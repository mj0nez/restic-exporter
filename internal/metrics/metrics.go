package metrics

import (
	"github.com/prometheus/client_golang/prometheus/collectors"
	promVersion "github.com/prometheus/client_golang/prometheus/collectors/version"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	commonLabels = [...]string{"repo", "client_hostname", "client_username", "client_version", "snapshot_hash", "snapshot_tag", "snapshot_tags", "snapshot_paths"}

	// metrics
	CheckSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "restic_check_success",
		Namespace: "restic",
		Help:      "Result of restic check operation in the repository",
	}, []string{"repo"})
	LocksTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "restic_locks_total",
		Namespace: "restic",
		Help:      "Total number of locks in the repository",
	}, []string{"repo"})
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

func NewRegistry() *prometheus.Registry {

	// Create non-global registry.
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		// TODO consider re-implementing this
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
