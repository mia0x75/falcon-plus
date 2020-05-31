package db

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/container/set"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// QueryStrategies 获取所有的Strategy列表
func QueryStrategies(tpls map[int]*cm.Template) (map[int]*cm.Strategy, error) {
	ret := make(map[int]*cm.Strategy)

	if tpls == nil || len(tpls) == 0 {
		return ret, fmt.Errorf("illegal argument")
	}

	now := time.Now().Format("15:04")
	q := fmt.Sprintf(
		"SELECT %s FROM strategies as s WHERE (s.run_begin='' and s.run_end='') "+
			"OR (s.run_begin <= '%s' AND s.run_end >= '%s')"+
			"OR (s.run_begin > s.run_end AND !(s.run_begin > '%s' and s.run_end < '%s'))",
		"s.id, s.metric, s.tags, s.func, s.op, s.right_value, s.max_step, s.priority, s.note, s.template_id",
		now,
		now,
		now,
		now,
	)

	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		s := cm.Strategy{}
		var tags string
		var tid int
		err = rows.Scan(&s.ID, &s.Metric, &tags, &s.Func, &s.Operator, &s.RightValue, &s.MaxStep, &s.Priority, &s.Note, &tid)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		tt := make(map[string]string)

		if tags != "" {
			arr := strings.Split(tags, ",")
			for _, tag := range arr {
				kv := strings.SplitN(tag, "=", 2)
				if len(kv) != 2 {
					continue
				}
				tt[kv[0]] = kv[1]
			}
		}

		s.Tags = tt
		s.Template = tpls[tid]
		if s.Template == nil {
			log.Warnf("[W] tpl is nil. strategy id=%d, tpl id=%d", s.ID, tid)
			// 如果Strategy没有对应的Tpl，那就没有action，就没法报警，无需往后传递了
			continue
		}

		ret[s.ID] = &s
	}

	return ret, nil
}

// QueryBuiltinMetrics TODO:
func QueryBuiltinMetrics(tids string) ([]*cm.BuiltinMetric, error) {
	q := fmt.Sprintf(
		"SELECT metric, tags FROM strategies WHERE template_id IN (%s) AND metric IN ('net.port.listen', 'proc.num', 'du.bs', 'url.check.health', 'fs.file.checksum')",
		tids,
	)

	ret := []*cm.BuiltinMetric{}

	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return ret, err
	}

	metricTagsSet := set.NewStringSet()

	defer rows.Close()
	for rows.Next() {
		builtinMetric := cm.BuiltinMetric{}
		err = rows.Scan(&builtinMetric.Metric, &builtinMetric.Tags)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		k := fmt.Sprintf("%s%s", builtinMetric.Metric, builtinMetric.Tags)
		if metricTagsSet.Exists(k) {
			continue
		}

		ret = append(ret, &builtinMetric)
		metricTagsSet.Add(k)
	}

	return ret, nil
}
