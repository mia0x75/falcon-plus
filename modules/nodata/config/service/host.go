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
	q := fmt.Sprintf("SELECT h.id, h.hostname FROM edges l "+
		" INNER JOIN hosts h ON h.id = l.descendant_id "+
		" INNER JOIN groups g ON g.id = l.ancestor_id "+
		" WHERE (h.maintain_begin > %d OR h.maintain_end < %d) AND l.type = 2 AND g.name = '%s'", now, now, grpName)

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
