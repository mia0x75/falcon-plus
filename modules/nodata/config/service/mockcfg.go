package service

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
)

type MockCfg struct {
	Id      int
	Name    string
	Obj     string
	ObjType string
	Metric  string
	Tags    map[string]string
	Type    string
	Step    int64
	Mock    float64
}

// 当 grp展开结果 与 host结果 存在冲突时, 优先选择 host结果
func GetMockCfgFromDB() map[string]*cm.NodataConfig {
	ret := make(map[string]*cm.NodataConfig)

	q := fmt.Sprintf("SELECT id,name,obj,obj_type,metric,tags,dstype,step,mock FROM mockcfg")
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] db.query error, mockcfg: %v", err)
		return ret
	}

	defer rows.Close()
	for rows.Next() {
		t := MockCfg{}
		tags := ""
		err := rows.Scan(&t.Id, &t.Name, &t.Obj, &t.ObjType, &t.Metric, &tags, &t.Type, &t.Step, &t.Mock)
		if err != nil {
			log.Errorf("[E] db.scan error, mockcfg: %v", err)
			continue
		}
		t.Tags = cu.DictedTagstring(tags)

		err = checkMockCfg(&t)
		if err != nil {
			log.Errorf("[E] check mockcfg, error: %v", err)
			continue
		}

		endpoints := getEndpoint(t.ObjType, t.Obj)
		if len(endpoints) < 1 {
			continue
		}

		for _, ep := range endpoints {
			uuid := cu.PK(ep, t.Metric, t.Tags)
			ncfg := cm.NewNodataConfig(t.Id, t.Name, t.ObjType, ep, t.Metric, t.Tags, t.Type, t.Step, t.Mock)

			val, found := ret[uuid]
			if !found { // so cute, it's the first one
				ret[uuid] = ncfg
				continue
			}

			if isSpuerNodataCfg(val, ncfg) {
				// val is spuer than ncfg, so drop ncfg
				log.Warnf("[W] nodata.mockcfg conflict, %s, used %s, drop %s", uuid, val.Name, ncfg.Name)
			} else {
				ret[uuid] = ncfg // overwrite the old one
				log.Warnf("[W] nodata.mockcfg conflict, %s, used %s, drop %s", uuid, ncfg.Name, val.Name)
			}
		}
	}

	return ret
}

func getEndpoint(objType string, obj string) []string {
	switch objType {
	case "host":
		return getEndpointFromHosts(obj)
	case "group":
		return getEndpointFromGroups(obj)
	case "other":
		return getEndpointFromOther(obj)
	default:
		return make([]string, 0)
	}
}

func getEndpointFromHosts(hosts string) []string {
	ret := make([]string, 0)

	hlist := strings.Split(hosts, "\n")
	for _, host := range hlist {
		nh := strings.TrimSpace(host)
		if nh != "" {
			ret = append(ret, nh)
		}
	}

	return ret
}

func getEndpointFromGroups(grps string) []string {
	grplist := strings.Split(grps, "\n")

	// get host map, avoid duplicating
	hosts := make(map[string]string)
	for _, grp := range grplist {
		ngrp := strings.TrimSpace(grp)
		if len(ngrp) < 1 {
			continue
		}

		hostmap := GetHostsFromGroup(grp)
		for hostname := range hostmap {
			if hostname != "" {
				hosts[hostname] = hostname
			}
		}
	}

	// get host slice
	ret := make([]string, 0)
	for key := range hosts {
		ret = append(ret, key)
	}

	return ret
}

func getEndpointFromOther(other string) []string {
	return getEndpointFromHosts(other)
}

func checkMockCfg(mc *MockCfg) error {
	// if len(mc.Endpoint) < 1 {
	// 	return fmt.Errorf("bad mockcfg, endpoint empty")
	// }

	if len(mc.Metric) < 1 {
		return fmt.Errorf("bad mockcfg, metric empty")
	}

	if mc.Type != "GAUGE" { // 只支持GAUGE类型
		return fmt.Errorf("bad mockcfg, type illegal, type=%s", mc.Type)
	}

	if mc.Step < 1 {
		return fmt.Errorf("bad mockcfg, step illegal, step=%d", mc.Step)
	}

	return nil
}

func isSpuerNodataCfg(A *cm.NodataConfig, B *cm.NodataConfig) bool {
	if A.ObjType == "group" && B.ObjType == "host" {
		return false
	}
	return true
}
