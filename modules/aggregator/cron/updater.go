package cron

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/aggregator/db"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
)

func UpdateItems() {
	go func() {
		d := time.Duration(g.Config().Database.Interval) * time.Second
		for range time.Tick(d) {
			updateItems()
		}
	}()
}

func updateItems() {
	items, err := db.ReadClusterMonitorItems()
	if err != nil {
		return
	}

	deleteNoUseWorker(items)
	createWorkerIfNeed(items)
}
