package model

import (
	"fmt"
)

// code == 0 => success
// code == 1 => bad request
type SimpleRPCResponse struct {
	Code int `json:"code"`
}

func (m *SimpleRPCResponse) String() string {
	return fmt.Sprintf("<Code: %d>", m.Code)
}

type NullRPCRequest struct {
}
