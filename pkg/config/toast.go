package config

// ToastConfig is the overall config struct used for running nomad-toast.
type ToastConfig struct {
	Nomad *NomadConfig
	Slack *SlackConfig
}
