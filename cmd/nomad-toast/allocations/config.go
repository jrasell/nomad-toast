package allocations

import (
	"github.com/jrasell/nomad-toast/pkg/config"
)

func getConfig() (*config.ToastConfig, error) {
	nomadConfig := config.GetNomadConfig()
	slackConfig := config.GetSlackConfig()

	return &config.ToastConfig{
		Nomad: &nomadConfig,
		Slack: &slackConfig,
	}, nil
}
