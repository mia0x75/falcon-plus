package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

// GetStrategies TODO:
func GetStrategies(c *gin.Context) {
	// var strategys []model.Strategy
	inputTemplateID := c.DefaultQuery("template_id", "")
	if inputTemplateID == "" {
		h.JSONR(c, h.HTTPBadRequest, "template id is missing")
		return
	}
	templateID, err := strconv.Atoi(inputTemplateID)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	strategies := cache.StrategiesMap.Filter(func(elem *model.Strategy) bool {
		if elem.TemplateID == int64(templateID) {
			return true
		}
		return false
	})

	// TODO:
	// if err := db.Where("template_id = ?", tid).Find(&strategys).Error; err != nil {
	// 	h.InternelError(c, "retrieving data", err)
	// 	return
	// }
	h.JSONR(c, strategies)
}

// APICreateStrategyInput TODO:
type APICreateStrategyInput struct {
	Metric     string `json:"metric"     binding:"required"`
	Tags       string `json:"tags"`
	MaxStep    int    `json:"maxStep"    binding:"required"`
	Priority   int    `json:"priority"   binding:"required"`
	Func       string `json:"func"       binding:"required"`
	Op         string `json:"op"         binding:"required"`
	RightValue string `json:"rightValue" binding:"required"`
	Note       string `json:"note"`
	RunBegin   string `json:"runBegin"`
	RunEnd     string `json:"runEnd"`
	TemplateID int64  `json:"templateID" binding:"required"`
}

// CheckFormat TODO:
func (s APICreateStrategyInput) CheckFormat() (err error) {
	validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
	validRightValue := regexp.MustCompile(`^\-?\d+(\.\d+)?$`)
	validTime := regexp.MustCompile(`^\d{2}:\d{2}$`)
	switch {
	case !validOp.MatchString(s.Op):
		err = errors.New("op's formating is not vaild")
	case !validRightValue.MatchString(s.RightValue):
		err = errors.New("right_value's formating is not vaild")
	case !validTime.MatchString(s.RunBegin) && s.RunBegin != "":
		err = errors.New("run_begin's formating is not vaild, please refer ex. 00:00")
	case !validTime.MatchString(s.RunEnd) && s.RunEnd != "":
		err = errors.New("run_end's formating is not vaild, please refer ex. 24:00")
	}
	return
}

// CreateStrategy TODO:
func CreateStrategy(c *gin.Context) {
	var inputs APICreateStrategyInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	strategy := model.Strategy{
		Metric:     inputs.Metric,
		Tags:       inputs.Tags,
		MaxStep:    inputs.MaxStep,
		Priority:   inputs.Priority,
		Func:       inputs.Func,
		Op:         inputs.Op,
		RightValue: inputs.RightValue,
		Note:       inputs.Note,
		RunBegin:   inputs.RunBegin,
		RunEnd:     inputs.RunEnd,
		TemplateID: inputs.TemplateID,
	}
	if err := db.Save(&strategy).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}
	h.JSONR(c, strategy)
	// 同步缓存
	go cache.StrategiesMap.Init()
}

// GetStrategy TODO:
func GetStrategy(c *gin.Context) {
	inputStrategyID := c.Params.ByName("id")
	if inputStrategyID == "" {
		h.JSONR(c, h.HTTPBadRequest, "strategy id is missing")
		return
	}
	strategyID, err := strconv.Atoi(inputStrategyID)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	// strategy := model.Strategy{}
	// if err := db.Where("id = ?", strategyID).Find(&strategy).Error; err != nil {
	// 	h.InternelError(c, "retrieving data", err)
	// 	return
	// }
	strategy := cache.StrategiesMap.Any(func(elem *model.Strategy) bool {
		if elem.ID == int64(strategyID) {
			return true
		}
		return false
	})
	h.JSONR(c, strategy)
}

// APIUpdateStrategyInput TODO:
type APIUpdateStrategyInput struct {
	ID         int64  `json:"id"         binding:"required"`
	Metric     string `json:"metric"     binding:"required"`
	Tags       string `json:"tags"`
	MaxStep    int    `json:"maxStep"    binding:"required"`
	Priority   int    `json:"priority"   binding:"required"`
	Func       string `json:"func"       binding:"required"`
	Op         string `json:"op"         binding:"required"`
	RightValue string `json:"rightValue" binding:"required"`
	Note       string `json:"note"`
	RunBegin   string `json:"runBegin"`
	RunEnd     string `json:"runEnd"`
}

// CheckFormat TODO:
func (s APIUpdateStrategyInput) CheckFormat() (err error) {
	validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
	validRightValue := regexp.MustCompile(`^\-?\d+(\.\d+)?$`)
	validTime := regexp.MustCompile(`^\d{2}:\d{2}$`)
	switch {
	case !validOp.MatchString(s.Op):
		err = errors.New("op's formating is not vaild")
	case !validRightValue.MatchString(s.RightValue):
		err = errors.New("right_value's formating is not vaild")
	case !validTime.MatchString(s.RunBegin) && s.RunBegin != "":
		err = errors.New("run_begin's formating is not vaild, please refer ex. 00:00")
	case !validTime.MatchString(s.RunEnd) && s.RunEnd != "":
		err = errors.New("run_end's formating is not vaild, please refer ex. 24:00")
	}
	return
}

// UpdateStrategy TODO:
func UpdateStrategy(c *gin.Context) {
	var inputs APIUpdateStrategyInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	strategy := model.Strategy{
		ID: inputs.ID,
	}
	if err := db.Find(&strategy).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	update := map[string]interface{}{
		"Metric":     inputs.Metric,
		"Tags":       inputs.Tags,
		"MaxStep":    inputs.MaxStep,
		"Priority":   inputs.Priority,
		"Func":       inputs.Func,
		"Op":         inputs.Op,
		"RightValue": inputs.RightValue,
		"Note":       inputs.Note,
		"RunBegin":   inputs.RunBegin,
		"RunEnd":     inputs.RunEnd}
	if err := db.Model(&strategy).Where("id = ?", strategy.ID).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("stragtegy: %d has been updated", strategy.ID))
	// 同步缓存
	go cache.StrategiesMap.Init()
}

// DeleteStrategy TODO:
func DeleteStrategy(c *gin.Context) {
	inputStrategyID := c.Params.ByName("id")
	if inputStrategyID == "" {
		h.JSONR(c, h.HTTPBadRequest, "strategy id is missing")
		return
	}
	strategyID, err := strconv.Atoi(inputStrategyID)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	strategy := model.Strategy{}
	if err := db.Where("id = ?", strategyID).Delete(&strategy).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("strategy: %d has been deleted", strategyID))
	// 同步缓存
	go cache.StrategiesMap.Init()
}

// MetricQuery TODO:
func MetricQuery(c *gin.Context) {
	filePath := g.Config().MetricListFile
	if filePath == "" {
		filePath = "./data/metric"
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	metrics := strings.Split(string(data), "\n")
	h.JSONR(c, metrics)
}
