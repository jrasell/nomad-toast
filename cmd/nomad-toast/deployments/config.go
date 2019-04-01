package deployments

import (
	"github.com/jrasell/nomad-toast/pkg/config"
)

func getConfig() (*config.ToastConfig, error) {
	nomadConfig := config.GetNomadConfig()
	slackConfig := config.GetSlackConfig()
	footerConfig := config.GetUIConfig()

	return &config.ToastConfig{
		Nomad: &nomadConfig,
		Slack: &slackConfig,
		UI:    &footerConfig,
	}, nil
}
