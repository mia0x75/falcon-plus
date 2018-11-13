package funcs

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/open-falcon/falcon-plus/modules/agent/hbs"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/radovskyb/watcher"
)

var (
	w      *watcher.Watcher
	files  map[string]int
	lock   = new(sync.RWMutex)
	events []*event
	once   sync.Once
)

type event struct {
	Op   string
	Path string
}

func Setup() {
	w = watcher.New()
	w.IgnoreHiddenFiles(true)
	d := time.Duration(g.COLLECT_INTERVAL) * time.Second
	w.Start(d)

	go func() {
		for {
			select {
			case e := <-w.Event:
				if !e.IsDir() {
					lock.Lock()
					events = append(events, &event{
						Op:   strings.ToLower(fmt.Sprintf("%s", e.Op)),
						Path: e.Path,
					})
					lock.Unlock()
				}
			case err := <-w.Error:
				log.Errorf("[E] incron error: %v", err)
			case <-w.Closed:
				return
			}
		}
	}()
}

func UpdateInconStats() {
	once.Do(Setup)

	files := hbs.ReportFiles()

	if files == nil || len(files) == 0 {
		for index, fi := range w.WatchedFiles() {
			if fi.IsDir() {
				w.RemoveRecursive(index)
			} else {
				w.Remove(index)
			}
		}
		return
	}

	for file, recursive := range files {
		if _, found := w.WatchedFiles()[file]; !found {
			if recursive > 0 {
				w.AddRecursive(file)
			} else {
				w.Add(file)
			}
		}
	}

	for index, fi := range w.WatchedFiles() {
		if _, found := files[index]; !found {
			if fi.IsDir() {
				w.RemoveRecursive(index)
			} else {
				w.Remove(index)
			}
		}
	}
}

func FilesMetrics() (L []*cmodel.MetricValue) {
	lock.Lock()
	defer lock.Unlock()
	if len(events) == 0 {
		return
	}
	for _, e := range events {
		tags := fmt.Sprintf("file=%s,op=%s", e.Path, e.Op)
		L = append(L, GaugeValue(g.FS_FILE_CHECKSUM, 1, tags))
	}
	events = make([]*event, 0)
	return
}
