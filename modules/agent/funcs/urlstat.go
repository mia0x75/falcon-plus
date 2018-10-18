package funcs

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

func UrlMetrics() (L []*cmodel.MetricValue) {
	reportUrls := g.ReportUrls()
	sz := len(reportUrls)
	if sz == 0 {
		return
	}
	hostname, err := g.Hostname()
	if err != nil {
		hostname = "None"
	}
	for furl, timeout := range reportUrls {
		tags := fmt.Sprintf("url=%v,timeout=%v,src=%v", furl, timeout, hostname)
		if ok, _ := probeUrl(furl, timeout); !ok {
			L = append(L, GaugeValue(g.URL_CHECK_HEALTH, 0, tags))
			continue
		}
		L = append(L, GaugeValue(g.URL_CHECK_HEALTH, 1, tags))
	}
	return
}

func probeUrl(furl string, timeout string) (bool, error) {
	bs, err := sys.CmdOutBytes("curl", "--max-filesize", "102400", "-I", "-m", timeout, "-o", "/dev/null", "-s", "-w", "%{http_code}", furl)
	if err != nil {
		log.Printf("probe url [%v] failed.the err is: [%v]\n", furl, err)
		return false, err
	}
	reader := bufio.NewReader(bytes.NewBuffer(bs))
	retcode, err := file.ReadLine(reader)
	if err != nil {
		log.Println("read retcode failed.err is:", err)
		return false, err
	}
	match, _ := regexp.MatchString("[20|30|10].*", strings.TrimSpace(string(retcode)))
	if !match {
		log.Printf("return code [%v] is not match regex. query url is [%v]", string(retcode), furl)
		return false, err
	}
	return true, err
}
