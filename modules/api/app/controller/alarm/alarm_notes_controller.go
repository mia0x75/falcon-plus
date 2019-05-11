package alarm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	alm "github.com/open-falcon/falcon-plus/modules/api/app/model/alarm"
)

// APIGetNotesOfAlarmInputs TODO:
type APIGetNotesOfAlarmInputs struct {
	StartTime int64 `json:"startTime" form:"startTime"`
	EndTime   int64 `json:"endTime" form:"endTime"`
	//id
	EventID string `json:"event_id" form:"event_id"`
	Status  string `json:"status" form:"status"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

func (input APIGetNotesOfAlarmInputs) checkInputsContain() error {
	if input.StartTime == 0 && input.EndTime == 0 {
		if input.EventID == "" {
			return errors.New("StartTime, endTime OR event_id, You have to at least pick one on the request")
		}
	}
	return nil
}

func (input APIGetNotesOfAlarmInputs) collectFilters() string {
	tmp := []string{}
	if input.StartTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", input.StartTime))
	}
	if input.EndTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", input.EndTime))
	}
	if input.Status != "" {
		tmp = append(tmp, fmt.Sprintf("status = '%s'", input.Status))
	}
	if input.EventID != "" {
		tmp = append(tmp, fmt.Sprintf("event_caseId = '%s'", input.EventID))
	}
	filterStrTmp := strings.Join(tmp, " AND ")
	if filterStrTmp != "" {
		filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
	}
	return filterStrTmp
}

func (input APIGetNotesOfAlarmInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
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
	if input.Status != "" {
		filterDB = filterDB.Where("status = ?", input.Status)
	}
	if input.EventID != "" {
		filterDB = filterDB.Where("event_caseId = ?", input.EventID)
	}
	return filterDB
}

// APIGetNotesOfAlarmOuput TODO:
type APIGetNotesOfAlarmOuput struct {
	EventCaseID string     `json:"event_caseId"`
	Note        string     `json:"note"`
	CaseID      string     `json:"case_id"`
	Status      string     `json:"status"`
	Timestamp   *time.Time `json:"timestamp"`
	UserName    string     `json:"user"`
}

// GetNotesOfAlarm TODO:
func GetNotesOfAlarm(c *gin.Context) {
	var inputs APIGetNotesOfAlarmInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, "binding input got error: "+err.Error())
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	//for get correct table name
	f := alm.EventNote{}
	noteDB := inputs.collectDBFilters(db.Alarm, f.TableName(), []string{"id", "event_caseId", "note", "case_id", "status", "timestamp", "user_id"})
	notes := []alm.EventNote{}
	if inputs.Limit <= 0 || inputs.Limit >= 50 {
		inputs.Limit = 50
	}
	step := (inputs.Page - 1) * inputs.Limit
	noteDB.Order("timestamp DESC").Offset(step).Limit(inputs.Limit).Scan(&notes)
	output := []APIGetNotesOfAlarmOuput{}
	for _, n := range notes {
		output = append(output, APIGetNotesOfAlarmOuput{
			EventCaseID: n.EventCaseId,
			Note:        n.Note,
			CaseID:      n.CaseId,
			Status:      n.Status,
			Timestamp:   n.Timestamp,
			UserName:    n.GetUserName(),
		})
	}
	h.JSONR(c, output)
}

// APIAddNotesToAlarmInputs TODO:
type APIAddNotesToAlarmInputs struct {
	EventID string `json:"event_id" form:"event_id" binding:"required"`
	Note    string `json:"note" form:"note" binding:"required"`
	Status  string `json:"status" form:"status" binding:"required"`
	CaseID  string `json:"case_id" form:"case_id"`
}

// CheckingFormating TODO:
func (input APIAddNotesToAlarmInputs) CheckingFormating() error {
	switch input.Status {
	case "in progress":
		return nil
	case "unresolved":
		return nil
	case "resolved":
		return nil
	case "ignored":
		return nil
	case "comment":
		return nil
	default:
		return errors.New(`Params status: only accepect ["in progress", "unresolved", "resolved", "ignored", "comment"]`)
	}
}

// AddNotesToAlarm TODO:
func AddNotesToAlarm(c *gin.Context) {
	var inputs APIAddNotesToAlarmInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckingFormating(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	Anote := alm.EventNote{
		UserId:      user.ID,
		Note:        inputs.Note,
		Status:      inputs.Status,
		EventCaseId: inputs.EventID,
		CaseId:      inputs.CaseID,
	}
	dt := db.Alarm.Begin()
	if err := dt.Save(&Anote); err.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, err.Error)
		return
	}
	if inputs.Status != "comment" {
		ecase := alm.EventCases{
			ProcessNote:   Anote.ID,
			ProcessStatus: Anote.Status,
		}
		if db := dt.Table(ecase.TableName()).Where("id = ?", Anote.EventCaseId).Update(&ecase); db.Error != nil {
			dt.Rollback()
			h.JSONR(c, badstatus, "update got error during update event_cases:"+db.Error.Error())
			return
		}
	}
	dt.Commit()
	h.JSONR(c, map[string]string{
		"id":      inputs.EventID,
		"message": fmt.Sprintf("add note to %s successfuled", inputs.EventID),
	})
	return
}
