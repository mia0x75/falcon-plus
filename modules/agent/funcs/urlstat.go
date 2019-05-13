package funcs

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
)

// URLMetrics TODO:
func URLMetrics() (L []*cmodel.MetricValue) {
	urls := hbs.ReportUrls()
	sz := len(urls)
	if sz == 0 {
		return
	}
	hostname, err := g.Hostname()
	if err != nil {
		hostname = "None"
	}
	for url, timeout := range urls {
		tags := fmt.Sprintf("url=%v,timeout=%v,src=%v", url, timeout, hostname)
		if ok, _ := probeURL(url, timeout); !ok {
			L = append(L, GaugeValue(g.URL_CHECK_HEALTH, 0, tags))
			continue
		}
		L = append(L, GaugeValue(g.URL_CHECK_HEALTH, 1, tags))
	}
	return
}

func probeURL(furl string, timeout string) (bool, error) {
	bs, err := sys.CmdOutBytes("curl", "--max-filesize", "102400", "-I", "-m", timeout, "-o", "/dev/null", "-s", "-w", "%{http_code}", furl)
	if err != nil {
		log.Errorf("[E] probe url %s fail, error: %v", furl, err)
		return false, err
	}
	reader := bufio.NewReader(bytes.NewBuffer(bs))
	retcode, err := file.ReadLine(reader)
	if err != nil {
		log.Errorf("[E] read retcode failed: %v", err)
		return false, err
	}
	match, _ := regexp.MatchString("[20|30|10].*", strings.TrimSpace(string(retcode)))
	if !match {
		log.Warnf("[W] return code %s is not match regex. query url is %s", string(retcode), furl)
		return false, err
	}
	return true, err
}
