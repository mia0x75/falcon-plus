package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// APIEventsGetInputs TODO:
type APIEventsGetInputs struct {
	StartTime int64  `json:"startTime" form:"startTime"`
	EndTime   int64  `json:"endTime"   form:"endTime"`
	Status    int    `json:"status"    form:"status"    binding:"gte=-1,lte=1"`
	EventID   string `json:"eventID"   form:"eventID"   binding:"required"` // event_caseId
	Limit     int    `json:"limit"     form:"limit"`                        // number of reacord's limit on each page
	Page      int    `json:"page"      form:"page"`                         // pagging
}

func (input APIEventsGetInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
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
	if input.EventID != "" {
		filterDB = filterDB.Where("event_caseId = ?", input.EventID)
	}
	if input.Status == 0 || input.Status == 1 {
		filterDB = filterDB.Where("status = ?", input.Status)
	}
	return filterDB
}

// GetEvents TODO:
func GetEvents(c *gin.Context) {
	var inputs APIEventsGetInputs
	inputs.Status = -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	// for get correct table name
	f := model.Event{}
	filterDB := inputs.collectDBFilters(db, f.TableName(), []string{"id", "step", "event_caseId", "cond", "status", "create_at"})
	events := []model.Event{}
	if inputs.Limit <= 0 || inputs.Limit >= 50 {
		inputs.Limit = 50
	}
	step := (inputs.Page - 1) * inputs.Limit
	if err := filterDB.Order("create_at DESC").Offset(step).Limit(inputs.Limit).Scan(&events).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
	} else {
		h.JSONR(c, events)
	}
}
