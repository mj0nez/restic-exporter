package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Global variable which allows sharing a cached Config instance.
var globalConfig *Config = nil

type (
	defaultValues = map[string]any
	defaultConfig = map[string]defaultValues
)

type ServerConfig struct {
	Addr            string
	ShutdownTimeout int32 `mapstructure:"shutdown_timeout"` // in seconds
}

var serverConfigDefaults = defaultValues{"addr": "0.0.0.0:8081"}

type LoggingConfig struct {
	Level string
}

var loggingConfigDefaults = defaultValues{"level": "info"}

type Repository struct {
	Name          string
	Url           string
	CheckInterval uint `mapstructure:"check_interval"` // in seconds
}

type Config struct {
	Logging      LoggingConfig
	Server       ServerConfig
	Repositories []Repository
}

func applyDefaults(vp *viper.Viper) *viper.Viper {

	defaults := defaultConfig{
		"logging": loggingConfigDefaults,
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
		if repo.CheckInterval == 0 {
			conf.Repositories[r] = Repository{
				Name:          repo.Name,
				Url:           repo.Url,
				CheckInterval: 60,
			}
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
