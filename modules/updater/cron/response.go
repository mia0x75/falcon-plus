package cron

import (
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/updater/g"
	"github.com/open-falcon/falcon-plus/modules/updater/model"
)

// HandleHeartbeatResponse TODO:
func HandleHeartbeatResponse(respone *model.HeartbeatResponse) {
	if respone.ErrorMessage != "" {
		log.Errorf("[E] receive error message: %s", respone.ErrorMessage)
		return
	}

	das := respone.DesiredAgents
	if das == nil || len(das) == 0 {
		return
	}

	for _, da := range das {
		da.FillAttrs(g.SelfDir)

		if g.Config().DesiredAgent == "" || g.Config().DesiredAgent == da.Name {
			HandleDesiredAgent(da)
		}
	}
}

// HandleDesiredAgent TODO:
func HandleDesiredAgent(da *model.DesiredAgent) {
	if da.Cmd == "start" {
		StartDesiredAgent(da)
	} else if da.Cmd == "stop" {
		StopDesiredAgent(da)
	} else {
		log.Warnf("[W] unknown cmd: %v", da)
	}
}
