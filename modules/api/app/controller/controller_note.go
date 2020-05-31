package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
)

// APIGetNotesOfAlarmInputs TODO:
type APIGetNotesOfAlarmInputs struct {
	StartTime int64  `json:"startTime" form:"startTime"`
	EndTime   int64  `json:"endTime"   form:"endTime"`
	EventID   string `json:"eventID"   form:"eventID"` // id
	Status    string `json:"status"    form:"status"`
	Limit     int    `json:"limit"     form:"limit"` // number of reacord's limit on each page
	Page      int    `json:"page"      form:"page"`  // pagging
}

func (input APIGetNotesOfAlarmInputs) checkInputsContain() error {
	if input.StartTime == 0 && input.EndTime == 0 {
		if input.EventID == "" {
			return errors.New("StartTime, endTime or eventID, You have to at least pick one on the request")
		}
	}
	return nil
}

func (input APIGetNotesOfAlarmInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
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
	EventCaseID string    `json:"eventCaseID"`
	Note        string    `json:"note"`
	CaseID      string    `json:"caseID"`
	Status      string    `json:"status"`
	CreateAt    time.Time `json:"createAt"`
	UserName    string    `json:"user"`
}

// GetNotesOfAlarm TODO:
func GetNotesOfAlarm(c *gin.Context) {
	var inputs APIGetNotesOfAlarmInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	// for get correct table name
	f := model.Note{}
	noteDB := inputs.collectDBFilters(db, f.TableName(), []string{"id", "event_caseId", "note", "case_id", "status", "create_at", "user_id"})
	notes := []model.Note{}
	if inputs.Limit <= 0 || inputs.Limit >= 50 {
		inputs.Limit = 50
	}
	step := (inputs.Page - 1) * inputs.Limit
	if err := noteDB.Order("create_at DESC").
		Offset(step).
		Limit(inputs.Limit).
		Scan(&notes).
		Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	output := []APIGetNotesOfAlarmOuput{}
	for _, n := range notes {
		user := cache.UsersMap.Any(func(elem *model.User) bool {
			if elem.ID == n.Creator {
				return true
			}
			return false
		})
		output = append(output, APIGetNotesOfAlarmOuput{
			EventCaseID: n.EventCaseID,
			Note:        n.Note,
			CaseID:      n.CaseID,
			Status:      n.Status,
			CreateAt:    time.Unix(n.CreateAt, 0),
			UserName:    user.Name,
		})
	}
	h.JSONR(c, output)
}

// APIAddNotesToAlarmInputs TODO:
type APIAddNotesToAlarmInputs struct {
	EventID string `json:"eventID" form:"eventID" binding:"required"`
	Note    string `json:"note"    form:"note"    binding:"required"`
	Status  string `json:"status"  form:"status"  binding:"required"`
	CaseID  string `json:"caseID"  form:"caseID"`
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
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if err := inputs.CheckingFormating(); err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	note := model.Note{
		Creator:     ctx.ID,
		Note:        inputs.Note,
		Status:      inputs.Status,
		EventCaseID: inputs.EventID,
		CaseID:      inputs.CaseID,
	}
	tx := db.Begin()
	if err := tx.Save(&note).Error; err != nil {
		h.InternelError(c, "creating data", err)
		tx.Rollback()
		return
	}
	if inputs.Status != "comment" {
		ecase := model.Case{
			ProcessNote:   note.ID,
			ProcessStatus: note.Status,
		}
		if err := tx.Where("id = ?", note.EventCaseID).Update(&ecase).Error; err != nil {
			h.InternelError(c, "updating data", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	resp := map[string]interface{}{
		"id":      inputs.EventID,
		"message": fmt.Sprintf("Note (id = %s) created", inputs.EventID),
	}
	h.JSONR(c, resp)
}
