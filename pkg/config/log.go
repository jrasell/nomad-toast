package config

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LogConfig is the nomad-toast logging config struct.
type LogConfig struct {
	LogLevel  string
	LogFormat string
	UseColor  bool
}

const (
	configKeyLogLevel  = "log-level"
	configKeyLogFormat = "log-format"
	configKeyUseColor  = "log-use-color"
)

// GetLogConfig uses viper to populate a LogConfig struct with values.
func GetLogConfig() LogConfig {
	return LogConfig{
		LogLevel:  viper.GetString(configKeyLogLevel),
		LogFormat: viper.GetString(configKeyLogFormat),
		UseColor:  viper.GetBool(configKeyUseColor),
	}
}

// RegisterLogConfig is used by commands to register the log config flags.
func RegisterLogConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyLogLevel
			longOpt      = "log-level"
			defaultValue = "info"
			description  = "Change the log level being sent to stderr"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyLogFormat
			longOpt      = "log-format"
			defaultValue = "auto"
			description  = `Specify the log format ("auto", "zerolog" or "human")`
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = configKeyUseColor
			longOpt     = "log-use-color"
			description = "Use ANSI colors in logging output"
		)
		defaultValue := false
		if isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd()) {
			defaultValue = true
		}

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
