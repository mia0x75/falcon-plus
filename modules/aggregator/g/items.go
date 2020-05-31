package g

import (
	"fmt"
	"sync"
	"time"
)

// Cluster TODO:
type Cluster struct {
	ID          int64
	GroupID     int64
	Numerator   string
	Denominator string
	Endpoint    string
	Metric      string
	Tags        string
	DsType      string
	Step        int
	LastUpdate  time.Time
}

// String TODO:
func (s *Cluster) String() string {
	return fmt.Sprintf(
		"<ID: %d, GroupID: %d, Numerator: %s, Denominator: %s, Endpoint: %s, Metric: %s, Tags: %s, DsType: %s, Step: %d, LastUpdate: %v>",
		s.ID,
		s.GroupID,
		s.Numerator,
		s.Denominator,
		s.Endpoint,
		s.Metric,
		s.Tags,
		s.DsType,
		s.Step,
		s.LastUpdate,
	)
}

// SafeClusterMonitorItems key: Id+LastUpdate
type SafeClusterMonitorItems struct {
	sync.RWMutex
	M map[string]*Cluster
}

// NewClusterMonitorItems TODO:
func NewClusterMonitorItems() *SafeClusterMonitorItems {
	return &SafeClusterMonitorItems{M: make(map[string]*Cluster)}
}

// Init TODO:
func (s *SafeClusterMonitorItems) Init(m map[string]*Cluster) {
	s.Lock()
	defer s.Unlock()
	s.M = m
}

// Get TODO:
func (s *SafeClusterMonitorItems) Get() map[string]*Cluster {
	s.RLock()
	defer s.RUnlock()
	return s.M
}
