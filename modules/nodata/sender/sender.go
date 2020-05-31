package sender

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	tsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/nmap"
	ttime "github.com/toolkits/time"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/nodata/g"
)

var (
	MockMap = nmap.NewSafeMap()
	sema    = tsema.NewSemaphore(1)
)

func AddMock(key string, endpoint string, metric string, tags string, ts int64, dstype string, step int64, value interface{}) {
	item := &cm.JSONMetaData{
		Metric:      metric,
		Endpoint:    endpoint,
		Timestamp:   ts,
		Step:        step,
		Value:       value,
		CounterType: dstype,
		Tags:        tags,
	}
	MockMap.Put(key, item)
}

func SendMockOnceAsync() {
	go SendMockOnce()
}

func SendMockOnce() int {
	if !sema.TryAcquire() {
		return -1
	}
	defer sema.Release()

	// not enabled
	if !g.Config().Transfer.Enabled {
		return 0
	}

	start := time.Now().Unix()
	cnt, _ := sendMock()
	end := time.Now().Unix()
	log.Debugf("[D] sender cron, cnt %d, time %ds, start %s", cnt, end-start, ttime.FormatTs(start))

	// Statistics
	g.SenderCronCnt.Incr()
	g.SenderLastTs.SetCnt(end - start)
	g.SenderCnt.IncrBy(int64(cnt))

	return cnt
}

func sendMock() (cnt int, errt error) {
	cfgTransfer := g.Config().Transfer
	batch := int(cfgTransfer.Batch)
	connTimeout := cfgTransfer.ConnectTimeout
	requTimeout := cfgTransfer.RequestTimeout

	// send mock to transfer
	mocks := MockMap.Slice()
	MockMap.Clear()
	mockSize := len(mocks)
	i := 0
	for i < mockSize {
		leftLen := mockSize - i
		sendSize := batch
		if leftLen < sendSize {
			sendSize = leftLen
		}
		fetchMocks := mocks[i : i+sendSize]
		i += sendSize

		items := make([]*cm.JSONMetaData, 0)
		for _, val := range fetchMocks {
			if val == nil {
				continue
			}
			items = append(items, val.(*cm.JSONMetaData))
		}
		cntonce, err := sendItemsToTransfer(items, len(items), "nodata.mock",
			time.Millisecond*time.Duration(connTimeout),
			time.Millisecond*time.Duration(requTimeout))
		if err == nil {
			log.Debugf("[D] send items: %v", items)
			cnt += cntonce
		}
	}

	return cnt, nil
}

//
func sendItemsToTransfer(items []*cm.JSONMetaData, size int, httpcliname string,
	connT, reqT time.Duration) (cnt int, err error) {
	if size < 1 {
		return
	}
	cnt = size

	cfg := g.Config()
	transUlr := fmt.Sprintf("http://%s/api/push", cfg.Transfer.Addr)

	// form request args
	var itemsBody []byte
	itemsBody, err = json.Marshal(items)
	if err != nil {
		log.Errorf("[E] %s, format body error: %v", transUlr, err)
		return
	}

	client := cu.NewHttp(transUlr)
	client.SetUserAgent(httpcliname)
	headers := map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		"Connection":   "close",
	}
	client.SetHeaders(headers)
	if _, err = client.Post(itemsBody); err != nil {
		log.Errorf("[E] %s, post to dest error: %v", transUlr, err)
		return
	}

	return
}
