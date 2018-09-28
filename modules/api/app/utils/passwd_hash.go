package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashIt(passwd string) (hashed string) {
	b, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		hashed = ""
	} else {
		hashed = string(b)
	}
	return
}
