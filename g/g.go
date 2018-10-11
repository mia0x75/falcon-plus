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
		"agent":      "/etc/mfp/agent.json",
		"aggregator": "/etc/mfp/aggregator.json",
		"graph":      "/etc/mfp/graph.json",
		"hbs":        "/etc/mfp/hbs.json",
		"judge":      "/etc/mfp/judge.json",
		"nodata":     "/etc/mfp/nodata.json",
		"transfer":   "/etc/mfp/transfer.json",
		"gateway":    "/etc/mfp/gateway.json",
		"api":        "/etc/mfp/api.json",
		"alarm":      "/etc/mfp/alarm.json",
		"updater":    "/etc/mfp/updater.json",
		"exporter":   "/etc/mfp/exporter.json",
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
		"agent":      "/var/log/mfp/agent.log",
		"aggregator": "/var/log/mfp/aggregator.log",
		"graph":      "/var/log/mfp/graph.log",
		"hbs":        "/var/log/mfp/hbs.log",
		"judge":      "/var/log/mfp/judge.log",
		"nodata":     "/var/log/mfp/nodata.log",
		"transfer":   "/var/log/mfp/transfer.log",
		"gateway":    "/var/log/mfp/gateway.log",
		"api":        "/var/log/mfp/api.log",
		"alarm":      "/var/log/mfp/alarm.log",
		"updater":    "/var/log/mfp/updater.log",
		"exporter":   "/var/log/mfp/exporter.log",
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

func Bin(name string) string {
	p, _ := filepath.Abs(BinOf[name])
	return p
}

func Cfg(name string) string {
	p, _ := filepath.Abs(cfgOf[name])
	return p
}

func LogPath(name string) string {
	p, _ := filepath.Abs(logpathOf[name])
	return p
}

func LogDir(name string) string {
	d, _ := filepath.Abs(filepath.Dir(logpathOf[name]))
	return d
}
