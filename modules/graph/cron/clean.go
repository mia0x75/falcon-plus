package cron

import (
	"strings"
	"time"

	pfc "github.com/mia0x75/gopfc/metric"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

// CleanCache TODO:
func CleanCache() {
	go clean()
}

func clean() {
	var ticker *time.Ticker
	// TODO: Move g.CLEAN_CACHE to cfg
	ticker = time.NewTicker(time.Duration(g.CLEAN_CACHE) * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		DeleteInvalidItems()   // 删除无效的GraphItems
		DeleteInvalidHistory() // 删除无效的HistoryCache
	}
}

/*

  概念定义及结构体简谱:
  ckey = md5_type_step
  uuid = endpoint/metric/tags/dstype/step
  md5  = md5(endpoint/metric/tags)

  GraphItems        [idx]  [ckey] [{timestamp, value}, {timestamp, value} ...]
  HistoryCache      [md5]  [itemFirst, itemSecond]
  IndexedItemCache  [md5]  {UUID, Item}

*/

// DeleteInvalidItems 删除长期不更新数据(依赖index)
func DeleteInvalidItems() int {
	var currentCnt, deleteCnt int
	graphItems := store.GraphItems

	for idx := 0; idx < graphItems.Size; idx++ {
		keys := graphItems.KeysByIndex(idx)

		for _, key := range keys {
			tmp := strings.Split(key, "_") // key = md5_type_step
			if len(tmp) == 3 && !index.IndexedItemCache.ContainsKey(tmp[0]) {
				graphItems.Remove(key)
				deleteCnt++
			}
		}
	}
	currentCnt = graphItems.Len()

	pfc.Gauge("cache.GraphItemsCacheCnt", int64(currentCnt))
	pfc.Gauge("cache.GraphItemsCacheInvalidCnt", int64(deleteCnt))
	log.Infof("[I] GraphItemsCache: Count=>%d, DeleteInvalid=>%d", currentCnt, deleteCnt)

	return deleteCnt
}

// DeleteInvalidHistory 删除长期不更新数据(依赖index)
func DeleteInvalidHistory() int {
	var currentCnt, deleteCnt int
	historyCache := store.HistoryCache

	keys := historyCache.Keys()
	for _, key := range keys {
		if !index.IndexedItemCache.ContainsKey(key) {
			historyCache.Remove(key)
			deleteCnt++
		}
	}
	currentCnt = historyCache.Size()

	pfc.Gauge("cache.HistoryCacheCnt", int64(currentCnt))
	pfc.Gauge("cache.HistoryCacheInvalidCnt", int64(deleteCnt))
	log.Infof("[I] HistoryCache: Count=>%d, DeleteInvalid=>%d", currentCnt, deleteCnt)

	return deleteCnt
}
