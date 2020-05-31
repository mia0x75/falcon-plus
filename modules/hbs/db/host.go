package db

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// QueryHosts TODO:
func QueryHosts() (map[string]int, error) {
	m := make(map[string]int)

	q := "SELECT id, hostname FROM hosts"
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id       int
			hostname string
		)

		err = rows.Scan(&id, &hostname)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		m[hostname] = id
	}

	return m, nil
}

// QueryMonitoredHosts TODO:
func QueryMonitoredHosts() (map[int]*cm.Host, error) {
	hosts := make(map[int]*cm.Host)
	now := time.Now().Unix()
	q := fmt.Sprintf("SELECT id, hostname FROM hosts WHERE maintain_begin > %d or maintain_end < %d", now, now)
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return hosts, err
	}

	defer rows.Close()
	for rows.Next() {
		t := cm.Host{}
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			log.Warnf("[W] %v", err)
			continue
		}
		hosts[t.ID] = &t
	}

	return hosts, nil
}
