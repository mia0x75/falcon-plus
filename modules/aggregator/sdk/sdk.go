package sdk

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/toolkits/net/httplib"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/sdk/requests"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/portal"
)

// HostnamesByID TODO:
func HostnamesByID(groupID int64) ([]string, error) {
	uri := fmt.Sprintf("%s/api/v1/hostgroup/%d", g.Config().API.API, groupID)
	req, err := requests.CurlPlus(uri, "GET", "aggregator", g.Config().API.Token,
		map[string]string{}, map[string]string{})

	if err != nil {
		return []string{}, err
	}

	type RESP struct {
		HostGroup f.HostGroup `json:"hostgroup"`
		Hosts     []f.Host    `json:"hosts"`
	}

	resp := &RESP{}
	err = req.ToJson(&resp)
	if err != nil {
		return []string{}, err
	}

	hosts := []string{}
	nowTs := time.Now().Unix()
	for _, host := range resp.Hosts {
		if host.MaintainBegin <= nowTs && nowTs <= host.MaintainEnd {
			continue
		}
		hosts = append(hosts, host.Hostname)
	}
	return hosts, nil
}

// QueryLastPoints TODO:
func QueryLastPoints(endpoints, counters []string) (resp []*cmodel.GraphLastResp, err error) {
	cfg := g.Config()
	uri := fmt.Sprintf("%s/api/v1/graph/lastpoint", cfg.API.API)

	var req *httplib.BeegoHttpRequest
	headers := map[string]string{"Content-type": "application/json"}
	req, err = requests.CurlPlus(uri, "POST", "aggregator", cfg.API.Token,
		headers, map[string]string{})

	if err != nil {
		return
	}

	req.SetTimeout(time.Duration(cfg.API.ConnectTimeout)*time.Millisecond,
		time.Duration(cfg.API.RequestTimeout)*time.Millisecond)

	body := []*cmodel.GraphLastParam{}
	for _, e := range endpoints {
		for _, c := range counters {
			body = append(body, &cmodel.GraphLastParam{e, c})
		}
	}

	b, err := json.Marshal(body)
	if err != nil {
		return
	}

	req.Body(b)

	err = req.ToJson(&resp)
	if err != nil {
		return
	}

	return resp, nil
}
