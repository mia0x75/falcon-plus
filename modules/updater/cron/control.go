package cron

import (
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

func Control(workdir, arg string) (string, error) {
	cmd := exec.Command("./control", arg)
	cmd.Dir = workdir
	bs, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("[E] cd %s; ./control %s fail %v. output: %s", workdir, arg, err, string(bs))
	}
	return string(bs), err
}

func ControlStatus(workdir string) (string, error) {
	return Control(workdir, "status")
}

func ControlStart(workdir string) (string, error) {
	return Control(workdir, "start")
}

func ControlStop(workdir string) (string, error) {
	return Control(workdir, "stop")
}
