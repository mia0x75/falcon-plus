package graph

import (
	"time"
)

type Endpoint struct {
	ID               uint              `gorm:"primary_key"`
	Endpoint         string            `json:"endpoint"`
	Ts               int               `json:"-"`
	TCreate          time.Time         `json:"-"`
	TModify          time.Time         `json:"-"`
	EndpointCounters []EndpointCounter `gorm:"ForeignKey:EndpointIDE"`
}

type Host struct {
	ID       uint   `gorm:"primary_key"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
}

func (Endpoint) TableName() string {
	return "endpoint"
}
