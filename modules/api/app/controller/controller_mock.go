package controller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// GetMockcfgList TODO:
func GetMockcfgList(c *gin.Context) {
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
	mockcfgs := []model.Mockcfg{}
	if limit != -1 && page != -1 {
		err = db.Raw("SELECT * FROM mockcfg LIMIT ?, ?", page, limit).
			Scan(&mockcfgs).
			Error
	} else {
		err = db.Find(&mockcfgs).Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, mockcfgs)
	return
}

// GetMockcfg TODO:
func GetMockcfg(c *gin.Context) {
	inputMockID := c.Params.ByName("id")
	if inputMockID == "" {
		log.Debug("[D] parameter `id` for mockcfg is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for mockcfg is missing")
		return
	}
	mockID, err := strconv.Atoi(inputMockID)
	if err != nil {
		log.Debugf("[D] parameter `id` for mockcfg is invalid, value = %v", inputMockID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for mockcfg is invalid, value = %v", inputMockID))
		return
	}
	mockcfg := model.Mockcfg{}
	if err := db.Where("id = ?", mockID).Find(&mockcfg).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, mockcfg)
}

// CreateMockcfg TODO:
func CreateMockcfg(c *gin.Context) {
	var inputs APICreateMockcfgInputs
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

	mockcfg := model.Mockcfg{
		Name:    inputs.Name,
		Obj:     inputs.Obj,
		ObjType: inputs.ObjType,
		Metric:  inputs.Metric,
		Tags:    inputs.Tags,
		DsType:  inputs.DsType,
		Step:    inputs.Step,
		Mock:    inputs.Mock,
		Creator: ctx.ID,
	}
	if err := db.Save(&mockcfg).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}
	h.JSONR(c, mockcfg)
}

// UpdateMockcfg TODO:
func UpdateMockcfg(c *gin.Context) {
	var inputs APIUpdateMockcfgInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	mockcfg := &model.Mockcfg{
		ID: inputs.ID,
	}
	update := map[string]interface{}{
		"Obj":     inputs.Obj,
		"ObjType": inputs.ObjType,
		"Metric":  inputs.Metric,
		"Tags":    inputs.Tags,
		"DsType":  inputs.DsType,
		"Step":    inputs.Step,
		"Mock":    inputs.Mock,
	}
	// TODO:
	if err := db.Model(&mockcfg).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("Mockcfg (id = %d) updated", inputs.ID))
}

// DeleteMockcfg TODO:
func DeleteMockcfg(c *gin.Context) {
	inputMockID := c.Params.ByName("id")
	if inputMockID == "" {
		h.JSONR(c, h.HTTPBadRequest, "nodata id is missing")
		return
	}
	mockID, err := strconv.Atoi(inputMockID)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	mockcfg := model.Mockcfg{}
	if err := db.Where("id = ?", mockID).Delete(&mockcfg).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("Mockcfg (id = %d) deleted", mockID))
}

// APICreateMockcfgInputs TODO:
type APICreateMockcfgInputs struct {
	Name    string  `json:"name"    binding:"required"`
	Obj     string  `json:"obj"     binding:"required"`
	ObjType string  `json:"objType" binding:"required"` // group, host, other
	Metric  string  `json:"metric"  binding:"required"`
	Tags    string  `json:"tags"    binding:"required"`
	DsType  string  `json:"dsType"  binding:"required"`
	Step    int     `json:"step"    binding:"required"`
	Mock    float64 `json:"mock"    binding:"required"`
}

// CheckFormat TODO:
func (s APICreateMockcfgInputs) CheckFormat() (err error) {
	switch {
	case s.ObjType != "group" && s.ObjType != "host" && s.ObjType != "other":
		err = errors.New("obj_type only accpect \"group, host, other\"")
	}
	return
}

// APIUpdateMockcfgInputs TODO:
type APIUpdateMockcfgInputs struct {
	ID      int64   `json:"ID"       binding:"required"`
	Obj     string  `json:"obj"      binding:"required"`
	ObjType string  `json:"objType" binding:"required"` // group, host, other
	Metric  string  `json:"metric"   binding:"required"`
	Tags    string  `json:"tags"     binding:"required"`
	DsType  string  `json:"dsType"   binding:"required"`
	Step    int     `json:"step"     binding:"required"`
	Mock    float64 `json:"mock"     binding:"required"`
}

// CheckFormat TODO:
func (s APIUpdateMockcfgInputs) CheckFormat() (err error) {
	switch {
	case s.ObjType != "group" && s.ObjType != "host" && s.ObjType != "other":
		err = errors.New("obj_type only accpect \"group, host, other\"")
	}
	return
}
