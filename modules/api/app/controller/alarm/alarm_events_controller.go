package alarm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	alm "github.com/open-falcon/falcon-plus/modules/api/app/model/alarm"
)

// APIGetAlarmListsInputs TODO:
type APIGetAlarmListsInputs struct {
	StartTime     int64  `json:"startTime" form:"startTime"`
	EndTime       int64  `json:"endTime" form:"endTime"`
	Priority      int    `json:"priority" form:"priority"`
	Status        string `json:"status" form:"status"`
	ProcessStatus string `json:"process_status" form:"process_status"`
	Metrics       string `json:"metrics" form:"metrics"`
	//id
	EventID string `json:"event_id" form:"event_id"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
	//endpoints strategy template
	Endpoints  []string `json:"endpoints" form:"endpoints"`
	StrategyID int      `json:"strategy_id" form:"strategy_id"`
	TemplateID int      `json:"template_id" form:"template_id"`
}

func (input APIGetAlarmListsInputs) checkInputsContain() error {
	if input.StartTime == 0 && input.EndTime == 0 {
		if input.EventID == "" && input.Endpoints == nil && input.StrategyID == 0 && input.TemplateID == 0 {
			return errors.New("StartTime, endTime, event_id, endpoints, strategy_id or template_id, You have to at least pick one on the request")
		}
	}
	return nil
}

// TODO: remove
func (input APIGetAlarmListsInputs) collectFilters() string {
	tmp := []string{}
	if input.StartTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", input.StartTime))
	}
	if input.EndTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", input.EndTime))
	}
	if input.Priority != -1 {
		tmp = append(tmp, fmt.Sprintf("priority = %d", input.Priority))
	}
	if input.Status != "" {
		status := ""
		statusTmp := strings.Split(input.Status, ",")
		for indx, n := range statusTmp {
			if indx == 0 {
				status = fmt.Sprintf(" status = '%s' ", n)
			} else {
				status = fmt.Sprintf(" %s OR status = '%s' ", status, n)
			}
		}
		status = fmt.Sprintf("( %s )", status)
		tmp = append(tmp, status)
	}
	if input.ProcessStatus != "" {
		pstatus := ""
		pstatusTmp := strings.Split(input.ProcessStatus, ",")
		for indx, n := range pstatusTmp {
			if indx == 0 {
				pstatus = fmt.Sprintf(" process_status = '%s' ", n)
			} else {
				pstatus = fmt.Sprintf(" %s OR process_status = '%s' ", pstatus, n)
			}
		}
		pstatus = fmt.Sprintf("( %s )", pstatus)
		tmp = append(tmp, pstatus)
	}
	if input.Metrics != "" {
		tmp = append(tmp, fmt.Sprintf("metrics regexp '%s'", input.Metrics))
	}
	if input.EventID != "" {
		tmp = append(tmp, fmt.Sprintf("id = '%s'", input.EventID))
	}
	if input.Endpoints != nil && len(input.Endpoints) != 0 {
		for i, ep := range input.Endpoints {
			input.Endpoints[i] = fmt.Sprintf("'%s'", ep)
		}
		tmp = append(tmp, fmt.Sprintf("endpoint in (%s)", strings.Join(input.Endpoints, ", ")))
	}
	if input.StrategyID != 0 {
		tmp = append(tmp, fmt.Sprintf("strategy_id = %d", input.StrategyID))
	}
	if input.TemplateID != 0 {
		tmp = append(tmp, fmt.Sprintf("template_id = %d", input.TemplateID))
	}
	filterStrTmp := strings.Join(tmp, " AND ")
	if filterStrTmp != "" {
		filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
	}
	return filterStrTmp
}

func (input APIGetAlarmListsInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
	filterDB := database.Table(tableName)
	// nil columns mean select all columns
	if columns != nil && len(columns) != 0 {
		filterDB = filterDB.Select(columns)
	}
	if input.StartTime != 0 {
		filterDB = filterDB.Where("timestamp >= FROM_UNIXTIME(?)", input.StartTime)
	}
	if input.EndTime != 0 {
		filterDB = filterDB.Where("timestamp <= FROM_UNIXTIME(?)", input.EndTime)
	}
	if input.Priority != -1 {
		filterDB = filterDB.Where("priority = ?", input.Priority)
	}
	if input.Status != "" {
		statusTmp := strings.Split(input.Status, ",")
		filterDB = filterDB.Where("status in (?)", statusTmp)
	}
	if input.ProcessStatus != "" {
		pstatusTmp := strings.Split(input.ProcessStatus, ",")
		filterDB = filterDB.Where("process_status in (?)", pstatusTmp)
	}
	if input.Metrics != "" {
		filterDB = filterDB.Where("metric regexp ?", input.Metrics)
	}
	if input.EventID != "" {
		filterDB = filterDB.Where("id = ?", input.EventID)
	}
	if input.Endpoints != nil && len(input.Endpoints) != 0 {
		filterDB = filterDB.Where("endpoint in (?)", input.Endpoints)
	}
	if input.StrategyID != 0 {
		filterDB = filterDB.Where("strategy_id = ?", input.StrategyID)
	}
	if input.TemplateID != 0 {
		filterDB = filterDB.Where("template_id = ?", input.TemplateID)
	}
	return filterDB
}

// Lists TODO:
func Lists(c *gin.Context) {
	var inputs APIGetAlarmListsInputs
	//set default
	inputs.Page = -1
	inputs.Limit = -1
	inputs.Priority = -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	//for get correct table name
	f := alm.EventCases{}
	alarmDB := inputs.collectDBFilters(db.Alarm, f.TableName(), nil)
	cevens := []alm.EventCases{}
	//if no specific, will give return first 2000 records
	if inputs.Page == -1 && inputs.Limit == -1 {
		inputs.Limit = 2000
		alarmDB = alarmDB.Order("timestamp DESC").Limit(inputs.Limit)
	} else if inputs.Limit == -1 {
		// set page but not set limit
		h.JSONR(c, badstatus, errors.New("You set page but skip limit params, please check your input"))
		return
	} else {
		// set limit but not set page
		if inputs.Page == -1 {
			// limit invalid
			if inputs.Limit <= 0 {
				h.JSONR(c, badstatus, errors.New("limit or page can not set to 0 or less than 0"))
				return
			}
			// set default page
			inputs.Page = 1
		} else {
			// set page and limit
			// page or limit invalid
			if inputs.Page <= 0 || inputs.Limit <= 0 {
				h.JSONR(c, badstatus, errors.New("limit or page can not set to 0 or less than 0"))
				return
			}
		}
		//set the max limit of each page
		if inputs.Limit >= 50 {
			inputs.Limit = 50
		}
		step := (inputs.Page - 1) * inputs.Limit
		alarmDB = alarmDB.Order("timestamp DESC").Offset(step).Limit(inputs.Limit)
	}
	alarmDB.Find(&cevens)
	h.JSONR(c, cevens)
}

// APIEventsGetInputs TODO:
type APIEventsGetInputs struct {
	StartTime int64 `json:"startTime" form:"startTime"`
	EndTime   int64 `json:"endTime" form:"endTime"`
	Status    int   `json:"status" form:"status" binding:"gte=-1,lte=1"`
	//event_caseId
	EventID string `json:"event_id" form:"event_id" binding:"required"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

func (input APIEventsGetInputs) collectFilters() string {
	tmp := []string{}
	filterStrTmp := ""
	if input.StartTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", input.StartTime))
	}
	if input.EndTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", input.EndTime))
	}
	if input.EventID != "" {
		tmp = append(tmp, fmt.Sprintf("event_caseId = '%s'", input.EventID))
	}
	if input.Status == 0 || input.Status == 1 {
		tmp = append(tmp, fmt.Sprintf("status = %d", input.Status))
	}
	if len(tmp) != 0 {
		filterStrTmp = strings.Join(tmp, " AND ")
		filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
	}
	return filterStrTmp
}

func (input APIEventsGetInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
	filterDB := database.Table(tableName)
	// nil columns mean select all columns
	if columns != nil && len(columns) != 0 {
		filterDB = filterDB.Select(columns)
	}
	if input.StartTime != 0 {
		filterDB = filterDB.Where("timestamp >= FROM_UNIXTIME(?)", input.StartTime)
	}
	if input.EndTime != 0 {
		filterDB = filterDB.Where("timestamp <= FROM_UNIXTIME(?)", input.EndTime)
	}
	if input.EventID != "" {
		filterDB = filterDB.Where("event_caseId = ?", input.EventID)
	}
	if input.Status == 0 || input.Status == 1 {
		filterDB = filterDB.Where("status = ?", input.Status)
	}
	return filterDB
}

// EventsGet TODO:
func EventsGet(c *gin.Context) {
	var inputs APIEventsGetInputs
	inputs.Status = -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	//for get correct table name
	f := alm.Events{}
	eventDB := inputs.collectDBFilters(db.Alarm, f.TableName(), []string{"id", "step", "event_caseId", "cond", "status", "timestamp"})
	evens := []alm.Events{}
	if inputs.Limit <= 0 || inputs.Limit >= 50 {
		inputs.Limit = 50
	}
	step := (inputs.Page - 1) * inputs.Limit
	eventDB.Order("timestamp DESC").Offset(step).Limit(inputs.Limit).Scan(&evens)
	h.JSONR(c, evens)
}
