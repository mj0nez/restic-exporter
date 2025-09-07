package cmd

import (
	"fmt"

	"github.com/mj0nez/restic-exporter/internal"
	"github.com/mj0nez/restic-exporter/internal/info"
	"github.com/mj0nez/restic-exporter/internal/metrics"
	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command = &cobra.Command{
		Use:     "restic-exporter",
		Short:   "A Prometheus metrics exporter for restic.",
		Long:    fmt.Sprintf("A Prometheus metrics exporter for restic. Version: %s Revision: %s", info.Version, info.Revision),
		Version: info.Version,
	}
	// rootLogger *zap.Logger = logger.New("setup")
)

func init() {

	rootCmd.AddCommand(&cobra.Command{
		Use: "server",
		Run: runServer,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "collect",
		Run: runCollect,
	})
}

func Execute() error {
	return rootCmd.Execute()
}

func runServer(cmd *cobra.Command, args []string) {

	metricsRegistry := metrics.NewRegistry()

	server := internal.NewHttpServer("0.0.0.0:8081", internal.NewRouter(metricsRegistry), internal.NewHttpServerOpts())
	err := internal.RunServer(server)

	if err != nil {
		fmt.Printf("%v", err)
	}

}

func runCollect(cmd *cobra.Command, args []string) {

	internal.Collect()
}
