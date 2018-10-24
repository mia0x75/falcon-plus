package event

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
)

const timeLayout = "2006-01-02 15:04:05"

func insertEvent(q orm.Ormer, eve *cmodel.Event) (res sql.Result, err error) {
	var status int
	if status = 0; eve.Status == "OK" {
		status = 1
	}
	sqltemplete := `INSERT INTO events (
		event_caseId,
		step,
		cond,
		status,
		timestamp
	) VALUES(?,?,?,?,?)`
	res, err = q.Raw(
		sqltemplete,
		eve.Id,
		eve.CurrentStep,
		fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
		status,
		time.Unix(eve.EventTime, 0).Format(timeLayout),
	).Exec()

	if err != nil {
		log.Errorf("insert event to db fail, error:%v", err)
	} else {
		lastid, _ := res.LastInsertId()
		log.Debug("insert event to db succ, last_insert_id:", lastid)
	}
	return
}

func InsertEvent(eve *cmodel.Event) {
	q := orm.NewOrm()
	db, err := orm.GetDB()
	if err != nil {
		log.Errorf("get database connection error:%v", err)
	} else {
		db.Ping()
	}
	var event []EventCases
	q.Raw("select * from event_cases where id = ?", eve.Id).QueryRows(&event)
	var sqlLog sql.Result
	var errRes error
	log.Debugf("events: %v", eve)
	log.Debugf("expression is null: %v", eve.Expression == nil)
	if len(event) == 0 {
		//create cases
		sqltemplete := `INSERT INTO event_cases (
					id,
					endpoint,
					metric,
					func,
					cond,
					note,
					max_step,
					current_step,
					priority,
					status,
					timestamp,
					update_at,
					tpl_creator,
					expression_id,
					strategy_id,
					template_id
					) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

		tpl_creator := ""
		if eve.Tpl() != nil {
			tpl_creator = eve.Tpl().Creator
		}
		sqlLog, errRes = q.Raw(
			sqltemplete,
			eve.Id,
			eve.Endpoint,
			counterGen(eve.Metric(), cutils.SortedTags(eve.PushedTags)),
			eve.Func(),
			//cond
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Note(),
			eve.MaxStep(),
			eve.CurrentStep,
			eve.Priority(),
			eve.Status,
			//start_at
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			//update_at
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			tpl_creator,
			eve.ExpressionId(),
			eve.StrategyId(),
			//template_id
			eve.TplId()).Exec()

	} else {
		sqltemplete := `UPDATE event_cases SET
				update_at = ?,
				max_step = ?,
				current_step = ?,
				note = ?,
				cond = ?,
				status = ?,
				func = ?,
				priority = ?,
				tpl_creator = ?,
				expression_id = ?,
				strategy_id = ?,
				template_id = ?`
		//reopen case
		if event[0].ProcessStatus == "resolved" || event[0].ProcessStatus == "ignored" {
			sqltemplete = fmt.Sprintf("%v ,process_status = '%s', process_note = %d", sqltemplete, "unresolved", 0)
		}

		tpl_creator := ""
		if eve.Tpl() != nil {
			tpl_creator = eve.Tpl().Creator
		}
		if eve.CurrentStep == 1 {
			//update start time of cases
			sqltemplete = fmt.Sprintf("%v ,timestamp = ? WHERE id = ?", sqltemplete)
			sqlLog, errRes = q.Raw(
				sqltemplete,
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.MaxStep(),
				eve.CurrentStep,
				eve.Note(),
				fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
				eve.Status,
				eve.Func(),
				eve.Priority(),
				tpl_creator,
				eve.ExpressionId(),
				eve.StrategyId(),
				eve.TplId(),
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.Id,
			).Exec()
		} else {
			sqltemplete = fmt.Sprintf("%v WHERE id = ?", sqltemplete)
			sqlLog, errRes = q.Raw(
				sqltemplete,
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.MaxStep(),
				eve.CurrentStep,
				eve.Note(),
				fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
				eve.Status,
				eve.Func(),
				eve.Priority(),
				tpl_creator,
				eve.ExpressionId(),
				eve.StrategyId(),
				eve.TplId(),
				eve.Id,
			).Exec()
		}
	}
	log.Debugf("%v, %v", sqlLog, errRes)
	//insert case
	insertEvent(q, eve)
}

func counterGen(metric string, tags string) (mycounter string) {
	mycounter = metric
	if tags != "" {
		mycounter = fmt.Sprintf("%s/%s", metric, tags)
	}
	return
}

func DeleteEventOlder(before time.Time, limit int) {
	t := before.Format(timeLayout)
	sqlTpl := `delete from events where timestamp<? limit ?`
	q := orm.NewOrm()
	resp, err := q.Raw(sqlTpl, t, limit).Exec()
	if err != nil {
		log.Errorf("delete event older than %v fail, error:%v", t, err)
	} else {
		affected, _ := resp.RowsAffected()
		log.Debugf("delete event older than %v, rows affected:%v", t, affected)
	}
}
