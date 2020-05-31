package api

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func init() {
	g.ParseConfig("../cfg.example.json")
}

func TestPortalAPI(t *testing.T) {
	Convey("Get action from api failed", t, func() {
		r := CurlAction(1)
		So(r.ID, ShouldEqual, 1)
	})

}
