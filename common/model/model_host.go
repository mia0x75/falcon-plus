package model

import (
	"fmt"
)

type Host struct {
	ID   int
	Name string
}

func (m *Host) String() string {
	return fmt.Sprintf(
		"<ID: %d, name: %s>",
		m.ID,
		m.Name,
	)
}
