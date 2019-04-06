package plugins

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// ListPlugins return: dict{sys/ntp/60_ntp.py : *Plugin}
func ListPlugins(scriptPath string) map[string]*Plugin {
	ret := make(map[string]*Plugin)
	if scriptPath == "" {
		return ret
	}

	absPath := filepath.Join(g.Config().Plugin.Dir, scriptPath)
	fs, err := ioutil.ReadDir(absPath)
	if err != nil {
		log.Errorf("[E] can not list files under: %s", absPath)
		return ret
	}

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		filename := f.Name()
		arr := strings.Split(filename, "_")
		if len(arr) < 2 {
			continue
		}

		// filename should be: $cycle_$xx
		var cycle int
		cycle, err = strconv.Atoi(arr[0])
		if err != nil {
			continue
		}

		fpath := filepath.Join(scriptPath, filename)
		plugin := &Plugin{FilePath: fpath, MTime: f.ModTime().Unix(), Cycle: cycle, Args: ""}
		ret[fpath] = plugin
	}
	return ret
}
