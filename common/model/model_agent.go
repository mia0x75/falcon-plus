package model

import (
	"fmt"
)

type AgentReportRequest struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
	Hostgroup     string
	MachineID     string
	Timezone      string
	OS            *system
	Kernel        *kernel
	Product       *product
	Mainboard     *mainboard
	Chassis       *chassis
	Bios          *bios
	CPU           *cpu
	Memory        *memory
	Storages      []*storage
	Networks      []*network
}

type system struct {
	Name         string
	Vendor       string
	Version      string
	Release      string
	Architecture string
}

type kernel struct {
	Release      string
	Version      string
	Architecture string
}

type product struct {
	Name    string
	Vendor  string
	Version string
	Serial  string
}

type mainboard struct {
	Name    string
	Vendor  string
	Version string
	Serial  string
}

type chassis struct {
	Type   int
	Vendor string
}

type bios struct {
	Vendor  string
	Version string
	Date    string
}

type cpu struct {
	Vendor  string
	Model   string
	Speed   int
	Cache   int
	Cpus    int
	Cores   int
	Threads int
}

type memory struct {
	Type  string
	Speed int
	Size  int
}

type storage struct {
	Name   string
	Driver string
	Vendor string
	Model  string
	Serial string
	Size   int
}

type network struct {
	Name   string
	Driver string
	MAC    string
	Port   string
	Speed  int
}

func (m *AgentReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, IP: %s, AgentVersion: %s, PluginVersion: %s, Hostgroup: %s>",
		m.Hostname,
		m.IP,
		m.AgentVersion,
		m.PluginVersion,
		m.Hostgroup,
	)
}

type AgentUpdateInfo struct {
	LastUpdate    int64
	ReportRequest *AgentReportRequest
}

type AgentHeartbeatRequest struct {
	Hostname string
	Checksum string
}

func (m *AgentHeartbeatRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, Checksum: %s>",
		m.Hostname,
		m.Checksum,
	)
}

type AgentPluginsResponse struct {
	Plugins   []string
	Timestamp int64
}

func (m *AgentPluginsResponse) String() string {
	return fmt.Sprintf(
		"<Plugins: %v, Timestamp: %v>",
		m.Plugins,
		m.Timestamp,
	)
}

// e.g. net.port.listen or proc.num
type BuiltinMetric struct {
	Metric string
	Tags   string
}

func (m *BuiltinMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		m.Metric,
		m.Tags,
	)
}

type BuiltinMetricResponse struct {
	Metrics   []*BuiltinMetric
	Checksum  string
	Timestamp int64
}

func (m *BuiltinMetricResponse) String() string {
	return fmt.Sprintf(
		"<Metrics: %v, Checksum: %s, Timestamp: %v>",
		m.Metrics,
		m.Checksum,
		m.Timestamp,
	)
}

type BuiltinMetricSlice []*BuiltinMetric

func (m BuiltinMetricSlice) Len() int {
	return len(m)
}

func (m BuiltinMetricSlice) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m BuiltinMetricSlice) Less(i, j int) bool {
	return m[i].String() < m[j].String()
}
