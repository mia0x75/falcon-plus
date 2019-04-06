package g

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Modules map[string]bool
var BinOf map[string]string
var cfgOf map[string]string
var ModuleApps map[string]string
var logpathOf map[string]string
var PidOf map[string]string
var AllModulesInOrder []string

func init() {
	Modules = map[string]bool{
		"agent":      true,
		"aggregator": true,
		"graph":      true,
		"hbs":        true,
		"judge":      true,
		"nodata":     true,
		"transfer":   true,
		"gateway":    true,
		"api":        true,
		"alarm":      true,
		"updater":    true,
		"exporter":   true,
	}

	BinOf = map[string]string{
		"agent":      "/usr/bin/falcon-agent",
		"aggregator": "/usr/bin/falcon-aggregator",
		"graph":      "/usr/bin/falcon-graph",
		"hbs":        "/usr/bin/falcon-hbs",
		"judge":      "/usr/bin/falcon-judge",
		"nodata":     "/usr/bin/falcon-nodata",
		"transfer":   "/usr/bin/falcon-transfer",
		"gateway":    "/usr/bin/falcon-gateway",
		"api":        "/usr/bin/falcon-api",
		"alarm":      "/usr/bin/falcon-alarm",
		"updater":    "/usr/bin/falcon-updater",
		"exporter":   "/usr/bin/falcon-exporter",
	}

	cfgOf = map[string]string{
		"agent":      "/etc/fpm/agent.json",
		"aggregator": "/etc/fpm/aggregator.json",
		"graph":      "/etc/fpm/graph.json",
		"hbs":        "/etc/fpm/hbs.json",
		"judge":      "/etc/fpm/judge.json",
		"nodata":     "/etc/fpm/nodata.json",
		"transfer":   "/etc/fpm/transfer.json",
		"gateway":    "/etc/fpm/gateway.json",
		"api":        "/etc/fpm/api.json",
		"alarm":      "/etc/fpm/alarm.json",
		"updater":    "/etc/fpm/updater.json",
		"exporter":   "/etc/fpm/exporter.json",
	}

	ModuleApps = map[string]string{
		"agent":      "falcon-agent",
		"aggregator": "falcon-aggregator",
		"graph":      "falcon-graph",
		"hbs":        "falcon-hbs",
		"judge":      "falcon-judge",
		"nodata":     "falcon-nodata",
		"transfer":   "falcon-transfer",
		"gateway":    "falcon-gateway",
		"api":        "falcon-api",
		"alarm":      "falcon-alarm",
		"updater":    "falcon-updater",
		"exporter":   "falcon-exporter",
	}

	logpathOf = map[string]string{
		"agent":      "/var/log/fpm/agent.log",
		"aggregator": "/var/log/fpm/aggregator.log",
		"graph":      "/var/log/fpm/graph.log",
		"hbs":        "/var/log/fpm/hbs.log",
		"judge":      "/var/log/fpm/judge.log",
		"nodata":     "/var/log/fpm/nodata.log",
		"transfer":   "/var/log/fpm/transfer.log",
		"gateway":    "/var/log/fpm/gateway.log",
		"api":        "/var/log/fpm/api.log",
		"alarm":      "/var/log/fpm/alarm.log",
		"updater":    "/var/log/fpm/updater.log",
		"exporter":   "/var/log/fpm/exporter.log",
	}

	PidOf = map[string]string{
		"agent":      "<NOT SET>",
		"aggregator": "<NOT SET>",
		"graph":      "<NOT SET>",
		"hbs":        "<NOT SET>",
		"judge":      "<NOT SET>",
		"nodata":     "<NOT SET>",
		"transfer":   "<NOT SET>",
		"gateway":    "<NOT SET>",
		"api":        "<NOT SET>",
		"alarm":      "<NOT SET>",
		"updater":    "<NOT SET>",
		"exporter":   "<NOT SET>",
	}

	// Modules are deployed in this order
	AllModulesInOrder = []string{
		"graph",
		"alarm",
		"judge",
		"api",
		"transfer",
		"nodata",
		"aggregator",
		"hbs",
		"agent",
		"gateway",
		"updater",
		"exporter",
	}
}

// Bin returns executable bin file path of the specify module
func Bin(name string) string {
	p, _ := filepath.Abs(BinOf[name])
	return p
}

// Cfg returns config file path of the specify module
func Cfg(name string) string {
	p, _ := filepath.Abs(cfgOf[name])
	return p
}

// LogPath returns the log file path of the specify module
func LogPath(name string) string {
	p, _ := filepath.Abs(logpathOf[name])
	return p
}

// LogDir return the log folder of the specify module
func LogDir(name string) string {
	d, _ := filepath.Abs(filepath.Dir(logpathOf[name]))
	return d
}

func HasLogfile(name string) bool {
	if _, err := os.Stat(LogPath(name)); err != nil {
		return false
	}
	return true
}

func PreqOrder(moduleArgs []string) []string {
	if len(moduleArgs) == 0 {
		return []string{}
	}

	var modulesInOrder []string

	// get arguments which are found in the order
	for _, nameOrder := range AllModulesInOrder {
		for _, nameArg := range moduleArgs {
			if nameOrder == nameArg {
				modulesInOrder = append(modulesInOrder, nameOrder)
			}
		}
	}
	// get arguments which are not found in the order
	for _, nameArg := range moduleArgs {
		end := 0
		for _, nameOrder := range modulesInOrder {
			if nameOrder == nameArg {
				break
			}
			end++
		}
		if end == len(modulesInOrder) {
			modulesInOrder = append(modulesInOrder, nameArg)
		}
	}
	return modulesInOrder
}

func Rel(p string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// filepath.Abs() returns an error only when os.Getwd() returns an error;
	abs, _ := filepath.Abs(p)

	r, err := filepath.Rel(wd, abs)
	if err != nil {
		return ""
	}

	return r
}

func HasCfg(name string) bool {
	if _, err := os.Stat(Cfg(name)); err != nil {
		return false
	}
	return true
}

func HasModule(name string) bool {
	return Modules[name]
}

func setPid(name string) {
	output, _ := exec.Command("pgrep", "-f", ModuleApps[name]).Output()
	pidStr := strings.TrimSpace(string(output))
	PidOf[name] = pidStr
}

func Pid(name string) string {
	if PidOf[name] == "<NOT SET>" {
		setPid(name)
	}
	return PidOf[name]
}

func IsRunning(name string) bool {
	setPid(name)
	return Pid(name) != ""
}

func RmDup(args []string) []string {
	if len(args) == 0 {
		return []string{}
	}
	if len(args) == 1 {
		return args
	}

	ret := []string{}
	isDup := make(map[string]bool)
	for _, arg := range args {
		if isDup[arg] == true {
			continue
		}
		ret = append(ret, arg)
		isDup[arg] = true
	}
	return ret
}
