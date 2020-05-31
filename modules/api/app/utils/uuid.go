package utils

import (
	"strings"

	"github.com/satori/go.uuid"
)

// GenerateUUID 生成一个UUID
func GenerateUUID() string {
	sig := ""
	id := uuid.NewV1()
	sig = id.String()
	sig = strings.Replace(sig, "-", "", -1)
	return sig
}
