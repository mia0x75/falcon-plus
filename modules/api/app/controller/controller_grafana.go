package controller

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
)

// APIGrafanaMainQueryInputs TODO:
type APIGrafanaMainQueryInputs struct {
	Limit int    `json:"limit" form:"limit"`
	Query string `json:"query" form:"query"`
}

// APIGrafanaMainQueryOutputs TODO:
type APIGrafanaMainQueryOutputs struct {
	Expandable bool   `json:"expandable"`
	Text       string `json:"text"`
}

// for return a host list for api test
func repsonseDefault(limit int) (result []APIGrafanaMainQueryOutputs) {
	result = []APIGrafanaMainQueryOutputs{}
	// for get right table name
	enps := []model.Endpoint{}
	enpip := model.Host{}
	db.Limit(limit).Find(&enps)
	for _, h := range enps {
		db.Select("hostname,ip").Where("hostname in (?)", h.Endpoint).Find(&enpip)
		result = append(result, APIGrafanaMainQueryOutputs{
			Expandable: true,
			Text:       h.Endpoint + "_" + enpip.IP,
		})
	}
	return
}

// for find host list & grafana template searching, regexp support
func responseHostsRegexp(limit int, regexpKey string) (result []APIGrafanaMainQueryOutputs) {
	result = []APIGrafanaMainQueryOutputs{}
	// for get right table name
	enps := []model.Endpoint{}
	enpip := model.Host{}
	// TODO: Error ??
	db.Where("endpoint regexp ?", regexpKey).Limit(limit).Find(&enps)
	for _, h := range enps {
		// TODO: Error ??
		db.Select("hostname,ip").Where("hostname = ?", h.Endpoint).First(&enpip)
		result = append(result, APIGrafanaMainQueryOutputs{
			Expandable: true,
			Text:       h.Endpoint + "_" + enpip.IP,
		})
	}
	return
}

// for find hostgroup list & grafana template searching
func responseHostGroupRegexp(limit int, pattern string) (result []APIGrafanaMainQueryOutputs) {
	result = []APIGrafanaMainQueryOutputs{}
	groups := []model.Group{}
	db.Where("name regexp ?", pattern).Limit(limit).Find(&groups)
	for _, g := range groups {
		result = append(result, APIGrafanaMainQueryOutputs{
			Expandable: true,
			Text:       g.Name,
		})
	}
	return
}

// for find host list by HostGroup & grafana template searching
func responseHostsByHostGroup(limit int, pattern string) (result []APIGrafanaMainQueryOutputs) {
	result = []APIGrafanaMainQueryOutputs{}
	groups := []model.Group{}
	db.Where("name REGEXP ?", pattern).Limit(limit).Find(&groups)
	if 0 == len(groups) {
		return
	}
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		for _, g := range groups {
			if elem.AncestorID == g.ID && elem.Type == 2 {
				return true
			}
		}
		return false
	})
	if edges == nil {
		return
	}
	hosts := cache.HostsMap.Filter(func(elem *model.Host) bool {
		for _, r := range edges {
			if elem.ID == r.DescendantID {
				return true
			}
		}
		return false
	})
	if hosts == nil {
		return
	}
	c := 0
	for _, host := range hosts {
		c = c + 1
		result = append(result, APIGrafanaMainQueryOutputs{
			Expandable: true,
			Text:       host.Hostname,
		})
		if c >= limit {
			return
		}
	}
	return
}

// for resolve mixed query with endpoint & counter of query string
func cutEndpointCounterHelp(input string) (hosts []string, counter string) {
	r, _ := regexp.Compile("^{?([^#}]+)}?#(.+)")
	matchedList := r.FindAllStringSubmatch(input, 1)
	if len(matchedList) != 0 {
		if len(matchedList[0]) > 1 {
			// get hosts
			hostsTmp := matchedList[0][1]
			counterTmp := matchedList[0][2]
			hosts = strings.Split(hostsTmp, ",")
			counter = strings.Replace(counterTmp, "#", "\\.", -1)
		}
	} else {
		log.Errorf("[E] grafana query inputs error: %v", input)
	}
	return
}

func expandableChecking(counter string, counterSearchKeyWord string) (expsub string, needexp bool) {
	re := regexp.MustCompile("(\\.\\+|\\.\\*)\\s*$")
	counterSearchKeyWord = re.ReplaceAllString(counterSearchKeyWord, "")
	counterSearchKeyWord = strings.Replace(counterSearchKeyWord, "\\.", ".", -1)
	expCheck := strings.Replace(counter, counterSearchKeyWord, "", -1)
	if expCheck == "" {
		needexp = false
		expsub = expCheck
	} else {
		needexp = true
		re = regexp.MustCompile("^\\.")
		expsubArr := strings.Split(re.ReplaceAllString(expCheck, ""), ".")

		switch len(expsubArr) {
		case 0:
			expsub = ""
		case 1:
			expsub = expsubArr[0]
			needexp = false
		default:
			expsub = expsubArr[0]
			needexp = true
			// if counter like switch.if.In/ifIndex=177,ifName=Eth-Trunk3.1000
			// not Split, return it
			if strings.Contains(expsub, "ifName") {
				expsub = expCheck
				needexp = false
			}
		}
	}
	return
}

/* add additional items (ex. $ & %)
   $ means metric is stop on here. no need expand any more.
   % means a wirecard string.
   also clean defecate metrics
*/
func addAddItionalItems(items []APIGrafanaMainQueryOutputs, regexpKey string) (result []APIGrafanaMainQueryOutputs) {
	flag := false
	mapset := hashmap.New()
	for _, i := range items {
		if !i.Expandable {
			flag = true
		}
		if val, exist := mapset.Get(i.Text); exist {
			if val != i.Expandable && i.Expandable {
				mapset.Put(i.Text, i.Expandable)
			}
		} else {
			mapset.Put(i.Text, i.Expandable)
		}
	}
	result = make([]APIGrafanaMainQueryOutputs, mapset.Size())
	for idx, ctmp := range mapset.Keys() {
		c := ctmp.(string)
		val, _ := mapset.Get(c)
		result[idx] = APIGrafanaMainQueryOutputs{
			Text:       c,
			Expandable: val.(bool),
		}
	}
	if flag {
		result = append(result, APIGrafanaMainQueryOutputs{
			Text:       "$",
			Expandable: false,
		})
	}
	if len(strings.Split(regexpKey, "\\.")) > 0 {
		result = append(result, APIGrafanaMainQueryOutputs{
			Text:       "%",
			Expandable: false,
		})
	}
	return
}

func findEndpointIDByEndpointList(hosts []string) []int64 {
	// for get right table name
	enps := []model.Endpoint{}
	var hostsNew []string
	for _, h := range hosts {
		hostsNew = append(hostsNew, strings.Split(h, "_")[0])
	}
	db.Where("endpoint in (?)", hosts).Find(&enps)
	hostIds := make([]int64, len(enps))
	for indx, h := range enps {
		hostIds[indx] = int64(h.ID)
	}
	return hostIds
}

// for reture counter list of endpoints
func responseCounterRegexp(regexpKey string) (result []APIGrafanaMainQueryOutputs) {
	result = []APIGrafanaMainQueryOutputs{}
	hosts, counter := cutEndpointCounterHelp(regexpKey)
	if len(hosts) == 0 || counter == "" {
		return
	}
	hostIds := findEndpointIDByEndpointList(hosts)
	// if not any endpoint matched
	if len(hostIds) == 0 {
		return
	}
	// for get right table name
	counters := []model.Counter{}
	db.Where("endpoint_id IN (?)", hostIds).Where("counter REGEXP ?", counter).Find(&counters)
	// if not any counter matched
	if len(counters) == 0 {
		return
	}
	for _, c := range counters {
		expsub, needexp := expandableChecking(c.Counter, counter)
		result = append(result, APIGrafanaMainQueryOutputs{
			Text:       expsub,
			Expandable: needexp,
		})
	}
	result = addAddItionalItems(result, regexpKey)
	return
}

// GrafanaMainQuery TODO:
func GrafanaMainQuery(c *gin.Context) {
	inputs := APIGrafanaMainQueryInputs{}
	inputs.Limit = 1000
	inputs.Query = "!N!"
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	log.Debugf("[D] got query string: %s", inputs.Query)
	output := []APIGrafanaMainQueryOutputs{}
	if inputs.Query == "!N!" {
		output = repsonseDefault(inputs.Limit)
	} else if strings.HasPrefix(inputs.Query, "!HG!") {
		output = responseHostGroupRegexp(inputs.Limit, inputs.Query[len("!HG!"):])
	} else if strings.HasPrefix(inputs.Query, "!HGH!") {
		output = responseHostsByHostGroup(inputs.Limit, inputs.Query[len("!HGH!"):])
	} else if !strings.Contains(inputs.Query, "#") {
		output = responseHostsRegexp(inputs.Limit, inputs.Query)
	} else if strings.Contains(inputs.Query, "#") && !strings.Contains(inputs.Query, "#select metric") {
		output = responseCounterRegexp(inputs.Query)
	}
	c.JSON(200, output)
	return
}

// APIGrafanaRenderInput TODO:
type APIGrafanaRenderInput struct {
	Target        []string `json:"target"        form:"target"        binding:"required"`
	From          int64    `json:"from"          form:"from"          binding:"required"`
	Until         int64    `json:"until"         form:"until"         binding:"required"`
	Format        string   `json:"format"        form:"format"`
	MaxDataPoints int64    `json:"maxDataPoints" form:"maxDataPoints"`
	Step          int      `json:"step"          form:"step"`
	ConsolFun     string   `json:"consolFun"     form:"consolFun"`
}

// GrafanaRender TODO:
func GrafanaRender(c *gin.Context) {
	inputs := APIGrafanaRenderInput{}
	// set default step is 60
	inputs.Step = 60
	inputs.ConsolFun = "AVERAGE"
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	respList := []*cm.GraphQueryResponse{}
	for _, target := range inputs.Target {
		hosts, counter := cutEndpointCounterHelp(target)
		// clean characters
		log.Debugf("[D] %s", counter)
		re := regexp.MustCompile("\\\\.\\$\\s*$")
		flag := re.MatchString(counter)
		counter = re.ReplaceAllString(counter, "")
		counter = strings.Replace(counter, "\\.%", ".+", -1)
		counters := []model.Counter{}
		hostIds := findEndpointIDByEndpointList(hosts)
		if flag {
			db.Select("distinct counter").Where("endpoint_id IN (?) AND counter = ?", hostIds, counter).Find(&counters)
		} else {
			db.Select("distinct counter").Where("endpoint_id IN (?) AND counter regexp ?", hostIds, counter).Find(&counters)
		}
		if len(counters) == 0 {
			// 没有匹配到的继续执行，避免当grafana graph有多个查询时，其他正常的查询也无法渲染视图
			continue
		}
		counterArr := make([]string, len(counters))
		for indx, c := range counters {
			counterArr[indx] = c.Counter
		}
		for _, host := range hosts {
			for _, c := range counterArr {
				resp, err := fetchData(host, c, inputs.ConsolFun, inputs.From, inputs.Until, inputs.Step)
				if err != nil {
					log.Debugf("[D] query graph got error with: %v", inputs)
				} else {
					respList = append(respList, resp)
				}
			}
		}
	}
	c.JSON(200, respList)
	return
}
