package cron

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/plugins"
)

// SyncMinePlugins TODO:
func SyncMinePlugins() {
	if !g.Config().Plugin.Enabled {
		return
	}

	if len(g.Config().Heartbeat.Addrs) == 0 {
		return
	}

	go syncMinePlugins()
}

func syncMinePlugins() {
	var (
		timestamp  int64 = -1
		pluginDirs []string
	)

	d := time.Duration(g.Config().Heartbeat.Interval) * time.Second
	for range time.Tick(d) {
		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req := cm.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var resp cm.AgentPluginsResponse
		err = g.HbsClient.Call("Agent.MinePlugins", req, &resp)
		if err != nil {
			log.Errorf("[E] Call Agent.MinePlugins fail, error: %v", err)
			continue
		}

		if resp.Timestamp <= timestamp {
			continue
		}

		pluginDirs = resp.Plugins
		timestamp = resp.Timestamp

		log.Debugf("[D] Call Agent.MinePlugin response: %v\n", &resp)

		if len(pluginDirs) == 0 {
			plugins.ClearAllPlugins()
			continue
		}

		desiredAll := make(map[string]*plugins.Plugin)
		filefmtScripts := [][]string{}
		dirfmtScripts := []string{}

		for _, dir := range pluginDirs {
			// script_path could be a DIR or a SCRIPT_FILE_WITH_OR_WITHOUT_ARGS
			// 比如： sys/ntp/60_ntp.py(arg1,arg2) 或者 sys/ntp/60_ntp.py 或者 sys/ntp
			// 1. 参数只对单个脚本文件生效，目录不支持参数
			// 2. 如果某个目录下的某个脚本被单独绑定到某个机器，那么再次绑定该目录时，该文件会不会再次执行
			args := ""

			re := regexp.MustCompile(`(.*)\((.*)\)`)
			pathArgs := re.FindAllStringSubmatch(dir, -1)
			if pathArgs != nil {
				dir = pathArgs[0][1]
				args = pathArgs[0][2]
			}

			absPath := filepath.Join(g.Config().Plugin.Dir, dir)
			if !file.IsExist(absPath) {
				continue
			}
			if file.IsFile(absPath) {
				filefmtScripts = append(filefmtScripts, []string{dir, args})
				continue
			}

			dirfmtScripts = append(dirfmtScripts, dir)
		}

		taken := make(map[string]struct{})
		for _, script := range filefmtScripts {
			absPath := filepath.Join(g.Config().Plugin.Dir, script[0])
			_, fileName := filepath.Split(absPath)
			arr := strings.Split(fileName, "_")
			var cycle int
			var err error
			cycle, err = strconv.Atoi(arr[0])
			if err == nil {
				fi, _ := os.Stat(absPath)
				plugin := &plugins.Plugin{FilePath: script[0], MTime: fi.ModTime().Unix(), Cycle: cycle, Args: script[1]}
				desiredAll[script[0]+"("+script[1]+")"] = plugin
			}
			// 针对某个 hostgroup 绑定了单个脚本后，再绑定该脚本的目录时，会忽略目录中的该文件
			taken[script[0]] = struct{}{}
		}

		for _, scriptPath := range dirfmtScripts {
			ps := plugins.ListPlugins(strings.Trim(scriptPath, "/"))
			for k, p := range ps {
				if _, ok := taken[k]; ok {
					continue
				}
				desiredAll[k] = p
			}
		}

		plugins.DelNoUsePlugins(desiredAll)
		plugins.AddNewPlugins(desiredAll)
		log.Debugf("[D] Current plugins: %v\n", plugins.Plugins)
	}
}
