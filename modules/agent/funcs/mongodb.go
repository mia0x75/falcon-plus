package funcs

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// AssertsStats has the assets metrics
type AssertsStats struct {
	Regular   float64 `bson:"regular"`
	Warning   float64 `bson:"warning"`
	Msg       float64 `bson:"msg"`
	User      float64 `bson:"user"`
	Rollovers float64 `bson:"rollovers"`
}

// DurTiming is the information about durability returned from the server.
type DurTiming struct {
	Dt               float64 `bson:"dt"`
	PrepLogBuffer    float64 `bson:"prepLogBuffer"`
	WriteToJournal   float64 `bson:"writeToJournal"`
	WriteToDataFiles float64 `bson:"writeToDataFiles"`
	RemapPrivateView float64 `bson:"remapPrivateView"`
}

// DurStats are the stats related to durability.
type DurStats struct {
	Commits            float64   `bson:"commits"`
	JournaledMB        float64   `bson:"journaledMB"`
	WriteToDataFilesMB float64   `bson:"writeToDataFilesMB"`
	Compression        float64   `bson:"compression"`
	CommitsInWriteLock float64   `bson:"commitsInWriteLock"`
	EarlyCommits       float64   `bson:"earlyCommits"`
	TimeMs             DurTiming `bson:"timeMs"`
}

// ConnectionStats are connections metrics
type ConnectionStats struct {
	Current      float64 `bson:"current"`
	Available    float64 `bson:"available"`
	TotalCreated float64 `bson:"totalCreated"`
}

// ExtraInfo has extra info metrics
type ExtraInfo struct {
	HeapUsageBytes float64 `bson:"heap_usage_bytes"`
	PageFaults     float64 `bson:"page_faults"`
}

// ReadWriteLockTimes information about the lock
type ReadWriteLockTimes struct {
	Read       float64 `bson:"R"`
	Write      float64 `bson:"W"`
	ReadLower  float64 `bson:"r"`
	WriteLower float64 `bson:"w"`
}

// LockStats lock stats
type LockStats struct {
	TimeLockedMicros    ReadWriteLockTimes `bson:"timeLockedMicros"`
	TimeAcquiringMicros ReadWriteLockTimes `bson:"timeAcquiringMicros"`
}

// LockStatsMap is a map of lock stats
type LockStatsMap map[string]LockStats

//IndexCounterStats index counter stats
type IndexCounterStats struct {
	Accesses  float64 `bson:"accesses"`
	Hits      float64 `bson:"hits"`
	Misses    float64 `bson:"misses"`
	Resets    float64 `bson:"resets"`
	MissRatio float64 `bson:"missRatio"`
}

// StorageEngineStats TODO:
type StorageEngineStats struct {
	Name string `bson:"name"`
}

// WTBlockManagerStats TODO:
type WTBlockManagerStats struct {
	MappedBytesRead  float64 `bson:"mapped bytes read"`
	BytesRead        float64 `bson:"bytes read"`
	BytesWritten     float64 `bson:"bytes written"`
	MappedBlocksRead float64 `bson:"mapped blocks read"`
	BlocksPreLoaded  float64 `bson:"blocks pre-loaded"`
	BlocksRead       float64 `bson:"blocks read"`
	BlocksWritten    float64 `bson:"blocks written"`
}

// RocksDbStatsCounters TODO:
type RocksDbStatsCounters struct {
	NumKeysWritten         float64 `bson:"num-keys-written"`
	NumKeysRead            float64 `bson:"num-keys-read"`
	NumSeeks               float64 `bson:"num-seeks"`
	NumForwardIter         float64 `bson:"num-forward-iterations"`
	NumBackwardIter        float64 `bson:"num-backward-iterations"`
	BlockCacheMisses       float64 `bson:"block-cache-misses"`
	BlockCacheHits         float64 `bson:"block-cache-hits"`
	BloomFilterUseful      float64 `bson:"bloom-filter-useful"`
	BytesWritten           float64 `bson:"bytes-written"`
	BytesReadPointLookup   float64 `bson:"bytes-read-point-lookup"`
	BytesReadIteration     float64 `bson:"bytes-read-iteration"`
	FlushBytesWritten      float64 `bson:"flush-bytes-written"`
	CompactionBytesRead    float64 `bson:"compaction-bytes-read"`
	CompactionBytesWritten float64 `bson:"compaction-bytes-written"`
}

// RocksDbStats TODO:
type RocksDbStats struct {
	NumImmutableMemTable       string                `bson:"num-immutable-mem-table"`
	MemTableFlushPending       string                `bson:"mem-table-flush-pending"`
	CompactionPending          string                `bson:"compaction-pending"`
	BackgroundErrors           string                `bson:"background-errors"`
	CurSizeMemTableActive      string                `bson:"cur-size-active-mem-table"`
	CurSizeAllMemTables        string                `bson:"cur-size-all-mem-tables"`
	NumEntriesMemTableActive   string                `bson:"num-entries-active-mem-table"`
	NumEntriesImmMemTables     string                `bson:"num-entries-imm-mem-tables"`
	EstimateTableReadersMem    string                `bson:"estimate-table-readers-mem"`
	NumSnapshots               string                `bson:"num-snapshots"`
	OldestSnapshotTime         string                `bson:"oldest-snapshot-time"`
	NumLiveVersions            string                `bson:"num-live-versions"`
	BlockCacheUsage            string                `bson:"block-cache-usage"`
	TotalLiveRecoveryUnits     float64               `bson:"total-live-recovery-units"`
	TransactionEngineKeys      float64               `bson:"transaction-engine-keys"`
	TransactionEngineSnapshots float64               `bson:"transaction-engine-snapshots"`
	Stats                      []string              `bson:"stats"`
	ThreadStatus               []string              `bson:"thread-status"`
	Counters                   *RocksDbStatsCounters `bson:"counters,omitempty"`
}

// RocksDbLevelStatsFiles TODO:
type RocksDbLevelStatsFiles struct {
	Num         float64
	CompThreads float64
}

// RocksDbLevelStats TODO:
type RocksDbLevelStats struct {
	Level    string
	Files    *RocksDbLevelStatsFiles
	Score    float64
	SizeMB   float64
	ReadGB   float64
	RnGB     float64
	Rnp1GB   float64
	WriteGB  float64
	WnewGB   float64
	MovedGB  float64
	WAmp     float64
	RdMBPSec float64
	WrMBPSec float64
	CompSec  float64
	CompCnt  float64
	AvgSec   float64
	KeyIn    float64
	KeyDrop  float64
}

// WTCacheStats cache stats
type WTCacheStats struct {
	BytesTotal         float64 `bson:"bytes currently in the cache"`
	BytesDirty         float64 `bson:"tracked dirty bytes in the cache"`
	BytesInternalPages float64 `bson:"tracked bytes belonging to internal pages in the cache"`
	BytesLeafPages     float64 `bson:"tracked bytes belonging to leaf pages in the cache"`
	MaxBytes           float64 `bson:"maximum bytes configured"`
	BytesReadInto      float64 `bson:"bytes read into cache"`
	BytesWrittenFrom   float64 `bson:"bytes written from cache"`
	EvictedUnmodified  float64 `bson:"unmodified pages evicted"`
	EvictedModified    float64 `bson:"modified pages evicted"`
	PercentOverhead    float64 `bson:"percentage overhead"`
	PagesTotal         float64 `bson:"pages currently held in the cache"`
	PagesReadInto      float64 `bson:"pages read into cache"`
	PagesWrittenFrom   float64 `bson:"pages written from cache"`
	PagesDirty         float64 `bson:"tracked dirty pages in the cache"`
}

// WTLogStats log stats
type WTLogStats struct {
	TotalBufferSize         float64 `bson:"total log buffer size"`
	TotalSizeCompressed     float64 `bson:"total size of compressed records"`
	BytesPayloadData        float64 `bson:"log bytes of payload data"`
	BytesWritten            float64 `bson:"log bytes written"`
	RecordsUncompressed     float64 `bson:"log records not compressed"`
	RecordsCompressed       float64 `bson:"log records compressed"`
	RecordsProcessedLogScan float64 `bson:"records processed by log scan"`
	MaxLogSize              float64 `bson:"maximum log file size"`
	LogFlushes              float64 `bson:"log flush operations"`
	LogReads                float64 `bson:"log read operations"`
	LogScansDouble          float64 `bson:"log scan records requiring two reads"`
	LogScans                float64 `bson:"log scan operations"`
	LogSyncs                float64 `bson:"log sync operations"`
	LogSyncDirs             float64 `bson:"log sync_dir operations"`
	LogWrites               float64 `bson:"log write operations"`
}

// WTSessionStats session stats
type WTSessionStats struct {
	Cursors  float64 `bson:"open cursor count"`
	Sessions float64 `bson:"open session count"`
}

// WTTransactionStats transaction stats
type WTTransactionStats struct {
	Begins               float64 `bson:"transaction begins"`
	Checkpoints          float64 `bson:"transaction checkpoints"`
	CheckpointsRunning   float64 `bson:"transaction checkpoint currently running"`
	CheckpointMaxMs      float64 `bson:"transaction checkpoint max time (msecs)"`
	CheckpointMinMs      float64 `bson:"transaction checkpoint min time (msecs)"`
	CheckpointLastMs     float64 `bson:"transaction checkpoint most recent time (msecs)"`
	CheckpointTotalMs    float64 `bson:"transaction checkpoint total time (msecs)"`
	Committed            float64 `bson:"transactions committed"`
	CacheOverflowFailure float64 `bson:"transaction failures due to cache overflow"`
	RolledBack           float64 `bson:"transactions rolled back"`
}

// WTConcurrentTransactionsTypeStats concurrenttransaction stats
type WTConcurrentTransactionsTypeStats struct {
	Out          float64 `bson:"out"`
	Available    float64 `bson:"available"`
	TotalTickets float64 `bson:"totalTickets"`
}

// WTConcurrentTransactionsStats TODO:
type WTConcurrentTransactionsStats struct {
	Write *WTConcurrentTransactionsTypeStats `bson:"read"`
	Read  *WTConcurrentTransactionsTypeStats `bson:"write"`
}

// WiredTigerStats WiredTiger stats
type WiredTigerStats struct {
	BlockManager           *WTBlockManagerStats           `bson:"block-manager"`
	Cache                  *WTCacheStats                  `bson:"cache"`
	Log                    *WTLogStats                    `bson:"log"`
	Session                *WTSessionStats                `bson:"session"`
	Transaction            *WTTransactionStats            `bson:"transaction"`
	ConcurrentTransactions *WTConcurrentTransactionsStats `bson:"concurrentTransactions"`
}

// Cursors are the cursor metrics
type Cursors struct {
	TotalOpen      float64 `bson:"totalOpen"`
	TimeOut        float64 `bson:"timedOut"`
	TotalNoTimeout float64 `bson:"totalNoTimeout"`
	Pinned         float64 `bson:"pinned"`
}

// MemStats tracks the mem stats metrics.
type MemStats struct {
	Bits              float64 `bson:"bits"`
	Resident          float64 `bson:"resident"`
	Virtual           float64 `bson:"virtual"`
	Mapped            float64 `bson:"mapped"`
	MappedWithJournal float64 `bson:"mappedWithJournal"`
}

// DocumentStats are the stats associated to a document.
type DocumentStats struct {
	Deleted  float64 `bson:"deleted"`
	Inserted float64 `bson:"inserted"`
	Returned float64 `bson:"returned"`
	Updated  float64 `bson:"updated"`
}

// BenchmarkStats is bechmark info about an operation.
type BenchmarkStats struct {
	Num         float64 `bson:"num"`
	TotalMillis float64 `bson:"totalMillis"`
}

// GetLastErrorStats are the last error stats.
type GetLastErrorStats struct {
	Wtimeouts float64         `bson:"wtimeouts"`
	Wtime     *BenchmarkStats `bson:"wtime"`
}

// OperationStats are the stats for some kind of operations.
type OperationStats struct {
	Fastmod      float64 `bson:"fastmod"`
	Idhack       float64 `bson:"idhack"`
	ScanAndOrder float64 `bson:"scanAndOrder"`
}

// QueryExecutorStats are the stats associated with a query execution.
type QueryExecutorStats struct {
	Scanned        float64 `bson:"scanned"`
	ScannedObjects float64 `bson:"scannedObjects"`
}

// RecordStats are stats associated with a record.
type RecordStats struct {
	Moves float64 `bson:"moves"`
}

// ApplyStats are the stats associated with the apply operation.
type ApplyStats struct {
	Batches *BenchmarkStats `bson:"batches"`
	Ops     float64         `bson:"ops"`
}

// BufferStats are the stats associated with the buffer
type BufferStats struct {
	Count        float64 `bson:"count"`
	MaxSizeBytes float64 `bson:"maxSizeBytes"`
	SizeBytes    float64 `bson:"sizeBytes"`
}

// ReplExecutorStats are the stats associated with replication execution
type ReplExecutorStats struct {
	Counters         map[string]float64 `bson:"counters"`
	Queues           map[string]float64 `bson:"queues"`
	EventWaiters     float64            `bson:"eventWaiters"`
	UnsignaledEvents float64            `bson:"unsignaledEvents"`
}

// MetricsNetworkStats are the network stats.
type MetricsNetworkStats struct {
	Bytes          float64         `bson:"bytes"`
	Ops            float64         `bson:"ops"`
	GetMores       *BenchmarkStats `bson:"getmores"`
	ReadersCreated float64         `bson:"readersCreated"`
}

// ReplStats are the stats associated with the replication process.
type ReplStats struct {
	Apply        *ApplyStats          `bson:"apply"`
	Buffer       *BufferStats         `bson:"buffer"`
	Executor     *ReplExecutorStats   `bson:"executor,omitempty"`
	Network      *MetricsNetworkStats `bson:"network"`
	PreloadStats *PreloadStats        `bson:"preload"`
}

// PreloadStats are the stats associated with preload operation.
type PreloadStats struct {
	Docs    *BenchmarkStats `bson:"docs"`
	Indexes *BenchmarkStats `bson:"indexes"`
}

// StorageStats are the stats associated with the storage.
type StorageStats struct {
	BucketExhausted float64 `bson:"freelist.search.bucketExhausted"`
	Requests        float64 `bson:"freelist.search.requests"`
	Scanned         float64 `bson:"freelist.search.scanned"`
}

// CursorStatsOpen are the stats for open cursors
type CursorStatsOpen struct {
	NoTimeout float64 `bson:"noTimeout"`
	Pinned    float64 `bson:"pinned"`
	Total     float64 `bson:"total"`
}

// CursorStats are the stats for cursors
type CursorStats struct {
	TimedOut float64          `bson:"timedOut"`
	Open     *CursorStatsOpen `bson:"open"`
}

// MetricsStats are all stats associated with metrics of the system
type MetricsStats struct {
	Document      *DocumentStats      `bson:"document"`
	GetLastError  *GetLastErrorStats  `bson:"getLastError"`
	Operation     *OperationStats     `bson:"operation"`
	QueryExecutor *QueryExecutorStats `bson:"queryExecutor"`
	Record        *RecordStats        `bson:"record"`
	Repl          *ReplStats          `bson:"repl"`
	Storage       *StorageStats       `bson:"storage"`
	Cursor        *CursorStats        `bson:"cursor"`
}

// OpcountersStats opcounters stats
type OpcountersStats struct {
	Insert  float64 `bson:"insert"`
	Query   float64 `bson:"query"`
	Update  float64 `bson:"update"`
	Delete  float64 `bson:"delete"`
	GetMore float64 `bson:"getmore"`
	Command float64 `bson:"command"`
}

// OpcountersReplStats opcounters stats
type OpcountersReplStats struct {
	Insert  float64 `bson:"insert"`
	Query   float64 `bson:"query"`
	Update  float64 `bson:"update"`
	Delete  float64 `bson:"delete"`
	GetMore float64 `bson:"getmore"`
	Command float64 `bson:"command"`
}

//NetworkStats network stats
type NetworkStats struct {
	BytesIn     float64 `bson:"bytesIn"`
	BytesOut    float64 `bson:"bytesOut"`
	NumRequests float64 `bson:"numRequests"`
}

// ClientStats metrics for client stats
type ClientStats struct {
	Total   float64 `bson:"total"`
	Readers float64 `bson:"readers"`
	Writers float64 `bson:"writers"`
}

// QueueStats queue stats
type QueueStats struct {
	Total   float64 `bson:"total"`
	Readers float64 `bson:"readers"`
	Writers float64 `bson:"writers"`
}

// GlobalLockStats global lock stats
type GlobalLockStats struct {
	TotalTime     float64      `bson:"totalTime"`
	LockTime      float64      `bson:"lockTime"`
	Ratio         float64      `bson:"ratio"`
	CurrentQueue  *QueueStats  `bson:"currentQueue"`
	ActiveClients *ClientStats `bson:"activeClients"`
}

// TopStatsMap is a map of top stats.
type TopStatsMap map[string]TopStats

// TopCounterStats represents top counter stats.
type TopCounterStats struct {
	Time  float64 `bson:"time"`
	Count float64 `bson:"count"`
}

// TopStats top collection stats
type TopStats struct {
	Total     TopCounterStats `bson:"total"`
	ReadLock  TopCounterStats `bson:"readLock"`
	WriteLock TopCounterStats `bson:"writeLock"`
	Queries   TopCounterStats `bson:"queries"`
	GetMore   TopCounterStats `bson:"getmore"`
	Insert    TopCounterStats `bson:"insert"`
	Update    TopCounterStats `bson:"update"`
	Remove    TopCounterStats `bson:"remove"`
	Commands  TopCounterStats `bson:"commands"`
}

// TopStatus represents top metrics
type TopStatus struct {
	TopStats TopStatsMap `bson:"totals,omitempty"`
}

// DatabaseStatList contains stats from all databases
type DatabaseStatList struct {
	Members []DatabaseStatus
}

// FlushStats is the flush stats metrics
type FlushStats struct {
	Flushes      float64   `bson:"flushes"`
	TotalMs      float64   `bson:"total_ms"`
	AverageMs    float64   `bson:"average_ms"`
	LastMs       float64   `bson:"last_ms"`
	LastFinished time.Time `bson:"last_finished"`
}

// DatabaseStatus represents stats about a database (mongod and raw from mongos)
type DatabaseStatus struct {
	Name        string `bson:"db,omitempty"`
	IndexSize   int    `bson:"indexSize,omitempty"`
	DataSize    int    `bson:"dataSize,omitempty"`
	Collections int    `bson:"collections,omitempty"`
	Objects     int    `bson:"objects,omitempty"`
	Indexes     int    `bson:"indexes,omitempty"`
}

// ServerStatus keeps the data returned by the serverStatus() method.
type ServerStatus struct {
	Version            string               `bson:"version"`            //
	Uptime             float64              `bson:"uptime"`             //
	UptimeEstimate     float64              `bson:"uptimeEstimate"`     //
	LocalTime          time.Time            `bson:"localTime"`          //
	Asserts            *AssertsStats        `bson:"asserts"`            //
	Connections        *ConnectionStats     `bson:"connections"`        //
	ExtraInfo          *ExtraInfo           `bson:"extra_info"`         //
	GlobalLock         *GlobalLockStats     `bson:"globalLock"`         //
	Locks              LockStatsMap         `bson:"locks,omitempty"`    //
	Network            *NetworkStats        `bson:"network"`            //
	Opcounters         *OpcountersStats     `bson:"opcounters"`         //
	OpcountersRepl     *OpcountersReplStats `bson:"opcountersRepl"`     //
	StorageEngine      *StorageEngineStats  `bson:"storageEngine"`      //
	Mem                *MemStats            `bson:"mem"`                //
	Metrics            *MetricsStats        `bson:"metrics"`            //
	Dur                *DurStats            `bson:"dur"`                //
	BackgroundFlushing *FlushStats          `bson:"backgroundFlushing"` //
	IndexCounter       *IndexCounterStats   `bson:"indexCounters"`      //
	Cursors            *Cursors             `bson:"cursors"`            //
	InMemory           *WiredTigerStats     `bson:"inMemory"`           //
	RocksDb            *RocksDbStats        `bson:"rocksdb"`            //
	WiredTiger         *WiredTigerStats     `bson:"wiredTiger"`         //
}

// NewMongoDB TODO:
func NewMongoDB() *mgo.Session {
	db, err := mgo.Dial("localhost")
	if err != nil {
		return nil
	}
	return db
}

var mgoTags string

// MongoDBMetrics TODO:
func MongoDBMetrics() (L []*cmodel.MetricValue) {
	if g.Config().Collector.MongoDB == nil {
		return nil
	}
	if !g.Config().Collector.MongoDB.Enabled {
		return nil
	}
	mgoTags = fmt.Sprintf("host=%s,port=%d",
		g.Config().Collector.MongoDB.Host,
		g.Config().Collector.MongoDB.Port)
	session := NewMongoDB()
	if session == nil {
		return
	}
	defer session.Close()
	L = append(L, MongoDBStatInfo(session)...)
	return
}

// MongoDBStatInfo TODO:
func MongoDBStatInfo(session *mgo.Session) (L []*cmodel.MetricValue) {
	var serverStatus = GetServerStatus(session)
	fmt.Println(serverStatus.Version)
	fmt.Println(serverStatus.Uptime)
	return nil
}

// GetServerStatus returns the server status info.
func GetServerStatus(session *mgo.Session) *ServerStatus {
	result := &ServerStatus{}
	d := bson.D{{"serverStatus", 1}, {"recordStats", 1}}
	err := session.DB("admin").Run(d, result)
	if err != nil {
		log.Errorf("[E] Failed to get server status: %s", err)
		return nil
	}

	return result
}
