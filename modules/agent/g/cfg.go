package g

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
)

type PluginConfig struct {
	Enabled bool   `json:"enabled"`
	Dir     string `json:"dir"`
	Git     string `json:"git"`
	LogDir  string `json:"logs"`
}

type HeartbeatConfig struct {
	Enabled  bool     `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

type TransferConfig struct {
	Enabled  bool     `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

type HttpConfig struct {
	Enabled  bool   `json:"enabled"`
	Listen   string `json:"listen"`
	Backdoor bool   `json:"backdoor"`
	Root     string `json:"root"`
}

type CollectorConfig struct {
	System  *SystemConfig  `json:"system"`
	MySQL   *MySQLConfig   `json:"mysql"`
	Redis   *RedisConfig   `json:"redis"`
	MongoDB *MongoDBConfig `json:"mongodb"`
	Jmx     *JmxConfig     `json:"jmx"`
	Nginx   *NginxConfig   `json:"nginx"`
}

type SystemConfig struct {
	IfacePrefix []string `json:"iface_prefix"`
	MountPoint  []string `json:"mount_point"`
}

type MySQLConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Passowrd string `json:"password"`
	Port     int    `json:"port"`
	Interval int    `json:"interval"`
}

type RedisConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"cluster_name"`
	Passowrd string `json:"password"`
	Port     int    `json:"port"`
	Interval int    `json:"interval"`
}

type MongoDBConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Passowrd string `json:"password"`
	Port     int    `json:"port"`
	Interval int    `json:"interval"`
}

type JmxConfig struct {
	Enabled bool `json:"enabled"`
}

type NginxConfig struct {
	Enabled bool `json:"enabled"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log           *LogConfig        `json:"log"`
	Hostname      string            `json:"hostname"`
	IP            string            `json:"ip"`
	Plugin        *PluginConfig     `json:"plugin"`
	Heartbeat     *HeartbeatConfig  `json:"heartbeat"`
	Transfer      *TransferConfig   `json:"transfer"`
	Http          *HttpConfig       `json:"http"`
	Collector     *CollectorConfig  `json:"collector"`
	DefaultTags   map[string]string `json:"default_tags"`
	IgnoreMetrics map[string]bool   `json:"ignore"`
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

func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	if os.Getenv("FALCON_ENDPOINT") != "" {
		hostname = os.Getenv("FALCON_ENDPOINT")
		return hostname, nil
	}

	if os.Getenv("ENDPOINT") != "" {
		hostname = os.Getenv("ENDPOINT")
		return hostname, nil
	}

	// parse /etc/endpoint.env ( ENDPOINT=xxx_xxx_1.2.3.4 )
	filePath := "/etc/endpoint.env"
	if _, err := os.Stat(filePath); err == nil {
		data, _ := ioutil.ReadFile(filePath)
		str := string(data)
		if !strings.ContainsAny(str, "ENDPOINT") {
			log.Panic("ERROR: /etc/endpoint.env missing ENDPOINT")
		}
		if !strings.ContainsAny(str, "=") {
			log.Panic("ERROR: /etc/endpoint.env missing =")
		}
		str = strings.Trim((strings.SplitAfter(str, "ENDPOINT"))[1], " ")
		hostname = strings.Trim((strings.SplitAfter(str, "="))[1], " ")
		if len(hostname) > 1 {
			return hostname, nil
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func IP() string {
	ip := Config().IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIp) > 0 {
		ip = LocalIp
	}

	return ip
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

	log.Debugln("read config file:", cfg, "successfully")
}
