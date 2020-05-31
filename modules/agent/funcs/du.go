package funcs

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/sys"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
)

var timeout = 30

// DuMetrics TODO:
func DuMetrics() (L []*cm.MetricValue) {
	paths := hbs.ReportPaths()
	result := make(chan *cm.MetricValue, len(paths))
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			var err error
			defer func() {
				if err != nil {
					log.Errorf("[E] %v", err)
					result <- GaugeValue(g.DU_BS, -1, "path="+path)
				}
				wg.Done()
			}()
			// 注意: Macos does not support -b.
			cmd := exec.Command("du", "-bs", path)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err = cmd.Start()
			if err != nil {
				return

			}
			err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Second)
			if isTimeout {
				err = fmt.Errorf(fmt.Sprintf("exec cmd : du -bs %s timeout", path))
				return
			}

			errStr := stderr.String()
			if errStr != "" {
				err = errors.New(errStr)
				return
			}

			if err != nil {
				err = fmt.Errorf(fmt.Sprintf("du -bs %s failed: %s", path, err.Error()))
				return
			}

			arr := strings.Fields(stdout.String())
			if len(arr) < 2 {
				err = fmt.Errorf(fmt.Sprintf("du -bs %s failed: %s", path, "return fields < 2"))
				return
			}

			size, err := strconv.ParseUint(arr[0], 10, 64)
			if err != nil {
				err = fmt.Errorf(fmt.Sprintf("cannot parse du -bs %s output", path))
				return
			}
			result <- GaugeValue(g.DU_BS, size, "path="+path)
		}(path)
	}
	wg.Wait()

	resultLen := len(result)
	for i := 0; i < resultLen; i++ {
		L = append(L, <-result)
	}
	return
}
