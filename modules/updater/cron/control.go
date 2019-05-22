package cron

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Control TODO:
func Control(workdir, arg string) (string, error) {
	cmd := exec.Command("./control", arg)
	cmd.Dir = workdir
	bs, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("[E] cd %s; ./control %s fail %v. output: %s", workdir, arg, err, string(bs))
	}
	return string(bs), err
}

// ControlStatus TODO:
func ControlStatus(workdir string) (string, error) {
	return Control(workdir, "status")
}

// ControlStart TODO:
func ControlStart(workdir string) (string, error) {
	return Control(workdir, "start")
}

// ControlStop TODO:
func ControlStop(workdir string) (string, error) {
	return Control(workdir, "stop")
}
