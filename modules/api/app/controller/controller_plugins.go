package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
)

// APICreatePluginInput TODO:
type APICreatePluginInput struct {
	GroupID int64  `json:"groupID" binding:"required"`
	DirPath string `json:"dirPath" binding:"required"`
}

// CreatePlugin TODO:
func CreatePlugin(c *gin.Context) {
	var inputs APICreatePluginInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	group := cache.GroupsMap.Any(func(elem *model.Group) bool {
		if elem.ID == inputs.GroupID {
			return true
		}
		return false
	})
	if !ctx.IsAdmin() && group.Creator != ctx.ID {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}
	plugin := model.Plugin{
		Dir:     inputs.DirPath,
		GroupID: inputs.GroupID,
		Creator: ctx.ID,
	}
	if err := db.Save(&plugin).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}
	h.JSONR(c, plugin)
}

// GetPluginOfGrp TODO:
func GetPluginOfGroup(c *gin.Context) {
	inputGroupID := c.Params.ByName("id")
	if inputGroupID == "" {
		log.Debug("[D] parameter `id` for group is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for group is missing")
		return
	}
	groupID, err := strconv.Atoi(inputGroupID)
	if err != nil {
		log.Debugf("[D] parameter `id` for group is invalid, value = %v", inputGroupID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for aggregator is invalid, value = %v", inputGroupID))
		return
	}
	plugins := []model.Plugin{}
	if err := db.Where("group_id = ?", groupID).Find(&plugins).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, plugins)
}

// DeletePlugin TODO:
func DeletePlugin(c *gin.Context) {
	inputPluginID := c.Params.ByName("id")
	if inputPluginID == "" {
		log.Debug("[D] parameter `id` for plugin is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for plugin is missing")
		return
	}
	pluginID, err := strconv.Atoi(inputPluginID)
	if err != nil {
		log.Debugf("[D] parameter `id` for group is invalid, value = %v", inputPluginID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for aggregator is invalid, value = %v", inputPluginID))
		return
	}
	plugin := model.Plugin{}
	if err := db.Where("id = ?", pluginID).Find(&plugin).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	group := cache.GroupsMap.Any(func(elem *model.Group) bool {
		if elem.ID == plugin.GroupID {
			return true
		}
		return false
	})
	if !ctx.IsAdmin() && group.Creator != ctx.ID && plugin.Creator != ctx.ID {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}

	if err := db.Where("id = ?", pluginID).Delete(&plugin).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("Plugin (id = %d) deleted", pluginID))
}
