package utils

import (
	"strings"

	"github.com/satori/go.uuid"
)

func GenerateUUID() string {
	sig := ""
	id := uuid.NewV1()
	sig = id.String()
	sig = strings.Replace(sig, "-", "", -1)
	return sig
}
