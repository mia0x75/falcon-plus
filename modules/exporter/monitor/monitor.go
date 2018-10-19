package monitor

import (
	"bytes"
	"fmt"
	"net/url"
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
)

type State struct {
	sync.RWMutex
	Queue  []string
	Name   string
	Host   string
	Errors int32
}

// AlarmItem Struct
type Alarm struct {
	Host      string
	Name      string
	Type      string
	Count     int32
	Timestamp int64
}

func (a *Alarm) String() string {
	switch a.Type {
	case "err":
		return fmt.Sprintf("PROBLEM\nP1\n%s\n%d\n%s", a.Name, a.Count, a.Host)
	case "ok":
		return fmt.Sprintf("OK\nP3\n%s\n%d\n%s", a.Name, a.Count, a.Host)
	default:
		return ""
	}
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

	HealthState = make(map[string]*State, 5)

	go startMonitor()
	go startJudge()

	log.Println("monitor.Start, ok")
}

func startMonitor() {
	cronHealth.AddFuncCC("*/30 * * * * *", func() { monitor() }, 1)
	cronHealth.Start()
}

// monitor
func monitor() {
	clusters := g.Config().Monitor.Hosts.Modules
	for _, node := range g.Config().Monitor.Hosts.Agents {
		clusters[fmt.Sprintf("agent-%s", node)] = node
	}
	for name, host := range clusters {
		go func(name, host string) {
			var state *State
			var found bool
			lock.Lock()
			if state, found = HealthState[host]; !found {
				state = &State{Errors: 0, Name: name, Host: host}
				state.Queue = []string{"ok", "ok"}
				HealthState[host] = state
			}
			lock.Unlock()
			url := fmt.Sprintf(g.Config().Monitor.Pattern, host)
			client := cutils.NewHttp(url)
			client.SetUserAgent("monitor.get")
			body, err := client.Get()
			if !(err == nil && len(body) >= 2 && string(body[:2]) == "ok") {
				state.Queue = append(state.Queue[1:], "err")
				state.RLock()
				state.Errors++
				if state.Errors >= 3 {
					// raise problem
					alarm := &Alarm{
						Host:      host,
						Name:      name,
						Type:      "err",
						Count:     state.Errors,
						Timestamp: time.Now().Unix(),
					}
					alarmCache.Put(host, alarm)
				}
				state.RUnlock()
				log.Errorf("%s, get health error.", host)
			} else {
				state.Queue = append(state.Queue[1:], "ok")
				state.RLock()
				if state.Errors >= 3 {
					// problem restore
					alarm := &Alarm{
						Host:      host,
						Name:      name,
						Type:      "ok",
						Count:     state.Errors,
						Timestamp: time.Now().Unix(),
					}
					alarmCache.Put(host, alarm)
				}
				state.Errors = 0
				state.RUnlock()
				log.Infof("%s, get health ok.", host)
			}
		}(name, host)
	}
}

// judge
func startJudge() {
	if !g.Config().Monitor.Alarm.Enabled {
		return
	}

	d := time.Duration(10) * time.Second
	for range time.Tick(d) {
		var content bytes.Buffer

		keys := alarmCache.Keys()
		if len(keys) == 0 {
			continue
		}
		for _, key := range keys {
			item, found := alarmCache.GetAndRemove(key)
			if !found {
				continue
			}
			content.WriteString(item.(*Alarm).String() + "\n")
		}

		if content.Len() < 6 {
			return
		}
		params := url.Values{}
		alarmUrl, _ := url.Parse(g.Config().Monitor.Alarm.Url)
		params.Add("tos", "tos")
		params.Add("subject", "subject")
		params.Add("content", content.String())
		params.Add("user", "exporter")
		alarmUrl.RawQuery = params.Encode()
		client := cutils.NewHttp(alarmUrl.String())
		client.SetUserAgent("monitor.alert")
		_, err := client.Post(nil)
		if err != nil {
			log.Infof("alarm send request for health check error, %s\n", err.Error())
		} else {
			log.Info("alarm send request for health check success\n")
			// statistics
		}
	}
}
