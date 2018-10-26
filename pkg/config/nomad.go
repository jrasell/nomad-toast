package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NomadConfig is the nomad-toast Nomad client config struct.
type NomadConfig struct {
	AllowStale   bool
	NomadAddress string
}

const (
	configKeyNomadAllowStale = "nomad-allow-stale"
	configKeyNomadAddress    = "nomad-address"
)

// GetNomadConfig uses viper to populate a NomadConfig struct with values.
func GetNomadConfig() NomadConfig {
	return NomadConfig{
		AllowStale:   viper.GetBool(configKeyNomadAllowStale),
		NomadAddress: viper.GetString(configKeyNomadAddress),
	}
}

// RegisterNomadConfig is used by commands to register the Nomad config flags.
func RegisterNomadConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyNomadAllowStale
			longOpt      = "nomad-allow-stale"
			defaultValue = true
			description  = "Allow stale Nomad consistency when making API calls"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyNomadAddress
			longOpt      = "nomad-address"
			defaultValue = "http://localhost:4646"
			description  = "The Nomad HTTP(S) API address"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

}
