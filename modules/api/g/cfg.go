package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type DbConfig struct {
	Portal    string
	Graph     string
	Uic       string
	Dashboard string
	Alarms    string
}

type GraphsConfig struct {
	MaxConnections int `json:"max_connections"`
	MaxIdle        int `json:"max_idle"`
	ConnectTimeout int `json:"connect_timeout"`
	ExecuteTimeout int `json:"execute_timeout"`
	Replicas       int `json:"replicas"`
	Cluster        map[string]string
}

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Address string `json:"addr"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log            *LogConfig    `json:"log"`
	Listen         string        `json:"listen"`
	AccessControl  bool          `json:"access_control"`
	SignupDisable  bool          `json:"signup_disable"`
	SkipAuth       bool          `json:"skip_auth"`
	DefaultToken   string        `json:"default_token"`
	GenDoc         bool          `json:"gen_doc"`
	GenDocPath     string        `json:"gen_doc_path"`
	MetricListFile string        `json:"metric_list_file"`
	Rpc            *RpcConfig    `json:"rpc"`
	DB             *DbConfig     `json:"db"`
	Graphs         *GraphsConfig `json:"graphs"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
