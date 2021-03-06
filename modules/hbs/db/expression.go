package db

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// QueryExpressions TODO:
func QueryExpressions() (ret []*cm.Expression, err error) {
	q := "SELECT id, expression, func, op, right_value, max_step, priority, note, action_id FROM expressions WHERE action_id > 0 AND pause = 0"
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		e := cm.Expression{}
		var exp string
		err = rows.Scan(
			&e.ID,
			&exp,
			&e.Func,
			&e.Operator,
			&e.RightValue,
			&e.MaxStep,
			&e.Priority,
			&e.Note,
			&e.ActionID,
		)

		if err != nil {
			log.Warnf("[W] row scan error:%v", err)
			continue
		}

		e.Metric, e.Tags, err = parseExpression(exp)
		if err != nil {
			log.Errorf("[E] parse expression error:%v", err)
			continue
		}

		ret = append(ret, &e)
	}

	return ret, nil
}

func parseExpression(exp string) (metric string, tags map[string]string, err error) {
	left := strings.Index(exp, "(")
	right := strings.Index(exp, ")")
	tagStrs := strings.TrimSpace(exp[left+1 : right])

	arr := strings.Fields(tagStrs)
	if len(arr) < 2 {
		err = fmt.Errorf("tag not enough. exp: %s", exp)
		return
	}

	tags = make(map[string]string)
	for _, item := range arr {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			err = fmt.Errorf("parse %s fail", exp)
			return
		}
		tags[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	metric, exists := tags["metric"]
	if !exists {
		err = fmt.Errorf("no metric give of %s", exp)
		return
	}

	delete(tags, "metric")
	return
}
