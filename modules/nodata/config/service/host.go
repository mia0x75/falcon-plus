package service

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// FIX ME: too many JOIN
func GetHostsFromGroup(grpName string) map[string]int {
	hosts := make(map[string]int)

	now := time.Now().Unix()
	q := fmt.Sprintf("SELECT host.id, host.hostname FROM grp_host AS gh "+
		" INNER JOIN host ON host.id=gh.host_id AND (host.maintain_begin > %d OR host.maintain_end < %d)"+
		" INNER JOIN grp ON grp.id=gh.grp_id AND grp.grp_name='%s'", now, now, grpName)

	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] %v", err)
		return hosts
	}

	defer rows.Close()
	for rows.Next() {
		hid := -1
		hostname := ""
		err = rows.Scan(&hid, &hostname)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}
		if hid < 0 || hostname == "" {
			continue
		}

		hosts[hostname] = hid
	}

	return hosts
}
