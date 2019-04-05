package g

import (
	"path/filepath"
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
