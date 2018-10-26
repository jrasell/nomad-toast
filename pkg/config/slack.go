package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SlackConfig is the nomad-toast Slack config struct.
type SlackConfig struct {
	AuthToken string
	Channel   string
}

const (
	configKeySlackAuthToken = "slack-auth-token"
	configKeySlackChannel   = "slack-channel"
)

// GetSlackConfig uses viper to populate a SlackConfig struct with values.
func GetSlackConfig() SlackConfig {
	return SlackConfig{
		AuthToken: viper.GetString(configKeySlackAuthToken),
		Channel:   viper.GetString(configKeySlackChannel),
	}
}

// RegisterSlackConfig is used by commands to register the Slack config flags.
func RegisterSlackConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeySlackAuthToken
			longOpt      = "slack-auth-token"
			defaultValue = ""
			description  = "The Slack API auth token"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeySlackChannel
			longOpt      = "slack-channel"
			defaultValue = ""
			description  = "The Slack channel to send notifications to"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

}
