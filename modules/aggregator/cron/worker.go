package cron

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
)

type Worker struct {
	Ticker      *time.Ticker
	ClusterItem *g.Cluster
	Quit        chan struct{}
}

func NewWorker(ci *g.Cluster) Worker {
	w := Worker{}
	w.Ticker = time.NewTicker(time.Duration(ci.Step) * time.Second)
	w.Quit = make(chan struct{})
	w.ClusterItem = ci
	return w
}

func (this Worker) Start() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				WorkerRun(this.ClusterItem)
			case <-this.Quit:
				log.Debugf("[D] drop worker %v", this.ClusterItem)
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this Worker) Drop() {
	close(this.Quit)
}

var Workers = make(map[string]Worker)

func deleteNoUseWorker(m map[string]*g.Cluster) {
	del := []string{}
	for key, worker := range Workers {
		if _, ok := m[key]; !ok {
			worker.Drop()
			del = append(del, key)
		}
	}

	for _, key := range del {
		delete(Workers, key)
	}
}

func createWorkerIfNeed(m map[string]*g.Cluster) {
	for key, item := range m {
		if _, ok := Workers[key]; !ok {
			if item.Step <= 0 {
				log.Warnf("[W] invalid cluster(step <= 0): %v", item)
				continue
			}
			worker := NewWorker(item)
			Workers[key] = worker
			worker.Start()
		}
	}
}
