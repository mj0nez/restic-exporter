package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

const (
	configFileName   = "config"
	configType       = "yaml"
	configPathEnvKey = "RESTIC_EXPORTER_CONFIG_PATH"
)

func NewViper(skipLoading bool) *viper.Viper {
	// Enable BindStruct to allow unmarshal env into a nested struct
	// this feature is experimental and requires v1.20.0-alpha.6
	vp := viper.NewWithOptions(viper.ExperimentalBindStruct())

	// replace - & . by _ for environment variable names and nested keys used in defaults
	// (eg: the env var for tls-server-name or tls.server.name is TLS_SERVER_NAME)
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	configPath := os.Getenv(configPathEnvKey)

	// file config
	// load config from env path or use common defaults
	if configPath != "" {
		vp.SetConfigFile(configPath)

	} else {
		// add common paths
		vp.AddConfigPath("/etc/restic-exporter/")

		if userHome, homeErr := os.UserHomeDir(); homeErr != nil {
			vp.AddConfigPath(path.Join(userHome, ".restic-exporter"))
		}

		vp.SetConfigName(configFileName)
		vp.SetConfigType(configType)

	}

	// return before loading anything
	if skipLoading {
		return vp
	}

	if err := vp.ReadInConfig(); err != nil {
		fmt.Printf("Reading config encountered an error:  %v\n", err)
	}

	// environment config
	vp.SetEnvPrefix("restic-exporter") // variables will be prefixed as RESTIC_EXPORTER_

	vp.AutomaticEnv() // read in env vars
	return vp
}

// func resolvePath(configPath string) (string, string, string) {
// 	dir, fileName := filepath.Split(configPath)
// 	ext := filepath.Ext(fileName)

// 	return dir, strings.Split(fileName, ".")[0], ext[1:]
// }
