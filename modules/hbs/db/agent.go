package db

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

// QueryAgentsInfo TODO:
func QueryAgentsInfo() (map[string]*cm.AgentUpdateInfo, error) {
	m := make(map[string]*cm.AgentUpdateInfo)
	q := "SELECT hostname, ip, agent_version, plugin_version, update_at FROM hosts"
	rows, err := DB.Query(q)
	if err != nil {
		log.Errorf("[E] exec %s fail: %v", q, err)
		return m, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			hostname      string
			ip            string
			agentVersion  string
			pluginVersion string
			updateAt      sql.NullInt64
		)
		err = rows.Scan(&hostname, &ip, &agentVersion, &pluginVersion, &updateAt)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}
		if !updateAt.Valid {
			updateAt.Int64 = 0
		}
		m[hostname] = &cm.AgentUpdateInfo{
			LastUpdate: updateAt.Int64,
			ReportRequest: &cm.AgentReportRequest{
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
func UpdateAgent(agentInfo *cm.AgentUpdateInfo) {
	var (
		hostname      string
		ip            string
		agentVersion  string
		pluginVersion string
	)

	q := fmt.Sprintf("SELECT hostname, ip, agent_version, plugin_version FROM hosts WHERE hostname = '%s'",
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
				"INSERT INTO hosts(hostname, ip, agent_version, plugin_version) VALUES ('%s', '%s', '%s', '%s')",
				agentInfo.ReportRequest.Hostname,
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
			)
		} else {
			q = fmt.Sprintf(
				"UPDATE hosts SET ip = '%s', agent_version = '%s', plugin_version = '%s' WHERE hostname = '%s'",
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
				agentInfo.ReportRequest.Hostname,
			)
		}
	} else {
		// sync, just update
		q = fmt.Sprintf(
			"UPDATE hosts SET ip = '%s', agent_version = '%s', plugin_version = '%s' WHERE hostname = '%s'",
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
