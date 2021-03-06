package index

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	tcache "github.com/toolkits/cache/localcache/timedcache"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/proc"
)

const (
	DefaultMaxCacheSize                     = 5000000 // 默认 最多500w个,太大了内存会耗尽
	DefaultCacheProcUpdateTaskSleepInterval = time.Duration(1) * time.Second
)

// item缓存
var (
	IndexedItemCache   = NewIndexCacheBase(DefaultMaxCacheSize)
	unIndexedItemCache = NewIndexCacheBase(DefaultMaxCacheSize)
)

// db本地缓存
var (
	// endpoints表的内存缓存, key:endpoint(string) / value:id(int64)
	dbEndpointCache = tcache.New(600*time.Second, 60*time.Second)
	// counters表的内存缓存, key:endpoint_id-counter(string) / val:dstype-step(string)
	dbEndpointCounterCache = tcache.New(600*time.Second, 60*time.Second)
)

// 初始化cache
func InitCache() {
	go startCacheProcUpdateTask()
}

// USED WHEN QUERY
func GetTypeAndStep(endpoint string, counter string) (dsType string, step int, found bool) {
	// get it from index cache
	pk := cu.Md5(fmt.Sprintf("%s/%s", endpoint, counter))
	if icitem := IndexedItemCache.Get(pk); icitem != nil {
		if item := icitem.(*IndexCacheItem).Item; item != nil {
			dsType = item.DsType
			step = item.Step
			found = true
			return
		}
	}

	// statistics
	proc.GraphLoadDbCnt.Incr()

	// get it from db, this should rarely happen
	var endpointId int64 = -1
	if endpointId, found = GetEndpointFromCache(endpoint); found {
		if dsType, step, found = GetCounterFromCache(endpointId, counter); found {
			// found = true
			return
		}
	}

	// do not find it, this must be a bad request
	found = false
	return
}

// GetEndpointFromCache returns EndpointId if found
func GetEndpointFromCache(endpoint string) (int64, bool) {
	// get from cache
	endpointId, found := dbEndpointCache.Get(endpoint)
	if found {
		return endpointId.(int64), true
	}

	// get from db
	var id int64 = -1
	err := g.DB.QueryRow("SELECT id FROM endpoints WHERE endpoint = ?", endpoint).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("[E] query endpoint id fail: %v", err)
		return -1, false
	}

	if err == sql.ErrNoRows || id < 0 {
		return -1, false
	}

	// update cache
	dbEndpointCache.Set(endpoint, id, 0)

	return id, true
}

// GetCounterFromCache returns DsType step if found
func GetCounterFromCache(endpointId int64, counter string) (dsType string, step int, found bool) {
	var err error
	// get from cache
	key := fmt.Sprintf("%d-%s", endpointId, counter)
	dsTypeStep, found := dbEndpointCounterCache.Get(key)
	if found {
		arr := strings.Split(dsTypeStep.(string), "_")
		step, err = strconv.Atoi(arr[1])
		if err != nil {
			found = false
			return
		}
		dsType = arr[0]
		return
	}

	// get from db
	err = g.DB.QueryRow("SELECT type, step FROM counters WHERE endpoint_id = ? and counter = ?",
		endpointId, counter).Scan(&dsType, &step)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("[E] query type and step fail: %v", err)
		return
	}

	if err == sql.ErrNoRows {
		return
	}

	// update cache
	dbEndpointCounterCache.Set(key, fmt.Sprintf("%s_%d", dsType, step), 0)

	found = true
	return
}

// 更新 cache的统计信息
func startCacheProcUpdateTask() {
	for range time.Tick(DefaultCacheProcUpdateTaskSleepInterval) {
		proc.IndexedItemCacheCnt.SetCnt(int64(IndexedItemCache.Size()))
		proc.UnIndexedItemCacheCnt.SetCnt(int64(unIndexedItemCache.Size()))
		proc.EndpointCacheCnt.SetCnt(int64(dbEndpointCache.Size()))
		proc.CounterCacheCnt.SetCnt(int64(dbEndpointCounterCache.Size()))
	}
}

// IndexCacheItem 索引缓存的元素数据结构
type IndexCacheItem struct {
	UUID string
	Item *cm.GraphItem
}

func NewIndexCacheItem(uuid string, item *cm.GraphItem) *IndexCacheItem {
	return &IndexCacheItem{UUID: uuid, Item: item}
}

// 索引缓存-基本缓存容器
type IndexCacheBase struct {
	sync.RWMutex
	maxSize int
	data    map[string]interface{}
}

func NewIndexCacheBase(max int) *IndexCacheBase {
	return &IndexCacheBase{maxSize: max, data: make(map[string]interface{})}
}

func (m *IndexCacheBase) GetMaxSize() int {
	return m.maxSize
}

func (m *IndexCacheBase) Put(key string, item interface{}) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = item
}

func (m *IndexCacheBase) Remove(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.data, key)
}

func (m *IndexCacheBase) Get(key string) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.data[key]
}

func (m *IndexCacheBase) ContainsKey(key string) bool {
	m.RLock()
	defer m.RUnlock()
	return m.data[key] != nil
}

func (m *IndexCacheBase) Size() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.data)
}

func (m *IndexCacheBase) Keys() []string {
	m.RLock()
	defer m.RUnlock()

	count := len(m.data)
	if count == 0 {
		return []string{}
	}

	keys := make([]string, 0, count)
	for key := range m.data {
		keys = append(keys, key)
	}

	return keys
}
