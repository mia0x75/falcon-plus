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

// GetHostBindToWhichHostGroup TODO:
func GetHostBindToWhichHostGroup(c *gin.Context) {
	inputHostID := c.Params.ByName("id")
	if inputHostID == "" {
		log.Debug("[D] parameter `id` for host is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for host is missing")
		return
	}
	hostID, err := strconv.Atoi(inputHostID)
	if err != nil {
		log.Debugf("[D] parameter `id` for host is invalid, value = %v", inputHostID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for host is invalid, value = %v", inputHostID))
		return
	}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.DescendantID == int64(hostID) && elem.Type == 2 {
			return true
		}
		return false
	})
	if edges == nil {
		h.JSONR(c, nil)
		return
	}
	groups := cache.GroupsMap.Filter(func(elem *model.Group) bool {
		for _, r := range edges {
			if elem.ID == r.AncestorID {
				return true
			}
		}
		return false
	})
	h.JSONR(c, groups)
}

// GetGroupWithTemplate 获取分组及分组上关联的模版信息
func GetGroupWithTemplate(c *gin.Context) {
	inputGroupID := c.Params.ByName("host_group")
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
	resp := map[string]interface{}{
		"group": group,
	}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.AncestorID == int64(groupID) && elem.Type == 2 {
			return true
		}
		return false
	})
	if edges != nil {
		hosts := cache.HostsMap.Filter(func(elem *model.Host) bool {
			for _, r := range edges {
				if elem.ID == r.DescendantID {
					return true
				}
			}
			return false
		})
		resp["hosts"] = hosts
	}
	h.JSONR(c, resp)
}

// GetHostRelatedGroups 获取主机关联的分组信息
func GetHostRelatedGroups(c *gin.Context) {
	inputHostID := c.Params.ByName("id")
	if inputHostID == "" {
		log.Debug("[D] parameter `id` for host is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for host is missing")
		return
	}
	hostID, err := strconv.Atoi(inputHostID)
	if err != nil {
		log.Debugf("[D] parameter `id` for host is invalid, value = %v", inputHostID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for host is invalid, value = %v", inputHostID))
		return
	}

	host := cache.HostsMap.Any(func(elem *model.Host) bool {
		if elem.ID == int64(hostID) {
			return true
		}
		return false
	})
	if host == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Host (id = %d) does not exist", hostID))
		return
	}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.DescendantID == host.ID && elem.Type == 2 {
			return true
		}
		return false
	})
	if edges == nil {
		h.JSONR(c, nil)
		return
	}
	groups := cache.GroupsMap.Filter(func(elem *model.Group) bool {
		for _, r := range edges {
			if r.AncestorID == elem.ID {
				return true
			}
		}
		return false
	})
	if groups == nil {
		h.JSONR(c, nil)
		return
	}
	h.JSONR(c, groups)
}

// GetHostRelatedTemplates 获取主机关联的模版信息
func GetHostRelatedTemplates(c *gin.Context) {
	inputHostID := c.Params.ByName("id")
	if inputHostID == "" {
		log.Debug("[D] parameter `id` for host is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for host is missing")
		return
	}
	hostID, err := strconv.Atoi(inputHostID)
	if err != nil {
		log.Debugf("[D] parameter `id` for host is invalid, value = %v", inputHostID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for host is invalid, value = %v", inputHostID))
		return
	}
	host := cache.HostsMap.Any(func(elem *model.Host) bool {
		if elem.ID == int64(hostID) {
			return true
		}
		return false
	})
	if host == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Host (id = %d) does not exist", hostID))
		return
	}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.DescendantID == host.ID && elem.Type == 3 {
			return true
		}
		return false
	})
	if edges == nil {
		h.JSONR(c, nil)
		return
	}
	templates := cache.TemplatesMap.Filter(func(elem *model.Template) bool {
		for _, r := range edges {
			if r.AncestorID == elem.ID {
				return true
			}
		}
		return false
	})
	if templates == nil {
		h.JSONR(c, nil)
		return
	}
	h.JSONR(c, templates)
}

// GetEndpointRelatedGroups 获取Endpoint关联的分组信息
func GetEndpointRelatedGroups(c *gin.Context) {
	inputHostName := c.Params.ByName("name")
	if inputHostName == "" {
		h.JSONR(c, h.HTTPBadRequest, "endpoint is missing")
		return
	}
	host := cache.HostsMap.Any(func(elem *model.Host) bool {
		if elem.Hostname == inputHostName {
			return true
		}
		return false
	})
	if host == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Host (name = %s) does not exist", inputHostName))
		return
	}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.DescendantID == host.ID && elem.Type == 2 {
			return true
		}
		return false
	})
	if edges == nil {
		h.JSONR(c, nil)
		return
	}
	groups := cache.GroupsMap.Filter(func(elem *model.Group) bool {
		for _, r := range edges {
			if r.AncestorID == elem.ID {
				return true
			}
		}
		return false
	})
	if groups == nil {
		h.JSONR(c, nil)
		return
	}
	h.JSONR(c, groups)
}

// GetHosts 分页获取主机信息
func GetHosts(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	inputPage := c.DefaultQuery("page", "")
	inputLimit := c.DefaultQuery("limit", "")
	q := c.DefaultQuery("q", ".+")
	page, limit, err = h.PageParser(inputPage, inputLimit)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err.Error())
		return
	}
	var hosts []model.Host
	if limit != -1 && page != -1 {
		err = db.Raw("SELECT * FROM hosts WHERE hostname regexp ? LIMIT ?, ?", q, page, limit).Scan(&hosts).Error
	} else {
		err = db.Where("hostname REGEXP ?", q).Find(&hosts).Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, hosts)
}

// GetMaintain TODO:
func GetMaintain(c *gin.Context) {
	var hosts = []model.Host{}
	if err := db.
		Where("maintain_begin != 0 and (select count(*) from edges l where l.descendant_id = hosts.id and l.type = 2) != 0").
		Find(&hosts).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, hosts)
}

// APIFindByMetricInput TODO:
type APIFindByMetricInput struct {
	Metric string `json:"metric"`
}

// APIMetricOutput TODO:
type APIMetricOutput struct {
	Strategy *model.Strategy `json:"strategy"`
	Hosts    []string        `json:"hosts"`
}

// FindByMetric TODO:
func FindByMetric(c *gin.Context) {
	var inputs APIFindByMetricInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	outputs := []APIMetricOutput{}
	strategies := cache.StrategiesMap.Filter(func(elem *model.Strategy) bool {
		if elem.Metric == inputs.Metric {
			return true
		}
		return false
	})
	for _, s := range strategies {
		edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
			if elem.DescendantID == s.TemplateID && elem.Type == 3 {
				return true
			}
			return false
		})
		if edges == nil {
			continue
		}
		var names = []string{}
		for _, r := range edges {
			edges2 := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
				if elem.AncestorID == r.AncestorID && elem.Type == 2 {
					return true
				}
				return false
			})
			if edges2 == nil {
				continue
			}
			hosts := cache.HostsMap.Filter(func(elem *model.Host) bool {
				for _, r := range edges2 {
					if elem.ID == r.DescendantID {
						return true
					}
				}
				return false
			})
			for _, h := range hosts {
				names = append(names, h.Hostname)
			}
		}
		if len(names) == 0 {
			continue
		}
		outputs = append(outputs, APIMetricOutput{
			Strategy: s,
			Hosts:    names,
		})
	}
	h.JSONR(c, outputs)
}
