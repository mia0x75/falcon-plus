package index

import (
	"time"

	log "github.com/sirupsen/logrus"
	cron "github.com/toolkits/cron"
	ntime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/modules/exporter/proc"
)

const (
	indexDeleteCronSpec = "0 0 2 ? * 6" // 每周6晚上22:00执行一次
	deteleStepInSec     = 7 * 24 * 3600 // 索引的最大生存周期, sec
)

var (
	indexDeleteCron = cron.New()
)

// 启动 索引全量更新 定时任务
func StartIndexDeleteTask() {
	indexDeleteCron.AddFuncCC(indexDeleteCronSpec, func() { DeleteIndex() }, 1)
	indexDeleteCron.Start()
}

// 索引的全量更新
func DeleteIndex() {
	startTs := time.Now().Unix()
	deleteIndex()
	endTs := time.Now().Unix()
	log.Infof("[I] deleteIndex, start %s, ts %ds", ntime.FormatTs(startTs), endTs-startTs)

	// statistics
	proc.IndexDeleteCnt.Incr()
}

// 先select 得到可能被删除的index的信息, 然后以相同的条件delete. select和delete不是原子操作,可能有一些不一致,但不影响正确性
func deleteIndex() error {
	ts := time.Now().Unix()
	lastTs := ts - deteleStepInSec
	log.Infof("[I] deleteIndex, lastTs %d", lastTs)

	// 复位 statistics
	proc.IndexDeleteCnt.PutOther("deleteCntEndpoint", 0)
	proc.IndexDeleteCnt.PutOther("deleteCntTagEndpoint", 0)
	proc.IndexDeleteCnt.PutOther("deleteCntEndpointCounter", 0)

	// endpoints 表
	{
		// Query
		rows, err := db.Query("SELECT count(*) as cnt FROM endpoints WHERE ts < ?", lastTs)
		if err != nil {
			log.Errorf("[E] %v", err)
			return err
		}

		cnt := 0
		if rows.Next() {
			err := rows.Scan(&cnt)
			if err != nil {
				log.Errorf("[E] %v", err)
				return err
			}
		}

		// Delete
		_, err = db.Exec("DELETE FROM endpoints WHERE ts < ?", lastTs)
		if err != nil {
			log.Errorf("[E] %v", err)
			return err
		}
		log.Infof("[I] delete endpoint, done, cnt %d", cnt)

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntEndpoint", cnt)
	}

	// tags表
	{
		// Query
		rows, err := db.Query("SELECT count(*) as cnt FROM tags WHERE ts < ?", lastTs)
		if err != nil {
			log.Errorf("[E] %v", err)
			return err
		}

		cnt := 0
		if rows.Next() {
			err := rows.Scan(&cnt)
			if err != nil {
				log.Errorf("[E] %v", err)
				return err
			}
		}

		// Delete
		_, err = db.Exec("DELETE FROM tags WHERE ts < ?", lastTs)
		if err != nil {
			log.Errorf("[E] %v", err)
			return err
		}
		log.Infof("[I] delete tag_endpoint, done, cnt %d", cnt)

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntTagEndpoint", cnt)
	}
	// counters表
	{
		// Query
		rows, err := db.Query("SELECT count(1) as cnt FROM counters WHERE ts < ?", lastTs)
		if err != nil {
			log.Errorf("[E] %v", err)
			return err
		}

		cnt := 0
		if rows.Next() {
			err := rows.Scan(&cnt)
			if err != nil {
				log.Errorf("[E] %v", err)
				return err
			}
		}

		// Delete
		_, err = db.Exec("DELETE FROM counters WHERE ts < ?", lastTs)
		if err != nil {
			log.Errorf("[E] %v", err)
			return err
		}
		log.Infof("[I] delete endpoint_counter, done, cnt %d", cnt)

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntEndpointCounter", cnt)
	}

	return nil
}
