package plugins

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// PluginScheduler TODO:
type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *Plugin
	Quit   chan struct{}
}

// NewPluginScheduler TODO:
func NewPluginScheduler(p *Plugin) *PluginScheduler {
	scheduler := PluginScheduler{Plugin: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Cycle) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

// Schedule TODO:
func (ps *PluginScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-ps.Ticker.C:
				PluginRun(ps.Plugin)
			case <-ps.Quit:
				ps.Ticker.Stop()
				return
			}
		}
	}()
}

// Stop TODO:
func (ps *PluginScheduler) Stop() {
	close(ps.Quit)
}

// PluginArgsParse using ',' as the seprator of args and '\,' to espace
func PluginArgsParse(rawArgs string) []string {
	ss := strings.Split(rawArgs, "\\,")

	out := [][]string{}
	for _, s := range ss {
		cleanArgs := []string{}
		for _, arg := range strings.Split(s, ",") {
			arg = strings.Trim(arg, " ")
			arg = strings.Trim(arg, "\"")
			arg = strings.Trim(arg, "'")
			cleanArgs = append(cleanArgs, arg)
		}
		out = append(out, cleanArgs)
	}

	ret := []string{}
	tail := ""

	for _, x := range out {
		for j, y := range x {
			if j == 0 {
				if tail != "" {
					ret = append(ret, tail+","+y)
					tail = ""
				} else {
					ret = append(ret, y)
				}
			} else if j == len(x)-1 {
				tail = y
			} else {
				ret = append(ret, y)
			}
		}
	}

	if tail != "" {
		ret = append(ret, tail)
	}

	return ret
}

// PluginRun TODO:
func PluginRun(plugin *Plugin) {
	timeout := plugin.Cycle*1000 - 500
	fpath := filepath.Join(g.Config().Plugin.Dir, plugin.FilePath)
	args := plugin.Args

	if !file.IsExist(fpath) {
		log.Warnf("[W] no such plugin: %s", fpath)
		return
	}

	log.Debugf("[D] %s running...", fpath)

	var cmd *exec.Cmd
	if args == "" {
		cmd = exec.Command(fpath)
	} else {
		argList := PluginArgsParse(args)
		cmd = exec.Command(fpath, argList...)
	}
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
		logFile := filepath.Join(g.Config().Plugin.LogDir, plugin.FilePath+"("+plugin.Args+")"+".stderr.log")
		if _, err = file.WriteString(logFile, errStr); err != nil {
			log.Errorf("[E] write log to %s fail, error: %s\n", logFile, err)
		}
	}

	if isTimeout {
		// has be killed
		if err == nil {
			log.Infof("[I] timeout and kill process %s(%s) successfully", fpath, args)
		}

		if err != nil {
			log.Errorf("[E] kill process %s(%s) occur error: %v", fpath, args, err)
		}

		return
	}

	if err != nil {
		log.Errorf("[E] exec plugin %s(%s) fail. error: %v", fpath, args, err)
		return
	}

	// exec successfully
	data := stdout.Bytes()
	if len(data) == 0 {
		log.Debugf("[D] stdout of %s(%s) is blank", fpath, args)
		return
	}

	var metrics []*cmodel.MetricValue
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Errorf("[E] json.Unmarshal stdout of %s(%s) fail. error:%s stdout: \n%s\n", fpath, args, err, stdout.String())
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
