package metrics

import (
	"github.com/mj0nez/restic-exporter/internal/info"
	promVersion "github.com/prometheus/client_golang/prometheus/collectors/version"

	"github.com/prometheus/client_golang/prometheus"
	promVersionInfo "github.com/prometheus/common/version"
)

var (
	commonLabels = [...]string{"repo", "client_hostname", "client_username", "client_version", "snapshot_hash", "snapshot_tag", "snapshot_tags", "snapshot_paths"}

	// metrics
	// this was restic_check_success a gauge, which was only increased on success,
	// which IMO neglects the case of an error
	CheckSuccess = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "restic_check_success",
		Namespace: "restic",
		Help:      "Number of successful integrity checks in the repository",
	}, []string{"repo"})
	// this metric is new compared to the python exporter
	CheckFailed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "restic_check_failed",
		Namespace: "restic",
		Help:      "Number of failed integrity checks in the repository",
	}, []string{"repo"})
	// this metric is new compared to the python exporter
	CheckSuggestRepairIndex = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "restic_check_suggest_repair_index",
		Namespace: "restic",
		Help:      "Metric indicating to repair the repository index",
	}, []string{"repo"})
	// this metric is new compared to the python exporter
	CheckSuggestPrune = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "restic_check_suggest_prune",
		Namespace: "restic",
		Help:      "Metric indicating to prune the repository",
	}, []string{"repo"})
	// this metric is new compared to the python exporter
	CheckErrorsTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "restic_errors_total",
		Namespace: "restic",
		Help:      "Total number of errors in the repository",
	}, []string{"repo"})

	LocksTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "restic_locks_total",
		Namespace: "restic",
		Help:      "Total number of locks in the repository",
	}, []string{"repo"})
	// this was a counter, but I assume with pruning the value might change
	// furthermore this avoids incrementing until we reach the current value
	SnapshotsTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "restic_snapshots_total",
		Namespace: "restic",
		Help:      "Total number of snapshots in the repository",
	}, []string{"repo"})
	BackupTimestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "restic_backup_timestamp",
		Namespace: "restic",
		Help:      "Timestamp of the last backup",
	}, commonLabels[:])
	BackupFilesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "restic_backup_files_total",
		Namespace: "restic",
		Help:      "Number of files in the backup",
	}, commonLabels[:])
	BackupSizeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "restic_backup_size_total",
		Namespace: "restic",
		Help:      "Total size of backup in bytes",
	}, commonLabels[:])
	BackupSnapshotsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "restic_backup_snapshots_total",
		Namespace: "restic",
		Help:      "Total number of snapshots",
	}, commonLabels[:])
	ScrapeDurationSeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "restic_scrape_duration_seconds",
		Namespace: "restic",
		Help:      "Amount of time each scrape takes",
	}, []string{"repo"})
)

func init() {
	promVersionInfo.Version = info.Version
	promVersionInfo.Revision = info.Revision
}

func NewRegistry() *prometheus.Registry {

	// Create non-global registry.
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		// collectors.NewGoCollector(),
		// collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		promVersion.NewCollector("restic"),
	)

	registerResticMetrics(reg)

	return reg
}

func registerResticMetrics(reg prometheus.Registerer) {
	reg.MustRegister(CheckSuccess)
	reg.MustRegister(LocksTotal)
	reg.MustRegister(SnapshotsTotal)
	reg.MustRegister(BackupTimestamp)
	reg.MustRegister(BackupFilesTotal)
	reg.MustRegister(BackupSizeTotal)
	reg.MustRegister(BackupSnapshotsTotal)
	reg.MustRegister(ScrapeDurationSeconds)
}
