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
	Op     string
	File   string
	Source string
}

func Setup() {
	w = watcher.New()
	w.IgnoreHiddenFiles(true)
	d := time.Duration(g.COLLECT_INTERVAL) * time.Second

	go func() {
		for {
			select {
			case e := <-w.Event:
				if !e.IsDir() {
					lock.Lock()
					ie := &event{
						Op:     strings.ToLower(fmt.Sprintf("%s", e.Op)),
						Source: e.WatchFileInfo.Source,
						File:   e.Path,
					}
					events = append(events, ie)
					lock.Unlock()
					log.Debugf("[D] fs.file.checksum, inotify event: %v", ie)
				}
			case err := <-w.Error:
				log.Errorf("[E] incron error: %v", err)
			case <-w.Closed:
				return
			}
		}
	}()
	w.Add("/etc/my.cnf.d/")
	w.Add("/etc/my.cnf")
	w.Start(d)
}

func UpdateInconStats() {
	once.Do(Setup)

	files := hbs.ReportSources()

	if files == nil || len(files) == 0 {
		for index, _ := range w.WatchedFiles() {
			w.Remove(index)
		}
		return
	}

	for _, file := range files {
		if _, found := w.WatchedFiles()[file]; !found {
			w.Add(file)
		}
	}

	for key, file := range w.WatchedFiles() {
		if _, found := files[key]; !found {
			w.Remove(file.Source)
		}
	}
	log.Debugf("[D] watched files: %v", w.WatchedFiles())
}

func IncronMetrics() (L []*cmodel.MetricValue) {
	lock.Lock()
	defer lock.Unlock()
	if len(events) == 0 {
		return
	}
	for _, e := range events {
		tags := fmt.Sprintf("source=%s,file=%s,op=%s", e.Source, e.File, e.Op)
		L = append(L, GaugeValue(g.FS_FILE_CHECKSUM, 1, tags))
		log.Debugf("[D] fs.file.checksum, tags: %s", tags)
	}
	events = make([]*event, 0)
	return
}
