package cron

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	eventmodel "github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

func CleanExpiredEvent() {
	go func() {
		d := time.Duration(10) * time.Minute
		for range time.Tick(d) {
			retention_days := g.Config().Housekeeper.EventRetentionDays
			delete_batch := g.Config().Housekeeper.EventDeleteBatch

			now := time.Now()
			before := now.Add(time.Duration(-retention_days*24) * time.Hour)
			eventmodel.DeleteEventOlder(before, delete_batch)
		}
	}()
}
