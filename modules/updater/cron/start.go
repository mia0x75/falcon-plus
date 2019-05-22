package cron

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"

	"github.com/open-falcon/falcon-plus/modules/updater/model"
	"github.com/open-falcon/falcon-plus/modules/updater/utils"
)

// StartDesiredAgent TODO:
func StartDesiredAgent(da *model.DesiredAgent) {
	if err := InsureDesiredAgentDirExists(da); err != nil {
		return
	}

	if err := InsureNewVersionFiles(da); err != nil {
		return
	}

	if err := Untar(da); err != nil {
		return
	}

	if err := StopAgentOf(da.Name, da.Version); err != nil {
		return
	}

	if err := ControlStartIn(da.AgentVersionDir); err != nil {
		return
	}

	file.WriteString(path.Join(da.AgentDir, ".version"), da.Version)
}

// Untar TODO:
func Untar(da *model.DesiredAgent) error {
	cmd := exec.Command("tar", "zxf", da.TarballFilename)
	cmd.Dir = da.AgentVersionDir
	err := cmd.Run()
	if err != nil {
		log.Errorf("[E] tar zxf %s fail: %v", da.TarballFilename, err)
		return err
	}

	return nil
}

// ControlStartIn TODO:
func ControlStartIn(workdir string) error {
	out, err := ControlStatus(workdir)
	if err == nil && strings.Contains(out, "started") {
		return nil
	}

	_, err = ControlStart(workdir)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 3)

	out, err = ControlStatus(workdir)
	if err == nil && strings.Contains(out, "started") {
		return nil
	}

	return err
}

// InsureNewVersionFiles TODO:
func InsureNewVersionFiles(da *model.DesiredAgent) error {
	if FilesReady(da) {
		return nil
	}
	content, err := ioutil.ReadFile("./password")
	password := strings.Trim(string(content), "\n")
	if err != nil {
		panic(err)
	}
	downloadTarballCmd := exec.Command("wget", "--no-check-certificate", "--auth-no-challenge", "--user=owl", "--password="+password, da.TarballUrl, "-O", da.TarballFilename)
	downloadTarballCmd.Dir = da.AgentVersionDir
	err = downloadTarballCmd.Run()
	if err != nil {
		log.Errorf(
			"[ERROR] wget -q --no-check-certificate --auth-no-challenge --user=owl --password=%s %s -O %s fail: %v",
			password,
			da.TarballUrl,
			da.TarballFilename,
			err,
		)
		return err
	}

	downloadMd5Cmd := exec.Command("wget", "--no-check-certificate", "--auth-no-challenge", "--user=owl", "--password="+password, da.Md5Url, "-O", da.Md5Filename)
	downloadMd5Cmd.Dir = da.AgentVersionDir
	err = downloadMd5Cmd.Run()
	if err != nil {
		log.Errorf(
			"[ERROR] wget -q --no-check-certificate --auth-no-challenge --user=owl --password=%s %s -O %s fail: %v",
			password,
			da.Md5Url,
			da.Md5Filename,
			err,
		)
		return err
	}

	if utils.Md5sumCheck(da.AgentVersionDir, da.Md5Filename) {
		return nil
	}
	return fmt.Errorf("md5sum -c fail")
}

// FilesReady TODO:
func FilesReady(da *model.DesiredAgent) bool {
	if !file.IsExist(da.Md5Filepath) {
		return false
	}

	if !file.IsExist(da.TarballFilepath) {
		return false
	}

	if !file.IsExist(da.ControlFilepath) {
		return false
	}

	return utils.Md5sumCheck(da.AgentVersionDir, da.Md5Filename)
}

// InsureDesiredAgentDirExists TODO:
func InsureDesiredAgentDirExists(da *model.DesiredAgent) error {
	err := file.InsureDir(da.AgentDir)
	if err != nil {
		log.Errorf("[E] insure dir %s fail: %v", da.AgentDir, err)
		return err
	}

	err = file.InsureDir(da.AgentVersionDir)
	if err != nil {
		log.Errorf("[E] insure dir %s fail: %v", da.AgentVersionDir, err)
	}
	return err
}
