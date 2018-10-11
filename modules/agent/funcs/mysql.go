package funcs

import (
	"bytes"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

const (
	TimeOut        = 30      //
	Origin         = "GAUGE" //
	ValueSplitChar = `\$`    // ValueSplitChar is the char splitting value and tag of ignore line
	TagSplitChar   = `/`     // TagSplitChar is the char splitting metric and tag of ignore line
)

// Global tag var
var (
	IsSlave    int
	IsReadOnly int
	Tag        string
)

//DataType all variables should be monitor
var DataType = map[string]string{
	"Innodb_buffer_pool_reads":           "COUNTER",
	"Innodb_buffer_pool_read_requests":   "COUNTER",
	"Innodb_buffer_pool_write_requests":  "COUNTER",
	"Innodb_compress_time":               "COUNTER",
	"Innodb_data_fsyncs":                 "COUNTER",
	"Innodb_data_read":                   "COUNTER",
	"Innodb_data_reads":                  "COUNTER",
	"Innodb_data_writes":                 "COUNTER",
	"Innodb_data_written":                "COUNTER",
	"Innodb_last_checkpoint_at":          "COUNTER",
	"Innodb_log_flushed_up_to":           "COUNTER",
	"Innodb_log_sequence_number":         "COUNTER",
	"Innodb_mutex_os_waits":              "COUNTER",
	"Innodb_mutex_spin_rounds":           "COUNTER",
	"Innodb_mutex_spin_waits":            "COUNTER",
	"Innodb_pages_flushed_up_to":         "COUNTER",
	"Innodb_rows_deleted":                "COUNTER",
	"Innodb_rows_inserted":               "COUNTER",
	"Innodb_rows_locked":                 "COUNTER",
	"Innodb_rows_modified":               "COUNTER",
	"Innodb_rows_read":                   "COUNTER",
	"Innodb_rows_updated":                "COUNTER",
	"Innodb_row_lock_time":               "COUNTER",
	"Innodb_row_lock_waits":              "COUNTER",
	"Innodb_uncompress_time":             "COUNTER",
	"Binlog_event_count":                 "COUNTER",
	"Binlog_number":                      "COUNTER",
	"Slave_count":                        "COUNTER",
	"Com_admin_commands":                 "COUNTER",
	"Com_assign_to_keycache":             "COUNTER",
	"Com_alter_db":                       "COUNTER",
	"Com_alter_db_upgrade":               "COUNTER",
	"Com_alter_event":                    "COUNTER",
	"Com_alter_function":                 "COUNTER",
	"Com_alter_procedure":                "COUNTER",
	"Com_alter_server":                   "COUNTER",
	"Com_alter_table":                    "COUNTER",
	"Com_alter_tablespace":               "COUNTER",
	"Com_analyze":                        "COUNTER",
	"Com_begin":                          "COUNTER",
	"Com_binlog":                         "COUNTER",
	"Com_call_procedure":                 "COUNTER",
	"Com_change_db":                      "COUNTER",
	"Com_change_master":                  "COUNTER",
	"Com_check":                          "COUNTER",
	"Com_checksum":                       "COUNTER",
	"Com_commit":                         "COUNTER",
	"Com_create_db":                      "COUNTER",
	"Com_create_event":                   "COUNTER",
	"Com_create_function":                "COUNTER",
	"Com_create_index":                   "COUNTER",
	"Com_create_procedure":               "COUNTER",
	"Com_create_server":                  "COUNTER",
	"Com_create_table":                   "COUNTER",
	"Com_create_trigger":                 "COUNTER",
	"Com_create_udf":                     "COUNTER",
	"Com_create_user":                    "COUNTER",
	"Com_create_view":                    "COUNTER",
	"Com_dealloc_sql":                    "COUNTER",
	"Com_delete":                         "COUNTER",
	"Com_delete_multi":                   "COUNTER",
	"Com_do":                             "COUNTER",
	"Com_drop_db":                        "COUNTER",
	"Com_drop_event":                     "COUNTER",
	"Com_drop_function":                  "COUNTER",
	"Com_drop_index":                     "COUNTER",
	"Com_drop_procedure":                 "COUNTER",
	"Com_drop_server":                    "COUNTER",
	"Com_drop_table":                     "COUNTER",
	"Com_drop_trigger":                   "COUNTER",
	"Com_drop_user":                      "COUNTER",
	"Com_drop_view":                      "COUNTER",
	"Com_empty_query":                    "COUNTER",
	"Com_execute_sql":                    "COUNTER",
	"Com_flush":                          "COUNTER",
	"Com_grant":                          "COUNTER",
	"Com_ha_close":                       "COUNTER",
	"Com_ha_open":                        "COUNTER",
	"Com_ha_read":                        "COUNTER",
	"Com_help":                           "COUNTER",
	"Com_insert":                         "COUNTER",
	"Com_insert_select":                  "COUNTER",
	"Com_install_plugin":                 "COUNTER",
	"Com_kill":                           "COUNTER",
	"Com_load":                           "COUNTER",
	"Com_lock_tables":                    "COUNTER",
	"Com_optimize":                       "COUNTER",
	"Com_preload_keys":                   "COUNTER",
	"Com_prepare_sql":                    "COUNTER",
	"Com_purge":                          "COUNTER",
	"Com_purge_before_date":              "COUNTER",
	"Com_release_savepoint":              "COUNTER",
	"Com_rename_table":                   "COUNTER",
	"Com_rename_user":                    "COUNTER",
	"Com_repair":                         "COUNTER",
	"Com_replace":                        "COUNTER",
	"Com_replace_select":                 "COUNTER",
	"Com_reset":                          "COUNTER",
	"Com_resignal":                       "COUNTER",
	"Com_revoke":                         "COUNTER",
	"Com_revoke_all":                     "COUNTER",
	"Com_rollback":                       "COUNTER",
	"Com_rollback_to_savepoint":          "COUNTER",
	"Com_savepoint":                      "COUNTER",
	"Com_select":                         "COUNTER",
	"Com_set_option":                     "COUNTER",
	"Com_signal":                         "COUNTER",
	"Com_show_authors":                   "COUNTER",
	"Com_show_binlog_events":             "COUNTER",
	"Com_show_binlogs":                   "COUNTER",
	"Com_show_charsets":                  "COUNTER",
	"Com_show_collations":                "COUNTER",
	"Com_show_contributors":              "COUNTER",
	"Com_show_create_db":                 "COUNTER",
	"Com_show_create_event":              "COUNTER",
	"Com_show_create_func":               "COUNTER",
	"Com_show_create_proc":               "COUNTER",
	"Com_show_create_table":              "COUNTER",
	"Com_show_create_trigger":            "COUNTER",
	"Com_show_databases":                 "COUNTER",
	"Com_show_engine_logs":               "COUNTER",
	"Com_show_engine_mutex":              "COUNTER",
	"Com_show_engine_status":             "COUNTER",
	"Com_show_events":                    "COUNTER",
	"Com_show_errors":                    "COUNTER",
	"Com_show_fields":                    "COUNTER",
	"Com_show_function_status":           "COUNTER",
	"Com_show_grants":                    "COUNTER",
	"Com_show_keys":                      "COUNTER",
	"Com_show_master_status":             "COUNTER",
	"Com_show_open_tables":               "COUNTER",
	"Com_show_plugins":                   "COUNTER",
	"Com_show_privileges":                "COUNTER",
	"Com_show_procedure_status":          "COUNTER",
	"Com_show_processlist":               "COUNTER",
	"Com_show_profile":                   "COUNTER",
	"Com_show_profiles":                  "COUNTER",
	"Com_show_relaylog_events":           "COUNTER",
	"Com_show_slave_hosts":               "COUNTER",
	"Com_show_slave_status":              "COUNTER",
	"Com_show_status":                    "COUNTER",
	"Com_show_storage_engines":           "COUNTER",
	"Com_show_table_status":              "COUNTER",
	"Com_show_tables":                    "COUNTER",
	"Com_show_triggers":                  "COUNTER",
	"Com_show_variables":                 "COUNTER",
	"Com_show_warnings":                  "COUNTER",
	"Com_slave_start":                    "COUNTER",
	"Com_slave_stop":                     "COUNTER",
	"Com_stmt_close":                     "COUNTER",
	"Com_stmt_execute":                   "COUNTER",
	"Com_stmt_fetch":                     "COUNTER",
	"Com_stmt_prepare":                   "COUNTER",
	"Com_stmt_reprepare":                 "COUNTER",
	"Com_stmt_reset":                     "COUNTER",
	"Com_stmt_send_long_data":            "COUNTER",
	"Com_truncate":                       "COUNTER",
	"Com_uninstall_plugin":               "COUNTER",
	"Com_unlock_tables":                  "COUNTER",
	"Com_update":                         "COUNTER",
	"Com_update_multi":                   "COUNTER",
	"Com_xa_commit":                      "COUNTER",
	"Com_xa_end":                         "COUNTER",
	"Com_xa_prepare":                     "COUNTER",
	"Com_xa_recover":                     "COUNTER",
	"Com_xa_rollback":                    "COUNTER",
	"Com_xa_start":                       "COUNTER",
	"Com_alter_user":                     "COUNTER",
	"Com_get_diagnostics":                "COUNTER",
	"Com_lock_tables_for_backup":         "COUNTER",
	"Com_lock_binlog_for_backup":         "COUNTER",
	"Com_purge_archived":                 "COUNTER",
	"Com_purge_archived_before_date":     "COUNTER",
	"Com_show_client_statistics":         "COUNTER",
	"Com_show_function_code":             "COUNTER",
	"Com_show_index_statistics":          "COUNTER",
	"Com_show_procedure_code":            "COUNTER",
	"Com_show_slave_status_nolock":       "COUNTER",
	"Com_show_table_statistics":          "COUNTER",
	"Com_show_thread_statistics":         "COUNTER",
	"Com_show_user_statistics":           "COUNTER",
	"Com_unlock_binlog":                  "COUNTER",
	"Aborted_clients":                    "COUNTER",
	"Aborted_connects":                   "COUNTER",
	"Access_denied_errors":               "COUNTER",
	"Binlog_bytes_written":               "COUNTER",
	"Binlog_cache_disk_use":              "COUNTER",
	"Binlog_cache_use":                   "COUNTER",
	"Binlog_stmt_cache_disk_use":         "COUNTER",
	"Binlog_stmt_cache_use":              "COUNTER",
	"Bytes_received":                     "COUNTER",
	"Bytes_sent":                         "COUNTER",
	"Connections":                        "COUNTER",
	"Created_tmp_disk_tables":            "COUNTER",
	"Created_tmp_files":                  "COUNTER",
	"Created_tmp_tables":                 "COUNTER",
	"Handler_delete":                     "COUNTER",
	"Handler_read_first":                 "COUNTER",
	"Handler_read_key":                   "COUNTER",
	"Handler_read_last":                  "COUNTER",
	"Handler_read_next":                  "COUNTER",
	"Handler_read_prev":                  "COUNTER",
	"Handler_read_rnd":                   "COUNTER",
	"Handler_read_rnd_next":              "COUNTER",
	"Handler_update":                     "COUNTER",
	"Handler_write":                      "COUNTER",
	"Opened_files":                       "COUNTER",
	"Opened_tables":                      "COUNTER",
	"Opened_table_definitions":           "COUNTER",
	"Qcache_hits":                        "COUNTER",
	"Qcache_inserts":                     "COUNTER",
	"Qcache_lowmem_prunes":               "COUNTER",
	"Qcache_not_cached":                  "COUNTER",
	"Queries":                            "COUNTER",
	"Questions":                          "COUNTER",
	"Select_full_join":                   "COUNTER",
	"Select_full_range_join":             "COUNTER",
	"Select_range_check":                 "COUNTER",
	"Select_scan":                        "COUNTER",
	"Slow_queries":                       "COUNTER",
	"Sort_merge_passes":                  "COUNTER",
	"Sort_range":                         "COUNTER",
	"Sort_rows":                          "COUNTER",
	"Sort_scan":                          "COUNTER",
	"Table_locks_immediate":              "COUNTER",
	"Table_locks_waited":                 "COUNTER",
	"Threads_created":                    "COUNTER",
	"Rpl_semi_sync_master_net_wait_time": "COUNTER",
	"Rpl_semi_sync_master_net_waits":     "COUNTER",
	"Rpl_semi_sync_master_no_times":      "COUNTER",
	"Rpl_semi_sync_master_no_tx":         "COUNTER",
	"Rpl_semi_sync_master_yes_tx":        "COUNTER",
	"Rpl_semi_sync_master_tx_wait_time":  "COUNTER",
	"Rpl_semi_sync_master_tx_waits":      "COUNTER",
}

// SlaveStatus not all slave status send to falcon-agent, this is a filter
var SlaveStatus = []string{
	"Exec_Master_Log_Pos",
	"Read_Master_log_Pos",
	"Relay_Log_Pos",
	"Seconds_Behind_Master",
	"Slave_IO_Running",
	"Slave_SQL_Running",
}

func dataType(k string) string {
	if v, ok := DataType[k]; ok {
		return v
	}
	return Origin
}

// NewMetric is the constructor of metric
func NewMetric(name string) *model.MetricValue {
	hostname, _ := g.Hostname()
	return &model.MetricValue{
		Metric:    "mysql." + name,
		Endpoint:  hostname,
		Type:      dataType(name),
		Tags:      Tag,
		Timestamp: time.Now().Unix(),
		Step:      60,
	}
}

// GetTag can get the tag to output
func GetTag() string {
	return fmt.Sprintf("host=%s,port=%d,read_only=%d",
		g.Config().Collector.MySQL.Host,
		g.Config().Collector.MySQL.Port,
		IsReadOnly)
}

// MySQLAlive checks if mysql can response
func MySQLAlive() []*model.MetricValue {
	if g.Config().Collector.MySQL != nil {
		data := NewMetric("alive")

		data.Value = 0
		if g.Config().Collector.MySQL.Enabled {
			data.Value = 1
		}
		return []*model.MetricValue{data}
	}
	return nil
}

// GetIsReadOnly get read_only variable of mysql
func GetIsReadOnly(db *sql.DB) (int, error) {
	readOnly := 0
	err := db.QueryRow("select @@read_only").Scan(&readOnly)
	if err != nil {
		return -1, err
	}
	return readOnly, nil
}

func MySQLMetrics() (L []*model.MetricValue) {
	if g.Config().Collector.MySQL == nil {
		return nil
	}
	if !g.Config().Collector.MySQL.Enabled {
		return nil
	}

	Tag = GetTag()
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql",
			g.Config().Collector.MySQL.User,
			g.Config().Collector.MySQL.Passowrd,
			g.Config().Collector.MySQL.Host,
			g.Config().Collector.MySQL.Port))

	if err != nil {
		fmt.Printf("Connect to mysql error: %s\n", err.Error())
		return nil
	}
	defer db.Close()

	defer func() {
		L = append(L, MySQLAlive()...)
	}()

	// Get GLOBAL variables
	IsReadOnly, err = GetIsReadOnly(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	// Get slave status and set IsSlave global var
	slaveState, err := ShowSlaveStatus(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	globalStatus, err := ShowGlobalStatus(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	L = append(L, globalStatus...)

	globalVars, err := ShowGlobalVariables(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	L = append(L, globalVars...)

	innodbState, err := ShowInnodbStatus(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	L = append(L, innodbState...)
	L = append(L, slaveState...)

	binaryLogStatus, err := ShowBinaryLogs(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	L = append(L, binaryLogStatus...)

	return
}

func Query(db *sql.DB, sql string) ([]map[string]interface{}, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))

	for i := range vals {
		scans[i] = &vals[i]
	}

	var results []map[string]interface{}

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	return results, nil
}

// ShowGlobalStatus execute mysql query `SHOW GLOBAL STATUS`
func ShowGlobalStatus(db *sql.DB) ([]*model.MetricValue, error) {
	return parseMySQLStatus(db, "SHOW /*!50001 GLOBAL */ STATUS")
}

// ShowGlobalVariables execute mysql query `SHOW GLOBAL VARIABLES`
func ShowGlobalVariables(db *sql.DB) ([]*model.MetricValue, error) {
	return parseMySQLStatus(db, "SHOW /*!50001 GLOBAL */ VARIABLES")
}

// ShowInnodbStatus execute mysql query `SHOW SHOW /*!50000 ENGINE*/ INNODB STATUS`
func ShowInnodbStatus(db *sql.DB) ([]*model.MetricValue, error) {
	type row struct {
		Type   string
		Name   string
		Status string
	}
	var r row
	err := db.QueryRow("SHOW /*!50000 ENGINE*/ INNODB STATUS").Scan(&r.Type, &r.Name, &r.Status)
	if err != nil {
		return nil, err
	}
	rows := strings.Split(r.Status, "\n")
	return parseInnodbStatus(rows)
}

// ShowBinaryLogs execute mysql query `SHOW BINARY LOGS`
func ShowBinaryLogs(db *sql.DB) ([]*model.MetricValue, error) {
	sum := 0

	binlogFileCounts := NewMetric("binlog_file_counts")
	binlogFileSize := NewMetric("binlog_file_size")

	results, err := Query(db, "SHOW BINARY LOGS")
	if err != nil {
		return []*model.MetricValue{binlogFileCounts, binlogFileSize}, err
	}

	for _, v := range results {
		size, err := strconv.Atoi(strings.TrimSpace(v["File_size"].(string)))
		if err == nil {
			sum += size
		}
	}

	binlogFileCounts.Value = len(results)
	binlogFileSize.Value = sum
	return []*model.MetricValue{binlogFileCounts, binlogFileSize}, err
}

// ShowSlaveStatus get all slave status of mysql serves
func ShowSlaveStatus(db *sql.DB) ([]*model.MetricValue, error) {
	// Check IsSlave
	results, err := Query(db, "SHOW SLAVE STATUS")
	if err != nil {
		IsSlave = -1
		return nil, err
	}

	if results != nil {
		IsSlave = 1
	} else {
		IsSlave = 0
	}

	isSlaveMetric := NewMetric("is_slave")
	isSlaveMetric.Value = IsSlave

	// be master
	if IsSlave == 0 {
		// Master_is_readonly VS master_is_read_only for version compatible, ugly
		masterReadOnly, err := ShowOtherMetric(db, "Master_is_readonly")
		if err != nil {
			return nil, err
		}
		masterReadOnly2, err := ShowOtherMetric(db, "master_is_read_only")
		if err != nil {
			return nil, err
		}
		innodbStatsOnMetadata, err := ShowOtherMetric(db, "innodb_stats_on_metadata")
		if err != nil {
			return nil, err
		}
		return []*model.MetricValue{isSlaveMetric, masterReadOnly, masterReadOnly2, innodbStatsOnMetadata}, nil
	}

	// be slave
	ioDelay, err := ShowOtherMetric(db, "io_thread_delay")
	if err != nil {
		return nil, err
	}
	slaveReadOnly, err := ShowOtherMetric(db, "slave_is_read_only")
	if err != nil {
		return nil, err
	}
	heartbeat, err := ShowOtherMetric(db, "Heartbeats_Behind_Master")
	if err != nil {
		// mysql.heartbeat table not necessary exist if you don't care about heartbeat
		// bypass heartbeat table not exist error
		err = nil
	}
	data := make([]*model.MetricValue, len(SlaveStatus))
	for i, s := range SlaveStatus {
		data[i] = NewMetric(s)
		switch s {
		case "Slave_SQL_Running", "Slave_IO_Running":
			data[i].Value = 0
			v := results[0][s]
			if v == "Yes" {
				data[i].Value = 1
			}
		default:
			if v, ok := results[0][s]; ok {
				data[i].Value = v
			} else {
				data[i].Value = -1
			}
		}
	}
	return append(data, []*model.MetricValue{isSlaveMetric, ioDelay, slaveReadOnly, heartbeat}...), nil
}

func GetLastNum(str string, split string) int {
	parts := strings.Split(str, split)
	if len(parts) < 2 {
		return -1
	}
	ans, err := strconv.ParseInt(parts[1], 10, 60)
	if err != nil {
		return -2
	}
	return int(ans)
}

// ShowOtherMetric all other metric will add in this func
func ShowOtherMetric(db *sql.DB, metric string) (*model.MetricValue, error) {
	var err error
	newMetaData := NewMetric(metric)
	switch metric {
	case "master_is_read_only", "slave_is_read_only", "Master_is_readonly":
		newMetaData.Value = IsReadOnly
	case "innodb_stats_on_metadata":
		v := 0
		err = db.QueryRow("SELECT /*!50504 @@GLOBAL.innodb_stats_on_metadata */;").Scan(&v)
		newMetaData.Value = v
	case "io_thread_delay":
		results, err := Query(db, "SHOW SLAVE STATUS")
		if err != nil {
			newMetaData.Value = -1
			return newMetaData, err
		}

		if bytes.Equal([]byte(results[0]["Master_Log_File"].(string)), []byte(results[0]["Relay_Master_Log_File"].(string))) {
			newMetaData.Value = 0
		} else {
			masterLogFile := GetLastNum(results[0]["Master_Log_File"].(string), ".")
			relayMasterLogFile := GetLastNum(results[0]["Relay_Master_Log_File"].(string), ".")
			newMetaData.Value = masterLogFile - relayMasterLogFile
			if masterLogFile < 0 || relayMasterLogFile < 0 {
				newMetaData.Value = -1
			}
		}
	case "Heartbeats_Behind_Master":
		var ts string
		err = db.QueryRow("select ts from mysql.heartbeat limit 1").Scan(&ts)
		// when row is empty, err is nil either
		if err != nil && err != sql.ErrNoRows {
			newMetaData.Value = -1
		} else {
			localTimezone, _ := time.LoadLocation("Local")
			heartbeatTimeStr := ts
			b := strings.Replace(heartbeatTimeStr, "T", " ", 1)
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", strings.Split(b, ".")[0], localTimezone)
			heartbeatTimestamp := t.Unix()
			currentTimestamp := time.Now().Unix()
			newMetaData.Value = currentTimestamp - heartbeatTimestamp
		}
	}

	return newMetaData, err
}

func parseMySQLStatus(db *sql.DB, sql string) ([]*model.MetricValue, error) {
	results, err := Query(db, sql)
	if err != nil {
		return nil, err
	}

	data := make([]*model.MetricValue, len(results))
	i := 0
	for _, row := range results {
		if key, ok := row["Variable_name"]; ok {
			if value, ok := row["Value"]; ok {
				d, err := strconv.Atoi(strings.TrimSpace(value.(string)))
				if err == nil {
					data[i] = NewMetric(key.(string))
					data[i].Value = d
					i++
				}
			}
		}
	}
	return data[:i], nil
}

func parseInnodbSection(
	row string, section string,
	pdata *[]*model.MetricValue, longTranTime *int) error {
	switch section {
	case "TRANSACTIONS":
		if strings.Contains(row, "ACTIVE") {
			tmpLongTransactionTime, err := strconv.Atoi(
				strings.Split(
					strings.Split(
						row, "ACTIVE ")[1],
					" sec")[0])
			if err != nil {
				return err
			}
			if tmpLongTransactionTime > *longTranTime {
				*longTranTime = tmpLongTransactionTime
			}
		}
		if strings.Contains(row, "History list length") {
			hisListLengthStr := strings.Split(row, "length ")[1]
			hisListLength, _ := strconv.Atoi(hisListLengthStr)
			HistoryListLength := NewMetric("History_list_length")
			HistoryListLength.Value = hisListLength
			*pdata = append(*pdata, HistoryListLength)
		}
	case "SEMAPHORES":
		matches := regexp.MustCompile(`^Mutex spin waits\s+(\d+),\s+rounds\s+(\d+),\s+OS waits\s+(\d+)`).FindStringSubmatch(row)
		if len(matches) == 4 {
			spinWaits, _ := strconv.Atoi(matches[1])
			innodbMutexSpinWaits := NewMetric("Innodb_mutex_spin_waits")
			innodbMutexSpinWaits.Value = spinWaits
			*pdata = append(*pdata, innodbMutexSpinWaits)

			spinRounds, _ := strconv.Atoi(matches[2])
			InnodbMutexSpinRounds := NewMetric("Innodb_mutex_spin_rounds")
			InnodbMutexSpinRounds.Value = spinRounds
			*pdata = append(*pdata, InnodbMutexSpinRounds)

			osWaits, _ := strconv.Atoi(matches[3])
			InnodbMutexOsWaits := NewMetric("Innodb_mutex_os_waits")
			InnodbMutexOsWaits.Value = osWaits
			*pdata = append(*pdata, InnodbMutexOsWaits)
		}
	}
	return nil
}

func parseInnodbStatus(rows []string) ([]*model.MetricValue, error) {
	var section string
	longTranTime := 0
	var err error
	var data []*model.MetricValue
	for _, row := range rows {
		switch row {
		case "BACKGROUND THREAD":
			section = row
			continue
		case "DEAD LOCK ERRORS":
			section = row
			continue
		case "LATEST DETECTED DEADLOCK":
			section = row
			continue
		case "FOREIGN KEY CONSTRAINT ERRORS", "LATEST FOREIGN KEY ERROR":
			section = row
			continue
		case "SEMAPHORES":
			section = row
			continue
		case "TRANSACTIONS":
			section = row
			continue
		case "FILE I/O":
			section = row
			continue
		case "INSERT BUFFER AND ADAPTIVE HASH INDEX":
			section = row
			continue
		case "LOG":
			section = row
			continue
		case "BUFFER POOL AND MEMORY":
			section = row
			continue
		case "ROW OPERATIONS":
			section = row
			continue
		}
		err = parseInnodbSection(row, section, &data, &longTranTime)
		if err != nil {
			return nil, err
		}
	}
	longTranMetric := NewMetric("longest_transaction")
	longTranMetric.Value = longTranTime
	data = append(data, longTranMetric)
	return data, nil
}
