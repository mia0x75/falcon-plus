package monitor

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
	"github.com/toolkits/container/nmap"
	"github.com/toolkits/cron"
)

var (
	lock       = new(sync.Mutex)
	cronHealth = cron.New()
	alarmCache = nmap.NewSafeMap()
	clusters   map[string]*HostInfo
)

type State struct {
	sync.Mutex
	Name     string
	Endpoint string
	Errors   int32
}

// AlarmDto Struct
type AlarmDto struct {
	Status     string
	Priority   int
	Endpoint   string
	Metric     string
	Tags       string
	Func       string
	LeftValue  string
	Operator   string
	RightValue string
	Note       string
	Max        int
	Current    int
	Timestamp  string
	Link       string
	Occur      int
	Subscriber []string
	Uic        string
}

func (this *AlarmDto) String() string {
	return fmt.Sprintf(
		"<Content:%s, Priority:P%d, Status:%s, Value:%s, Operator:%s Threshold:%s, Occur:%d, Uic:%s, Tos:%s>",
		this.Note,
		this.Priority,
		this.Status,
		this.LeftValue,
		this.Operator,
		this.RightValue,
		this.Occur,
		this.Uic,
		this.Subscriber,
	)
}

type HostInfo struct {
	Endpoint string
	Tag      string
}

var HealthState map[string]*State

func Start() {
	if !g.Config().Monitor.Enabled {
		log.Println("monitor.Start warning, not enable")
		return
	}
	// init url
	if g.Config().Monitor.Alarm.Url == "" {
		return
	}
	if g.Config().Monitor.Pattern == "" {
		return
	}

	initClusters()
	HealthState = make(map[string]*State, 5)

	go startMonitor()
	go startJudge()

	log.Println("monitor.Start, ok")
}

func initClusters() {
	clusters = make(map[string]*HostInfo)
	for k, v := range g.Config().Monitor.Hosts.Modules {
		clusters[k] = &HostInfo{
			Endpoint: v,
			Tag:      k,
		}
	}
	for _, node := range g.Config().Monitor.Hosts.Agents {
		clusters[fmt.Sprintf("agent-%s", node)] = &HostInfo{
			Endpoint: node,
			Tag:      "agent",
		}
	}
}

func createAlarm(endpoint, status, tag string, tos []string) *AlarmDto {
	return &AlarmDto{
		Status:     status,
		Priority:   0,
		Endpoint:   endpoint,
		Metric:     "health.ok",
		Tags:       tag,
		Func:       "all(#3)",
		LeftValue:  "0",
		Operator:   "==",
		RightValue: "0",
		Note:       "组件/health存活检查",
		Max:        3,
		Current:    1,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Link:       "",
		Subscriber: []string{"tos"},
	}
}

func startMonitor() {
	cronHealth.AddFuncCC("*/30 * * * * *", func() { monitor() }, 1)
	cronHealth.Start()
}

// monitor
func monitor() {
	if clusters == nil {
		return
	}
	for _, h := range clusters {
		go func(tag, endpoint string) {
			lock.Lock()
			defer lock.Unlock()
			var state *State
			if s, ok := HealthState[endpoint]; !ok {
				s = &State{Errors: 0, Name: tag, Endpoint: endpoint}
				HealthState[endpoint] = s
				state = s
			} else {
				state = s
			}
			host := strings.Split(endpoint, ":")[0]
			url := fmt.Sprintf(g.Config().Monitor.Pattern, endpoint)
			body, err := cutils.Get(url)
			state.Lock()
			defer state.Unlock()
			if !(err == nil && len(body) >= 2 && string(body[:2]) == "ok") {
				state.Errors++
				if state.Errors >= 3 {
					// raise problem
					alarm := createAlarm(host, "PROBLEM", fmt.Sprintf("module=%s", tag), nil)
					alarm.Occur = int(state.Errors)
					alarmCache.Put(endpoint, alarm)
				}
				log.Errorf("%s, get health error.", endpoint)
			} else {
				if state.Errors >= 3 {
					// problem restore
					alarm := createAlarm(host, "OK", fmt.Sprintf("module=%s", tag), nil)
					alarm.Occur = 0
					alarmCache.Put(endpoint, alarm)
				}
				state.Errors = 0
				log.Infof("%s, get health ok.", endpoint)
			}
		}(h.Tag, h.Endpoint)
	}
}

// judge
func startJudge() {
	if !g.Config().Monitor.Alarm.Enabled {
		return
	}

	d := time.Duration(10) * time.Second
	for range time.Tick(d) {
		keys := alarmCache.Keys()
		if len(keys) == 0 {
			continue
		}
		for _, key := range keys {
			item, found := alarmCache.GetAndRemove(key)
			if !found {
				continue
			}
			if data, err := json.Marshal(item.(*AlarmDto)); err != nil {
				log.Infof("json marshal error:%v", err)
			} else {
				_, err := cutils.Post(g.Config().Monitor.Alarm.Url, data)
				if err != nil {
					log.Infof("alarm send request for health check error:%v", err)
				} else {
					log.Info("alarm send request for health check success\n")
					// statistics
				}
			}
		}
	}
}
