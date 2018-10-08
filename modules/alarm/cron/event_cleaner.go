package cron

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	eventmodel "github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

func CleanExpiredEvent() {
	for {

		retention_days := g.Config().Housekeeper.EventRetentionDays
		delete_batch := g.Config().Housekeeper.EventDeleteBatch

		now := time.Now()
		before := now.Add(time.Duration(-retention_days*24) * time.Hour)
		eventmodel.DeleteEventOlder(before, delete_batch)

		time.Sleep(time.Minute * 10)
	}
}
