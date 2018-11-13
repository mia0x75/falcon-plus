package db

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

func QueryAgentsInfo() (map[string]*cmodel.AgentUpdateInfo, error) {
	m := make(map[string]*cmodel.AgentUpdateInfo)
	sql := "select hostname, ip, agent_version, plugin_version, update_at from host"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Errorf("[E] %v", err)
		return m, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			hostname       string
			ip             string
			agent_version  string
			plugin_version string
			update_at      time.Time
		)
		err = rows.Scan(&hostname, &ip, &agent_version, &plugin_version, &update_at)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}
		m[hostname] = &cmodel.AgentUpdateInfo{
			LastUpdate: update_at.UnixNano(),
			ReportRequest: &cmodel.AgentReportRequest{
				Hostname:      hostname,
				IP:            ip,
				AgentVersion:  agent_version,
				PluginVersion: plugin_version,
			},
		}
	}
	return m, nil
}

func UpdateAgent(agentInfo *cmodel.AgentUpdateInfo) {
	var (
		hostname       string
		ip             string
		agent_version  string
		plugin_version string
	)

	sql := fmt.Sprintf("select hostname, ip, agent_version, plugin_version from host where hostname = '%s'",
		agentInfo.ReportRequest.Hostname,
	)

	err := DB.QueryRow(sql).Scan(&hostname, &ip, &agent_version, &plugin_version)
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	if agentInfo.ReportRequest.Hostname == hostname &&
		agentInfo.ReportRequest.IP == ip &&
		agentInfo.ReportRequest.AgentVersion == agent_version &&
		agentInfo.ReportRequest.PluginVersion == plugin_version {
		return
	}

	sql = ""
	if g.Config().Hosts == "" {
		if hostname == "" && ip == "" && agent_version == "" && plugin_version == "" {
			sql = fmt.Sprintf(
				"insert into host(hostname, ip, agent_version, plugin_version) values ('%s', '%s', '%s', '%s')",
				agentInfo.ReportRequest.Hostname,
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
			)
		} else {
			sql = fmt.Sprintf(
				"update host set ip='%s', agent_version='%s', plugin_version='%s' where hostname='%s'",
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
				agentInfo.ReportRequest.Hostname,
			)
		}
	} else {
		// sync, just update
		sql = fmt.Sprintf(
			"update host set ip='%s', agent_version='%s', plugin_version='%s' where hostname='%s'",
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.Hostname,
		)
	}

	_, err = DB.Exec(sql)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", sql, err)
	}
}
