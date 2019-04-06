package api

import (
	"testing"

	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func init() {
	log.SetLevel(log.DebugLevel)
	g.ParseConfig("../cfg.example.json")
}

func TestUicAPI(t *testing.T) {
	Convey("Get team users from api failed", t, func() {
		r := CurlUic("plus-dev")
		for _, x := range r {
			log.Debug("[D] %#v", x)
		}
		So(len(r), ShouldEqual, 1)
	})
}
