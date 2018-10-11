package gopfc

import (
	"github.com/mia0x75/gopfc/g"
	"github.com/mia0x75/gopfc/http"
	"github.com/mia0x75/gopfc/metric"
)

func Start() {
	if g.Config() == nil {
		return
	}

	// http
	if g.Config().Http.Enabled {
		go http.Start()
	}

	// base collector cron
	if len(g.Config().Bases) > 0 {
		go metric.CollectBase(g.Config().Bases)
	}

	// push cron
	if g.Config().Push.Enabled {
		go metric.PushToFalcon()
	}
}
