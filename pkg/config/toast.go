package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ToastConfig is the overall config struct used for running nomad-toast.
type ToastConfig struct {
	Nomad *NomadConfig
	Slack *SlackConfig
	UI    *UI
}

const (
	configKeyHashiUIEnabled = "hashiui-enabled"
	configKeyHashiUIHost    = "hashiui-host"
)

// UI represents the configuration to provides links to a Nomad UI inside a notification.
type UI struct {
	HashiUIEnabled bool
	HashiUIHost    string
}

// GetUIConfig uses viper to populate a UI config struct with values.
func GetUIConfig() UI {
	return UI{
		HashiUIEnabled: viper.GetBool(configKeyHashiUIEnabled),
		HashiUIHost:    viper.GetString(configKeyHashiUIHost),
	}
}

// RegisterUIConfig is used by commands to register the UI config flags.
func RegisterUIConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyHashiUIEnabled
			longOpt      = "hashiui-enabled"
			defaultValue = false
			description  = "This flag enables adding a link to HashiUI within a notification."
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyHashiUIHost
			longOpt      = "hashiui-host"
			defaultValue = "http://localhost:8000"
			description  = "The base URL where the HashiUI lives."
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
