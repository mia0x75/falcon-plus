package g

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	pfcg "github.com/mia0x75/gopfc/g"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

// HTTPConfig HTTP配置
type HTTPConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

// RPCConfig RPC配置
type RPCConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

// SocketConfig TCP配置
type SocketConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout int    `json:"timeout"`
}

// JudgeConfig Judge配置
type JudgeConfig struct {
	Enabled        bool                    `json:"enabled"`
	Batch          int                     `json:"batch"`
	ConnectTimeout int                     `json:"connect_timeout"`
	ExecuteTimeout int                     `json:"execute_timeout"`
	MaxConnections int                     `json:"max_connections"`
	MaxIdle        int                     `json:"max_idle"`
	Replicas       int                     `json:"replicas"`
	Cluster        map[string]string       `json:"cluster"`
	ClusterList    map[string]*ClusterNode `json:"cluster_list"`
}

// GraphConfig Graph配置
type GraphConfig struct {
	Enabled        bool                    `json:"enabled"`
	Batch          int                     `json:"batch"`
	ConnectTimeout int                     `json:"connect_timeout"`
	ExecuteTimeout int                     `json:"execute_timeout"`
	MaxConnections int                     `json:"max_connections"`
	MaxIdle        int                     `json:"max_idle"`
	Replicas       int                     `json:"replicas"`
	Cluster        map[string]string       `json:"cluster"`
	ClusterList    map[string]*ClusterNode `json:"cluster_list"`
}

// TSDBConfig TSDB配置
type TSDBConfig struct {
	Enabled        bool   `json:"enabled"`
	Batch          int    `json:"batch"`
	ConnectTimeout int    `json:"connect_timeout"`
	ExecuteTimeout int    `json:"execute_timeout"`
	MaxConnections int    `json:"max_connections"`
	MaxIdle        int    `json:"max_idle"`
	MaxRetry       int    `json:"retry"`
	Address        string `json:"address"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
}

// TransferConfig Transfer配置
type TransferConfig struct {
	Enabled     bool              `json:"enabled"`
	Batch       int               `json:"batch"`
	ConnTimeout int               `json:"connTimeout"`
	CallTimeout int               `json:"callTimeout"`
	MaxConns    int               `json:"maxConns"`
	MaxIdle     int               `json:"maxIdle"`
	MaxRetry    int               `json:"retry"`
	Cluster     map[string]string `json:"cluster"`
}

// GlobalConfig Transfer模块配置
type GlobalConfig struct {
	Log           *LogConfig         `json:"log"`
	MinStep       int                `json:"min_step"` // 最小周期,单位sec
	HTTP          *HTTPConfig        `json:"http"`
	RPC           *RPCConfig         `json:"rpc"`
	Socket        *SocketConfig      `json:"socket"`
	Judge         *JudgeConfig       `json:"judge"`
	Graph         *GraphConfig       `json:"graph"`
	TSDB          *TSDBConfig        `json:"tsdb"`
	IgnoreMetrics map[string]bool    `json:"ignore"`
	Transfer      *TransferConfig    `json:"transfer"`
	PerfCounter   *pfcg.GlobalConfig `json:"pfc"`
}

// 变量定义
var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 获取Transfer模块的配置
func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// ParseConfig 解析配置文件
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

	// split cluster config
	c.Judge.ClusterList = formatClusterItems(c.Judge.Cluster)
	c.Graph.ClusterList = formatClusterItems(c.Graph.Cluster)

	configLock.Lock()
	defer configLock.Unlock()

	if c.PerfCounter != nil {
		c.PerfCounter.Debug = c.Log.Level == "debug"
		port := strings.Split(c.HTTP.Listen, ":")[1]
		if c.PerfCounter.Tags == "" {
			c.PerfCounter.Tags = fmt.Sprintf("port=%s", port)
		} else {
			c.PerfCounter.Tags = fmt.Sprintf("%s,port=%s", c.PerfCounter.Tags, port)
		}
	}

	config = &c

	log.Debugf("[D] read config file: %s successfully", cfg)
}

// ClusterNode TODO:
type ClusterNode struct {
	Addrs []string `json:"addrs"`
}

// NewClusterNode TODO:
func NewClusterNode(addrs []string) *ClusterNode {
	return &ClusterNode{addrs}
}

// map["node"]="host1,host2" --> map["node"]=["host1", "host2"]
func formatClusterItems(cluster map[string]string) map[string]*ClusterNode {
	ret := make(map[string]*ClusterNode)
	for node, clusterStr := range cluster {
		items := strings.Split(clusterStr, ",")
		nitems := make([]string, 0)
		for _, item := range items {
			nitems = append(nitems, strings.TrimSpace(item))
		}
		ret[node] = NewClusterNode(nitems)
	}

	return ret
}
