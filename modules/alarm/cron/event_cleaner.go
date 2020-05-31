package cron

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	log "github.com/sirupsen/logrus"
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
			DeleteEventOlder(before, deleteBatch)
		}
	}()
}

// DeleteEventOlder 删除过时记录
func DeleteEventOlder(before time.Time, limit int) {
	sqlTpl := `DELETE FROM events WHERE create_at < ? LIMIT ?`
	db := g.Con()
	if err := db.Exec(sqlTpl, before.Unix(), limit).Error; err != nil {
		log.Errorf("[E] delete events fail, error: %v", err)
	} else {
		log.Debugf("[D] event older than %v deleted", before.Unix())
	}
}
