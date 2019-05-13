package g

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

// PluginConfig 插件配置信息
type PluginConfig struct {
	Enabled bool   `json:"enabled"`
	Dir     string `json:"dir"`
	Git     string `json:"git"`
	LogDir  string `json:"logs"`
}

// HeartbeatConfig 心跳服务器配置信息
type HeartbeatConfig struct {
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

// TransferConfig 数据中转服务器配置信息
type TransferConfig struct {
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

// HTTPConfig TODO:
type HTTPConfig struct {
	Listen   string `json:"listen"`
	Backdoor bool   `json:"backdoor"`
	Root     string `json:"root"`
}

// CollectorConfig TODO:
type CollectorConfig struct {
	System  *SystemConfig  `json:"system"`
	MySQL   *MySQLConfig   `json:"mysql"`
	Redis   *RedisConfig   `json:"redis"`
	MongoDB *MongoDBConfig `json:"mongodb"`
	Jmx     *JmxConfig     `json:"jmx"`
	Nginx   *NginxConfig   `json:"nginx"`
}

// SystemConfig TODO:
type SystemConfig struct {
	IfacePrefix []string `json:"iface_prefix"`
	MountPoint  []string `json:"mount_point"`
	Interval    int      `json:"interval"`
}

// MySQLConfig TODO:
type MySQLConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Passowrd string `json:"password"`
	Port     int    `json:"port"`
	Interval int    `json:"interval"`
}

// RedisConfig TODO:
type RedisConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"cluster_name"`
	Passowrd string `json:"password"`
	Port     int    `json:"port"`
	Interval int    `json:"interval"`
}

// MongoDBConfig TODO:
type MongoDBConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Passowrd string `json:"password"`
	Port     int    `json:"port"`
	Interval int    `json:"interval"`
}

// JmxConfig TODO:
type JmxConfig struct {
	Enabled  bool `json:"enabled"`
	Interval int  `json:"interval"`
}

// NginxConfig TODO:
type NginxConfig struct {
	Enabled  bool `json:"enabled"`
	Interval int  `json:"interval"`
}

// LogConfig TODO:
type LogConfig struct {
	Level string `json:"level"`
}

// GlobalConfig TODO:
type GlobalConfig struct {
	Log           *LogConfig        `json:"log"`
	Hostname      string            `json:"hostname"`
	IP            string            `json:"ip"`
	Plugin        *PluginConfig     `json:"plugin"`
	Heartbeat     *HeartbeatConfig  `json:"heartbeat"`
	Transfer      *TransferConfig   `json:"transfer"`
	HTTP          *HTTPConfig       `json:"http"`
	Collector     *CollectorConfig  `json:"collector"`
	DefaultTags   map[string]string `json:"default_tags"`
	IgnoreMetrics map[string]bool   `json:"ignore"`
}

// TODO:
var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

// Config 获取配置信息
func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

// Hostname TODO:
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
			log.Panic("[P] /etc/endpoint.env missing ENDPOINT")
		}
		if !strings.ContainsAny(str, "=") {
			log.Panic("[P] /etc/endpoint.env missing =")
		}
		str = strings.Trim((strings.SplitAfter(str, "ENDPOINT"))[1], " ")
		hostname = strings.Trim((strings.SplitAfter(str, "="))[1], " ")
		if len(hostname) > 1 {
			return hostname, nil
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("[E] os.Hostname() fail: %v", err)
	}
	return hostname, err
}

// IP TODO:
func IP() string {
	ip := Config().IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIP) > 0 {
		ip = LocalIP
	}

	return ip
}

// ParseConfig 读取并解析配置
func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatal("[F] use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalf("[F] config file: %s is not existent. maybe you need `mv cfg.example.json cfg.json`", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalf("[F] read config file: %s fail: %v", cfg, err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalf("[F] parse config file: %s fail: %v", cfg, err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Debugf("[D] read config file: %s successfully", cfg)
}
