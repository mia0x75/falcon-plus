package judge

import (
	log "github.com/Sirupsen/logrus"
)

func Start() {
	go StartJudgeCron()
	log.Println("judge.Start ok")
}
