package cache

// 每个agent心跳上来的时候立马更新一下数据库是没必要的
// 缓存起来，每隔一个小时写一次DB
// 提供http接口查询机器信息，排查重名机器的时候比较有用

import (
	"sync"
	"time"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

type SafeAgents struct {
	sync.RWMutex
	M map[string]*cm.AgentUpdateInfo
}

var (
	Agents = NewSafeAgents()
)

func NewSafeAgents() *SafeAgents {
	return &SafeAgents{M: make(map[string]*cm.AgentUpdateInfo)}
}

func (m *SafeAgents) Init() {
	agents, err := db.QueryAgentsInfo()
	if err != nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	m.M = agents
}

func (m *SafeAgents) Put(req *cm.AgentReportRequest) {
	val := &cm.AgentUpdateInfo{
		LastUpdate:    time.Now().Unix(),
		ReportRequest: req,
	}

	if agentInfo, exists := m.Get(req.Hostname); !exists ||
		agentInfo.ReportRequest.AgentVersion != req.AgentVersion ||
		agentInfo.ReportRequest.IP != req.IP ||
		agentInfo.ReportRequest.PluginVersion != req.PluginVersion {

		db.UpdateAgent(val)
	}
	// 更新hbs时间
	m.Lock()
	defer m.Unlock()
	m.M[req.Hostname] = val
}

func (m *SafeAgents) Get(hostname string) (*cm.AgentUpdateInfo, bool) {
	m.RLock()
	defer m.RUnlock()
	val, exists := m.M[hostname]
	return val, exists
}

func (m *SafeAgents) Delete(hostname string) {
	m.Lock()
	defer m.Unlock()
	delete(m.M, hostname)
}

func (m *SafeAgents) Keys() []string {
	m.RLock()
	defer m.RUnlock()
	count := len(m.M)
	keys := make([]string, count)
	i := 0
	for hostname := range m.M {
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
