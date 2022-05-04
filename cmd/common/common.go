package common

const (
	FlagNameProfile      = "profile"
	FlagNameProfileShort = "p"

	FlagNameEmailAddressFrom = "email-address-from"
	FlagNameEmailAddressTo   = "email-address-to"
	FlagNameEmailPassword    = "email-password"
	FlagNameEmailHost        = "email-host"
	FlagNameEmailHostPort    = "email-host-port"

	FlagNameServiceURL      = "service-url"
	FlagNameServiceUsername = "service-username"
	FlagNameServicePassword = "service-password"

	FlagNameMonitorEveryMinutes            = "monitor-every-minutes"
	FlagNameMonitorTempDifferenceThreshold = "monitor-temp-difference-threshold"
	ViperEnvPrefix                         = "PUMP"
)

const (
	DefaultValueProfile       = "pump-monitor.conf"
	DefaultValueEmailHost     = "smtp.gmail.com"
	DefaultValueEmailHostPort = 587

	DefaultValueServiceURL = "https://remotecontrol.at/index.php/"

	DefaultValueMonitorEveryMinutes            = 5
	DefaultValueMonitorTempDifferenceThreshold = 12
)

var (
	FlagsConfig CommandConfigFlags
)

type CommandConfigFlags struct {
	ProfileName string

	EmailAddressFrom string
	EmailAddressTo   string
	EmailPassword    string
	EmailHost        string
	EmailHostPort    int

	ServiceURL      string
	ServiceUsername string
	ServicePassword string

	MonitorEveryMinutes            int
	MonitorTempDifferenceThreshold int
}
