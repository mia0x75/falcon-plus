package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// APISetMaintainInput 设定主机状态维护对应的结构体
type APISetMaintainInput struct {
	Hosts []string `json:"hosts"         binding:"required"`
	Ids   []int64  `json:"ids"           binding:"required"`
	Begin int64    `json:"maintainBegin" binding:"required"`
	End   int64    `json:"maintainEnd"   binding:"required"`
}

// SetMaintain 设定主机状态维护
func SetMaintain(c *gin.Context) {
	var (
		inputs APISetMaintainInput
		method string
		err    error
	)

	if err = c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	maintainInfo := map[string]int64{
		"maintain_begin": inputs.Begin,
		"maintain_end":   inputs.End,
	}

	if len(inputs.Hosts) > 0 {
		method = "hosts"
		err = db.Model(model.Host{}).
			Where("hostname IN (?)", inputs.Hosts).
			Updates(maintainInfo).Error
	} else if len(inputs.Ids) > 0 {
		method = "ids"
		err = db.Model(model.Host{}).
			Where("id IN (?)", inputs.Ids).
			Updates(maintainInfo).Error
	} else {
		h.JSONR(c, h.HTTPBadRequest, "hosts or ids is required")
		return
	}

	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("Through: %s", method))
}

// APIUnsetMaintainInput 解除主机维护状态对应的结构体
type APIUnsetMaintainInput struct {
	Hosts []string `json:"hosts" binding:"required"`
	Ids   []int64  `json:"ids"   binding:"required"`
}

// UnsetMaintain 解除主机维护状态
func UnsetMaintain(c *gin.Context) {
	var inputs APIUnsetMaintainInput
	var method string
	var err error

	if err = c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	maintainInfo := map[string]int64{
		"maintain_begin": 0,
		"maintain_end":   0,
	}
	if len(inputs.Hosts) > 0 {
		method = "hosts"
		err = db.Model(model.Host{}).
			Where("hostname IN (?)", inputs.Hosts).
			Updates(maintainInfo).
			Error
	} else if len(inputs.Ids) > 0 {
		method = "ids"
		err = db.Model(model.Host{}).
			Where("id IN (?)", inputs.Ids).
			Updates(maintainInfo).
			Error
	} else {
		h.JSONR(c, h.HTTPBadRequest, "hosts or ids is required")
		return
	}

	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("Through: %s", method))
}
