package collector

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/sdk/requests"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/nodata/config"
	"github.com/open-falcon/falcon-plus/modules/nodata/g"
	tsema "github.com/toolkits/concurrent/semaphore"
	tcron "github.com/toolkits/cron"
	"github.com/toolkits/net/httplib"
	ttime "github.com/toolkits/time"
)

var (
	collectorCron = tcron.New()
)

func StartCollectorCron() {
	collectorCron.AddFuncCC("*/20 * * * * ?", func() {
		start := time.Now().Unix()
		cnt := collectDataOnce()
		end := time.Now().Unix()
		log.Debugf("collect cron, cnt %d, time %ds, start %s\n", cnt, end-start, ttime.FormatTs(start))

		// statistics
		g.CollectorCronCnt.Incr()
		g.CollectorLastTs.SetCnt(end - start)
		g.CollectorLastCnt.SetCnt(int64(cnt))
		g.CollectorCnt.IncrBy(int64(cnt))
	}, 1)
	collectorCron.Start()
}

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
		batch = 200 //batch不能太小, 否则channel将会很大
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
				log.Printf("fetchItemAndStore fail, size:%v, error:%v", size, err)
			}
			log.Debugf("fetchItemAndStore keys:%v, key_size:%v, ret_size:%v", keys, keySize, size)
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
	args := make([]*cmodel.GraphLastParam, 0)
	for _, key := range fetchKeys {
		ndcfg, found := config.GetNdConfig(key)
		if !found {
			continue
		}

		endpoint := ndcfg.Endpoint
		counter := cutils.Counter(ndcfg.Metric, ndcfg.Tags)
		arg := &cmodel.GraphLastParam{endpoint, counter}
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
		//log.Printf("collect:%v\n", glr)
		if glr == nil || glr.Value == nil {
			continue
		}
		AddItem(cutils.PK2(glr.Endpoint, glr.Counter), NewDataItem(glr.Value.Timestamp, float64(glr.Value.Value), "OK", fts))
	}

	return len(resp), nil
}

func queryLastPoints(param []*cmodel.GraphLastParam) (resp []*cmodel.GraphLastResp, err error) {
	cfg := g.Config()
	uri := fmt.Sprintf("%s/api/v1/graph/lastpoint", cfg.Api.Addr)

	var req *httplib.BeegoHttpRequest
	headers := map[string]string{"Content-type": "application/json"}
	req, err = requests.CurlPlus(uri, "POST", "nodata", cfg.Api.Token,
		headers, map[string]string{})

	if err != nil {
		return
	}
	req.SetTimeout(time.Duration(cfg.Api.ConnectTimeout)*time.Millisecond,
		time.Duration(cfg.Api.RequestTimeout)*time.Millisecond)

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