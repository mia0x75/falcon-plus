package plugins

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *Plugin
	Quit   chan struct{}
}

func NewPluginScheduler(p *Plugin) *PluginScheduler {
	scheduler := PluginScheduler{Plugin: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Cycle) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

func (this *PluginScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				PluginRun(this.Plugin)
			case <-this.Quit:
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this *PluginScheduler) Stop() {
	close(this.Quit)
}

func PluginRun(plugin *Plugin) {
	timeout := plugin.Cycle*1000 - 500
	fpath := filepath.Join(g.Config().Plugin.Dir, plugin.FilePath)

	if !file.IsExist(fpath) {
		log.Warnf("[W] no such plugin: %s", fpath)
		return
	}

	log.Debugf("[D] %s running...", fpath)

	cmd := exec.Command(fpath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	if err != nil {
		log.Errorf("[E] plugin start fail, error: %s\n", err)
		return
	}
	log.Debugf("[D] plugin started: %s", fpath)

	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)

	errStr := stderr.String()
	if errStr != "" {
		logFile := filepath.Join(g.Config().Plugin.LogDir, plugin.FilePath+".stderr.log")
		if _, err = file.WriteString(logFile, errStr); err != nil {
			log.Errorf("[E] write log to %s fail, error: %s\n", logFile, err)
		}
	}

	if isTimeout {
		// has be killed
		if err == nil {
			log.Infof("[I] timeout and kill process %s successfully", fpath)
		}

		if err != nil {
			log.Errorf("[E] kill process %s occur error: %v", fpath, err)
		}

		return
	}

	if err != nil {
		log.Errorf("[E] exec plugin %s fail. error: %v", fpath, err)
		return
	}

	// exec successfully
	data := stdout.Bytes()
	if len(data) == 0 {
		log.Debugf("[D] stdout of %s is blank", fpath)
		return
	}

	var metrics []*cmodel.MetricValue
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Errorf("[E] json.Unmarshal stdout of %s fail. error: %s stdout: \n%s\n", fpath, err, stdout.String())
		return
	}

	hostname, err := g.Hostname()
	if err != nil {
		log.Errorf("[E] get hostname fail: %v", err)
		return
	}
	// 如果插件中没有配置Endpoint 则使用当前agent 的Endpint
	// 适用于需要统一agent 和插件Endpoint的情况，只需将插件
	// Endpoint 置空即可
	for _, metric := range metrics {
		if metric.Endpoint == "" {
			metric.Endpoint = hostname
		}
	}

	g.SendToTransfer(metrics)
}
