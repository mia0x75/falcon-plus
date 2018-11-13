package judge

import (
	log "github.com/Sirupsen/logrus"
)

func Start() {
	go StartJudgeCron()
	log.Info("[I] judge.Start ok")
}
