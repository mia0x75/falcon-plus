package model

import (
	"fmt"
)

type Template struct {
	ID       int    `json:"ID"`
	Name     string `json:"name"`
	ParentID int    `json:"parentID"`
	ActionID int    `json:"actionID"`
	Creator  int    `json:"creator"`
}

func (m *Template) String() string {
	return fmt.Sprintf(
		"<ID: %d, Name: %s, ParentID: %d, ActionID: %d, Creator: %d>",
		m.ID,
		m.Name,
		m.ParentID,
		m.ActionID,
		m.Creator,
	)
}
