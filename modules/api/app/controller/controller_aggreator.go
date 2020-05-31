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

// GetAggregatorListOfGroup TODO:
func GetAggregatorListOfGroup(c *gin.Context) {
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

	aggregators := []model.Cluster{}
	if limit != -1 && page != -1 {
		err = db.Raw("SELECT * FROM clusters WHERE group_id = ? LIMIT ?, ?", groupID, page, limit).Scan(&aggregators).Error
	} else {
		err = db.Where("group_id = ?", groupID).Find(&aggregators).Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	if len(aggregators) != 0 {
		resp["aggregators"] = aggregators
	}
	h.JSONR(c, resp)
}

// GetAggregator TODO:
func GetAggregator(c *gin.Context) {
	inputAggregatorID := c.Params.ByName("id")
	if inputAggregatorID == "" {
		log.Debug("[D] parameter `id` for aggregator is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for aggregator is missing")
		return
	}
	aggregatorID, err := strconv.Atoi(inputAggregatorID)
	if err != nil {
		log.Debugf("[D] parameter `id` for aggregator is invalid, value = %v", inputAggregatorID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for aggregator is invalid, value = %v", inputAggregatorID))
		return
	}

	aggregator := cache.ClustersMap.Any(func(elem *model.Cluster) bool {
		if elem.ID == int64(aggregatorID) {
			return true
		}
		return false
	})
	h.JSONR(c, aggregator)
}

// APICreateAggregatorInput TODO:
type APICreateAggregatorInput struct {
	GroupID     int64  `json:"groupID"     binding:"required"`
	Numerator   string `json:"numerator"   binding:"required"`
	Denominator string `json:"denominator" binding:"required"`
	Endpoint    string `json:"endpoint"    binding:"required"`
	Metric      string `json:"metric"      binding:"required"`
	Tags        string `json:"tags"        binding:"required"`
	Step        int    `json:"step"        binding:"required"`
}

// CreateAggregator TODO:
func CreateAggregator(c *gin.Context) {
	var inputs APICreateAggregatorInput
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
		if elem.ID == int64(inputs.GroupID) {
			return true
		}
		return false
	})
	if group == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Invalid group (id = %d)", inputs.GroupID))
		return
	}
	if !ctx.IsAdmin() {
		if group.Creator != ctx.ID {
			h.JSONR(c, h.HTTPBadRequest, "Permission denied")
			return
		}
	}
	aggregator := model.Cluster{
		GroupID:     inputs.GroupID,
		Numerator:   inputs.Numerator,
		Denominator: inputs.Denominator,
		Endpoint:    inputs.Endpoint,
		Metric:      inputs.Metric,
		Tags:        inputs.Tags,
		DsType:      "GAUGE",
		Step:        inputs.Step,
		Creator:     ctx.ID,
	}
	if err := db.Create(&aggregator).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}
	h.JSONR(c, aggregator)
	go cache.ClustersMap.Init()
}

// APIUpdateAggregatorInput TODO:
type APIUpdateAggregatorInput struct {
	ID          int64  `json:"id"          binding:"required"`
	Numerator   string `json:"numerator"   binding:"required"`
	Denominator string `json:"denominator" binding:"required"`
	Endpoint    string `json:"endpoint"    binding:"required"`
	Metric      string `json:"metric"      binding:"required"`
	Tags        string `json:"tags"        binding:"required"`
	Step        int    `json:"step"        binding:"required"`
}

// UpdateAggregator TODO:
func UpdateAggregator(c *gin.Context) {
	var inputs APIUpdateAggregatorInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	aggregator := cache.ClustersMap.Any(func(elem *model.Cluster) bool {
		if elem.ID == inputs.ID {
			return true
		}
		return false
	})
	if aggregator == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Aggregator (id = %d) does not exist", inputs.ID))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}
	if !ctx.IsAdmin() && aggregator.Creator != ctx.ID {
		group := cache.GroupsMap.Any(func(elem *model.Group) bool {
			if elem.ID == aggregator.GroupID {
				return true
			}
			return false
		})
		// only admin & aggregator creator can update it
		if group == nil || group.Creator != ctx.ID {
			h.JSONR(c, h.HTTPBadRequest, "Permission denied")
			return
		}
	}
	update := map[string]interface{}{
		"Numerator":   inputs.Numerator,
		"Denominator": inputs.Denominator,
		"Endpoint":    inputs.Endpoint,
		"Metric":      inputs.Metric,
		"Tags":        inputs.Tags,
		"Step":        inputs.Step,
	}
	if err := db.Model(aggregator).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, aggregator)
	go cache.ClustersMap.Init()
}

// DeleteAggregator TODO:
func DeleteAggregator(c *gin.Context) {
	inputAggregatorID := c.Params.ByName("id")
	if inputAggregatorID == "" {
		log.Debug("[D] parameter `id` for aggregator is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for aggregator is missing")
		return
	}
	aggregatorID, err := strconv.Atoi(inputAggregatorID)
	if err != nil {
		log.Debugf("[D] parameter `id` for aggregator is invalid, value = %v", inputAggregatorID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for aggregator is invalid, value = %v", inputAggregatorID))
		return
	}
	aggregator := cache.ClustersMap.Any(func(elem *model.Cluster) bool {
		if elem.ID == int64(aggregatorID) {
			return true
		}
		return false
	})
	if aggregator == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Aggregator (id = %d) does not exist", aggregatorID))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}
	if !ctx.IsAdmin() && aggregator.Creator != ctx.ID {
		group := cache.GroupsMap.Any(func(elem *model.Group) bool {
			if elem.ID == aggregator.GroupID {
				return true
			}
			return false
		})
		if group == nil || group.Creator != ctx.ID {
			h.JSONR(c, h.HTTPBadRequest, "Permission denied")
			return
		}
	}

	if err := db.Delete(aggregator).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("Aggregator (id = %d) deleted", aggregatorID))
	go cache.ClustersMap.Init()
}
