package host

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
)

// APISetMaintainInput TODO:
type APISetMaintainInput struct {
	Hosts []string `json:"hosts"`
	Ids   []int64  `json:"ids"`
	Begin int64    `json:"maintain_begin" binding:"required"`
	End   int64    `json:"maintain_end" binding:"required"`
}

// SetMaintain TODO:
func SetMaintain(c *gin.Context) {
	var dt *gorm.DB
	var inputs APISetMaintainInput
	var method string

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	if len(inputs.Hosts) > 0 {
		method = "hosts"
		dt = db.Falcon.Table("host").Where("hostname IN (?)", inputs.Hosts).Updates(map[string]interface{}{"maintain_begin": inputs.Begin, "maintain_end": inputs.End})
	} else if len(inputs.Ids) > 0 {
		method = "ids"
		dt = db.Falcon.Table("host").Where("id IN (?)", inputs.Ids).Updates(map[string]interface{}{"maintain_begin": inputs.Begin, "maintain_end": inputs.End})
	} else {
		h.JSONR(c, badstatus, "hosts or ids is required")
		return
	}

	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("Through: %s, Affect row: %d", method, dt.RowsAffected))
}

// APIUnsetMaintainInput TODO:
type APIUnsetMaintainInput struct {
	Hosts []string `json:"hosts"`
	Ids   []int64  `json:"ids"`
}

// UnsetMaintain TODO:
func UnsetMaintain(c *gin.Context) {
	var dt *gorm.DB
	var inputs APIUnsetMaintainInput
	var method string

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	if len(inputs.Hosts) > 0 {
		method = "hosts"
		dt = db.Falcon.Table("host").Where("hostname IN (?)", inputs.Hosts).Updates(map[string]interface{}{"maintain_begin": 0, "maintain_end": 0})
	} else if len(inputs.Ids) > 0 {
		method = "ids"
		dt = db.Falcon.Table("host").Where("id IN (?)", inputs.Ids).Updates(map[string]interface{}{"maintain_begin": 0, "maintain_end": 0})
	} else {
		h.JSONR(c, badstatus, "hosts or ids is required")
		return
	}

	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("Through: %s, Affect row: %d", method, dt.RowsAffected))
}
