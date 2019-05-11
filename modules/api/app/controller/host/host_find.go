package host

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/portal"
)

// FindByMaintain TODO:
func FindByMaintain(c *gin.Context) {
	var dt *gorm.DB
	var hosts = []f.Host{}
	if dt = db.Falcon.Where("maintain_begin != 0 and (select count(*) from grp_host where grp_host.host_id = host.id) != 0").Find(&hosts); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, hosts)
}

// APIFindByMetricInput TODO:
type APIFindByMetricInput struct {
	Metric string `json:"metric"`
}

// RMetric TODO:
type RMetric struct {
	Strategy f.Strategy `json:"strategy"`
	Hosts    []string   `json:"hosts"`
}

// FindByMetric TODO:
func FindByMetric(c *gin.Context) {
	var inputs APIFindByMetricInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	ret := []RMetric{}
	stgs := []f.Strategy{}
	var dt *gorm.DB
	if dt = db.Falcon.Where("metric = ?", inputs.Metric).Find(&stgs); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	for _, stg := range stgs {
		var grpTpls = []f.GrpTpl{}
		if dt = db.Falcon.Where("tpl_id = ?", stg.TplId).Find(&grpTpls); dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			return
		}
		if len(grpTpls) == 0 {
			continue
		}
		var hosts = []string{}
		for _, tpl := range grpTpls {
			var tmpHosts = []f.Host{}
			if dt = db.Falcon.Joins("JOIN grp_host on host.id = grp_host.host_id AND grp_host.grp_id = ?", tpl.GrpID).Find(&tmpHosts); dt.Error != nil {
				h.JSONR(c, badstatus, dt.Error)
				return
			}
			for _, host := range tmpHosts {
				hosts = append(hosts, host.Hostname)
			}
		}
		if len(hosts) == 0 {
			continue
		}
		ret = append(ret, RMetric{
			Strategy: stg,
			Hosts:    hosts,
		})
	}
	h.JSONR(c, ret)
}
