package utils

import (
	"strings"

	"github.com/satori/go.uuid"
)

func GenerateUUID() string {
	sig := ""
	if id, err := uuid.NewV1(); err != nil {
		sig = id.String()
		sig = strings.Replace(sig, "-", "", -1)
	}
	return sig
}
