package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	defaultIntervalCheck    int32 = 3600
	defaultIntervalSnapshot int32 = 60
)

// Global variable which allows sharing a cached Config instance.
var globalConfig *Config = nil

type (
	defaultValues = map[string]any
	defaultConfig = map[string]defaultValues
)

type LoggingConfig struct {
	Level string
}

var loggingConfigDefaults = defaultValues{"level": "info"}

type AppConfig struct {
	BinaryPath string `mapstructure:"binary_path"`
	Prefetch   bool
	// ShutdownTimeout int32 `mapstructure:"shutdown_timeout"` // in seconds
}

var appConfigDefaults = defaultValues{"binary_path": "/usr/bin/restic", "prefetch": true}

type ServerConfig struct {
	Addr            string
	ShutdownTimeout int32 `mapstructure:"shutdown_timeout"` // in seconds
}

var serverConfigDefaults = defaultValues{"addr": "0.0.0.0:8081"}

type CollectionIntervalsConfig struct {
	Check    int32
	Snapshot int32
}

type ResticConfig struct {
	Repo     string
	Password string
}

type Repository struct {
	Name                string
	Restic              ResticConfig
	CollectionIntervals CollectionIntervalsConfig `mapstructure:"collection_intervals"` // in seconds
}

type Config struct {
	Logging      LoggingConfig
	App          AppConfig
	Server       ServerConfig
	Repositories []Repository
}

func applyDefaults(vp *viper.Viper) *viper.Viper {

	defaults := defaultConfig{
		"logging": loggingConfigDefaults,
		"app":     appConfigDefaults,
		"server":  serverConfigDefaults,
	}

	for confKey, confValues := range defaults {
		for valKey, value := range confValues {
			viperKey := fmt.Sprintf("%v.%v", confKey, valKey)
			vp.SetDefault(viperKey, value)
		}

	}
	return vp
}

// allows setting of defaults for our array configs, namely the repositories
func postProcessConfig(conf *Config) {
	for r, repo := range conf.Repositories {

		collectionConfig := &repo.CollectionIntervals

		switch {
		case collectionConfig.Check == 0:
			collectionConfig.Check = defaultIntervalCheck
		case collectionConfig.Check < 0:
			// ensure downstream can process this
			collectionConfig.Check = -1
		}

		switch {
		case collectionConfig.Snapshot == 0:
			collectionConfig.Snapshot = defaultIntervalSnapshot
		case collectionConfig.Snapshot < 0:
			// ensure downstream can process this
			collectionConfig.Snapshot = -1
		}

		conf.Repositories[r] = Repository{
			Name:                repo.Name,
			Restic:              repo.Restic,
			CollectionIntervals: repo.CollectionIntervals,
		}
	}

}

func LoadConfig(onlyDefaults bool) (*Config, error) {

	if globalConfig != nil && !onlyDefaults {
		return globalConfig, nil
	}

	vp := NewViper(onlyDefaults)
	applyDefaults(vp)

	globalConfig = &Config{}
	if err := vp.Unmarshal(globalConfig); err != nil {
		return globalConfig, err
	}

	postProcessConfig(globalConfig)

	return globalConfig, nil
}

func MustLoadConfig(onlyDefaults bool) *Config {

	config, err := LoadConfig(onlyDefaults)

	if err != nil {
		fmt.Printf("Unable to load configuration because of %+v \n", err)
		os.Exit(1)
	}

	return config
}

func ExportDefaultConfig() error {
	vp := NewViper(true)
	applyDefaults(vp)
	var buffer bytes.Buffer
	if err := vp.WriteConfigTo(&buffer); err != nil {
		return err
	}

	fmt.Println(buffer.String())
	return nil

}
