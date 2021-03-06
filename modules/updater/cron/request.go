package cron

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	f "github.com/toolkits/file"

	"github.com/open-falcon/falcon-plus/modules/updater/g"
	"github.com/open-falcon/falcon-plus/modules/updater/model"
)

// BuildHeartbeatRequest TODO:
func BuildHeartbeatRequest(hostname string, agentDirs []string) model.HeartbeatRequest {
	req := model.HeartbeatRequest{Hostname: hostname}

	realAgents := []*model.RealAgent{}
	now := time.Now().Unix()

	for _, agentDir := range agentDirs {
		// 如果目录下没有.version，我们认为这根本不是一个agent
		versionFile := path.Join(g.SelfDir, agentDir, ".version")
		if !f.IsExist(versionFile) {
			continue
		}

		version, err := f.ToTrimString(versionFile)
		if err != nil {
			log.Errorf("[E] read %s/.version fail: %v", agentDir, err)
			continue
		}

		controlFile := path.Join(g.SelfDir, agentDir, version, "control")
		if !f.IsExist(controlFile) {
			log.Warnf("[W] %s is nonexistent", controlFile)
			continue
		}

		cmd := exec.Command("./control", "status")
		cmd.Dir = path.Join(g.SelfDir, agentDir, version)
		bs, err := cmd.CombinedOutput()

		status := ""
		if err != nil {
			status = fmt.Sprintf("exec `./control status` fail: %s", err)
		} else {
			status = strings.TrimSpace(string(bs))
		}

		realAgent := &model.RealAgent{
			Name:      agentDir,
			Version:   version,
			Status:    status,
			Timestamp: now,
		}

		realAgents = append(realAgents, realAgent)
	}

	req.RealAgents = realAgents
	return req
}

// ListAgentDirs TODO:
func ListAgentDirs() ([]string, error) {
	agentDirs, err := f.DirsUnder(g.SelfDir)
	if err != nil {
		log.Errorf("[E] list dirs under %s fail: %v", g.SelfDir, err)
	}
	return agentDirs, err
}
