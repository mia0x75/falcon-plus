package cron

import (
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"

	"github.com/open-falcon/falcon-plus/modules/updater/g"
	"github.com/open-falcon/falcon-plus/modules/updater/model"
)

// StopDesiredAgent TODO:
func StopDesiredAgent(da *model.DesiredAgent) {
	if !file.IsExist(da.ControlFilepath) {
		return
	}

	ControlStopIn(da.AgentVersionDir)
}

// StopAgentOf TODO:
func StopAgentOf(agentName, newVersion string) error {
	agentDir := path.Join(g.SelfDir, agentName)
	versionFile := path.Join(agentDir, ".version")

	if !file.IsExist(versionFile) {
		log.Warnf("[W] %s is nonexistent", versionFile)
		return nil
	}

	version, err := file.ToTrimString(versionFile)
	if err != nil {
		log.Errorf("[E] read %s fail %s", version, err)
		return nil
	}

	if version == newVersion {
		// do nothing
		return nil
	}

	versionDir := path.Join(agentDir, version)
	if !file.IsExist(versionDir) {
		log.Warnf("[W] %s nonexistent", versionDir)
		return nil
	}

	return ControlStopIn(versionDir)
}

// ControlStopIn TODO:
func ControlStopIn(workdir string) error {
	if !file.IsExist(workdir) {
		return nil
	}

	out, err := ControlStatus(workdir)
	if err == nil && strings.Contains(out, "stopped") {
		return nil
	}

	_, err = ControlStop(workdir)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 3)

	out, err = ControlStatus(workdir)
	if err == nil && strings.Contains(out, "stopped") {
		return nil
	}

	return err
}
