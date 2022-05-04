package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	c "github.com/achuchev/pump-monitor/cmd/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func GetURL(suffix string) string {
	return c.FlagsConfig.ServiceURL + suffix
}

// BindFlags Bind each cobra flag to its associated viper configuration (config file and environment variable)
func BindFlags(cmd *cobra.Command) {
	var envFlagName string
	var envNameChanged bool
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		envNameChanged = false
		envFlagName = strings.ToUpper(f.Name)

		// If command is present then add double underscore to the flag name
		if strings.Contains(envFlagName, ".") {
			envFlagName = strings.Replace(envFlagName, ".", "__", 1)
			envNameChanged = true
		}

		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores
		if strings.Contains(envFlagName, "-") {
			envFlagName = strings.ReplaceAll(envFlagName, "-", "_")
			envNameChanged = true
		}

		if envNameChanged {
			bindStr := fmt.Sprintf("%s_%s", c.ViperEnvPrefix, envFlagName)
			err := viper.BindEnv(f.Name, bindStr)
			if err != nil {
				log.Errorf("Environment variable '%s' could not be binded to '%s'. Error: %s", f.Name, bindStr, err.Error())
			}
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				log.Errorf("Default value '%s' could not be set to flag '%s'. Error: %s", val, f.Name, err.Error())
			}
		}
	})
}

func GetApplicationDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	log.Debug("Application dir: %s", exPath)
	return exPath
}
