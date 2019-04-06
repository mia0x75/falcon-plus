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
	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
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
	"Innodb_buffer_pool_read_ahead_rnd":                            "COUNTER",
	"Innodb_buffer_pool_read_ahead":                                "COUNTER",
	"Innodb_buffer_pool_read_ahead_evicted":                        "COUNTER",
	"Innodb_buffer_pool_pages_flushed":                             "COUNTER",
	"Innodb_buffer_pool_reads":                                     "COUNTER",
	"Innodb_buffer_pool_read_requests":                             "COUNTER",
	"Innodb_buffer_pool_write_requests":                            "COUNTER",
	"Innodb_compress_time":                                         "COUNTER",
	"Innodb_data_fsyncs":                                           "COUNTER",
	"Innodb_data_pending_fsyncs":                                   "COUNTER",
	"Innodb_data_pending_reads":                                    "COUNTER",
	"Innodb_data_pending_writes":                                   "COUNTER",
	"Innodb_data_read":                                             "COUNTER",
	"Innodb_data_reads":                                            "COUNTER",
	"Innodb_data_writes":                                           "COUNTER",
	"Innodb_data_written":                                          "COUNTER",
	"Innodb_dblwr_pages_written":                                   "COUNTER",
	"Innodb_dblwr_writes":                                          "COUNTER",
	"Innodb_log_waits":                                             "COUNTER",
	"Innodb_log_write_requests":                                    "COUNTER",
	"Innodb_log_writes":                                            "COUNTER",
	"Innodb_os_log_fsyncs":                                         "COUNTER",
	"Innodb_os_log_pending_fsyncs":                                 "COUNTER",
	"Innodb_os_log_pending_writes":                                 "COUNTER",
	"Innodb_os_log_written":                                        "COUNTER",
	"Innodb_pages_created":                                         "COUNTER",
	"Innodb_pages_read":                                            "COUNTER",
	"Innodb_pages0_read":                                           "COUNTER", // MariaDB
	"Innodb_pages_written":                                         "COUNTER",
	"Innodb_row_lock_time":                                         "COUNTER",
	"Innodb_row_lock_waits":                                        "COUNTER",
	"Innodb_rows_deleted":                                          "COUNTER",
	"Innodb_rows_inserted":                                         "COUNTER",
	"Innodb_rows_read":                                             "COUNTER",
	"Innodb_rows_updated":                                          "COUNTER",
	"Innodb_system_rows_deleted":                                   "COUNTER", // MariaDB
	"Innodb_system_rows_inserted":                                  "COUNTER", // MariaDB
	"Innodb_system_rows_read":                                      "COUNTER", // MariaDB
	"Innodb_system_rows_updated":                                   "COUNTER", // MariaDB
	"Innodb_num_open_files":                                        "COUNTER",
	"Innodb_truncated_status_writes":                               "COUNTER",
	"Innodb_page_compression_saved":                                "COUNTER", // MariaDB
	"Innodb_num_index_pages_written":                               "COUNTER", // MariaDB
	"Innodb_num_non_index_pages_written":                           "COUNTER", // MariaDB
	"Innodb_num_pages_page_compressed":                             "COUNTER", // MariaDB
	"Innodb_num_page_compressed_trim_op":                           "COUNTER", // MariaDB
	"Innodb_num_pages_page_decompressed":                           "COUNTER", // MariaDB
	"Innodb_num_pages_page_compression_error":                      "COUNTER", // MariaDB
	"Innodb_num_pages_encrypted":                                   "COUNTER", // MariaDB
	"Innodb_num_pages_decrypted":                                   "COUNTER", // MariaDB
	"Innodb_defragment_compression_failures":                       "COUNTER", // MariaDB
	"Innodb_defragment_failures":                                   "COUNTER", // MariaDB
	"Innodb_defragment_count":                                      "COUNTER", // MariaDB
	"Innodb_instant_alter_column":                                  "COUNTER", // MariaDB
	"Innodb_secondary_index_triggered_cluster_reads":               "COUNTER", // MariaDB
	"Innodb_secondary_index_triggered_cluster_reads_avoided":       "COUNTER", // MariaDB
	"Innodb_encryption_rotation_pages_read_from_cache":             "COUNTER", // MariaDB
	"Innodb_encryption_rotation_pages_read_from_disk":              "COUNTER", // MariaDB
	"Innodb_encryption_rotation_pages_modified":                    "COUNTER", // MariaDB
	"Innodb_encryption_rotation_pages_flushed":                     "COUNTER", // MariaDB
	"Innodb_encryption_rotation_estimated_iops":                    "COUNTER", // MariaDB
	"Innodb_encryption_key_rotation_list_length":                   "COUNTER", // MariaDB
	"Innodb_encryption_n_merge_blocks_encrypted":                   "COUNTER", // MariaDB
	"Innodb_encryption_n_merge_blocks_decrypted":                   "COUNTER", // MariaDB
	"Innodb_encryption_n_rowlog_blocks_encrypted":                  "COUNTER", // MariaDB
	"Innodb_encryption_n_rowlog_blocks_decrypted":                  "COUNTER", // MariaDB
	"Innodb_scrub_background_page_reorganizations":                 "COUNTER", // MariaDB
	"Innodb_scrub_background_page_splits":                          "COUNTER", // MariaDB
	"Innodb_scrub_background_page_split_failures_underflow":        "COUNTER", // MariaDB
	"Innodb_scrub_background_page_split_failures_out_of_filespace": "COUNTER", // MariaDB
	"Innodb_scrub_background_page_split_failures_missing_index":    "COUNTER", // MariaDB
	"Innodb_scrub_background_page_split_failures_unknown":          "COUNTER", // MariaDB
	"Innodb_encryption_num_key_requests":                           "COUNTER", // MariaDB
	"Binlog_event_count":                                           "COUNTER",
	"Binlog_number":                                                "COUNTER",
	"Slave_open_temp_tables":                                       "COUNTER",
	"Slave_received_heartbeats":                                    "COUNTER", // MariaDB
	"Slave_retried_transactions":                                   "COUNTER", // MariaDB
	"Slave_skipped_errors":                                         "COUNTER", // MariaDB
	"Com_admin_commands":                                           "COUNTER",
	"Com_alter_db":                                                 "COUNTER",
	"Com_alter_db_upgrade":                                         "COUNTER",
	"Com_alter_event":                                              "COUNTER",
	"Com_alter_function":                                           "COUNTER",
	"Com_alter_instance":                                           "COUNTER", // MySQL
	"Com_alter_procedure":                                          "COUNTER",
	"Com_alter_server":                                             "COUNTER",
	"Com_alter_sequence":                                           "COUNTER",
	"Com_alter_table":                                              "COUNTER",
	"Com_alter_tablespace":                                         "COUNTER",
	"Com_alter_user":                                               "COUNTER",
	"Com_analyze":                                                  "COUNTER",
	"Com_assign_to_keycache":                                       "COUNTER",
	"Com_begin":                                                    "COUNTER",
	"Com_binlog":                                                   "COUNTER",
	"Com_call_procedure":                                           "COUNTER",
	"Com_change_db":                                                "COUNTER",
	"Com_change_master":                                            "COUNTER",
	"Com_change_repl_filter":                                       "COUNTER", // MySQL
	"Com_check":                                                    "COUNTER",
	"Com_checksum":                                                 "COUNTER",
	"Com_commit":                                                   "COUNTER",
	"Com_compound_sql":                                             "COUNTER", // MariaDB
	"Com_create_db":                                                "COUNTER",
	"Com_create_event":                                             "COUNTER",
	"Com_create_function":                                          "COUNTER",
	"Com_create_index":                                             "COUNTER",
	"Com_create_package":                                           "COUNTER", // MariaDB
	"Com_create_package_body":                                      "COUNTER", // MariaDB
	"Com_create_procedure":                                         "COUNTER",
	"Com_create_role":                                              "COUNTER", // MariaDB
	"Com_create_sequence":                                          "COUNTER", // MariaDB
	"Com_create_server":                                            "COUNTER",
	"Com_create_table":                                             "COUNTER",
	"Com_create_temporary_table":                                   "COUNTER", // MariaDB
	"Com_create_trigger":                                           "COUNTER",
	"Com_create_udf":                                               "COUNTER",
	"Com_create_user":                                              "COUNTER",
	"Com_create_view":                                              "COUNTER",
	"Com_dealloc_sql":                                              "COUNTER",
	"Com_delete":                                                   "COUNTER",
	"Com_delete_multi":                                             "COUNTER",
	"Com_do":                                                       "COUNTER",
	"Com_drop_db":                                                  "COUNTER",
	"Com_drop_event":                                               "COUNTER",
	"Com_drop_function":                                            "COUNTER",
	"Com_drop_index":                                               "COUNTER",
	"Com_drop_procedure":                                           "COUNTER",
	"Com_drop_package":                                             "COUNTER", // MariaDB
	"Com_drop_package_body":                                        "COUNTER", // MariaDB
	"Com_drop_role":                                                "COUNTER", // MariaDB
	"Com_drop_server":                                              "COUNTER",
	"Com_drop_sequence":                                            "COUNTER", // MariaDB
	"Com_drop_table":                                               "COUNTER",
	"Com_drop_temporary_table":                                     "COUNTER", // MariaDB
	"Com_drop_trigger":                                             "COUNTER",
	"Com_drop_user":                                                "COUNTER",
	"Com_drop_view":                                                "COUNTER",
	"Com_empty_query":                                              "COUNTER",
	"Com_execute_immediate":                                        "COUNTER", // MariaDB
	"Com_execute_sql":                                              "COUNTER",
	"Com_explain_other":                                            "COUNTER", // MySQL
	"Com_flush":                                                    "COUNTER",
	"Com_get_diagnostics":                                          "COUNTER",
	"Com_grant":                                                    "COUNTER",
	"Com_grant_role":                                               "COUNTER", // MariaDB
	"Com_group_replication_start":                                  "COUNTER", // MySQL
	"Com_group_replication_stop":                                   "COUNTER", // MySQL
	"Com_ha_close":                                                 "COUNTER",
	"Com_ha_open":                                                  "COUNTER",
	"Com_ha_read":                                                  "COUNTER",
	"Com_help":                                                     "COUNTER",
	"Com_insert":                                                   "COUNTER",
	"Com_insert_select":                                            "COUNTER",
	"Com_install_plugin":                                           "COUNTER",
	"Com_kill":                                                     "COUNTER",
	"Com_load":                                                     "COUNTER",
	"Com_lock_tables":                                              "COUNTER",
	"Com_multi":                                                    "COUNTER", // MariaDB
	"Com_optimize":                                                 "COUNTER",
	"Com_preload_keys":                                             "COUNTER",
	"Com_prepare_sql":                                              "COUNTER",
	"Com_purge":                                                    "COUNTER",
	"Com_purge_before_date":                                        "COUNTER",
	"Com_release_savepoint":                                        "COUNTER",
	"Com_rename_table":                                             "COUNTER",
	"Com_rename_user":                                              "COUNTER",
	"Com_repair":                                                   "COUNTER",
	"Com_replace":                                                  "COUNTER",
	"Com_replace_select":                                           "COUNTER",
	"Com_reset":                                                    "COUNTER",
	"Com_resignal":                                                 "COUNTER",
	"Com_revoke":                                                   "COUNTER",
	"Com_revoke_all":                                               "COUNTER",
	"Com_revoke_role":                                              "COUNTER", // MariaDB
	"Com_rollback":                                                 "COUNTER",
	"Com_rollback_to_savepoint":                                    "COUNTER",
	"Com_savepoint":                                                "COUNTER",
	"Com_select":                                                   "COUNTER",
	"Com_set_option":                                               "COUNTER",
	"Com_show_authors":                                             "COUNTER", // MariaDB
	"Com_show_binlog_events":                                       "COUNTER",
	"Com_show_binlogs":                                             "COUNTER",
	"Com_show_charsets":                                            "COUNTER",
	"Com_show_collations":                                          "COUNTER",
	"Com_show_contributors":                                        "COUNTER",
	"Com_show_create_db":                                           "COUNTER",
	"Com_show_create_event":                                        "COUNTER",
	"Com_show_create_func":                                         "COUNTER",
	"Com_show_create_package":                                      "COUNTER", // MariaDB
	"Com_show_create_package_body":                                 "COUNTER", // MariaDB
	"Com_show_create_proc":                                         "COUNTER",
	"Com_show_create_table":                                        "COUNTER",
	"Com_show_create_trigger":                                      "COUNTER",
	"Com_show_create_user":                                         "COUNTER",
	"Com_show_databases":                                           "COUNTER",
	"Com_show_engine_logs":                                         "COUNTER",
	"Com_show_engine_mutex":                                        "COUNTER",
	"Com_show_engine_status":                                       "COUNTER",
	"Com_show_errors":                                              "COUNTER",
	"Com_show_events":                                              "COUNTER",
	"Com_show_explain":                                             "COUNTER",
	"Com_show_fields":                                              "COUNTER",
	"Com_show_function_code":                                       "COUNTER", // MySQL
	"Com_show_function_status":                                     "COUNTER",
	"Com_show_generic":                                             "COUNTER", // MariaDB
	"Com_show_grants":                                              "COUNTER",
	"Com_show_keys":                                                "COUNTER",
	"Com_show_master_status":                                       "COUNTER",
	"Com_show_open_tables":                                         "COUNTER",
	"Com_show_package_status":                                      "COUNTER", // MariaDB
	"Com_show_package_body_status":                                 "COUNTER", // MariaDB
	"Com_show_plugins":                                             "COUNTER",
	"Com_show_privileges":                                          "COUNTER",
	"Com_show_procedure_code":                                      "COUNTER", // MySQL
	"Com_show_procedure_status":                                    "COUNTER",
	"Com_show_processlist":                                         "COUNTER",
	"Com_show_profile":                                             "COUNTER",
	"Com_show_profiles":                                            "COUNTER",
	"Com_show_relaylog_events":                                     "COUNTER",
	"Com_show_slave_hosts":                                         "COUNTER",
	"Com_show_slave_status":                                        "COUNTER",
	"Com_show_status":                                              "COUNTER",
	"Com_show_storage_engines":                                     "COUNTER",
	"Com_show_table_status":                                        "COUNTER",
	"Com_show_tables":                                              "COUNTER",
	"Com_show_triggers":                                            "COUNTER",
	"Com_show_variables":                                           "COUNTER",
	"Com_show_warnings":                                            "COUNTER",
	"Com_shutdown":                                                 "COUNTER",
	"Com_signal":                                                   "COUNTER",
	"Com_start_all_slaves":                                         "COUNTER",
	"Com_start_slave":                                              "COUNTER", // MariaDB
	"Com_stmt_close":                                               "COUNTER",
	"Com_stmt_execute":                                             "COUNTER",
	"Com_stmt_fetch":                                               "COUNTER",
	"Com_stmt_prepare":                                             "COUNTER",
	"Com_stmt_reprepare":                                           "COUNTER",
	"Com_stmt_reset":                                               "COUNTER",
	"Com_stmt_send_long_data":                                      "COUNTER",
	"Com_stop_all_slaves":                                          "COUNTER",
	"Com_stop_slave":                                               "COUNTER", // MariaDB
	"Com_truncate":                                                 "COUNTER",
	"Com_uninstall_plugin":                                         "COUNTER",
	"Com_unlock_tables":                                            "COUNTER",
	"Com_update":                                                   "COUNTER",
	"Com_update_multi":                                             "COUNTER",
	"Com_xa_commit":                                                "COUNTER",
	"Com_xa_end":                                                   "COUNTER",
	"Com_xa_prepare":                                               "COUNTER",
	"Com_xa_recover":                                               "COUNTER",
	"Com_xa_rollback":                                              "COUNTER",
	"Com_xa_start":                                                 "COUNTER",
	"Com_slave_start":                                              "COUNTER", // Mysql
	"Com_slave_stop":                                               "COUNTER", // Mysql
	"Aborted_clients":                                              "COUNTER",
	"Aborted_connects":                                             "COUNTER",
	"Access_denied_errors":                                         "COUNTER",
	"Binlog_bytes_written":                                         "COUNTER",
	"Binlog_cache_disk_use":                                        "COUNTER",
	"Binlog_cache_use":                                             "COUNTER",
	"Binlog_stmt_cache_disk_use":                                   "COUNTER",
	"Binlog_stmt_cache_use":                                        "COUNTER",
	"Bytes_received":                                               "COUNTER",
	"Bytes_sent":                                                   "COUNTER",
	"Connections":                                                  "COUNTER",
	"Created_tmp_disk_tables":                                      "COUNTER",
	"Created_tmp_files":                                            "COUNTER",
	"Created_tmp_tables":                                           "COUNTER",
	"Handler_commit":                                               "COUNTER",
	"Handler_delete":                                               "COUNTER",
	"Handler_discover":                                             "COUNTER",
	"Handler_external_lock":                                        "COUNTER",
	"Handler_icp_attempts":                                         "COUNTER",
	"Handler_icp_match":                                            "COUNTER",
	"Handler_mrr_init":                                             "COUNTER",
	"Handler_mrr_key_refills":                                      "COUNTER",
	"Handler_mrr_rowid_refills":                                    "COUNTER",
	"Handler_prepare":                                              "COUNTER",
	"Handler_read_first":                                           "COUNTER",
	"Handler_read_key":                                             "COUNTER",
	"Handler_read_last":                                            "COUNTER",
	"Handler_read_next":                                            "COUNTER",
	"Handler_read_prev":                                            "COUNTER",
	"Handler_read_retry":                                           "COUNTER",
	"Handler_read_rnd":                                             "COUNTER",
	"Handler_read_rnd_deleted":                                     "COUNTER",
	"Handler_read_rnd_next":                                        "COUNTER",
	"Handler_rollback":                                             "COUNTER",
	"Handler_savepoint":                                            "COUNTER",
	"Handler_savepoint_rollback":                                   "COUNTER",
	"Handler_tmp_delete":                                           "COUNTER",
	"Handler_tmp_update":                                           "COUNTER",
	"Handler_tmp_write":                                            "COUNTER",
	"Handler_update":                                               "COUNTER",
	"Handler_write":                                                "COUNTER",
	"Opened_files":                                                 "COUNTER",
	"Opened_tables":                                                "COUNTER",
	"Opened_table_definitions":                                     "COUNTER",
	"Qcache_hits":                                                  "COUNTER",
	"Qcache_inserts":                                               "COUNTER",
	"Qcache_lowmem_prunes":                                         "COUNTER",
	"Qcache_not_cached":                                            "COUNTER",
	"Queries":                                                      "COUNTER",
	"Questions":                                                    "COUNTER",
	"Select_full_join":                                             "COUNTER",
	"Select_full_range_join":                                       "COUNTER",
	"Select_range":                                                 "COUNTER",
	"Select_range_check":                                           "COUNTER",
	"Select_scan":                                                  "COUNTER",
	"Slow_queries":                                                 "COUNTER",
	"Sort_merge_passes":                                            "COUNTER",
	"Sort_range":                                                   "COUNTER",
	"Sort_rows":                                                    "COUNTER",
	"Sort_scan":                                                    "COUNTER",
	"Table_locks_immediate":                                        "COUNTER",
	"Table_locks_waited":                                           "COUNTER",
	"Threads_created":                                              "COUNTER",
	"Rpl_semi_sync_master_net_wait_time":                           "COUNTER",
	"Rpl_semi_sync_master_net_waits":                               "COUNTER",
	"Rpl_semi_sync_master_no_times":                                "COUNTER",
	"Rpl_semi_sync_master_no_tx":                                   "COUNTER",
	"Rpl_semi_sync_master_yes_tx":                                  "COUNTER",
	"Rpl_semi_sync_master_tx_wait_time":                            "COUNTER",
	"Rpl_semi_sync_master_tx_waits":                                "COUNTER",
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
func NewMetric(name string) *cmodel.MetricValue {
	hostname, _ := g.Hostname()
	return &cmodel.MetricValue{
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
	role := "master"
	if IsSlave == 1 {
		role = "slave"
	}
	return fmt.Sprintf("role=%s,port=%d", role, g.Config().Collector.MySQL.Port)
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

func MySQLMetrics() (L []*cmodel.MetricValue) {
	if g.Config().Collector.MySQL == nil {
		return nil
	}
	if !g.Config().Collector.MySQL.Enabled {
		return nil
	}

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql",
			g.Config().Collector.MySQL.User,
			g.Config().Collector.MySQL.Passowrd,
			g.Config().Collector.MySQL.Host,
			g.Config().Collector.MySQL.Port))

	if err != nil {
		log.Errorf("[E] connect to mysql error: %v", err)
		return nil
	}
	defer db.Close()

	// Get slave status and set IsSlave global var
	slaveState, err := ShowSlaveStatus(db)
	if err != nil {
		log.Errorf("[E] show slave status error: %v", err)
		return
	}
	Tag = GetTag()

	globalStatus, err := ShowGlobalStatus(db)
	if err != nil {
		log.Errorf("[E] show status error: %v", err)
		return
	}
	L = append(L, globalStatus...)

	globalVars, err := ShowGlobalVariables(db)
	if err != nil {
		log.Errorf("[E] show variables error: %v", err)
		return
	}
	L = append(L, globalVars...)

	innodbState, err := ShowInnodbStatus(db)
	if err != nil {
		log.Errorf("[E] show innodb status error: %v", err)
		return
	}
	L = append(L, innodbState...)
	L = append(L, slaveState...)

	binaryLogStatus, err := ShowBinaryLogs(db)
	if err != nil {
		log.Errorf("[E] show bin log error: %v", err)
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
func ShowGlobalStatus(db *sql.DB) ([]*cmodel.MetricValue, error) {
	return parseMySQLStatus(db, "SHOW /*!50001 GLOBAL */ STATUS")
}

// ShowGlobalVariables execute mysql query `SHOW GLOBAL VARIABLES`
func ShowGlobalVariables(db *sql.DB) ([]*cmodel.MetricValue, error) {
	return parseMySQLStatus(db, "SHOW /*!50001 GLOBAL */ VARIABLES")
}

// ShowInnodbStatus execute mysql query `SHOW SHOW /*!50000 ENGINE*/ INNODB STATUS`
func ShowInnodbStatus(db *sql.DB) ([]*cmodel.MetricValue, error) {
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
func ShowBinaryLogs(db *sql.DB) ([]*cmodel.MetricValue, error) {
	sum := 0

	binlogFileCounts := NewMetric("binlog_file_counts")
	binlogFileSize := NewMetric("binlog_file_size")

	results, err := Query(db, "SHOW BINARY LOGS")
	if err != nil {
		return []*cmodel.MetricValue{binlogFileCounts, binlogFileSize}, err
	}

	for _, v := range results {
		size, err := strconv.Atoi(strings.TrimSpace(v["File_size"].(string)))
		if err == nil {
			sum += size
		}
	}

	binlogFileCounts.Value = len(results)
	binlogFileSize.Value = sum
	return []*cmodel.MetricValue{binlogFileCounts, binlogFileSize}, err
}

// ShowSlaveStatus get all slave status of mysql serves
func ShowSlaveStatus(db *sql.DB) ([]*cmodel.MetricValue, error) {
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
		return []*cmodel.MetricValue{masterReadOnly, masterReadOnly2, innodbStatsOnMetadata}, nil
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
	data := make([]*cmodel.MetricValue, len(SlaveStatus))
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
	return append(data, []*cmodel.MetricValue{ioDelay, slaveReadOnly}...), nil
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
func ShowOtherMetric(db *sql.DB, metric string) (*cmodel.MetricValue, error) {
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
	}

	return newMetaData, err
}

func parseMySQLStatus(db *sql.DB, sql string) ([]*cmodel.MetricValue, error) {
	results, err := Query(db, sql)
	if err != nil {
		return nil, err
	}

	data := make([]*cmodel.MetricValue, len(results))
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
	pdata *[]*cmodel.MetricValue, longTranTime *int) error {
	switch section {
	case "TRANSACTIONS":
		if strings.Contains(row, "ACTIVE") {
			// TRANSACTION 97779198, ACTIVE 2 sec
			// TRANSACTION 97779198, ACTIVE (PREPARED) 2 sec
			// TRANSACTION 97869332, ACTIVE (PREPARED) 0 sec committing
			matches := regexp.MustCompile(`^---TRANSACTION\s+(?P<T>\d+),\s+ACTIVE\s+(?:\(PREPARED\)\s+){0,1}(?P<D>\d+)\s+sec`).FindStringSubmatch(row)
			if len(matches) == 3 {
				tmpLongTransactionTime, _ := strconv.Atoi(matches[2])
				if tmpLongTransactionTime > *longTranTime {
					*longTranTime = tmpLongTransactionTime
				}
			}
		}
		if strings.Contains(row, "History list length") {
			hisListLengthStr := strings.Split(row, "length ")[1]
			hisListLength, err := strconv.Atoi(hisListLengthStr)
			if err != nil {
				log.Errorf("[E] extract history list length from %s error: %v", row, err)
			} else {
				HistoryListLength := NewMetric("History_list_length")
				HistoryListLength.Value = hisListLength
				*pdata = append(*pdata, HistoryListLength)
			}
		}
	case "SEMAPHORES":
		matches := regexp.MustCompile(`^Mutex spin waits\s+(\d+),\s+rounds\s+(\d+),\s+OS waits\s+(\d+)`).FindStringSubmatch(row)
		if len(matches) == 4 {
			if spinWaits, err := strconv.Atoi(matches[1]); err != nil {
				log.Errorf("[E] extract spin waits from %s error: %v", matches[1], err)
			} else {
				innodbMutexSpinWaits := NewMetric("Innodb_mutex_spin_waits")
				innodbMutexSpinWaits.Value = spinWaits
				*pdata = append(*pdata, innodbMutexSpinWaits)
			}

			if spinRounds, err := strconv.Atoi(matches[2]); err != nil {
				log.Errorf("[E] extract spin rounds from %s error: %v", matches[2], err)
			} else {
				InnodbMutexSpinRounds := NewMetric("Innodb_mutex_spin_rounds")
				InnodbMutexSpinRounds.Value = spinRounds
				*pdata = append(*pdata, InnodbMutexSpinRounds)
			}

			if osWaits, err := strconv.Atoi(matches[3]); err != nil {
				log.Errorf("[E] extract os waits from %s error: %v", matches[3], err)
			} else {
				InnodbMutexOsWaits := NewMetric("Innodb_mutex_os_waits")
				InnodbMutexOsWaits.Value = osWaits
				*pdata = append(*pdata, InnodbMutexOsWaits)
			}
		}
	}
	return nil
}

func parseInnodbStatus(rows []string) ([]*cmodel.MetricValue, error) {
	var section string
	longTranTime := 0
	var err error
	var data []*cmodel.MetricValue
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
