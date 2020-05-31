package controller

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
)

// GetGroups 分页获取主机分组信息表列
func GetGroups(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	inputPage := c.DefaultQuery("page", "")
	inputLimit := c.DefaultQuery("limit", "")
	q := c.DefaultQuery("q", ".+")
	q = template.HTMLEscapeString(q)
	page, limit, err = h.PageParser(inputPage, inputLimit)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err.Error())
		return
	}
	var groups []model.Group
	if limit != -1 && page != -1 {
		err = db.Raw("SELECT * FROM groups WHERE name REGEXP ? LIMIT ?, ?", q, page, limit).Scan(&groups).Error
	} else {
		err = db.Where("name REGEXP ?", q).Find(&groups).Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, groups)
}

// APICreateGroupInput 创建主机分组时对应的输入信息
type APICreateGroupInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateGroup 创建一个主机分组信息
func CreateGroup(c *gin.Context) {
	var inputs APICreateGroupInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	group := model.Group{
		Name:    inputs.Name,
		Creator: ctx.ID,
	}
	if cache.GroupsMap.Has(group) {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (name = %s) already exists", inputs.Name))
		return
	}
	if err := db.Create(&group).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}
	h.JSONR(c, group)
	go cache.GroupsMap.Init()
}

// APIAddHostsToGroupInput 将一个或者多个主机添加到某一个主机组对应的输入信息
type APIAddHostsToGroupInput struct {
	Hosts   []string `json:"hosts"   binding:"required"`
	GroupID int64    `json:"groupID" binding:"required"`
}

// AddHostsToGroup 将一个或者多个主机添加到某一个主机组中
func AddHostsToGroup(c *gin.Context) {
	var inputs APIAddHostsToGroupInput
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
	if group == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (id = %d) does not exist", inputs.GroupID))
		return
	}
	if !ctx.IsAdmin() && group.Creator != ctx.ID {
		h.JSONR(c, h.HTTPExpectationFailed, "Permission denied")
		return
	}
	tx := db.Begin()
	if err := tx.Where("ancestor_id = ? AND type = 2", group.ID).Delete(&model.Edge{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	for _, name := range inputs.Hosts {
		host := model.Host{
			Hostname: name,
		}
		if !cache.HostsMap.Has(host) {
			continue
		}
		edge := model.Edge{
			AncestorID:   group.ID,
			DescendantID: host.ID,
			Type:         2,
			Creator:      ctx.ID,
		}
		// 如果数据本身有问题则这里目前无法完全避免
		if cache.EdgesMap.Has(edge) {
			continue
		}
		if err := tx.Create(&edge).Error; err != nil {
			h.InternelError(c, "creating data", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	h.JSONR(c, "%v bind to hostgroup: %v")
	go cache.EdgesMap.Init()
}

// APIRemoveHostFromGroup 从一个主机组中移除一个主机对应的输入信息
type APIRemoveHostFromGroup struct {
	HostID  int64 `json:"hostID"  binding:"required"`
	GroupID int64 `json:"groupID" binding:"required"`
}

// RemoveHostFromGroup 从一个主机组中移除一个主机
func RemoveHostFromGroup(c *gin.Context) {
	var inputs APIRemoveHostFromGroup
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
	if group == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (id = %d) does not exist", inputs.GroupID))
		return
	}
	if !ctx.IsAdmin() && group.Creator != ctx.ID {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}
	if err := db.Where("ancestor_id = ? AND descendant_id = ? AND type = 2", inputs.GroupID, inputs.HostID).Delete(&model.Edge{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("unbind host: %v of hostgroup: %v", inputs.HostID, inputs.GroupID))
	go cache.EdgesMap.Init()
}

// DeleteGroup 根据分组名称删除一个主机分组
func DeleteGroup(c *gin.Context) {
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
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	group := cache.GroupsMap.Any(func(elem *model.Group) bool {
		if elem.ID == int64(groupID) {
			return true
		}
		return false
	})
	if group == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (id = %d) does not exist", groupID))
		return
	}
	if !ctx.IsAdmin() && group.Creator != ctx.ID {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}
	tx := db.Begin()
	// Delete groups referance of edges table
	if err := tx.Where("ancestor_id = ? AND type = 2", groupID).Delete(&model.Edge{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	// Delete plugins of group
	if err := tx.Where("group_id = ?", groupID).Delete(&model.Plugin{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	// Delete aggregators of groups
	if err := tx.Where("group_id = ?", groupID).Delete(&model.Cluster{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	// Delete from groups
	if err := tx.Where("id = ?", groupID).Delete(&model.Group{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("Group (id = %d) deleted", groupID))
	go cache.GroupsMap.Init()
	go cache.EdgesMap.Init()
	go cache.ClustersMap.Init()
}

// GetGroup 根据主机分组标识获取主机分组信息
func GetGroup(c *gin.Context) {
	inputGroupID := c.Params.ByName("id")
	q := c.DefaultQuery("q", ".+")
	if inputGroupID == "" {
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for group is missing")
		return
	}
	groupID, err := strconv.Atoi(inputGroupID)
	if err != nil {
		log.Debugf("[D] inputGroupID: %v", inputGroupID)
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	group := cache.GroupsMap.Any(func(elem *model.Group) bool {
		if elem.ID == int64(groupID) {
			return true
		}
		return false
	})
	if group == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (id = %d) does not exist", groupID))
		return
	}
	resp := map[string]interface{}{
		"group": group,
	}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.AncestorID == int64(groupID) && elem.Type == 2 {
			return true
		}
		return false
	})
	for _, r := range edges {
		hosts := cache.HostsMap.Filter(func(elem *model.Host) bool {
			if elem.ID == r.DescendantID {
				if ok, err := regexp.MatchString(q, elem.Hostname); ok == true && err == nil {
					return true
				}
			}
			return false
		})
		if hosts != nil {
			resp["hosts"] = hosts
		}
	}
	h.JSONR(c, resp)
}

// APIPutGroupInputs 修改一个已存在的主机分组对应的输入信息
type APIPutGroupInputs struct {
	ID   int64  `json:"id"   binding:"required"`
	Name string `json:"name" binding:"required"`
}

// PutGroup 修改一个已存在的主机分组
func PutGroup(c *gin.Context) {
	var inputs APIPutGroupInputs
	err := c.Bind(&inputs)
	switch {
	case err != nil:
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	case utils.HasDangerousCharacters(inputs.Name):
		h.JSONR(c, h.HTTPBadRequest, "name is invalid")
		return
	}
	if !cache.GroupsMap.Has(model.Group{ID: inputs.ID}) {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (id = %d) does not exist", inputs.ID))
		return
	}
	group := cache.GroupsMap.Any(func(elem *model.Group) bool {
		if elem.Name == inputs.Name && elem.ID != inputs.ID {
			return true
		}
		return false
	})
	if group != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (name = %s) already exists", inputs.Name))
		return
	}
	update := map[string]interface{}{
		"Name": inputs.Name,
	}

	if err := db.Model(model.Group{}).Where("id = ?", inputs.ID).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("Group (id = %d) updated", inputs.ID))
	go cache.GroupsMap.Init()
}

// APIBindTemplateToGroupInputs TODO:
type APIBindTemplateToGroupInputs struct {
	TemplateID int64 `json:"templateID"`
	GroupID    int64 `json:"groupID"`
}

// BindTemplateToGroup TODO:
func BindTemplateToGroup(c *gin.Context) {
	var inputs APIBindTemplateToGroupInputs
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
	if group == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (id = %d) does not exist", inputs.GroupID))
		return
	}
	template := cache.TemplatesMap.Any(func(elem *model.Template) bool {
		if elem.ID == inputs.TemplateID {
			return true
		}
		return false
	})
	if template == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Template (id = %d) does not exist", inputs.TemplateID))
		return
	}
	edge := model.Edge{
		AncestorID:   inputs.GroupID,
		DescendantID: inputs.TemplateID,
		Type:         3,
		Creator:      ctx.ID,
	}
	if cache.EdgesMap.Has(edge) {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Edge between group (id = %d) and template (id = %d) already exists", inputs.GroupID, inputs.TemplateID))
		return
	}
	if err := db.Create(&edge).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}
	h.JSONR(c, edge)
	go cache.EdgesMap.Init()
}

// APIUnBindTemplateToGroupInputs TODO:
type APIUnBindTemplateToGroupInputs struct {
	TemplateID int64 `json:"templateID"`
	GroupID    int64 `json:"groupID"`
}

// UnBindTemplateToGroup 解绑分组和模版
func UnBindTemplateToGroup(c *gin.Context) {
	var inputs APIUnBindTemplateToGroupInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	edge := cache.EdgesMap.Any(func(elem *model.Edge) bool {
		if elem.AncestorID == inputs.GroupID && elem.DescendantID == inputs.TemplateID && elem.Type == 3 {
			return true
		}
		return false
	})
	if edge == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Edge between group (id = %d) and template (id = %d) does not exist", inputs.GroupID, inputs.TemplateID))
		return
	}
	if !ctx.IsAdmin() && edge.Creator != ctx.ID {
		h.JSONR(c, h.HTTPBadRequest, errors.New("Permission denied"))
		return
	}
	if err := db.Delete(&edge).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("template: %v is unbind of HostGroup: %v", inputs.TemplateID, inputs.GroupID))
	go cache.EdgesMap.Init()
}

// GetTemplateOfGroup 获取分组及分组关联的模版信息
func GetTemplateOfGroup(c *gin.Context) {
	inputGroupID := c.Params.ByName("id")
	if inputGroupID == "" {
		log.Debug("[D] parameter `id` for group is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for group is missing")
		return
	}
	groupID, err := strconv.Atoi(inputGroupID)
	if err != nil {
		log.Debugf("[D] parameter `id` for group is invalid, value = %v", inputGroupID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for group is invalid, value = %v", inputGroupID))
		return
	}
	group := cache.GroupsMap.Any(func(elem *model.Group) bool {
		if elem.ID == int64(groupID) {
			return true
		}
		return false
	})
	if group == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (id = %d) does not exist", groupID))
		return
	}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.AncestorID == int64(groupID) && elem.Type == 3 {
			return true
		}
		return false
	})
	resp := map[string]interface{}{
		"group": group,
	}
	if edges != nil {
		templates := cache.TemplatesMap.Filter(func(elem *model.Template) bool {
			for _, r := range edges {
				if elem.ID == r.DescendantID {
					return true
				}
			}
			return false
		})
		resp["templates"] = templates
	}
	h.JSONR(c, resp)
}

// APIPatchHostInGroup TODO:
type APIPatchHostInGroup struct {
	Op    string   `json:"op"    binding:"required"`
	Hosts []string `json:"hosts" binding:"required"`
}

// PatchHostInGroup TODO:
func PatchHostInGroup(c *gin.Context) {
	var inputs APIPatchHostInGroup
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	inputGroupID := c.Params.ByName("id")
	if inputGroupID == "" {
		log.Debug("[D] parameter `id` for group is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for group is missing")
		return
	}
	groupID, err := strconv.Atoi(inputGroupID)
	if err != nil {
		log.Debugf("[D] parameter `id` for group is invalid, value = %v", inputGroupID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for group is invalid, value = %v", inputGroupID))
		return
	}

	if inputs.Op != "add" && inputs.Op != "remove" {
		h.JSONR(c, h.HTTPBadRequest, "Op must be add or remove")
		return
	}

	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	group := cache.GroupsMap.Any(func(elem *model.Group) bool {
		if elem.ID == int64(groupID) {
			return true
		}
		return false
	})
	if group == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Group (id = %d) does not exist", groupID))
		return
	}
	if !ctx.IsAdmin() && group.Creator != ctx.ID {
		h.JSONR(c, h.HTTPExpectationFailed, "Permission denied")
		return
	}

	switch inputs.Op {
	case "add":
		bindHostToGroup(c, ctx, group, inputs.Hosts)
	case "remove":
		unbindHostToGroup(c, group, inputs.Hosts)
	}
}

// bindHostToGroup 绑定主机和分组
func bindHostToGroup(c *gin.Context, ctx model.User, group *model.Group, hosts []string) {
	tx := db.Begin()
	for _, name := range hosts {
		host := cache.HostsMap.Any(func(elem *model.Host) bool {
			if elem.Hostname == name {
				return true
			}
			return false
		})
		// 主机不存在则忽略
		if host == nil {
			continue
		}

		edge := model.Edge{
			AncestorID:   group.ID,
			DescendantID: host.ID,
			Type:         2,
			Creator:      ctx.ID,
		}
		// 关联已存在则忽略
		if cache.EdgesMap.Has(edge) {
			continue
		}
		if err := tx.Create(&edge).Error; err != nil {
			h.InternelError(c, "creating data", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	h.JSONR(c, "ok")
	go cache.EdgesMap.Init()
}

// unbindHostToGroup 解绑主机和分组
func unbindHostToGroup(c *gin.Context, group *model.Group, hosts []string) {
	tx := db.Begin()
	for _, name := range hosts {
		host := cache.HostsMap.Any(func(elem *model.Host) bool {
			if elem.Hostname == name {
				return true
			}
			return false
		})
		if host == nil {
			continue
		}
		edge := cache.EdgesMap.Any(func(elem *model.Edge) bool {
			if elem.AncestorID == group.ID && elem.DescendantID == host.ID && elem.Type == 2 {
				return true
			}
			return false
		})
		if edge == nil {
			continue
		}
		if err := db.Delete(edge).Error; err != nil {
			h.InternelError(c, "deleting data", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	h.JSONR(c, "ok")
	go cache.EdgesMap.Init()
}
