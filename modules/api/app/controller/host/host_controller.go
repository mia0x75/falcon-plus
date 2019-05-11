package host

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/portal"
)

// GetHostBindToWhichHostGroup TODO:
func GetHostBindToWhichHostGroup(c *gin.Context) {
	HostIDTmp := c.Params.ByName("host_id")
	if HostIDTmp == "" {
		h.JSONR(c, badstatus, "host id is missing")
		return
	}
	hostID, err := strconv.Atoi(HostIDTmp)
	if err != nil {
		log.Debugf("[D] HostId: %v", HostIDTmp)
		h.JSONR(c, badstatus, err)
		return
	}
	grpHostMap := []f.GrpHost{}
	db.Falcon.Select("grp_id").Where("host_id = ?", hostID).Find(&grpHostMap)
	grpIDs := []int64{}
	for _, g := range grpHostMap {
		grpIDs = append(grpIDs, g.GrpID)
	}
	hostgroups := []f.HostGroup{}
	if len(grpIDs) != 0 {
		db.Falcon.Where("id in (？)", grpIDs).Find(&hostgroups)
	}
	h.JSONR(c, hostgroups)
	return
}

// GetHostGroupWithTemplate TODO:
func GetHostGroupWithTemplate(c *gin.Context) {
	grpIDtmp := c.Params.ByName("host_group")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		log.Debugf("[D] grpIDtmp: %v", grpIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	hostgroup := f.HostGroup{ID: int64(grpID)}
	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	hosts := []f.Host{}
	grpHosts := []f.GrpHost{}
	if dt := db.Falcon.Where("grp_id = ?", grpID).Find(&grpHosts); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	for _, grph := range grpHosts {
		var host f.Host
		db.Falcon.Find(&host, grph.HostID)
		if host.ID != 0 {
			hosts = append(hosts, host)
		}
	}
	h.JSONR(c, map[string]interface{}{
		"hostgroup": hostgroup,
		"hosts":     hosts,
	})
	return
}

// GetGrpsRelatedHost TODO:
func GetGrpsRelatedHost(c *gin.Context) {
	hostIDtmp := c.Params.ByName("host_id")
	if hostIDtmp == "" {
		h.JSONR(c, badstatus, "host id is missing")
		return
	}
	hostID, err := strconv.Atoi(hostIDtmp)
	if err != nil {
		log.Debugf("[D] host id: %v", hostIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}

	host := f.Host{}
	if dt := db.Falcon.Where("id = ?", hostID).Find(&host); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	grps := host.RelatedGrp()
	h.JSONR(c, grps)
	return
}

// GetTplsRelatedHost TODO:
func GetTplsRelatedHost(c *gin.Context) {
	hostIDtmp := c.Params.ByName("host_id")
	if hostIDtmp == "" {
		h.JSONR(c, badstatus, "host id is missing")
		return
	}
	hostID, err := strconv.Atoi(hostIDtmp)
	if err != nil {
		log.Debugf("[D] host id: %v", hostIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	host := f.Host{}
	if dt := db.Falcon.Where("id = ?", hostID).Find(&host); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	tpls := host.RelatedTpl()
	log.Debugf("[D] hostid: %d", host.ID)
	h.JSONR(c, tpls)
	return
}

// GetGrpsRelatedEndpoint TODO:
func GetGrpsRelatedEndpoint(c *gin.Context) {
	hostNameTmp := c.Params.ByName("endpoint_name")
	if hostNameTmp == "" {
		h.JSONR(c, badstatus, "endpoint is missing")
		return
	}
	ahost := f.Host{Hostname: hostNameTmp}
	var hostID int64
	var ok bool
	if hostID, ok = ahost.Existing(); ok {
		host := f.Host{ID: int64(hostID)}
		grps := host.RelatedGrp()
		h.JSONR(c, grps)
		return
	}
	h.JSONR(c, badstatus, "endpoint is missing")
	return
}

// GetHosts TODO:
func GetHosts(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	pageTmp := c.DefaultQuery("page", "")
	limitTmp := c.DefaultQuery("limit", "")
	q := c.DefaultQuery("q", ".+")
	page, limit, err = h.PageParser(pageTmp, limitTmp)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	var hosts []f.Host
	var dt *gorm.DB
	if limit != -1 && page != -1 {
		dt = db.Falcon.Raw("SELECT * from host where hostname regexp ? limit ?,?", q, page, limit).Scan(&hosts)
	} else {
		dt = db.Falcon.Table("host").Where("hostname regexp ?", q).Find(&hosts)
	}
	if dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, hosts)
	return
}
