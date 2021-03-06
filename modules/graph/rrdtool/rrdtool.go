package rrdtool

import (
	"errors"
	"math"
	"sync/atomic"
	"time"

	"github.com/open-falcon/rrdlite"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

var (
	disk_counter uint64
	net_counter  uint64
)

type fetch_t struct {
	filename string
	cf       string
	start    int64
	end      int64
	step     int
	data     []*cm.RRDData
}

type flushfile_t struct {
	filename string
	items    []*cm.GraphItem
}

type readfile_t struct {
	filename string
	data     []byte
}

func Start() {
	cfg := g.Config()
	var err error
	// check data dir
	if err = file.EnsureDirRW(cfg.RRD.Storage); err != nil {
		log.Fatalf("[F] rrdtool.Start error, bad data dir %s, error: %v", cfg.RRD.Storage, err)
	}

	migrate_start(cfg)

	// sync disk
	go syncDisk()
	go ioWorker()
	log.Info("[I] rrdtool.Start ok")
}

// RRA.Point.Size
const (
	RRA1PointCnt    = 8640 // 原始值 - 5s / 12h
	RRA6PointCnt    = 8640 // 采样值 - 30s / 3d
	RRA30PointCnt   = 4032 // 采样值 - 2m30s / 1w
	RRA180PointCnt  = 2880 // 采样值 - 15m / 1m
	RRA720PointCnt  = 2160 // 采样值 - 1h / 3m
	RRA4320PointCnt = 1460 // 采样值 - 6h / 1y
)

func create(filename string, item *cm.GraphItem) error {
	now := time.Now()
	start := now.Add(time.Duration(-24) * time.Hour)
	step := uint(item.Step)

	c := rrdlite.NewCreator(filename, start, step)
	c.DS("metric", item.DsType, item.Heartbeat, item.Min, item.Max)

	// 设置各种归档策略
	// 原始值
	c.RRA("AVERAGE", 0, 1, RRA1PointCnt)

	// 采样值
	c.RRA("AVERAGE", 0, 6, RRA6PointCnt)
	c.RRA("MAX", 0, 6, RRA6PointCnt)
	c.RRA("MIN", 0, 6, RRA6PointCnt)

	// 采样值
	c.RRA("AVERAGE", 0, 30, RRA30PointCnt)
	c.RRA("MAX", 0, 30, RRA30PointCnt)
	c.RRA("MIN", 0, 30, RRA30PointCnt)

	// 采样值
	c.RRA("AVERAGE", 0, 180, RRA180PointCnt)
	c.RRA("MAX", 0, 180, RRA180PointCnt)
	c.RRA("MIN", 0, 180, RRA180PointCnt)

	// 采样值
	c.RRA("AVERAGE", 0, 4320, RRA4320PointCnt)
	c.RRA("MAX", 0, 4320, RRA4320PointCnt)
	c.RRA("MIN", 0, 4320, RRA4320PointCnt)

	// 采样值
	c.RRA("AVERAGE", 0, 720, RRA720PointCnt)
	c.RRA("MAX", 0, 720, RRA720PointCnt)
	c.RRA("MIN", 0, 720, RRA720PointCnt)

	return c.Create(true)
}

func update(filename string, items []*cm.GraphItem) error {
	u := rrdlite.NewUpdater(filename)

	for _, item := range items {
		v := math.Abs(item.Value)
		if v > 1e+300 || (v < 1e-300 && v > 0) {
			continue
		}
		if item.DsType == "DERIVE" || item.DsType == "COUNTER" {
			u.Cache(item.Timestamp, int(item.Value))
		} else {
			u.Cache(item.Timestamp, item.Value)
		}
	}

	return u.Update()
}

// flush to disk from memory
// 最新的数据在列表的最后面
// TODO fix me, filename fmt from item[0], it's hard to keep consistent
func flushrrd(filename string, items []*cm.GraphItem) error {
	if items == nil || len(items) == 0 {
		return errors.New("empty items")
	}

	if !g.IsRrdFileExist(filename) {
		baseDir := file.Dir(filename)

		err := file.InsureDir(baseDir)
		if err != nil {
			return err
		}

		err = create(filename, items[0])
		if err != nil {
			return err
		}
	}

	return update(filename, items)
}

func ReadFile(filename, md5 string) ([]byte, error) {
	done := make(chan error, 1)
	task := &io_task_t{
		method: IO_TASK_M_READ,
		args:   &readfile_t{filename: filename},
		done:   done,
	}

	io_task_chans[getIndex(md5)] <- task
	err := <-done
	return task.args.(*readfile_t).data, err
}

func FlushFile(filename, md5 string, items []*cm.GraphItem) error {
	done := make(chan error, 1)
	io_task_chans[getIndex(md5)] <- &io_task_t{
		method: IO_TASK_M_FLUSH,
		args: &flushfile_t{
			filename: filename,
			items:    items,
		},
		done: done,
	}
	atomic.AddUint64(&disk_counter, 1)
	return <-done
}

func Fetch(filename string, md5 string, cf string, start, end int64, step int) ([]*cm.RRDData, error) {
	done := make(chan error, 1)
	task := &io_task_t{
		method: IO_TASK_M_FETCH,
		args: &fetch_t{
			filename: filename,
			cf:       cf,
			start:    start,
			end:      end,
			step:     step,
		},
		done: done,
	}
	io_task_chans[getIndex(md5)] <- task
	err := <-done
	return task.args.(*fetch_t).data, err
}

func fetch(filename string, cf string, start, end int64, step int) ([]*cm.RRDData, error) {
	start_t := time.Unix(start, 0)
	end_t := time.Unix(end, 0)
	step_t := time.Duration(step) * time.Second

	fetchRes, err := rrdlite.Fetch(filename, cf, start_t, end_t, step_t)
	if err != nil {
		return []*cm.RRDData{}, err
	}

	defer fetchRes.FreeValues()

	values := fetchRes.Values()
	size := len(values)
	ret := make([]*cm.RRDData, size)

	start_ts := fetchRes.Start.Unix()
	step_s := fetchRes.Step.Seconds()

	for i, val := range values {
		ts := start_ts + int64(i+1)*int64(step_s)
		d := &cm.RRDData{
			Timestamp: ts,
			Value:     cm.JSONFloat(val),
		}
		ret[i] = d
	}

	return ret, nil
}

func FlushAll(force bool) {
	n := store.GraphItems.Size / 10
	for i := 0; i < store.GraphItems.Size; i++ {
		FlushRRD(i, force)
		if i%n == 0 {
			log.Debugf("[D] flush hash idx: %03d size: %03d disk: %08d net: %08d\n",
				i, store.GraphItems.Size, disk_counter, net_counter)
		}
	}
	log.Infof("[I] flush hash done (disk: %08d net: %08d)", disk_counter, net_counter)
}

func CommitByKey(key string) {
	md5, dsType, step, err := g.SplitRrdCacheKey(key)
	if err != nil {
		return
	}
	filename := g.RrdFileName(g.Config().RRD.Storage, md5, dsType, step)

	items := store.GraphItems.PopAll(key)
	if len(items) == 0 {
		return
	}
	FlushFile(filename, md5, items)
}

func PullByKey(key string) {
	done := make(chan error)

	item := store.GraphItems.First(key)
	if item == nil {
		return
	}
	node, err := Consistent.Get(item.PrimaryKey())
	if err != nil {
		return
	}
	Net_task_ch[node] <- &Net_task_t{
		Method: NET_TASK_M_PULL,
		Key:    key,
		Done:   done,
	}
	// net_task slow, shouldn't block syncDisk() or FlushAll()
	// warning: recev sigout when migrating, maybe lost memory data
	go func() {
		err := <-done
		if err != nil {
			log.Errorf("[E] get %s from remote error: %v", key, err)
			return
		}
		atomic.AddUint64(&net_counter, 1)
		// TODO: flushfile after getfile? not yet
	}()
}

func FlushRRD(idx int, force bool) {
	begin := time.Now()
	atomic.StoreInt32(&flushrrd_timeout, 0)

	keys := store.GraphItems.KeysByIndex(idx)
	if len(keys) == 0 {
		return
	}

	for _, key := range keys {
		flag, _ := store.GraphItems.GetFlag(key)

		// write err data to local filename
		if force == false && g.Config().Migrate.Enabled && flag&g.GRAPH_F_MISS != 0 {
			if time.Since(begin) > time.Millisecond*g.FLUSH_DISK_STEP {
				atomic.StoreInt32(&flushrrd_timeout, 1)
			}
			PullByKey(key)
		} else if force || shouldFlush(key) {
			CommitByKey(key)
		}
	}
}

func shouldFlush(key string) bool {

	if store.GraphItems.ItemCnt(key) >= g.FLUSH_MIN_COUNT {
		return true
	}

	deadline := time.Now().Unix() - int64(g.FLUSH_MAX_WAIT)
	back := store.GraphItems.Back(key)
	if back != nil && back.Timestamp <= deadline {
		return true
	}

	return false
}
