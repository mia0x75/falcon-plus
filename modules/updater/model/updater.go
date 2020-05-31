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

func (m *DesiredAgent) String() string {
	return fmt.Sprintf(
		"<Name:%s, Version:%s, RunUser:%s, WorkDir: %s, Md5:%s, Cmd:%s>",
		m.Name,
		m.Version,
		m.RunUser,
		m.WorkDir,
		m.Md5,
		m.Cmd,
	)
}

func (m *DesiredAgent) FillAttrs(workdir string) {
	m.AgentDir = path.Join(workdir, m.Name)
	m.AgentVersionDir = path.Join(m.AgentDir, m.Version)
	m.TarballFilename = fmt.Sprintf("%s-%s.tar.gz", m.Name, m.Version)
	m.Md5Filename = fmt.Sprintf("%s.md5", m.TarballFilename)
	m.TarballFilepath = path.Join(m.AgentVersionDir, m.TarballFilename)
	m.Md5Filepath = path.Join(m.AgentVersionDir, m.Md5Filename)
	m.ControlFilepath = path.Join(m.AgentVersionDir, "control")

	if m.Md5 == "" {
		m.Md5 = m.Tarball
	}

	m.TarballUrl = fmt.Sprintf("%s/%s", m.Tarball, m.TarballFilename)
	m.Md5Url = fmt.Sprintf("%s/%s", m.Md5, m.Md5Filename)
}

type RealAgent struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	RunUser   string `json:"runUser"`
	WorkDir   string `json:"workDir"`
}

func (m *RealAgent) String() string {
	return fmt.Sprintf(
		"<Name:%s, Version:%s, Status:%s, Timestamp:%v, RunUser: %s, WorkDir: %s>",
		m.Name,
		m.Version,
		m.Status,
		m.Timestamp,
		m.RunUser,
		m.WorkDir,
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
