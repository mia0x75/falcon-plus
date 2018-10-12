package g

import (
	"os"
	"sync"
)

type GlobalConfig struct {
	Debug    bool        `json:"debug"`
	Hostname string      `json:"hostname"`
	Prefix   string      `json:"prefix"`
	Tags     string      `json:"tags"`
	Step     int64       `json:"step"`
	Bases    []string    `json:"bases"`
	Push     *PushConfig `json:"push"`
	Http     *HttpConfig `json:"http"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type PushConfig struct {
	Enabled bool   `json:"enabled"`
	Api     string `json:"api"`
}

var (
	defaultPrefix = ""
	defaultTags   = ""
	defaultStep   = int64(60) //time in sec
	defaultBases  = []string{"debug", "runtime"}
	defaultPush   = &PushConfig{Enabled: true, Api: "http://127.0.0.1:1988/v1/push"}
	defaultHttp   = &HttpConfig{Enabled: false, Listen: ""}
)

var (
	cfg     *GlobalConfig
	cfgLock = new(sync.RWMutex)
)

var (
	DefaultConfig = &GlobalConfig{
		Debug:    false,
		Hostname: defaultHostname(),
		Prefix:   defaultPrefix,
		Tags:     defaultTags,
		Step:     defaultStep,
		Bases:    defaultBases,
		Push:     defaultPush,
		Http:     defaultHttp,
	}
)

func PFC() {
	PFCWithConfig(DefaultConfig)
}

func PFCWithConfig(c *GlobalConfig) {
	if c.Hostname == "" {
		c.Hostname = DefaultConfig.Hostname
	}
	if c.Push == nil {
		c.Push = DefaultConfig.Push
	}
	if c.Http == nil {
		c.Http = DefaultConfig.Http
	}
	cfg = c
}

//
func Config() *GlobalConfig {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}

func defaultHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
