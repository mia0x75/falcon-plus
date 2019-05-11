package dashboard_graph

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	m "github.com/open-falcon/falcon-plus/modules/api/app/model/dashboard"
)

// APITmpGraphCreateReqData TODO:
type APITmpGraphCreateReqData struct {
	Endpoints []string `json:"endpoints" binding:"required"`
	Counters  []string `json:"counters" binding:"required"`
}

// DashboardTmpGraphCreate TODO:
func DashboardTmpGraphCreate(c *gin.Context) {
	var inputs APITmpGraphCreateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	es := inputs.Endpoints
	cs := inputs.Counters
	sort.Strings(es)
	sort.Strings(cs)

	esString := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	csString := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
	ck := cutils.Md5(esString + ":" + csString)

	dt := db.Dashboard.Exec("insert ignore into `tmp_graph` (endpoints, counters, ck) values(?, ?, ?) on duplicate key update time_=?", esString, csString, ck, time.Now())
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	tmpGraph := m.DashboardTmpGraph{}
	dt = db.Dashboard.Table("tmp_graph").Where("ck=?", ck).First(&tmpGraph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, map[string]int{"id": int(tmpGraph.ID)})
}

// DashboardTmpGraphQuery TODO:
func DashboardTmpGraphQuery(c *gin.Context) {
	id := c.Param("id")

	tmpGraph := m.DashboardTmpGraph{}
	dt := db.Dashboard.Table("tmp_graph").Where("id = ?", id).First(&tmpGraph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	es := strings.Split(tmpGraph.Endpoints, TMP_GRAPH_FILED_DELIMITER)
	cs := strings.Split(tmpGraph.Counters, TMP_GRAPH_FILED_DELIMITER)

	ret := map[string][]string{
		"endpoints": es,
		"counters":  cs,
	}

	h.JSONR(c, ret)
}

// APIGraphCreateReqData TODO:
type APIGraphCreateReqData struct {
	ScreenID    int      `json:"screen_id" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Endpoints   []string `json:"endpoints" binding:"required"`
	Counters    []string `json:"counters" binding:"required"`
	TimeSpan    int      `json:"timespan"`
	RelativeDay int      `json:"relativeday"`
	GraphType   string   `json:"graph_type"`
	Method      string   `json:"method"`
	Position    int      `json:"position"`
	FalconTags  string   `json:"falcon_tags"`
}

// DashboardGraphCreate TODO:
func DashboardGraphCreate(c *gin.Context) {
	var inputs APIGraphCreateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	es := inputs.Endpoints
	cs := inputs.Counters
	sort.Strings(es)
	sort.Strings(cs)
	esString := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	csString := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)

	d := m.DashboardGraph{
		Title:       inputs.Title,
		Hosts:       esString,
		Counters:    csString,
		ScreenId:    int64(inputs.ScreenID),
		TimeSpan:    inputs.TimeSpan,
		RelativeDay: inputs.RelativeDay,
		GraphType:   inputs.GraphType,
		Method:      inputs.Method,
		Position:    inputs.Position,
	}
	if d.TimeSpan == 0 {
		d.TimeSpan = 3600
	}
	if d.GraphType == "" {
		d.GraphType = "h"
	}

	tx := db.Dashboard.Begin()
	dt := tx.Table("dashboard_graph").Create(&d)
	if dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	var lid []int
	dt = tx.Table("dashboard_graph").Raw("select LAST_INSERT_ID() as id").Pluck("id", &lid)
	if dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	tx.Commit()
	aid := lid[0]

	h.JSONR(c, map[string]int{"id": aid})

}

// APIGraphUpdateReqData TODO:
type APIGraphUpdateReqData struct {
	ScreenID    int      `json:"screen_id"`
	Title       string   `json:"title"`
	Endpoints   []string `json:"endpoints"`
	Counters    []string `json:"counters"`
	TimeSpan    int      `json:"timespan"`
	RelativeDay int      `json:"relativeday"`
	GraphType   string   `json:"graph_type"`
	Method      string   `json:"method"`
	Position    int      `json:"position"`
	FalconTags  string   `json:"falcon_tags"`
}

// DashboardGraphUpdate TODO:
func DashboardGraphUpdate(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	var inputs APIGraphUpdateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	d := m.DashboardGraph{}

	if len(inputs.Endpoints) != 0 {
		es := inputs.Endpoints
		sort.Strings(es)
		esString := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
		d.Hosts = esString
	}
	if len(inputs.Counters) != 0 {
		cs := inputs.Counters
		sort.Strings(cs)
		csString := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
		d.Counters = csString
	}
	if inputs.Title != "" {
		d.Title = inputs.Title
	}
	if inputs.ScreenID != 0 {
		d.ScreenId = int64(inputs.ScreenID)
	}
	if inputs.TimeSpan != 0 {
		d.TimeSpan = inputs.TimeSpan
	}
	if inputs.RelativeDay >= 0 {
		d.RelativeDay = inputs.RelativeDay
	}
	if inputs.GraphType != "" {
		d.GraphType = inputs.GraphType
	}
	if inputs.Method != "" {
		d.Method = inputs.Method
	}
	if inputs.Position != 0 {
		d.Position = inputs.Position
	}
	if inputs.FalconTags != "" {
		d.FalconTags = inputs.FalconTags
	}

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Model(&graph).Where("id = ?", gid).Updates(d)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, map[string]int{"id": gid})
}

// DashboardGraphGet TODO:
func DashboardGraphGet(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("id = ?", gid).First(&graph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
	cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)

	h.JSONR(c, map[string]interface{}{
		"graph_id":    graph.ID,
		"title":       graph.Title,
		"endpoints":   es,
		"counters":    cs,
		"screen_id":   graph.ScreenId,
		"graph_type":  graph.GraphType,
		"timespan":    graph.TimeSpan,
		"relativeday": graph.RelativeDay,
		"method":      graph.Method,
		"position":    graph.Position,
		"falcon_tags": graph.FalconTags,
	})
}

// DashboardGraphDelete TODO:
func DashboardGraphDelete(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("id = ?", gid).Delete(&graph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, map[string]int{"id": gid})
}

// DashboardGraphGetsByScreenID TODO:
func DashboardGraphGetsByScreenID(c *gin.Context) {
	id := c.Param("screen_id")
	sid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen id")
		return
	}
	limit := c.DefaultQuery("limit", "500")

	graphs := []m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("screen_id = ?", sid).Limit(limit).Find(&graphs)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	ret := []map[string]interface{}{}
	for _, graph := range graphs {
		es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
		cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)

		r := map[string]interface{}{
			"graph_id":    graph.ID,
			"title":       graph.Title,
			"endpoints":   es,
			"counters":    cs,
			"screen_id":   graph.ScreenId,
			"graph_type":  graph.GraphType,
			"timespan":    graph.TimeSpan,
			"relativeday": graph.RelativeDay,
			"method":      graph.Method,
			"position":    graph.Position,
			"falcon_tags": graph.FalconTags,
		}
		ret = append(ret, r)
	}

	h.JSONR(c, ret)
}
