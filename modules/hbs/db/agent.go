package db

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

// QueryAgentsInfo TODO:
func QueryAgentsInfo() (map[string]*cmodel.AgentUpdateInfo, error) {
	m := make(map[string]*cmodel.AgentUpdateInfo)
	sql := "select hostname, ip, agent_version, plugin_version, update_at from host"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", sql, err)
		return m, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			hostname      string
			ip            string
			agentVersion  string
			pluginVersion string
			updateAt      time.Time
		)
		err = rows.Scan(&hostname, &ip, &agentVersion, &pluginVersion, &updateAt)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}
		m[hostname] = &cmodel.AgentUpdateInfo{
			LastUpdate: updateAt.UnixNano(),
			ReportRequest: &cmodel.AgentReportRequest{
				Hostname:      hostname,
				IP:            ip,
				AgentVersion:  agentVersion,
				PluginVersion: pluginVersion,
			},
		}
	}
	return m, nil
}

// UpdateAgent TODO:
func UpdateAgent(agentInfo *cmodel.AgentUpdateInfo) {
	var (
		hostname      string
		ip            string
		agentVersion  string
		pluginVersion string
	)

	q := fmt.Sprintf("select hostname, ip, agent_version, plugin_version from host where hostname = '%s'",
		agentInfo.ReportRequest.Hostname,
	)

	err := DB.QueryRow(q).Scan(&hostname, &ip, &agentVersion, &pluginVersion)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("[E] exec %s fail: %v", q, err)
			return
		}
	}

	if agentInfo.ReportRequest.Hostname == hostname &&
		agentInfo.ReportRequest.IP == ip &&
		agentInfo.ReportRequest.AgentVersion == agentVersion &&
		agentInfo.ReportRequest.PluginVersion == pluginVersion {
		return
	}

	q = ""
	if g.Config().Hosts == "" {
		if hostname == "" && ip == "" && agentVersion == "" && pluginVersion == "" {
			q = fmt.Sprintf(
				"insert into host(hostname, ip, agent_version, plugin_version) values ('%s', '%s', '%s', '%s')",
				agentInfo.ReportRequest.Hostname,
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
			)
		} else {
			q = fmt.Sprintf(
				"update host set ip = '%s', agent_version = '%s', plugin_version = '%s' where hostname = '%s'",
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
				agentInfo.ReportRequest.Hostname,
			)
		}
	} else {
		// sync, just update
		q = fmt.Sprintf(
			"update host set ip = '%s', agent_version = '%s', plugin_version = '%s' where hostname = '%s'",
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.Hostname,
		)
	}

	_, err = DB.Exec(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
	}
}
