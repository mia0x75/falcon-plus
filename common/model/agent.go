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

func (this *AgentReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, IP: %s, AgentVersion: %s, PluginVersion: %s, Hostgroup: %s>",
		this.Hostname,
		this.IP,
		this.AgentVersion,
		this.PluginVersion,
		this.Hostgroup,
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

func (this *AgentHeartbeatRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, Checksum: %s>",
		this.Hostname,
		this.Checksum,
	)
}

type AgentPluginsResponse struct {
	Plugins   []string
	Timestamp int64
}

func (this *AgentPluginsResponse) String() string {
	return fmt.Sprintf(
		"<Plugins: %v, Timestamp: %v>",
		this.Plugins,
		this.Timestamp,
	)
}

// e.g. net.port.listen or proc.num
type BuiltinMetric struct {
	Metric string
	Tags   string
}

func (this *BuiltinMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		this.Metric,
		this.Tags,
	)
}

type BuiltinMetricResponse struct {
	Metrics   []*BuiltinMetric
	Checksum  string
	Timestamp int64
}

func (this *BuiltinMetricResponse) String() string {
	return fmt.Sprintf(
		"<Metrics: %v, Checksum: %s, Timestamp: %v>",
		this.Metrics,
		this.Checksum,
		this.Timestamp,
	)
}

type BuiltinMetricSlice []*BuiltinMetric

func (this BuiltinMetricSlice) Len() int {
	return len(this)
}

func (this BuiltinMetricSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this BuiltinMetricSlice) Less(i, j int) bool {
	return this[i].String() < this[j].String()
}
