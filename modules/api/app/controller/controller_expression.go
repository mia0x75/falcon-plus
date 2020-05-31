package controller

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
)

// GetExpressionList TODO:
func GetExpressionList(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	inputPage := c.DefaultQuery("page", "")
	inputLimit := c.DefaultQuery("limit", "")
	page, limit, err = h.PageParser(inputPage, inputLimit)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err.Error())
		return
	}
	expressions := []model.Expression{}
	if limit != -1 && page != -1 {
		err = db.Raw("SELECT * FROM expressions LIMIT ?, ?", page, limit).Find(&expressions).Error
	} else {
		err = db.Find(&expressions).Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, expressions)
}

// GetExpression TODO:
func GetExpression(c *gin.Context) {
	inputExpressionID := c.Params.ByName("id")
	if inputExpressionID == "" {
		log.Debug("[D] parameter `id` for expression is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for expression is missing")
		return
	}
	expressionID, err := strconv.Atoi(inputExpressionID)
	if err != nil {
		log.Debugf("[D] parameter `id` for expression is invalid, value = %v", inputExpressionID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for expression is invalid, value = %v", inputExpressionID))
		return
	}
	expression := cache.ExpressionsMap.Any(func(elem *model.Expression) bool {
		if elem.ID == int64(expressionID) {
			return true
		}
		return false
	})
	if expression == nil {
		h.JSONR(c, h.HTTPBadRequest, "...") // TODO:
		return
	}
	resp := map[string]interface{}{
		"expression": expression,
	}
	action := cache.ActionsMap.Any(func(elem *model.Action) bool {
		if elem.ID == expression.ActionID {
			return true
		}
		return false
	})
	if action != nil {
		resp["action"] = action
	}
	h.JSONR(c, resp)
}

// APICreateExrpessionInput TODO:
type APICreateExrpessionInput struct {
	Expression string    `json:"expression" binding:"required"`
	Func       string    `json:"func"       binding:"required"`
	Op         string    `json:"op"         binding:"required"`
	RightValue string    `json:"rightValue" binding:"required"`
	MaxStep    int       `json:"maxStep"    binding:"required"`
	Priority   int       `json:"priority"   binding:"required"`
	Note       string    `json:"note"       binding:"required"`
	Pause      int       `json:"pause"      binding:"required"`
	Action     ActionTmp `json:"action"     binding:"required"`
}

// ActionTmp TODO:
type ActionTmp struct {
	UIC                []string `json:"uic"                binding:"required"`
	URL                string   `json:"url"                binding:"required"`
	Callback           int      `json:"callback"           binding:"required"`
	BeforeCallbackSMS  int      `json:"beforeCallbackSMS"  binding:"required"`
	AfterCallbackSMS   int      `json:"afterCallbackSMS"   binding:"required"`
	BeforeCallbackMail int      `json:"beforeCallbackMail" binding:"required"`
	AfterCallbackMail  int      `json:"afterCallbackMail"  binding:"required"`
}

// CheckFormat TODO:
func (input APICreateExrpessionInput) CheckFormat() (err error) {
	validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
	validRightValue := regexp.MustCompile(`^\-?\d+(\.\d+)?$`)
	switch {
	case !validOp.MatchString(input.Op):
		err = errors.New("op's formating is not vaild")
	case !validRightValue.MatchString(input.RightValue):
		err = errors.New("right_value's formating is not vaild")
	}
	return
}

// CreateExrpession TODO:
func CreateExrpession(c *gin.Context) {
	var inputs APICreateExrpessionInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	action := model.Action{
		UIC:                strings.Join(inputs.Action.UIC, ","),
		URL:                inputs.Action.URL,
		Callback:           inputs.Action.Callback,
		BeforeCallbackSMS:  inputs.Action.BeforeCallbackSMS,
		BeforeCallbackMail: inputs.Action.BeforeCallbackMail,
		AfterCallbackSMS:   inputs.Action.AfterCallbackSMS,
		AfterCallbackMail:  inputs.Action.AfterCallbackMail,
	}
	tx := db.Begin()
	if err := tx.Save(&action).Error; err != nil {
		h.InternelError(c, "creating data", err)
		tx.Rollback()
		return
	}
	expression := model.Expression{
		Expression: inputs.Expression,
		Func:       inputs.Func,
		Op:         inputs.Op,
		RightValue: inputs.RightValue,
		MaxStep:    inputs.MaxStep,
		Priority:   inputs.Priority,
		Note:       inputs.Note,
		Pause:      inputs.Pause,
		Creator:    ctx.ID,
		ActionID:   action.ID,
	}
	if err := tx.Save(&expression).Error; err != nil {
		h.InternelError(c, "creating data", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, expression)
	go cache.ExpressionsMap.Init()
	go cache.ActionsMap.Init()
}

// APIUpdateExrpessionInput TODO:
type APIUpdateExrpessionInput struct {
	ID         int64      `json:"id"         binding:"required"`
	Expression string     `json:"expression" binding:"required"`
	Func       string     `json:"func"       binding:"required"`
	Op         string     `json:"op"         binding:"required"`
	RightValue string     `json:"rightValue" binding:"required"`
	MaxStep    int        `json:"maxStep"    binding:"required"`
	Priority   int        `json:"priority"   binding:"required"`
	Note       string     `json:"note"       binding:"required"`
	Pause      int        `json:"pause"      binding:"required"`
	Action     ActionTmpU `json:"action"     binding:"required"`
}

// ActionTmpU TODO:
type ActionTmpU struct {
	UIC                []string `json:"uic"                binding:"required"`
	URL                string   `json:"url"                binding:"required"`
	Callback           int      `json:"callback"           binding:"required"`
	BeforeCallbackSMS  int      `json:"beforeCallbackSMS"  binding:"required"`
	AfterCallbackSMS   int      `json:"afterCallbackSMS"   binding:"required"`
	BeforeCallbackMail int      `json:"beforeCallbackMail" binding:"required"`
	AfterCallbackMail  int      `json:"afterCallbackMail"  binding:"required"`
}

// CheckFormat TODO:
func (input APIUpdateExrpessionInput) CheckFormat() (err error) {
	validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
	validRightValue := regexp.MustCompile(`^\d+$`)
	switch {
	case !validOp.MatchString(input.Op):
		err = errors.New("op's formating is not vaild")
	case !validRightValue.MatchString(input.RightValue):
		err = errors.New("right_value's formating is not vaild")
	}
	return
}

// UpdateExrpession TODO:
func UpdateExrpession(c *gin.Context) {
	var inputs APIUpdateExrpessionInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	expression := cache.ExpressionsMap.Any(func(elem *model.Expression) bool {
		if elem.ID == inputs.ID {
			return true
		}
		return false
	})
	if expression == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Expression (id = %d) does not exist", inputs.ID))
		return
	}
	if !ctx.IsAdmin() && expression.Creator != ctx.ID {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}
	action := cache.ActionsMap.Any(func(elem *model.Action) bool {
		if elem.ID == expression.ActionID {
			return true
		}
		return false
	})
	update := map[string]interface{}{
		"ID":         expression.ID,
		"Expression": inputs.Expression,
		"Func":       inputs.Func,
		"Op":         inputs.Op,
		"RightValue": inputs.RightValue,
		"MaxStep":    inputs.MaxStep,
		"Priority":   inputs.Priority,
		"Note":       inputs.Note,
		"Pause":      inputs.Pause,
	}
	tx := db.Begin()
	if err := tx.Model(expression).Update(update).Error; err != nil {
		h.InternelError(c, "updating ", err)
		tx.Rollback()
		return
	}
	if action != nil {
		update := map[string]interface{}{
			"ID":                 expression.ActionID,
			"UIC":                strings.Join(inputs.Action.UIC, ","),
			"URL":                inputs.Action.URL,
			"Callback":           inputs.Action.Callback,
			"BeforeCallbackSMS":  inputs.Action.BeforeCallbackSMS,
			"BeforeCallbackMail": inputs.Action.BeforeCallbackMail,
			"AfterCallbackSMS":   inputs.Action.AfterCallbackSMS,
			"AfterCallbackMail":  inputs.Action.AfterCallbackMail,
		}
		if err := tx.Model(action).Update(update).Error; err != nil {
			h.InternelError(c, "updating data", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("Expression (id = %d) updated", inputs.ID))
	go cache.ExpressionsMap.Init()
	go cache.ActionsMap.Init()
}

// DeleteExpression TODO:
func DeleteExpression(c *gin.Context) {
	inputExpressionID := c.Params.ByName("id")
	if inputExpressionID == "" {
		log.Debug("[D] parameter `id` for expression is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for expression is missing")
		return
	}
	expressionID, err := strconv.Atoi(inputExpressionID)
	if err != nil {
		log.Debugf("[D] parameter `id` for expression is invalid, value = %v", inputExpressionID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for expression is invalid, value = %v", inputExpressionID))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	expression := cache.ExpressionsMap.Any(func(elem *model.Expression) bool {
		if elem.ID == int64(expressionID) {
			return true
		}
		return false
	})
	if expression == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Expression (id = %d) does not exist", expressionID))
		return
	}
	if !ctx.IsAdmin() && expression.Creator != ctx.ID {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}
	tx := db.Begin()
	if err := tx.Where("id = ?", expression.ActionID).Delete(&model.Action{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	if err := tx.Delete(expression).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("Expression (id = %d) deleted", expressionID))
	go cache.ExpressionsMap.Init()
	go cache.ActionsMap.Init()
}
