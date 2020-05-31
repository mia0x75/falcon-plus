package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// CreateScreen TODO:
func CreateScreen(c *gin.Context) {
	pid := c.DefaultPostForm("pid", "0")
	name := c.DefaultPostForm("name", "")
	if name == "" {
		h.JSONR(c, h.HTTPBadRequest, "empty name")
		return
	}

	ipid, err := strconv.Atoi(pid)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, "invalid screen pid")
		return
	}

	if err := db.Exec("insert ignore into screens (pid, name) values(?, ?)", ipid, name).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}

	var lid []int
	screen := model.Screen{}
	if err := db.Table(screen.TableName()).Select("id").Where("pid = ? and name = ?", ipid, name).Limit(1).Pluck("id", &lid).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	if len(lid) == 0 {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("no such screen where name=%s", name))
		return
	}
	screenID := lid[0]
	resp := map[string]interface{}{
		"pid":  ipid,
		"id":   screenID,
		"name": name,
	}
	h.JSONR(c, resp)
}

// GetScreen TODO:
func GetScreen(c *gin.Context) {
	id := c.Param("screen_id")

	screenID, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, "invalid screen id")
		return
	}

	screen := model.Screen{}
	if err := db.Where("id = ?", screenID).First(&screen).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	h.JSONR(c, screen)
}

// GetScreensByPid TODO:
func GetScreensByPid(c *gin.Context) {
	id := c.Param("pid")

	pid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, "invalid screen pid")
		return
	}

	screens := []model.Screen{}
	if err := db.Where("pid = ?", pid).Find(&screens).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	h.JSONR(c, screens)
}

// GetScreensAll TODO:
func GetScreensAll(c *gin.Context) {
	limit := c.DefaultQuery("limit", "500")
	screens := []model.Screen{}
	if err := db.Limit(limit).Find(&screens).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	h.JSONR(c, screens)
}

// DeleteScreen TODO:
func DeleteScreen(c *gin.Context) {
	inputScreenID := c.Param("id")
	if inputScreenID == "" {
		log.Debug("[D] parameter `id` for screen is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for screen is missing")
		return
	}

	screenID, err := strconv.Atoi(inputScreenID)
	if err != nil {
		log.Debugf("[D] parameter `id` for screen is invalid, value = %v", inputScreenID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for screen is invalid, value = %v", inputScreenID))
		return
	}

	screen := model.Screen{}
	if err := db.Where("id = ?", screenID).Delete(&screen).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}

	h.JSONR(c, "ok")
}

// UpdateScreen TODO:
func UpdateScreen(c *gin.Context) {
	inputScreenID := c.Param("id")
	if inputScreenID == "" {
		log.Debug("[D] parameter `id` for screen is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for screen is missing")
		return
	}

	screenID, err := strconv.Atoi(inputScreenID)
	if err != nil {
		log.Debugf("[D] parameter `id` for screen is invalid, value = %v", inputScreenID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for screen is invalid, value = %v", inputScreenID))
		return
	}

	data := map[string]interface{}{}
	pid := c.PostForm("pid")
	name := c.PostForm("name")
	if name != "" {
		data["name"] = name
	}

	if pid != "" {
		ipid, err := strconv.Atoi(pid)
		if err != nil {
			h.JSONR(c, h.HTTPBadRequest, "invalid screen pid")
			return
		}
		data["pid"] = ipid
	}

	if err := db.Model(model.Screen{}).Where("id = ?", screenID).Update(data).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}

	h.JSONR(c, "ok")
}
