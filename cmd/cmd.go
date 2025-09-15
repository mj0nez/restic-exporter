package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mj0nez/restic-exporter/internal"
	"github.com/mj0nez/restic-exporter/internal/collector"
	"github.com/mj0nez/restic-exporter/internal/config"
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
	cfgCmd := &cobra.Command{
		Use: "config",
		Run: exportDefault,
	}
	cfgCmd.Flags().Bool("defaults", false, "get only the default config values")
	rootCmd.AddCommand(cfgCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func runServer(cmd *cobra.Command, args []string) {
	conf := config.MustLoadConfig(false)

	metricsRegistry := metrics.NewRegistry()
	server := internal.NewHttpServer(conf.Server.Addr, internal.NewRouter(metricsRegistry), internal.NewHttpServerOpts())
	err := internal.RunServer(server)

	if err != nil {
		fmt.Printf("%v", err)
	}

}

func runCollect(cmd *cobra.Command, args []string) {

	collector.Collect()
}

func exportDefault(cmd *cobra.Command, args []string) {

	onlyDefaults, err := cmd.Flags().GetBool("defaults")

	if err != nil {
		fmt.Println("Usage error in config flags")
	}

	// view the configuration
	conf := config.MustLoadConfig(onlyDefaults)
	confb, _ := json.MarshalIndent(*conf, "", "\t")
	fmt.Println(string(confb))

}
