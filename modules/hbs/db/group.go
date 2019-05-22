package db

import (
	log "github.com/sirupsen/logrus"
)

// QueryHostGroups TODO:
func QueryHostGroups() (map[int][]int, error) {
	m := make(map[int][]int)

	q := "select grp_id, host_id from grp_host"
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var gid, hid int
		err = rows.Scan(&gid, &hid)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		if _, exists := m[hid]; exists {
			m[hid] = append(m[hid], gid)
		} else {
			m[hid] = []int{gid}
		}
	}

	return m, nil
}
