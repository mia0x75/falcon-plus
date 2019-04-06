package judge

import (
	log "github.com/sirupsen/logrus"
)

func Start() {
	go StartJudgeCron()
	log.Info("[I] judge.Start ok")
}
