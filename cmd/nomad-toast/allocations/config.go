package allocations

import (
	"github.com/jrasell/nomad-toast/pkg/config"
	"github.com/spf13/viper"
)

func getConfig() (*config.ToastConfig, error) {
	nomadConfig := config.GetNomadConfig()
	nomadConfig.AllocCfg = getAllocConfig()
	slackConfig := config.GetSlackConfig()
	uiConfig := config.GetUIConfig()

	return &config.ToastConfig{
		Nomad: &nomadConfig,
		Slack: &slackConfig,
		UI:    &uiConfig,
	}, nil
}

func getAllocConfig() *config.AllocConfig {
	return &config.AllocConfig{
		IncludeStates: viper.GetStringSlice(cfgKeyAllocIncludeCStates),
		ExcludeStates: viper.GetStringSlice(cfgKeyAllocExcludeCStates),
	}
}
