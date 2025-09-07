package internal

import (
	"github.com/prometheus/client_golang/prometheus/collectors"
	promVersion "github.com/prometheus/client_golang/prometheus/collectors/version"
	promVersionInfo "github.com/prometheus/common/version"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricsPath = "/metrics"
)

func init() {
	// pass build information to prometheus - avoids additional linker flags
	promVersionInfo.Version = Version
	promVersionInfo.Revision = Revision
}

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

	return reg
}
