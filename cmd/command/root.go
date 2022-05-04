package command

import (
	c "github.com/achuchev/pump-monitor/cmd/common"
	"github.com/achuchev/pump-monitor/cmd/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	RootCmd = &cobra.Command{
		Use:               "pump-monitor",
		Short:             "CLI utility which helps to monitor the status of the pump",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
	}
)

// Execute executes the root command.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&c.FlagsConfig.ProfileName, c.FlagNameProfile, c.FlagNameProfileShort, c.DefaultValueProfile, "the name of the configuration profile")
	RootCmd.PersistentFlags().StringVar(&c.FlagsConfig.EmailAddressFrom, c.FlagNameEmailAddressFrom, "", "SMTP sender")
	RootCmd.PersistentFlags().StringVar(&c.FlagsConfig.EmailPassword, c.FlagNameEmailPassword, "", "SMTP sender's password")
	RootCmd.PersistentFlags().StringVar(&c.FlagsConfig.EmailAddressTo, c.FlagNameEmailAddressTo, "", "email address of the recipient")
	RootCmd.PersistentFlags().StringVar(&c.FlagsConfig.EmailHost, c.FlagNameEmailHost, c.DefaultValueEmailHost, "SMTP host")
	RootCmd.PersistentFlags().IntVar(&c.FlagsConfig.EmailHostPort, c.FlagNameEmailHostPort, c.DefaultValueEmailHostPort, "SMTP host port")

	RootCmd.PersistentFlags().StringVar(&c.FlagsConfig.ServiceURL, c.FlagNameServiceURL, c.DefaultValueServiceURL, "URL of the service")
	RootCmd.PersistentFlags().StringVar(&c.FlagsConfig.ServiceUsername, c.FlagNameServiceUsername, "", "username of the service")
	RootCmd.PersistentFlags().StringVar(&c.FlagsConfig.ServicePassword, c.FlagNameServicePassword, "", "password of the service")

	RootCmd.PersistentFlags().IntVar(&c.FlagsConfig.MonitorEveryMinutes, c.FlagNameMonitorEveryMinutes, c.DefaultValueMonitorEveryMinutes, "monitor every minutes")
	RootCmd.PersistentFlags().IntVar(&c.FlagsConfig.MonitorTempDifferenceThreshold, c.FlagNameMonitorTempDifferenceThreshold, c.DefaultValueMonitorTempDifferenceThreshold, "monitor temp difference threshold")

	RootCmd.AddCommand(debugCmd)
	RootCmd.AddCommand(showCmd)
}

func initConfig() {
	viper.AddConfigPath(utils.GetApplicationDir())
	viper.AddConfigPath("./")
	viper.SetConfigType("json")
	viper.SetEnvPrefix(c.ViperEnvPrefix)
	viper.SetConfigPermissions(0600)
	//viper.SetConfigName("pump-monitor.config")

	if c.FlagsConfig.ProfileName == "" {
		c.FlagsConfig.ProfileName = c.DefaultValueProfile
	}
	viper.SetConfigName(c.FlagsConfig.ProfileName)

	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using configuration from file: %v", viper.ConfigFileUsed())
	} else {
		log.Debugf("No configuration file found")
	}

	viper.AutomaticEnv()
	utils.BindFlags(RootCmd)
}
