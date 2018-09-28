falcon-agent
===

This is a linux monitor agent. Just like zabbix-agent and tcollector.


## Installation

It is a golang classic project

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/falcon-plus.git
cd falcon-plus/modules/agent
go get
./control build
./control start

# goto http://localhost:1988
```

I use [linux-dash](https://github.com/afaqurk/linux-dash) as the page theme.

## Configuration

- heartbeat: heartbeat server rpc address
- transfer: transfer rpc address
- ignore: the metrics should ignore

# Auto deployment

Just look at https://github.com/open-falcon/ops-updater

## 采集的Redis指标
----------------------------------------

| Counters | Type | Notes|
|-----|------|------|
|aof_current_rewrite_time_sec  |GAUGE|当前AOF重写持续的耗时|
|aof_enabled                   |GAUGE| appenonly是否开启,appendonly为yes则为1,no则为0|
|aof_last_bgrewrite_status     |GAUGE|最近一次AOF重写操作是否成功|
|aof_last_rewrite_time_sec     |GAUGE|最近一次AOF重写操作耗时|
|aof_last_write_status         |GAUGE|最近一次AOF写入操作是否成功|
|aof_rewrite_in_progress       |GAUGE|AOF重写是否正在进行|
|aof_rewrite_scheduled         |GAUGE|AOF重写是否被RDB save操作阻塞等待|
|blocked_clients               |GAUGE|正在等待阻塞命令（BLPOP、BRPOP、BRPOPLPUSH）的客户端的数量|
|client_biggest_input_buf      |GAUGE|当前客户端连接中，最大输入缓存|
|client_longest_output_list    |GAUGE|当前客户端连接中，最长的输出列表|
|cluster_enabled               |GAUGE|是否启用Redis集群模式，cluster_enabled|
|cluster_known_nodes           |GAUGE|集群中节点的个数|
|cluster_size                  |GAUGE|集群的大小，即集群的分区数个数|
|cluster_slots_assigned        |GAUGE|集群中已被指派slot个数，正常是16385个|
|cluster_slots_fail            |GAUGE|集群中已下线（客观失效）的slot个数|
|cluster_slots_ok              |GAUGE|集群中正常slots个数|
|cluster_slots_pfail           |GAUGE|集群中疑似下线（主观失效）的slot个数|
|cluster_state                 |GAUGE|集群的状态是否正常|
|cmdstat_auth                  |COUNTER|auth命令每秒执行次数|
|cmdstat_config                |COUNTER|config命令每秒执行次数|
|cmdstat_get                   |COUNTER|get命令每秒执行次数|
|cmdstat_info                  |COUNTER|info命令每秒执行次数|
|cmdstat_ping                  |COUNTER|ping命令每秒执行次数|
|cmdstat_set                   |COUNTER|set命令每秒执行次数|
|cmdstat_slowlog               |COUNTER|slowlog命令每秒执行次数|
|connected_clients             |GAUGE|当前已连接的客户端个数|
|connected_clients_pct         |GAUGE|已使用的连接数百分比，connected_clients/maxclients |
|connected_slaves              |GAUGE|已连接的Redis从库个数|
|evicted_keys                  |COUNTER|因内存used_memory达到maxmemory后，每秒被驱逐的key个数|
|expired_keys                  |COUNTER|因键过期后，被惰性和主动删除清理key的每秒个数|
|hz			       |GAUGE|serverCron执行的频率，默认10，表示100ms执行一次，建议不要大于120|
|instantaneous_input_kbps      |GAUGE|瞬间的Redis输入网络流量(kbps)|
|instantaneous_ops_per_sec     |GAUGE|瞬间的Redis操作QPS|
|instantaneous_output_kbps     |GAUGE|瞬间的Redis输出网络流量(kbps)|
|keys                          |GAUGE|当前Redis实例的key总数|
|keyspace_hit_ratio            |GAUGE|查找键的命中率（每个周期60sec精确计算)|
|keyspace_hits                 |COUNTER|查找键命中的次数|
|keyspace_misses               |COUNTER|查找键未命中的次数|
|latest_fork_usec              |GAUGE|最近一次fork操作的耗时的微秒数(BGREWRITEAOF,BGSAVE,SYNC等都会触发fork),当并发场景fork耗时过长对服务影响较大|
|loading		       |GAUGE|标志位，是否在载入数据文件|
|master_repl_offset            |GAUGE|master复制的偏移量，除了写入aof外，Redis定期为自动增加|
|mem_fragmentation_ratio       |GAUGE|内存碎片率，used_memory_rss/used_memory|
|pubsub_channels               |GAUGE|目前被订阅的频道数量|
|pubsub_patterns               |GAUGE|目前被订阅的模式数量|
|rdb_bgsave_in_progress        |GAUGE|标志位，记录当前是否在创建RDB快照|
|rdb_current_bgsave_time_sec   |GAUGE|当前bgsave执行耗时秒数|
|rdb_last_bgsave_status        |GAUGE|标志位，记录最近一次bgsave操作是否创建成功|
|rdb_last_bgsave_time_sec      |GAUGE|最近一次bgsave操作耗时秒数|
|rdb_last_save_time            |GAUGE|最近一次创建RDB快照文件的Unix时间戳|
|rdb_changes_since_last_save   |GAUGE|从最近一次dump快照后，未被dump的变更次数(和save里变更计数器类似)|
|alive                   |GAUGE|当前Redis是否存活，ping监控socket_time默认500ms|
|rejected_connections          |COUNTER|因连接数达到maxclients上限后，被拒绝的连接个数|
|repl_backlog_active           |GAUGE|标志位，master是否开启了repl_backlog,有效地psync(2.8+)|
|repl_backlog_first_byte_offset|GAUGE|repl_backlog中首字节的复制偏移位|
|repl_backlog_histlen          |GAUGE|repl_backlog当前使用的字节数|
|repl_backlog_size             |GAUGE|repl_backlog的长度(repl-backlog-size)，网络环境不稳定的，建议调整大些
|role                          |GAUGE|当前实例的角色：master 1， slave 0|
|master_link_status            |GAUGE|标志位,从库复制是否正常，正常1，断开0|
|master_link_down_since_seconds|GAUGE|从库断开复制的秒数|
|slave_read_only	       |GAUGE|从库是否设置为只读状态，避免写入|
|slowlog_len                   |COUNTER|slowlog的个数(因未转存slowlog实例，每次采集不会slowlog reset，所以当slowlog占满后，此值无意义)|
|sync_full                     |GAUGE|累计Master full sync的次数;如果值比较大，说明常常出现全量复制，就得分析原因，或调整repl-backlog-size|
|sync_partial_err              |GAUGE|累计Master pysync 出错失败的次数|
|sync_partial_ok               |GAUGE|累计Master psync成功的次数|
|total_commands_processed      |COUNTER|每秒执行的命令数，比较准确的QPS|
|total_connections_received    |COUNTER|每秒新创建的客户端连接数|
|total_net_input_bytes         |COUNTER|Redis每秒网络输入的字节数|
|total_net_output_bytes        |COUNTER|Redis每秒网络输出的字节数|
|uptime_in_days                |GAUGE|Redis运行时长天数|
|uptime_in_seconds	       |GAUGE|Redis运行时长的秒数|
|used_cpu_sys                  |COUNTER|Redis进程消耗的sys cpu|
|used_cpu_user                 |COUNTER|Redis进程消耗的user cpu|
|used_memory                   |GAUGE|由Redis分配的内存的总量，字节数|
|used_memory_lua               |GAUGE|lua引擎使用的内存总量，字节数；有使用lua脚本的注意监控|
|used_memory_pct               |GAUGE|最大内存已使用百分比,used_memory/maxmemory; 存储场景无淘汰key注意监控.(如果maxmemory=0表示未设置限制,pct永远为0)|
|used_memory_peak              |GAUGE|Redis使用内存的峰值，字节数|
|used_memory_rss               |GAUGE|Redis进程从OS角度分配的物理内存，如key被删除后，malloc不一定把内存归还给OS,但可以Redis进程复用|

## 建议设置监控告警项
-----------------------------

| 告警项 | 触发条件 | 备注|
|-----|------|------|
|load.1min|all(#3)>10|Redis服务器过载，处理能力下降|
|cpu.idle |all(#3)<10|CPU idle过低，处理能力下降|
|df.bytes.free.percent|all(#3)<20|磁盘可用空间百分比低于20%，影响从库RDB和AOF持久化|
|mem.memfree.percent|all(#3)<15|内存剩余低于15%，Redis有OOM killer和使用swap的风险|
|mem.swapfree.percent|all(#3)<80|使用20% swap,Redis性能下降或OOM风险|
|net.if.out.bytes|all(#3)>94371840|网络出口流量超90MB,影响Redis响应|
|net.if.in.bytes|all(#3)>94371840|网络入口流量超90MB,影响Redis响应|
|disk.io.util|all(#3)>90|磁盘IO可能存负载，影响从库持久化和阻塞写|
|redis.alive|all(#2)=0|Redis实例存活有问题，可能不可用|
|used_memory|all(#2)>32212254720|单实例使用30G，建议拆分扩容；对fork卡停，full_sync时长都有明显性能影响|
|used_memory_pct|all(#3)>85|(存储场景)使用内存达85%,存储场景会写入失败|
|mem_fragmentation_ratio|all(#3)>2|内存碎片过高(如果实例比较小，这个指标可能比较大，不实用)|
|connected_clients|all(#3)>5000|客户端连接数超5000|
|connected_clients_pct|all(#3)>85|客户端连接数占最大连接数超85%|
|rejected_connections|all(#1)>0|连接数达到maxclients后，创建新连接失败|
|total_connections_received|每秒新创建连接数超5000，对Redis性能有明显影响，常见于PHP短连接场景|
|master_link_status>|all(#1)=0|主从同步断开；会全量同步，HA/备份/读写分离数据最终一致性受影响|
|slave_read_only|all(#1)=0|从库非只读状态|
|repl_backlog_active|all(#1)=0|repl_backlog关闭，对网络闪断场景不能psync|
|keys|all(#1)>50000000|keyspace key总数5千万，建议拆分扩容|
|instantaneous_ops_per_sec|all(#2)>30000|整体QPS 30000,建议拆分扩容|
|slowlog_len|all(#1)>10|1分钟中内，出现慢查询个数(一般Redis命令响应大于1ms，记录为slowlog)|
|latest_fork_usec|all(#1)>1000000|最近一次fork耗时超1秒(其间Redis不能响应任何请求)|
|keyspace_hit_ratio|all(#2)<80|命中率低于80%|
|cluster_state|all(#1)=0|Redis集群处理于FAIL状态，不能提供读写|
|cluster_slots_assigned|all(#1)<16384|keyspace的所有数据槽未被全部指派，集群处理于FAIL状态|
|cluster_slots_fail|all(#1)>0|集群中有槽处于失败，集群处理于FAIL状态|


## 采集的MongoDB指标
----------------------------------------

| Counters | Type | Notes|
|-----|------|------|
|mongo_local_alive|                        GAUGE   |mongodb存活本地监控，如果开启Auth，要求连接认证成功
|asserts_msg|                              COUNTER |消息断言数/秒
|asserts_regular|                          COUNTER |常规断言数/秒
|asserts_rollovers|                        COUNTER |计数器roll over的次数/秒,计数器每2^30个断言就会清零
|asserts_user|                             COUNTER |用户断言数/秒
|asserts_warning|                          COUNTER |警告断言数/秒
|page_faults|                              COUNTER |页缺失次数/秒
|connections_available|                    GAUGE   |未使用的可用连接数
|connections_current|                      GAUGE   |当前所有客户端的已连接的连接数
|connections_used_percent|                 GAUGE   |已使用连接数百分比
|connections_totalCreated|                 COUNTER |创建的新连接数/秒
|globalLock_currentQueue_total|            GAUGE   |当前队列中等待锁的操作数
|globalLock_currentQueue_readers|          GAUGE   |当前队列中等待读锁的操作数
|globalLock_currentQueue_writers|          GAUGE   |当前队列中等待写锁的操作数
|locks_Global_acquireCount_ISlock|         COUNTER |实例级意向共享锁获取次数
|locks_Global_acquireCount_IXlock|         COUNTER |实例级意向排他锁获取次数
|locks_Global_acquireCount_Slock|          COUNTER |实例级共享锁获取次数
|locks_Global_acquireCount_Xlock|          COUNTER |实例级排他锁获取次数
|locks_Global_acquireWaitCount_ISlock|     COUNTER |实例级意向共享锁等待次数
|locks_Global_acquireWaitCount_IXlock|     COUNTER |实例级意向排他锁等待次数
|locks_Global_timeAcquiringMicros_ISlock|  COUNTER |实例级共享锁获取耗时 单位:微秒
|locks_Global_timeAcquiringMicros_IXlock|  COUNTER |实例级共排他获取耗时 单位:微秒
|locks_Database_acquireCount_ISlock|       COUNTER |数据库级意向共享锁获取次数
|locks_Database_acquireCount_IXlock|       COUNTER |数据库级意向排他锁获取次数
|locks_Database_acquireCount_Slock|        COUNTER |数据库级共享锁获取次数
|locks_Database_acquireCount_Xlock|        COUNTER |数据库级排他锁获取次数
|locks_Collection_acquireCount_ISlock|     COUNTER |集合级意向共享锁获取次数
|locks_Collection_acquireCount_IXlock|     COUNTER |集合级意向排他锁获取次数
|locks_Collection_acquireCount_Xlock|      COUNTER |集合级排他锁获取次数
|opcounters_command|                       COUNTER |数据库执行的所有命令/秒
|opcounters_insert|                        COUNTER |数据库执行的插入操作次数/秒
|opcounters_delete|                        COUNTER |数据库执行的删除操作次数/秒
|opcounters_update|                        COUNTER |数据库执行的更新操作次数/秒
|opcounters_query|                         COUNTER |数据库执行的查询操作次数/秒
|opcounters_getmore|                       COUNTER |数据库执行的getmore操作次数/秒
|opcountersRepl_command|                   COUNTER |数据库复制执行的所有命令次数/秒
|opcountersRepl_insert|                    COUNTER |数据库复制执行的插入命令次数/秒
|opcountersRepl_delete|                    COUNTER |数据库复制执行的删除命令次数/秒
|opcountersRepl_update|                    COUNTER |数据库复制执行的更新命令次数/秒
|opcountersRepl_query|                     COUNTER |数据库复制执行的查询命令次数/秒
|opcountersRepl_getmore|                   COUNTER |数据库复制执行的gtemore命令次数/秒
|network_bytesIn|                          COUNTER |数据库接受的网络传输字节数/秒
|network_bytesOut|                         COUNTER |数据库发送的网络传输字节数/秒
|network_numRequests|                      COUNTER |数据库接收到的请求的总次数/秒
|mem_virtual|                              GAUGE   |数据库进程使用的虚拟内存
|mem_resident|                             GAUGE   |数据库进程使用的物理内存
|mem_mapped|                               GAUGE   |mapped的内存,只用于MMAPv1 存储引擎
|mem_bits|                                 GAUGE   |64 or 32bit
|mem_mappedWithJournal|                    GAUGE   |journal日志消耗的映射内存，只用于MMAPv1 存储引擎
|backgroundFlushing_flushes|               COUNTER |数据库刷新写操作到磁盘的次数/秒
|backgroundFlushing_average_ms|            GAUGE   |数据库刷新写操作到磁盘的平均耗时，单位ms
|backgroundFlushing_last_ms|               COUNTER |当前最近一次数据库刷新写操作到磁盘的耗时，单位ms
|backgroundFlushing_total_ms|              GAUGE   |数据库刷新写操作到磁盘的总耗时/秒，单位ms
|cursor_open_total|                        GAUGE   |当前数据库为客户端维护的游标总数
|cursor_timedOut|                          COUNTER |数据库timout的游标个数/秒
|cursor_open_noTimeout|                    GAUGE   |设置DBQuery.Option.noTimeout的游标数
|cursor_open_pinned|                       GAUGE   |打开的pinned的游标数
|repl_health|                              GAUGE   |复制的健康状态
|repl_myState|                             GAUGE   |当前节点的副本集状态
|repl_oplog_window|                        GAUGE   |oplog的窗口大小
|repl_optime|                              GAUGE   |上次执行的时间戳
|replication_lag_percent|                  GAUGE   |延时占比(lag/oplog_window)
|repl_lag|                                 GAUGE   |Secondary复制延时，单位秒
|shards_size|                              GAUGE   |数据库集群的分片个数; config.shards.count
|shards_mongosSize|                        GAUGE   |数据库集群中mongos节点个数；config.mongos.count
|shards_chunkSize|                         GAUGE   |数据库集群的chunksize大小设置，以config.settings集合中获取
|shards_activeWindow|                      GAUGE   |数据库集群的数据均衡器是否设置了时间窗口，1/0
|shards_activeWindow_start|                GAUGE   |数据库集群的数据均衡器时间窗口开始时间，格式23.30表示 23：30分
|shards_activeWindow_stop|                 GAUGE   |数据库集群的数据均衡器时间窗口结束时间，格式23.30表示 23：30分
|shards_BalancerState|                     GAUGE   |数据库集群的数据均衡器的状态，是否为打开
|shards_isBalancerRunning|                 GAUGE   |数据库集群的数据均衡器是否正在运行块迁移
|wt_cache_used_total_bytes|                GAUGE   |wiredTiger cache的字节数
|wt_cache_dirty_bytes|                     GAUGE   |wiredTiger cache中"dirty"数据的字节数
|wt_cache_readinto_bytes|                  COUNTER |数据库写入wiredTiger cache的字节数/秒
|wt_cache_writtenfrom_bytes|               COUNTER |数据库从wiredTiger cache写入到磁盘的字节数/秒
|wt_concurrentTransactions_write|          GAUGE   |write tickets available to the WiredTiger storage engine
|wt_concurrentTransactions_read|           GAUGE   |read tickets available to the WiredTiger storage engine
|wt_bm_bytes_read|                         COUNTER |block-manager read字节数/秒
|wt_bm_bytes_written|                      COUNTER |block-manager write字节数/秒
|wt_bm_blocks_read|                        COUNTER |block-manager read块数/秒
|wt_bm_blocks_written|                     COUNTER |block-manager write块数/秒
|rocksdb_num_immutable_mem_table|
|rocksdb_mem_table_flush_pending|
|rocksdb_compaction_pending|
|rocksdb_background_errors|
|rocksdb_num_entries_active_mem_table|
|rocksdb_num_entries_imm_mem_tables|
|rocksdb_num_snapshots|
|rocksdb_oldest_snapshot_time|
|rocksdb_num_live_versions|
|rocksdb_total_live_recovery_units|
|PerconaFT_cachetable_size_current|
|PerconaFT_cachetable_size_limit|
|PerconaFT_cachetable_size_writing|
|PerconaFT_checkpoint_count|
|PerconaFT_checkpoint_time|
|PerconaFT_checkpoint_write_leaf_bytes_compressed|
|PerconaFT_checkpoint_write_leaf_bytes_uncompressed|
|PerconaFT_checkpoint_write_leaf_count|
|PerconaFT_checkpoint_write_leaf_time|
|PerconaFT_checkpoint_write_nonleaf_bytes_compressed|
|PerconaFT_checkpoint_write_nonleaf_bytes_uncompressed|
|PerconaFT_checkpoint_write_nonleaf_count|
|PerconaFT_checkpoint_write_nonleaf_time|
|PerconaFT_compressionRatio_leaf|
|PerconaFT_compressionRatio_nonleaf|
|PerconaFT_compressionRatio_overall|
|PerconaFT_fsync_count|
|PerconaFT_fsync_time|
|PerconaFT_log_bytes|
|PerconaFT_log_count|
|PerconaFT_log_time|
|PerconaFT_serializeTime_leaf_compress|
|PerconaFT_serializeTime_leaf_decompress|
|PerconaFT_serializeTime_leaf_deserialize|
|PerconaFT_serializeTime_leaf_serialize|
|PerconaFT_serializeTime_nonleaf_compress|
|PerconaFT_serializeTime_nonleaf_decompress|
|PerconaFT_serializeTime_nonleaf_deserialize|
|PerconaFT_serializeTime_nonleaf_serialize|


## 建议设置监控告警项
-----------------------------

|告警项|
|-----|
|load.1min>10|
|cpu.idle<10|
|df.bytes.free.percent<30|
|df.bytes.free.percent<10|
|mem.memfree.percent<20|
|mem.memfree.percent<10|
|mem.memfree.percent<5|
|mem.swapfree.percent<50|
|mem.memused.percent>=50|
|mem.memused.percent>=10|
|net.if.out.bytes>94371840|
|net.if.in.bytes>94371840|
|disk.io.util>90|
|mongo_local_alive=0|
|page_faults>100|
|connections_current>5000|
|connections_used_percent>60|
|connections_used_percent>80|
|connections_totalCreated>1000|
|globalLock_currentQueue_total>10|
|globalLock_currentQueue_readers>10|
|globalLock_currentQueue_writers>10|
|opcounters_command|
|opcounters_insert|
|opcounters_delete|
|opcounters_update|
|opcounters_query|
|opcounters_getmore|
|opcountersRepl_command|
|opcountersRepl_insert|
|opcountersRepl_delete|
|opcountersRepl_update|
|opcountersRepl_query|
|opcountersRepl_getmore|
|network_bytesIn|
|network_bytesOut|
|network_numRequests|
|repl_health=0|
|repl_myState not 1/2/7|
|repl_oplog_window<168|
|repl_oplog_window<48|
|replication_lag_percent>50|
|repl_lag>60|
|repl_lag>300|
|shards_mongosSize|
