package api

import (
	"fmt"
	"time"

	"github.com/toolkits/net/httplib"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func LinkToSMS(content string) (string, error) {
	uri := fmt.Sprintf("%s/portal/links/store", g.Config().Api.Dashboard)
	req := httplib.Post(uri).SetTimeout(3*time.Second, 10*time.Second)
	req.Body([]byte(content))
	return req.String()
}
