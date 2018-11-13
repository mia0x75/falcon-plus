package db

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

func QueryHosts() (map[string]int, error) {
	m := make(map[string]int)

	sql := "select id, hostname from host"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Errorf("[E] %v", err)
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

func QueryMonitoredHosts() (map[int]*cmodel.Host, error) {
	hosts := make(map[int]*cmodel.Host)
	now := time.Now().Unix()
	sql := fmt.Sprintf("select id, hostname from host where maintain_begin > %d or maintain_end < %d", now, now)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Errorf("[E] %v", err)
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
