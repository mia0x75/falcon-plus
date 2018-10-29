package cache

// 每个agent心跳上来的时候立马更新一下数据库是没必要的
// 缓存起来，每隔一个小时写一次DB
// 提供http接口查询机器信息，排查重名机器的时候比较有用

import (
	"sync"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

type SafeAgents struct {
	sync.RWMutex
	M map[string]*cmodel.AgentUpdateInfo
}

var (
	Agents = NewSafeAgents()
)

func NewSafeAgents() *SafeAgents {
	return &SafeAgents{M: make(map[string]*cmodel.AgentUpdateInfo)}
}

func (this *SafeAgents) Init() {
	m, err := db.QueryAgentsInfo()
	if err != nil {
		return
	}
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeAgents) Put(req *cmodel.AgentReportRequest) {
	val := &cmodel.AgentUpdateInfo{
		LastUpdate:    time.Now().Unix(),
		ReportRequest: req,
	}

	if agentInfo, exists := this.Get(req.Hostname); !exists ||
		agentInfo.ReportRequest.AgentVersion != req.AgentVersion ||
		agentInfo.ReportRequest.IP != req.IP ||
		agentInfo.ReportRequest.PluginVersion != req.PluginVersion {

		db.UpdateAgent(val)
	}
	// 更新hbs时间
	this.Lock()
	defer this.Unlock()
	this.M[req.Hostname] = val
}

func (this *SafeAgents) Get(hostname string) (*cmodel.AgentUpdateInfo, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[hostname]
	return val, exists
}

func (this *SafeAgents) Delete(hostname string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, hostname)
}

func (this *SafeAgents) Keys() []string {
	this.RLock()
	defer this.RUnlock()
	count := len(this.M)
	keys := make([]string, count)
	i := 0
	for hostname := range this.M {
		keys[i] = hostname
		i++
	}
	return keys
}

func DeleteStaleAgents() {
	go func() {
		d := time.Hour * time.Duration(24)
		for range time.Tick(d) {
			deleteStaleAgents()
		}
	}()
}

func deleteStaleAgents() {
	// 一天都没有心跳的Agent，从内存中干掉
	before := time.Now().Unix() - 3600*24
	keys := Agents.Keys()
	count := len(keys)
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		curr, _ := Agents.Get(keys[i])
		if curr.LastUpdate < before {
			Agents.Delete(curr.ReportRequest.Hostname)
		}
	}
}
