package index

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	nsema "github.com/toolkits/concurrent/semaphore"
	ntime "github.com/toolkits/time"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	proc "github.com/open-falcon/falcon-plus/modules/graph/proc"
)

// 常量定义
const (
	DefaultUpdateStepInSec     = 2 * 24 * 3600 // 更新步长,一定不能大于删除步长. 两天内的数据,都可以用来建立索引
	ConcurrentOfUpdateIndexAll = 1
)

var (
	semaIndexUpdateAllTask = nsema.NewSemaphore(ConcurrentOfUpdateIndexAll) // 全量同步任务 并发控制器
	semaIndexUpdateAll     = nsema.NewSemaphore(4)                          // 索引全量更新时的mysql操作并发控制
)

// GetConcurrentOfUpdateIndexAll 索引全量更新的当前并行数
func GetConcurrentOfUpdateIndexAll() int {
	return ConcurrentOfUpdateIndexAll - semaIndexUpdateAllTask.AvailablePermits()
}

// UpdateIndexAllByDefaultStep 索引的全量更新
func UpdateIndexAllByDefaultStep() {
	UpdateIndexAll(DefaultUpdateStepInSec)
}

// UpdateIndexAll 索引的全量更新
func UpdateIndexAll(updateStepInSec int64) {
	// 减少任务积压,但高并发时可能无效(AvailablePermits不是线程安全的)
	if semaIndexUpdateAllTask.AvailablePermits() <= 0 {
		log.Info("[I] updateIndexAll, concurrent not available")
		return
	}

	semaIndexUpdateAllTask.Acquire()
	defer semaIndexUpdateAllTask.Release()

	startTs := time.Now().Unix()
	cnt := updateIndexAll(updateStepInSec)
	endTs := time.Now().Unix()
	log.Infof("[I] UpdateIndexAll, lastStartTs %s, updateStepInSec %d, lastTimeConsumingInSec %d\n",
		ntime.FormatTs(startTs), updateStepInSec, endTs-startTs)

	// statistics
	proc.IndexUpdateAllCnt.SetCnt(int64(cnt))
	proc.IndexUpdateAll.Incr()
	proc.IndexUpdateAll.PutOther("lastStartTs", ntime.FormatTs(startTs))
	proc.IndexUpdateAll.PutOther("updateStepInSec", updateStepInSec)
	proc.IndexUpdateAll.PutOther("lastTimeConsumingInSec", endTs-startTs)
	proc.IndexUpdateAll.PutOther("updateCnt", cnt)
}

// UpdateIndexOne 更新一条监控数据对应的索引. 用于手动添加索引,一般情况下不会使用
func UpdateIndexOne(endpoint string, metric string, tags map[string]string, dstype string, step int) error {
	itemDemo := &cm.GraphItem{
		Endpoint: endpoint,
		Metric:   metric,
		Tags:     tags,
		DsType:   dstype,
		Step:     step,
	}
	md5 := itemDemo.Checksum()
	uuid := itemDemo.UUID()

	cached := IndexedItemCache.Get(md5)
	if cached == nil {
		return fmt.Errorf("not found")
	}

	icitem := cached.(*IndexCacheItem)
	if icitem.UUID != uuid {
		return fmt.Errorf("bad type or step")
	}
	gitem := icitem.Item

	return updateIndexFromOneItem(gitem, g.DB)
}

func updateIndexAll(updateStepInSec int64) int {
	var ret int = 0
	if IndexedItemCache == nil || IndexedItemCache.Size() <= 0 {
		return ret
	}

	// lastTs for update index
	ts := time.Now().Unix()
	lastTs := ts - updateStepInSec

	keys := IndexedItemCache.Keys()
	for _, key := range keys {
		icitem := IndexedItemCache.Get(key)
		if icitem == nil {
			continue
		}

		gitem := icitem.(*IndexCacheItem).Item
		if gitem.Timestamp < lastTs { // 缓存中的数据太旧了,不能用于索引的全量更新
			IndexedItemCache.Remove(key) // 在这里做个删除,有点恶心
			continue
		}
		// 并发写mysql
		semaIndexUpdateAll.Acquire()
		go func(gitem *cm.GraphItem, db *sql.DB) {
			defer semaIndexUpdateAll.Release()
			err := updateIndexFromOneItem(gitem, db)
			if err != nil {
				proc.IndexUpdateAllErrorCnt.Incr()
			}
		}(gitem, g.DB)

		ret++
	}

	return ret
}

// 根据item,更新mysql
func updateIndexFromOneItem(item *cm.GraphItem, conn *sql.DB) error {
	if item == nil {
		return nil
	}

	ts := item.Timestamp
	var endpointId int64 = -1

	// endpoint表
	sqlStr := `INSERT INTO endpoints(endpoint, ts, create_at)
		VALUES (?, ?, UNIX_TIMESTAMP())
		ON DUPLICATE KEY UPDATE ts=?, update_at = UNIX_TIMESTAMP()`

	_, err := conn.Exec(sqlStr, item.Endpoint, ts, ts)
	if err != nil {
		log.Errorf("[E] %v", err)
		return err
	}
	proc.IndexUpdateIncrDbEndpointInsertCnt.Incr()

	err = conn.QueryRow("SELECT id FROM endpoints WHERE endpoint = ?", item.Endpoint).Scan(&endpointId)
	if err != nil {
		log.Errorf("[E] %v", err)
		return err
	}
	if endpointId <= 0 {
		log.Errorf("[E] no such endpoint in db, endpoint = %s", item.Endpoint)
		return errors.New("no such endpoint")
	}

	// tag_endpoint表
	for tagKey, tagVal := range item.Tags {
		tag := fmt.Sprintf("%s=%s", tagKey, tagVal)
		sqlStr := `INSERT INTO tags(tag, endpoint_id, ts, create_at)
			VALUES (?, ?, ?, UNIX_TIMESTAMP())
			ON DUPLICATE KEY UPDATE ts=?, update_at=UNIX_TIMESTAMP()`

		_, err := conn.Exec(sqlStr, tag, endpointId, ts, ts)
		if err != nil {
			log.Errorf("[E] %v", err)
			return err
		}
		proc.IndexUpdateIncrDbTagEndpointInsertCnt.Incr()
	}

	// endpoint_counter表
	counter := item.Metric
	if len(item.Tags) > 0 {
		counter = fmt.Sprintf("%s/%s", counter, cu.SortedTags(item.Tags))
	}

	sqlStr = `INSERT INTO counters(endpoint_id,counter,step,type,ts,create_at)
		VALUES (?,?,?,?,?,UNIX_TIMESTAMP())
		ON DUPLICATE KEY UPDATE ts=?,step=?,type=?,update_at=UNIX_TIMESTAMP()`

	_, err = conn.Exec(sqlStr, endpointId, counter, item.Step, item.DsType, ts, ts, item.Step, item.DsType)
	if err != nil {
		log.Errorf("[E] %v", err)
		return err
	}
	proc.IndexUpdateIncrDbEndpointCounterInsertCnt.Incr()

	return nil
}
