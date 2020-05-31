package controller

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	tcache "github.com/toolkits/cache/localcache/timedcache"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	grh "github.com/open-falcon/falcon-plus/modules/api/graph"
)

var (
	localStepCache = tcache.New(600*time.Second, 60*time.Second)
)

// APITmpGraphCreateReqData TODO:
type APITmpGraphCreateReqData struct {
	Endpoints []string `json:"endpoints" binding:"required"`
	Counters  []string `json:"counters"  binding:"required"`
}

// CreateTempGraph TODO:
func CreateTempGraph(c *gin.Context) {
	var inputs APITmpGraphCreateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	es := inputs.Endpoints
	cs := inputs.Counters
	sort.Strings(es)
	sort.Strings(cs)

	esString := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	csString := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
	sign := cu.Md5(esString + ":" + csString)

	draft := model.Draft{
		Endpoints: esString,
		Counters:  csString,
		Sign:      sign,
	}
	if err := db.Where("sign = ?", sign).First(&draft).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			h.InternelError(c, "retrieving data", err)
			return
		}
	}
	if draft.ID == 0 {
		if err := db.Create(&draft).Error; err != nil {
			h.InternelError(c, "creating data", err)
			return
		}
	}
	resp := map[string]interface{}{
		"id": draft.ID,
	}
	h.JSONR(c, resp)
}

// GetTempGraph TODO:
func GetTempGraph(c *gin.Context) {
	inputDraftID := c.Param("id")
	if inputDraftID == "" {
		log.Debug("[D] parameter `id` for draft is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for draft is missing")
		return
	}
	draftID, err := strconv.Atoi(inputDraftID)
	if err != nil {
		log.Debugf("[D] parameter `id` for draft is invalid, value = %v", inputDraftID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for draft is invalid, value = %v", inputDraftID))
		return
	}

	draft := model.Draft{}
	if err := db.Where("id = ?", draftID).First(&draft).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	es := strings.Split(draft.Endpoints, TMP_GRAPH_FILED_DELIMITER)
	cs := strings.Split(draft.Counters, TMP_GRAPH_FILED_DELIMITER)

	ret := map[string][]string{
		"endpoints": es,
		"counters":  cs,
	}

	h.JSONR(c, ret)
}

// APIGraphCreateReqData TODO:
type APIGraphCreateReqData struct {
	ScreenID    int      `json:"screenID"  binding:"required"`
	Title       string   `json:"title"     binding:"required"`
	Endpoints   []string `json:"endpoints" binding:"required"`
	Counters    []string `json:"counters"  binding:"required"`
	TimeSpan    int      `json:"timespan"`
	RelativeDay int      `json:"relativeDay"`
	Type        string   `json:"type"`
	Method      string   `json:"method"`
	Position    int      `json:"position"`
	Tags        string   `json:"tags"`
}

// CreateGraph TODO:
func CreateGraph(c *gin.Context) {
	var inputs APIGraphCreateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	es := inputs.Endpoints
	cs := inputs.Counters
	sort.Strings(es)
	sort.Strings(cs)
	esString := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	csString := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)

	graph := model.Graph{
		Title:    inputs.Title,
		Hosts:    esString,
		Counters: csString,
		ScreenID: int64(inputs.ScreenID),
		TimeSpan: inputs.TimeSpan,
		Type:     inputs.Type,
		Method:   inputs.Method,
		Position: inputs.Position,
	}
	if graph.TimeSpan == 0 {
		graph.TimeSpan = 3600
	}
	if graph.Type == "" {
		graph.Type = "h"
	}

	if err := db.Create(&graph).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}

	resp := map[string]interface{}{
		"id": graph.ID,
	}
	h.JSONR(c, resp)

}

// APIGraphUpdateReqData TODO:
type APIGraphUpdateReqData struct {
	ScreenID    int      `json:"screen_id"`
	Title       string   `json:"title"`
	Endpoints   []string `json:"endpoints"`
	Counters    []string `json:"counters"`
	TimeSpan    int      `json:"timespan"`
	RelativeDay int      `json:"relativeday"`
	Type        string   `json:"type"`
	Method      string   `json:"method"`
	Position    int      `json:"position"`
	Tags        string   `json:"tags"`
}

// UpdateGraph TODO:
func UpdateGraph(c *gin.Context) {
	inputGraphID := c.Param("id")
	if inputGraphID == "" {
		log.Debug("[D] parameter `id` for graph is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for graph is missing")
		return
	}
	graphID, err := strconv.Atoi(inputGraphID)
	if err != nil {
		log.Debugf("[D] parameter `id` for graph is invalid, value = %v", inputGraphID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for graph is invalid, value = %v", inputGraphID))
		return
	}

	var inputs APIGraphUpdateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	update := map[string]interface{}{}

	if len(inputs.Endpoints) != 0 {
		es := inputs.Endpoints
		sort.Strings(es)
		esString := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
		update["endpoints"] = esString
	}
	if len(inputs.Counters) != 0 {
		cs := inputs.Counters
		sort.Strings(cs)
		csString := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
		update["counters"] = csString
	}
	if inputs.Title != "" {
		update["title"] = inputs.Title
	}
	if inputs.ScreenID != 0 {
		update["screen_id"] = int64(inputs.ScreenID)
	}
	if inputs.TimeSpan != 0 {
		update["timespan"] = inputs.TimeSpan
	}
	if inputs.Type != "" {
		update["type"] = inputs.Type
	}
	if inputs.Method != "" {
		update["method"] = inputs.Method
	} else {
		update["method"] = " "
	}
	if inputs.Position != 0 {
		update["position"] = inputs.Position
	}
	if inputs.Tags != "" {
		update["tags"] = inputs.Tags
	}

	if err := db.Model(model.Graph{}).Where("id = ?", graphID).Updates(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	resp := map[string]interface{}{
		"id": graphID,
	}
	h.JSONR(c, resp)
}

// GetGraph TODO:
func GetGraph(c *gin.Context) {
	inputGraphID := c.Param("id")
	if inputGraphID == "" {
		log.Debug("[D] parameter `id` for graph is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for graph is missing")
		return
	}
	graphID, err := strconv.Atoi(inputGraphID)
	if err != nil {
		log.Debugf("[D] parameter `id` for graph is invalid, value = %v", inputGraphID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for graph is invalid, value = %v", inputGraphID))
		return
	}

	graph := model.Graph{}
	if err := db.Where("id = ?", graphID).First(&graph).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
	cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)
	resp := map[string]interface{}{
		"graph_id":  graph.ID,
		"title":     graph.Title,
		"endpoints": es,
		"counters":  cs,
		"screen_id": graph.ScreenID,
		"type":      graph.Type,
		"timespan":  graph.TimeSpan,
		"method":    graph.Method,
		"position":  graph.Position,
		"tags":      graph.Tags,
	}
	h.JSONR(c, resp)
}

// DeleteGraph TODO:
func DeleteGraph(c *gin.Context) {
	inputGraphID := c.Param("id")
	if inputGraphID == "" {
		log.Debug("[D] parameter `id` for graph is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for graph is missing")
		return
	}
	graphID, err := strconv.Atoi(inputGraphID)
	if err != nil {
		log.Debugf("[D] parameter `id` for graph is invalid, value = %v", inputGraphID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for graph is invalid, value = %v", inputGraphID))
		return
	}

	graph := model.Graph{}
	if err := db.Where("id = ?", graphID).Delete(&graph).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}
	resp := map[string]interface{}{
		"id": graphID,
	}
	h.JSONR(c, resp)
}

// GetGraphsByScreenID TODO:
func GetGraphsByScreenID(c *gin.Context) {
	inputScreenID := c.Param("screen_id")
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

	limit := c.DefaultQuery("limit", "500")

	graphs := []model.Graph{}
	if err := db.Where("screen_id = ?", screenID).Limit(limit).Find(&graphs).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	ret := []map[string]interface{}{}
	for _, graph := range graphs {
		es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
		cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)

		r := map[string]interface{}{
			"graph_id":  graph.ID,
			"title":     graph.Title,
			"endpoints": es,
			"counters":  cs,
			"screen_id": graph.ScreenID,
			"type":      graph.Type,
			"timespan":  graph.TimeSpan,
			"method":    graph.Method,
			"position":  graph.Position,
			"tags":      graph.Tags,
			// "relativeday": graph.RelativeDay,
		}
		ret = append(ret, r)
	}

	h.JSONR(c, ret)
}

// APIEndpointObjGetInputs TODO:
type APIEndpointObjGetInputs struct {
	Endpoints []string `json:"endpoints" form:"endpoints" binding:"required"`
	Deadline  int64    `json:"deadline"  form:"deadline"`
}

// EndpointObjGet TODO:
func EndpointObjGet(c *gin.Context) {
	inputs := APIEndpointObjGetInputs{
		Deadline: 0,
	}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	// TODO: endpoints missing with required binding ??
	if len(inputs.Endpoints) == 0 {
		h.JSONR(c, h.HTTPBadRequest, "endpoints missing")
		return
	}

	result := []model.Endpoint{}
	// TODO:
	if err := db.
		Where("endpoint in (?) and ts >= ?", inputs.Endpoints, inputs.Deadline).
		Find(&result).Error; err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}

	endpoints := []map[string]interface{}{}
	for _, r := range result {
		endpoints = append(endpoints, map[string]interface{}{"id": r.ID, "endpoint": r.Endpoint, "ts": r.Ts})
	}

	h.JSONR(c, endpoints)
}

// APIEndpointRegexpQueryInputs TODO:
type APIEndpointRegexpQueryInputs struct {
	Q     string `json:"q"     form:"q"`
	Label string `json:"tags"  form:"tags"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page"  form:"page"`
}

// IsIP TODO:
func IsIP(ip string) (b bool) {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return false
	}
	return true
}

// EndpointRegexpQuery TODO:
func EndpointRegexpQuery(c *gin.Context) {
	inputs := APIEndpointRegexpQueryInputs{
		// set default is 500
		Limit: 500,
		Page:  1,
	}
	var err error
	if err = c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if inputs.Q == "" && inputs.Label == "" {
		h.JSONR(c, h.HTTPBadRequest, "q and labels are all missing")
		return
	}

	labels := []string{}
	if inputs.Label != "" {
		labels = strings.Split(inputs.Label, ",")
	}
	qsNew := []string{}
	qs := []string{}
	enphostname := model.Host{}
	if inputs.Q != "" {
		qsNew = strings.Split(inputs.Q, " ")
		for _, each := range qsNew {
			if IsIP(each) {
				db.Select("hostname").Where("ip IN (?)", each).Find(&enphostname)
				qs = append(qs, enphostname.Hostname)
			} else {
				qs = append(qs, each)
			}
		}
	}

	var offset int
	if inputs.Page > 1 {
		offset = (inputs.Page - 1) * inputs.Limit
	}

	var endpoint []model.Endpoint
	var endpointIDs []int
	if len(labels) != 0 {
		dt := db.Select("distinct endpoint_id")
		for _, trem := range labels {
			dt = dt.Where(" counter like ? ", "%"+strings.TrimSpace(trem)+"%")
		}
		if err = dt.Model(model.Counter{}).Limit(inputs.Limit).Offset(offset).Pluck("distinct endpoint_id", &endpointIDs).Error; err != nil {
			h.JSONR(c, h.HTTPBadRequest, err)
			return
		}
	}
	if len(qs) != 0 {
		dt := db.Select("endpoint, id")
		if len(endpointIDs) != 0 {
			dt = dt.Where("id in (?)", endpointIDs)
		}

		for _, trem := range qs {
			dt = dt.Where(" endpoint regexp ? ", strings.TrimSpace(trem))
		}
		err = dt.Limit(inputs.Limit).Offset(offset).Find(&endpoint).Error
	} else if len(endpointIDs) != 0 {
		err = db.Select("endpoint, id").
			Where("id in (?)", endpointIDs).
			Find(&endpoint).Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	endpoints := []map[string]interface{}{}
	for _, e := range endpoint {
		enpip := model.Host{}
		db.Select("ip").Where("hostname in (?)", e.Endpoint).First(&enpip)
		endpoints = append(endpoints, map[string]interface{}{"id": e.ID, "endpoint": e.Endpoint, "ip": enpip.IP})
	}

	h.JSONR(c, endpoints)
}

// EndpointCounterRegexpQuery TODO:
func EndpointCounterRegexpQuery(c *gin.Context) {
	inputEndpointID := c.DefaultQuery("eid", "")
	metricQuery := c.DefaultQuery("metricQuery", ".+")
	limitTmp := c.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitTmp)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	pageTmp := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageTmp)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	var offset int
	if page > 1 {
		offset = (page - 1) * limit
	}
	eidArray := []string{}
	if inputEndpointID == "" {
		h.JSONR(c, h.HTTPBadRequest, "eid is missing")
		return
	} else {
		eids := utils.ConvertIntStringToList(inputEndpointID)
		if eids == "" {
			h.JSONR(c, h.HTTPBadRequest, "input error, please check your input info.")
			return
		}
		eidArray = strings.Split(eids, ",")

		var counters []model.Counter
		dt := db.Select("endpoint_id, counter, step, type").Where("endpoint_id IN (?)", eidArray)
		if metricQuery != "" {
			qs := strings.Split(metricQuery, " ")
			if len(qs) > 0 {
				for _, term := range qs {
					t := strings.TrimSpace(term)
					if t != "" {
						if strings.HasPrefix(term, "!") {
							dt = dt.Where("NOT counter REGEXP ?", term[1:])
						} else {
							dt = dt.Where("counter REGEXP ?", term)
						}
					}
				}
			}
		}
		if err := dt.Limit(limit).Offset(offset).Find(&counters).Error; err != nil {
			h.InternelError(c, "retrieving data", err)
			return
		}

		countersResp := []interface{}{}
		for _, c := range counters {
			countersResp = append(countersResp, map[string]interface{}{
				"endpoint_id": c.EndpointID,
				"counter":     c.Counter,
				"step":        c.Step,
				"type":        c.Type,
			})
		}
		h.JSONR(c, countersResp)
	}
	return
}

// APIQueryGraphDrawData TODO:
type APIQueryGraphDrawData struct {
	HostNames []string `json:"hostnames" binding:"required"`
	Counters  []string `json:"counters"  binding:"required"`
	ConsolFun string   `json:"consolFun" binding:"required"`
	StartTime int64    `json:"startTime" binding:"required"`
	EndTime   int64    `json:"endTime"   binding:"required"`
	Step      int      `json:"step"`
}

// QueryGraphDrawData TODO:
func QueryGraphDrawData(c *gin.Context) {
	var inputs APIQueryGraphDrawData
	var err error
	if err = c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	respData := []*cm.GraphQueryResponse{}
	for _, host := range inputs.HostNames {
		for _, counter := range inputs.Counters {
			var step int
			if inputs.Step > 0 && inputs.StartTime != -1 {
				step = inputs.Step
			} else {
				step, err = getCounterStep(host, counter)
				if err != nil {
					h.InternelError(c, "retrieving data", err)
					continue
				}
			}
			if data, err := fetchData(host, counter, inputs.ConsolFun, inputs.StartTime, inputs.EndTime, step); err != nil {
				h.InternelError(c, "retrieving data", err)
				return
			} else {
				respData = append(respData, data)
			}
		}
	}
	h.JSONR(c, respData)
}

// QueryGraphLastPoint TODO:
func QueryGraphLastPoint(c *gin.Context) {
	var inputs []cm.GraphLastParam
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	respData := []*cm.GraphLastResp{}

	for _, param := range inputs {
		resp, err := grh.Last(param)
		if err != nil {
			log.Errorf("[E] query last point from graph fail: %v", err)
		} else {
			respData = append(respData, resp)
		}
	}

	h.JSONR(c, respData)
}

// DeleteGraphEndpoint TODO:
func DeleteGraphEndpoint(c *gin.Context) {
	inputs := []string{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	type DBRows struct {
		Endpoint  string
		CounterID int
		Counter   string
		Type      string
		Step      int
	}

	rows := []DBRows{}
	if err := db.Raw(
		`select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoints as a, counters as b
		where b.endpoint_id = a.id
		AND a.endpoint in (?)`, inputs).Scan(&rows).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	var affectedCounter int64
	var affectedEndpoint int64

	if len(rows) > 0 {
		params := []*cm.GraphDeleteParam{}
		for _, row := range rows {
			param := &cm.GraphDeleteParam{
				Endpoint: row.Endpoint,
				DsType:   row.Type,
				Step:     row.Step,
			}
			fields := strings.SplitN(row.Counter, "/", 2)
			if len(fields) == 1 {
				param.Metric = fields[0]
			} else if len(fields) == 2 {
				param.Metric = fields[0]
				param.Tags = fields[1]
			} else {
				log.Errorf("[E] invalid counter %s", row.Counter)
				continue
			}
			params = append(params, param)
		}
		grh.Delete(params)
	}

	tx := db.Begin()

	if len(rows) > 0 {
		cids := make([]int, len(rows))
		for i, row := range rows {
			cids[i] = row.CounterID
		}

		if err := tx.Where("id IN (?)", cids).Delete(&model.Counter{}).Error; err != nil {
			h.InternelError(c, "deleting data", err)
			tx.Rollback()
			return
		} else {
			affectedCounter = db.RowsAffected
		}

		if err := tx.Exec(`delete from tags where endpoint_id in 
			(select id from endpoints where endpoint in (?))`, inputs).Error; err != nil {
			h.InternelError(c, "deleting data", err)
			tx.Rollback()
			return
		}
	}

	if err := tx.Where("endpoint in (?)", inputs).Delete(&model.Endpoint{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	} else {
		affectedEndpoint = db.RowsAffected
	}
	tx.Commit()
	resp := map[string]interface{}{
		"affected_endpoint": affectedEndpoint,
		"affected_counter":  affectedCounter,
	}
	h.JSONR(c, resp)
}

// APIGraphDeleteCounterInputs TODO:
type APIGraphDeleteCounterInputs struct {
	Endpoints []string `json:"endpoints" binding:"required"`
	Counters  []string `json:"counters"  binding:"required"`
}

// DeleteGraphCounter TODO:
func DeleteGraphCounter(c *gin.Context) {
	inputs := APIGraphDeleteCounterInputs{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	type DBRows struct {
		Endpoint  string
		CounterID int
		Counter   string
		Type      string
		Step      int
	}

	rows := []DBRows{}
	if err := db.Raw(`select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoints as a,
		counters as b
		where b.endpoint_id = a.id 
		AND a.endpoint in (?)
		AND b.counter in (?)`, inputs.Endpoints, inputs.Counters).Scan(&rows).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	}
	if len(rows) == 0 {
		resp := map[string]interface{}{
			"affected_counter": 0,
		}
		h.JSONR(c, resp)
		return
	}

	params := []*cm.GraphDeleteParam{}
	for _, row := range rows {
		param := &cm.GraphDeleteParam{
			Endpoint: row.Endpoint,
			DsType:   row.Type,
			Step:     row.Step,
		}
		fields := strings.SplitN(row.Counter, "/", 2)
		if len(fields) == 1 {
			param.Metric = fields[0]
		} else if len(fields) == 2 {
			param.Metric = fields[0]
			param.Tags = fields[1]
		} else {
			log.Errorf("[E] invalid counter %s", row.Counter)
			continue
		}
		params = append(params, param)
	}
	grh.Delete(params)

	tx := db.Begin()
	cids := make([]int, len(rows))
	for i, row := range rows {
		cids[i] = row.CounterID
	}

	var affectedRows int64
	if err := tx.Where("id in (?)", cids).Delete(&model.Counter{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	} else {
		affectedRows = db.RowsAffected
	}
	tx.Commit()

	resp := map[string]interface{}{
		"affected_counter": affectedRows,
	}
	h.JSONR(c, resp)
}

func fetchData(hostname string, counter string, consolFun string, startTime int64, endTime int64, step int) (resp *cm.GraphQueryResponse, err error) {
	hostnameNew := strings.Split(hostname, "_")[0]
	qparm := grh.GenQParam(hostnameNew, counter, consolFun, startTime, endTime, step)
	log.Debugf("[D] qparm: %v", qparm)
	resp, err = grh.QueryOne(qparm)
	return
}

func getCounterStep(endpoint, counter string) (step int, err error) {
	key := fmt.Sprintf("step:%s/%s", endpoint, counter)
	s, found := localStepCache.Get(key)
	if found && s != nil {
		step = s.(int)
		return
	}

	var rows []int
	if err = db.Raw(`select a.step from counters as a, endpoints as b
		 where b.endpoint = ? and a.endpoint_id = b.id and a.counter = ? limit 1`, endpoint, counter).
		Scan(&rows).Error; err != nil {
		return
	}
	if len(rows) == 0 {
		err = errors.New("empty result")
		return
	}
	step = rows[0]
	localStepCache.Set(key, step, tcache.DefaultExpiration)

	return
}
