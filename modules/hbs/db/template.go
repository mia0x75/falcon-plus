package db

import (
	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// QueryGroupTemplates TODO:
func QueryGroupTemplates() (map[int][]int, error) {
	m := make(map[int][]int)

	q := "SELECT ancestor_id, descendant_id FROM edges WHERE type = 3"
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var gid, tid int
		err = rows.Scan(&gid, &tid)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		if _, exists := m[gid]; exists {
			m[gid] = append(m[gid], tid)
		} else {
			m[gid] = []int{tid}
		}
	}

	return m, nil
}

// QueryTemplates 获取所有的策略模板列表
func QueryTemplates() (map[int]*cm.Template, error) {
	templates := make(map[int]*cm.Template)

	q := "select id, name, parent_id, action_id, creator from templates"
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return templates, err
	}

	defer rows.Close()
	for rows.Next() {
		t := cm.Template{}
		err = rows.Scan(&t.ID, &t.Name, &t.ParentID, &t.ActionID, &t.Creator)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}
		templates[t.ID] = &t
	}

	return templates, nil
}

// QueryHostTemplateIDs 一个机器ID对应了多个模板ID
func QueryHostTemplateIDs() (map[int][]int, error) {
	ret := make(map[int][]int)
	q := "SELECT a.descendant_id, b.descendant_id FROM edges a INNER JOIN edges b ON a.ancestor_id = b.ancestor_id AND a.type = 3 AND b.type = 2"
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		var tid, hid int

		err = rows.Scan(&tid, &hid)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		if _, ok := ret[hid]; ok {
			ret[hid] = append(ret[hid], tid)
		} else {
			ret[hid] = []int{tid}
		}
	}

	return ret, nil
}
