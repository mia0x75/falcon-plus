package db

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

// QueryHosts TODO:
func QueryHosts() (map[string]int, error) {
	m := make(map[string]int)

	q := "select id, hostname from host"
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
func QueryMonitoredHosts() (map[int]*cmodel.Host, error) {
	hosts := make(map[int]*cmodel.Host)
	now := time.Now().Unix()
	q := fmt.Sprintf("select id, hostname from host where maintain_begin > %d or maintain_end < %d", now, now)
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return hosts, err
	}

	defer rows.Close()
	for rows.Next() {
		t := cmodel.Host{}
		err = rows.Scan(&t.Id, &t.Name)
		if err != nil {
			log.Warnf("[W] %v", err)
			continue
		}
		hosts[t.Id] = &t
	}

	return hosts, nil
}
