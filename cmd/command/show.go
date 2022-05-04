package command

import (
	"reflect"

	c "github.com/achuchev/pump-monitor/cmd/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
}

var (
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows the configuration",
		Long:  `...`,
		RunE:  doShow,
	}
)

func doShow(cmd *cobra.Command, args []string) error {
	settings := viper.AllSettings()
	log.Infof("Showing settings of '%s' profile", c.FlagsConfig.ProfileName)
	for key, value := range settings {
		if reflect.ValueOf(value).Kind() == reflect.Map {
			log.Infof("%s:", key)
			nestedMap := value.(map[string]interface{})
			for k, v := range nestedMap {
				log.Infof("    %s: %v\n", k, v)
			}
		} else {
			log.Infof("  %v: %v\n", key, value)
		}
	}
	return nil
}
