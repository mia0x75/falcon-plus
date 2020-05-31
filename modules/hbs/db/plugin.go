package db

import (
	log "github.com/sirupsen/logrus"
)

// QueryPlugins TODO:
func QueryPlugins() (map[int][]string, error) {
	m := make(map[int][]string)

	q := "SELECT group_id, dir FROM plugin_dir"
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id  int
			dir string
		)

		err = rows.Scan(&id, &dir)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		if _, exists := m[id]; exists {
			m[id] = append(m[id], dir)
		} else {
			m[id] = []string{dir}
		}
	}

	return m, nil
}
