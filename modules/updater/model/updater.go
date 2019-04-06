package model

import (
	"fmt"
	"path"
)

type DesiredAgent struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Tarball         string `json:"tarball"`
	Md5             string `json:"md5"`
	Cmd             string `json:"cmd"`
	RunUser         string `json:"runUser"`
	WorkDir         string `json:"workDir"`
	ConfigFileName  string `json:"configFileName"`
	ConfigRemoteUrl string `json:"configRemoteUrl"`
	AgentDir        string `json:"-"`
	AgentVersionDir string `json:"-"`
	TarballFilename string `json:"-"`
	Md5Filename     string `json:"-"`
	TarballFilepath string `json:"-"`
	Md5Filepath     string `json:"-"`
	ControlFilepath string `json:"-"`
	TarballUrl      string `json:"-"`
	Md5Url          string `json:"-"`
}

func (this *DesiredAgent) String() string {
	return fmt.Sprintf(
		"<Name:%s, Version:%s, RunUser:%s, WorkDir: %s, Md5:%s, Cmd:%s>",
		this.Name,
		this.Version,
		this.RunUser,
		this.WorkDir,
		this.Md5,
		this.Cmd,
	)
}

func (this *DesiredAgent) FillAttrs(workdir string) {
	this.AgentDir = path.Join(workdir, this.Name)
	this.AgentVersionDir = path.Join(this.AgentDir, this.Version)
	this.TarballFilename = fmt.Sprintf("%s-%s.tar.gz", this.Name, this.Version)
	this.Md5Filename = fmt.Sprintf("%s.md5", this.TarballFilename)
	this.TarballFilepath = path.Join(this.AgentVersionDir, this.TarballFilename)
	this.Md5Filepath = path.Join(this.AgentVersionDir, this.Md5Filename)
	this.ControlFilepath = path.Join(this.AgentVersionDir, "control")

	if this.Md5 == "" {
		this.Md5 = this.Tarball
	}

	this.TarballUrl = fmt.Sprintf("%s/%s", this.Tarball, this.TarballFilename)
	this.Md5Url = fmt.Sprintf("%s/%s", this.Md5, this.Md5Filename)
}

type RealAgent struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	RunUser   string `json:"runUser"`
	WorkDir   string `json:"workDir"`
}

func (this *RealAgent) String() string {
	return fmt.Sprintf(
		"<Name:%s, Version:%s, Status:%s, Timestamp:%v, RunUser: %s, WorkDir: %s>",
		this.Name,
		this.Version,
		this.Status,
		this.Timestamp,
		this.RunUser,
		this.WorkDir,
	)
}

type HeartbeatRequest struct {
	Hostname       string       `json:"hostname"`
	Ip             string       `json:"ip"`
	RunUser        string       `json:"runUser"`
	UpdaterVersion string       `json:"updaterVersion"`
	RealAgents     []*RealAgent `json:"realAgents"`
}

type HeartbeatResponse struct {
	ErrorMessage  string          `json:"errorMessage"`
	DesiredAgents []*DesiredAgent `json:"desiredAgents"`
}
