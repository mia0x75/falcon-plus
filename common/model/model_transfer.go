package model

import (
	"fmt"
)

type TransferResponse struct {
	Message string
	Total   int
	Invalid int
	Latency int64
}

func (m *TransferResponse) String() string {
	return fmt.Sprintf(
		"<Total=%v, Invalid: %v, Latency=%vms, Message: %s>",
		m.Total,
		m.Invalid,
		m.Latency,
		m.Message,
	)
}
