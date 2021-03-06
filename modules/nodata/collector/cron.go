package collector

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	tsema "github.com/toolkits/concurrent/semaphore"
	tcron "github.com/toolkits/cron"
	"github.com/toolkits/net/httplib"
	ttime "github.com/toolkits/time"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cr "github.com/open-falcon/falcon-plus/common/sdk/requests"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/nodata/config"
	"github.com/open-falcon/falcon-plus/modules/nodata/g"
)

var (
	collectorCron = tcron.New()
)

// StartCollectorCron TODO:
func StartCollectorCron() {
	collectorCron.AddFuncCC("*/20 * * * * ?", func() {
		start := time.Now().Unix()
		cnt := collectDataOnce()
		end := time.Now().Unix()
		log.Debugf("[D] collect cron, cnt %d, time %ds, start %s\n", cnt, end-start, ttime.FormatTs(start))

		// Statistics
		g.CollectorCronCnt.Incr()
		g.CollectorLastTs.SetCnt(end - start)
		g.CollectorLastCnt.SetCnt(int64(cnt))
		g.CollectorCnt.IncrBy(int64(cnt))
	}, 1)
	collectorCron.Start()
}

// CollectDataOnce TODO:
func CollectDataOnce() int {
	return collectDataOnce()
}

func collectDataOnce() int {
	keys := config.Keys()
	keysLen := len(keys)

	// 并发+同步控制
	cfg := g.Config().Collector
	concurrent := int(cfg.Concurrent)
	if concurrent < 1 || concurrent > 50 {
		concurrent = 10
	}
	sema := tsema.NewSemaphore(concurrent)

	batch := int(cfg.Batch)
	if batch < 100 || batch > 1000 {
		batch = 200 // batch不能太小, 否则channel将会很大
	}

	batchCnt := (keysLen + batch - 1) / batch
	rch := make(chan int, batchCnt+1)

	i := 0
	for i < keysLen {
		leftLen := keysLen - i
		fetchSize := batch // 每次处理batch个配置
		if leftLen < fetchSize {
			fetchSize = leftLen
		}
		fetchKeys := keys[i : i+fetchSize]

		// 并发collect数据
		sema.Acquire()
		go func(keys []string, keySize int) {
			defer sema.Release()
			size, err := fetchItemsAndStore(keys, keySize)
			if err != nil {
				log.Errorf("[E] fetchItemAndStore fail, size: %v, error: %v", size, err)
			}
			log.Debugf("[D] fetchItemAndStore keys: %v, key_size: %v, ret_size: %v", keys, keySize, size)
			rch <- size
		}(fetchKeys, fetchSize)

		i += fetchSize
	}

	collectCnt := 0
	for i := 0; i < batchCnt; i++ {
		select {
		case cnt := <-rch:
			collectCnt += cnt
		}
	}

	return collectCnt
}

func fetchItemsAndStore(fetchKeys []string, fetchSize int) (size int, errt error) {
	if fetchSize < 1 {
		return
	}

	// form request args
	args := make([]*cm.GraphLastParam, 0)
	for _, key := range fetchKeys {
		ndcfg, found := config.GetNdConfig(key)
		if !found {
			continue
		}

		endpoint := ndcfg.Endpoint
		counter := cu.Counter(ndcfg.Metric, ndcfg.Tags)
		arg := &cm.GraphLastParam{
			Endpoint: endpoint,
			Counter:  counter,
		}
		args = append(args, arg)
	}
	if len(args) < 1 {
		return
	}

	resp, err := queryLastPoints(args)
	if err != nil {
		return 0, err
	}

	// store items
	fts := time.Now().Unix()
	for _, glr := range resp {
		log.Infof("[I] collect: %v", glr)
		if glr == nil || glr.Value == nil {
			continue
		}
		AddItem(cu.PK2(glr.Endpoint, glr.Counter), NewDataItem(glr.Value.Timestamp, float64(glr.Value.Value), "OK", fts))
	}

	return len(resp), nil
}

func queryLastPoints(param []*cm.GraphLastParam) (resp []*cm.GraphLastResp, err error) {
	cfg := g.Config()
	uri := fmt.Sprintf("%s/api/v1/graph/lastpoint", cfg.API.Addr)

	var req *httplib.BeegoHttpRequest
	headers := map[string]string{"Content-type": "application/json"}
	req, err = cr.CurlPlus(uri, "POST", "nodata", cfg.API.Token,
		headers, map[string]string{})

	if err != nil {
		return
	}
	req.SetTimeout(time.Duration(cfg.API.ConnectTimeout)*time.Millisecond,
		time.Duration(cfg.API.RequestTimeout)*time.Millisecond)

	b, err := json.Marshal(param)
	if err != nil {
		return
	}

	req.Body(b)

	err = req.ToJson(&resp)
	if err != nil {
		return
	}

	return resp, nil
}
