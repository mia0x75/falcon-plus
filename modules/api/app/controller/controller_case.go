package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// APIGetAlarmListsInputs TODO:
type APIGetAlarmListsInputs struct {
	StartTime     int64    `json:"startTime"     form:"startTime"`
	EndTime       int64    `json:"endTime"       form:"endTime"`
	Priority      int      `json:"priority"      form:"priority"`
	Status        string   `json:"status"        form:"status"`
	ProcessStatus string   `json:"processStatus" form:"processStatus"`
	Metrics       string   `json:"metrics"       form:"metrics"`
	EventID       string   `json:"eventID"       form:"eventID"`   // id
	Limit         int      `json:"limit"         form:"limit"`     // number of reacord's limit on each page
	Page          int      `json:"page"          form:"page"`      // pagging
	Endpoints     []string `json:"endpoints"     form:"endpoints"` // endpoints strategy template
	StrategyID    int      `json:"strategyID"    form:"strategyID"`
	TemplateID    int      `json:"templateID"    form:"templateID"`
}

func (input APIGetAlarmListsInputs) checkInputsContain() error {
	if input.StartTime == 0 && input.EndTime == 0 {
		if input.EventID == "" && input.Endpoints == nil && input.StrategyID == 0 && input.TemplateID == 0 {
			return errors.New("StartTime, endTime, eventID, endpoints, strategyID or templateID, You have to at least pick one on the request")
		}
	}
	return nil
}

func (input APIGetAlarmListsInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
	filterDB := database.Table(tableName)
	// nil columns mean select all columns
	if columns != nil && len(columns) != 0 {
		filterDB = filterDB.Select(columns)
	}
	if input.StartTime != 0 {
		filterDB = filterDB.Where("create_at >= FROM_UNIXTIME(?)", input.StartTime)
	}
	if input.EndTime != 0 {
		filterDB = filterDB.Where("create_at <= FROM_UNIXTIME(?)", input.EndTime)
	}
	if input.Priority != -1 {
		filterDB = filterDB.Where("priority = ?", input.Priority)
	}
	if input.Status != "" {
		statusTmp := strings.Split(input.Status, ",")
		filterDB = filterDB.Where("status IN (?)", statusTmp)
	}
	if input.ProcessStatus != "" {
		pstatusTmp := strings.Split(input.ProcessStatus, ",")
		filterDB = filterDB.Where("process_status IN (?)", pstatusTmp)
	}
	if input.Metrics != "" {
		filterDB = filterDB.Where("metric REGEXP ?", input.Metrics)
	}
	if input.EventID != "" {
		filterDB = filterDB.Where("id = ?", input.EventID)
	}
	if input.Endpoints != nil && len(input.Endpoints) != 0 {
		filterDB = filterDB.Where("endpoint IN (?)", input.Endpoints)
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
	// set default
	inputs.Page = -1
	inputs.Limit = -1
	inputs.Priority = -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	// for get correct table name
	f := model.Case{}
	filterDB := inputs.collectDBFilters(db, f.TableName(), nil)
	cases := []model.Case{}
	// if no specific, will give return first 2000 records
	if inputs.Page == -1 && inputs.Limit == -1 {
		inputs.Limit = 2000
		filterDB = filterDB.Order("create_at DESC").Limit(inputs.Limit)
	} else if inputs.Limit == -1 {
		// set page but not set limit
		h.JSONR(c, h.HTTPBadRequest, errors.New("You set page but skip limit params, please check your input"))
		return
	} else {
		// set limit but not set page
		if inputs.Page == -1 {
			// limit invalid
			if inputs.Limit <= 0 {
				h.JSONR(c, h.HTTPBadRequest, errors.New("limit or page can not set to 0 or less than 0"))
				return
			}
			// set default page
			inputs.Page = 1
		} else {
			// set page and limit
			// page or limit invalid
			if inputs.Page <= 0 || inputs.Limit <= 0 {
				h.JSONR(c, h.HTTPBadRequest, errors.New("limit or page can not set to 0 or less than 0"))
				return
			}
		}
		// set the max limit of each page
		if inputs.Limit >= 50 {
			inputs.Limit = 50
		}
		step := (inputs.Page - 1) * inputs.Limit
		filterDB = filterDB.Order("create_at DESC").Offset(step).Limit(inputs.Limit)
	}
	if err := filterDB.Find(&cases).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
	} else {
		h.JSONR(c, cases)
	}
}
