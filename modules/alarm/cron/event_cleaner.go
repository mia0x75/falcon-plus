package cron

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	eventmodel "github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

// CleanExpiredEvent 清理过期的事件
func CleanExpiredEvent() {
	go func() {
		d := time.Duration(10) * time.Minute
		for range time.Tick(d) {
			retentionDays := g.Config().Housekeeper.EventRetentionDays
			deleteBatch := g.Config().Housekeeper.EventDeleteBatch

			now := time.Now()
			before := now.Add(time.Duration(-retentionDays*24) * time.Hour)
			eventmodel.DeleteEventOlder(before, deleteBatch)
		}
	}()
}
