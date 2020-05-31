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

// APIGetTemplatesOutput TODO:
type APIGetTemplatesOutput struct {
	Templates []CTemplate `json:"templates"`
}

// CTemplate TODO:
type CTemplate struct {
	Template   model.Template `json:"template"`
	ParentName string         `json:"parent_name"`
}

// GetTemplates TODO:
func GetTemplates(c *gin.Context) {
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
	var templates []model.Template
	q := c.DefaultQuery("q", ".+")
	if limit != -1 && page != -1 {
		err = db.Raw("SELECT * FROM templates WHERE name REGEXP ? LIMIT ?, ?", q, page, limit).Scan(&templates).Error
	} else {
		err = db.Where("name REGEXP ?", q).Find(&templates).Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	f := func(t model.Template) (name string) {
		if t.ParentID <= 0 {
			return
		}
		template := cache.TemplatesMap.Any(func(elem *model.Template) bool {
			if elem.ID == t.ParentID {
				return true
			}
			return false
		})
		name = template.Name
		return
	}
	output := APIGetTemplatesOutput{}
	output.Templates = []CTemplate{}
	for _, t := range templates {
		output.Templates = append(output.Templates, CTemplate{
			Template:   t,
			ParentName: f(t),
		})
	}
	h.JSONR(c, output)
}

// GetTemplatesSimple TODO:
func GetTemplatesSimple(c *gin.Context) {
	templates := []model.Template{}
	q := c.DefaultQuery("q", ".+")
	if err := db.Select("id, name").Where("name REGEXP ?", q).Find(&templates).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, templates)
}

// GetTemplate TODO:
func GetTemplate(c *gin.Context) {
	inputTemplateID := c.Params.ByName("id")
	if inputTemplateID == "" {
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for template is missing")
		return
	}
	templateID, err := strconv.Atoi(inputTemplateID)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	template := cache.TemplatesMap.Any(func(elem *model.Template) bool {
		if elem.ID == int64(templateID) {
			return true
		}
		return false
	})
	stratgies := cache.StrategiesMap.Filter(func(elem *model.Strategy) bool {
		if elem.TemplateID == int64(templateID) {
			return true
		}
		return false
	})
	action := cache.ActionsMap.Any(func(elem *model.Action) bool {
		if elem.ID == template.ActionID {
			return true
		}
		return false
	})

	parentTemplateName := ""
	if p := cache.TemplatesMap.Any(func(elem *model.Template) bool {
		if elem.ID == template.ParentID {
			return true
		}
		return false
	}); p != nil {
		parentTemplateName = p.Name
	}
	resp := map[string]interface{}{
		"template":    template,
		"stratgies":   stratgies,
		"action":      action,
		"parent_name": parentTemplateName,
	}
	h.JSONR(c, resp)
}

// APICreateTemplateInput TODO:
type APICreateTemplateInput struct {
	Name     string `json:"name"     binding:"required"`
	ParentID int64  `json:"parentID" binding:"required"`
	ActionID int64  `json:"actionID"`
}

// CreateTemplate TODO:
func CreateTemplate(c *gin.Context) {
	var inputs APICreateTemplateInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	// TODO: inputs.Name == "" ??
	if inputs.Name == "" {
		h.JSONR(c, h.HTTPBadRequest, "input name is empty, please check it")
		return
	}
	template := model.Template{
		Name:     inputs.Name,
		ParentID: inputs.ParentID,
		ActionID: inputs.ActionID,
		Creator:  ctx.ID,
	}
	if err := db.Save(&template).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}
	h.JSONR(c, template)
	go cache.TemplatesMap.Init()
}

// APIUpdateTemplateInput TODO:
type APIUpdateTemplateInput struct {
	Name       string `json:"name"       binding:"required"`
	ParentID   int64  `json:"parentID"   binding:"required"`
	TemplateID int64  `json:"templateID" binding:"required"`
}

// UpdateTemplate TODO:
func UpdateTemplate(c *gin.Context) {
	var inputs APIUpdateTemplateInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	template := cache.TemplatesMap.Any(func(elem *model.Template) bool {
		if elem.ID == int64(inputs.TemplateID) {
			return true
		}
		return false
	})
	if template == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Template (id = %d) does not exist", inputs.TemplateID))
		return
	}
	if template.Creator != ctx.ID && !ctx.IsAdmin() {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}

	update := map[string]interface{}{
		"Name":     inputs.Name,
		"ParentID": inputs.ParentID,
	}
	if err := db.Model(template).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, template)
	go cache.TemplatesMap.Init()
}

// DeleteTemplate TODO:
func DeleteTemplate(c *gin.Context) {
	inputTemplateID, _ := c.Params.Get("id")
	if inputTemplateID == "" {
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for template is missing")
		return
	}
	templateID, err := strconv.Atoi(inputTemplateID)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	tx := db.Begin()
	var template model.Template
	if err := tx.Find(&template, templateID).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		tx.Rollback()
		return
	}
	// delete template
	actionID := template.ActionID
	if err := tx.Delete(&template).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	// delete action
	if actionID != 0 {
		if err := tx.Delete(&model.Action{}, actionID).Error; err != nil {
			h.InternelError(c, "deleting data", err)
			tx.Rollback()
			return
		}
	}
	// delete strategy
	if err := tx.Where("template_id = ?", templateID).Delete(&model.Strategy{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	// delete link
	if err := tx.Where("descendant_id = ? AND type = 3", templateID).Delete(&model.Edge{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("Template (id = %d) deleted", templateID))
	go cache.TemplatesMap.Init()
	go cache.EdgesMap.Init()
}

// GetATemplateHostgroup TODO:
func GetATemplateHostgroup(c *gin.Context) {
	inputTemplateID := c.Params.ByName("id")
	if inputTemplateID == "" {
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for template is missing")
		return
	}
	templateID, err := strconv.Atoi(inputTemplateID)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	template := cache.TemplatesMap.Any(func(elem *model.Template) bool {
		if elem.ID == int64(templateID) {
			return true
		}
		return false
	})
	groups := []*model.Group{}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.DescendantID == int64(templateID) && elem.Type == 3 {
			return true
		}
		return false
	})
	if len(edges) != 0 {
		groups = cache.GroupsMap.Filter(func(elem *model.Group) bool {
			for _, r := range edges {
				if elem.ID == r.AncestorID {
					return true
				}
			}
			return false

		})
	}
	resp := map[string]interface{}{
		"template": template,
		"groups":   groups,
	}
	h.JSONR(c, resp)
}

// APICreateActionToTmplateInput TODO:
type APICreateActionToTmplateInput struct {
	UIC                string `json:"uic"                binding:"required"`
	URL                string `json:"url"                binding:"required"`
	Callback           int    `json:"callback"           binding:"required"`
	BeforeCallbackSMS  int    `json:"beforeCallbackSMS"  binding:"required"`
	AfterCallbackSMS   int    `json:"afterCallbackSMS"   binding:"required"`
	BeforeCallbackMail int    `json:"beforeCallbackMail" binding:"required"`
	AfterCallbackMail  int    `json:"afterCallbackMail"  binding:"required"`
	TemplateID         int64  `json:"templateID"         binding:"required"`
}

// CreateActionToTmplate TODO:
func CreateActionToTmplate(c *gin.Context) {
	var inputs APICreateActionToTmplateInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	action := model.Action{
		UIC:                inputs.UIC,
		URL:                inputs.URL,
		Callback:           inputs.Callback,
		BeforeCallbackSMS:  inputs.BeforeCallbackSMS,
		BeforeCallbackMail: inputs.BeforeCallbackMail,
		AfterCallbackMail:  inputs.AfterCallbackMail,
		AfterCallbackSMS:   inputs.AfterCallbackSMS,
	}
	tx := db.Begin()
	if err := tx.Save(&action).Error; err != nil {
		h.InternelError(c, "creating data", err)
		tx.Rollback()
		return
	}
	// TODO: lid[0] == action.ID ?
	var lid []int
	tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &lid)
	inputActionID := lid[0]
	var tpl model.Template
	if err := tx.Find(&tpl, inputs.TemplateID).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		tx.Rollback()
		return
	}

	if err := tx.Model(&tpl).UpdateColumns(model.Template{ActionID: int64(inputActionID)}).Error; err != nil {
		h.InternelError(c, "updating data", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("action is created and bind to template: %d", inputs.TemplateID))
}

// APIUpdateActionToTmplateInput TODO:
type APIUpdateActionToTmplateInput struct {
	ID                 int64  `json:"id"                 binding:"required"`
	UIC                string `json:"uic"                binding:"required"`
	URL                string `json:"url"                binding:"required"`
	Callback           int    `json:"callback"           binding:"required"`
	BeforeCallbackSMS  int    `json:"beforeCallbackSMS"  binding:"required"`
	AfterCallbackSMS   int    `json:"afterCallbackSMS"   binding:"required"`
	BeforeCallbackMail int    `json:"beforeCallbackMail" binding:"required"`
	AfterCallbackMail  int    `json:"afterCallbackMail"  binding:"required"`
}

// UpdateActionToTmplate TODO：
func UpdateActionToTmplate(c *gin.Context) {
	var inputs APIUpdateActionToTmplateInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	var action model.Action
	tx := db.Begin()
	if err := tx.Find(&action, inputs.ID).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		tx.Rollback()
		return
	}

	update := map[string]interface{}{
		"UIC":                inputs.UIC,
		"URL":                inputs.URL,
		"Callback":           inputs.Callback,
		"BeforeCallbackSMS":  inputs.BeforeCallbackSMS,
		"BeforeCallbackMail": inputs.BeforeCallbackMail,
		"AfterCallbackMail":  inputs.AfterCallbackMail,
		"AfterCallbackSMS":   inputs.AfterCallbackSMS,
	}
	if err := tx.Model(&action).Where("id = ?", inputs.ID).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("Action (id = %d) updated", inputs.ID))
}

// GetAction TODO:
func GetAction(c *gin.Context) {
	inputActionID := c.Param("id")
	if inputActionID == "" {
		log.Debug("[D] parameter `id` for action is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for action is missing")
		return
	}
	actionID, err := strconv.Atoi(inputActionID)
	if err != nil {
		log.Debugf("[D] parameter `id` for action is invalid, value = %v", inputActionID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for action is invalid, value = %v", inputActionID))
		return
	}

	action := cache.ActionsMap.Any(func(elem *model.Action) bool {
		if elem.ID == int64(actionID) {
			return true
		}
		return false
	})

	h.JSONR(c, action)
}
