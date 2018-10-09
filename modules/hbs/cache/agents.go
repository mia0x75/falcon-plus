package cache

// 每个agent心跳上来的时候立马更新一下数据库是没必要的
// 缓存起来，每隔一个小时写一次DB
// 提供http接口查询机器信息，排查重名机器的时候比较有用

import (
	"sync"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

type SafeAgents struct {
	sync.RWMutex
	M map[string]*model.AgentUpdateInfo
}

var (
	Agents = NewSafeAgents()
)

func NewSafeAgents() *SafeAgents {
	return &SafeAgents{M: make(map[string]*model.AgentUpdateInfo)}
}

func (this *SafeAgents) Put(req *model.AgentReportRequest) {
	val := &model.AgentUpdateInfo{
		LastUpdate:    time.Now().Unix(),
		ReportRequest: req,
	}

	if agentInfo, exists := this.Get(req.Hostname); !exists ||
		agentInfo.ReportRequest.AgentVersion != req.AgentVersion ||
		agentInfo.ReportRequest.IP != req.IP ||
		agentInfo.ReportRequest.PluginVersion != req.PluginVersion {

		db.UpdateAgent(val)
	}

	this.Lock()
	this.M[req.Hostname] = val
	this.Unlock()
}

func (this *SafeAgents) Get(hostname string) (*model.AgentUpdateInfo, bool) {
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
	duration := time.Hour * time.Duration(24)
	go func() {
		for {
			time.Sleep(duration)
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

// CheckAgentHbs 检查agent hbs 间隔时间，超过AgentMaxIdle 则agent.alive = -1
/*
func CheckAgentHbs() {
	duration := time.Second * time.Duration(g.Config().AgentMaxIdle)
	for {
		time.Sleep(duration)
		checkAgentHbs()
	}
}

func checkAgentHbs() {
	before := time.Now().Unix() - g.Config().AgentMaxIdle
	keys := Agents.Keys()
	count := len(keys)
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		curr, _ := Agents.Get(keys[i])
		if curr.LastUpdate < before {
			key := cutils.PK(curr.ReportRequest.Hostname, "agent.alive", nil)
			genMock(genTs(time.Now().Unix(), g.Config().AgentStep), key, curr.ReportRequest.Hostname)
		}
	}

	sender.SendMockOnceAsync()
}

func genMock(ts int64, key, hostname string) {
	sender.AddMock(key, hostname, "agent.alive", "", ts, "GAUGE", g.Config().AgentStep, -1)
}

//mock的数据,要前移1+个周期、防止覆盖正常值
func genTs(nowTs int64, step int64) int64 {
	if step < 1 {
		step = 60
	}

	return nowTs - nowTs%step - 2*step
}
*/
