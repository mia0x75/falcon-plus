package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
)

// ReadClusterMonitorItems TODO:
func ReadClusterMonitorItems() (M map[string]*g.Cluster, err error) {
	M = make(map[string]*g.Cluster)
	sql := "SELECT `id`, `group_id`, `numerator`, `denominator`, `endpoint`, `metric`, `tags`, `ds_type`, `step`, `update_at` FROM `clusters`"

	cfg := g.Config()
	ids := cfg.Database.Ids
	if len(ids) != 2 {
		log.Fatal("[F] ids configuration error")
	}

	if ids[0] != -1 && ids[1] != -1 {
		sql = fmt.Sprintf("%s WHERE `id` >= %d and `id` <= %d", sql, ids[0], ids[1])
	} else {
		if ids[0] != -1 {
			sql = fmt.Sprintf("%s WHERE `id` >= %d", sql, ids[0])
		}

		if ids[1] != -1 {
			sql = fmt.Sprintf("%s WHERE `id` <= %d", sql, ids[1])
		}
	}

	log.Debugf("[D] %s", sql)

	rows, err := DB.Query(sql)
	if err != nil {
		log.Errorf("[E] %v", err)
		return M, err
	}

	defer rows.Close()
	for rows.Next() {
		var c g.Cluster
		err = rows.Scan(&c.ID, &c.GroupID, &c.Numerator, &c.Denominator, &c.Endpoint, &c.Metric, &c.Tags, &c.DsType, &c.Step, &c.LastUpdate)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		M[fmt.Sprintf("%d%v", c.ID, c.LastUpdate)] = &c
	}

	return M, err
}
