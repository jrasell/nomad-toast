package allocations

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jrasell/nomad-toast/pkg/notifier"
	"github.com/jrasell/nomad-toast/pkg/watcher"
	"github.com/rs/zerolog/log"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgKeyAllocIncludeCStates = "include-states"
	cfgKeyAllocExcludeCStates = "exclude-states"
)

// RegisterCommand is used to register the deployments nomad-toast command.
func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "allocations",
		Short: "Run the nomad-toast allocations watcher and notifier",
		Run: func(cmd *cobra.Command, args []string) {
			runDeployments(cmd, args)
		},
	}

	flags := cmd.Flags()
	{
		const (
			key         = cfgKeyAllocIncludeCStates
			longOpt     = cfgKeyAllocIncludeCStates
			description = "Comma-separated list of allocation client states that will be whitelisted for notifications. \n" +
				"If specified, *only* these states will be included in notifications. The set of client states is defined by the \n" +
				"Nomad API - see the AllocClientStatus* constants in the docs: https://godoc.org/github.com/hashicorp/nomad/api#pkg-constants."
		)

		flags.StringSlice(longOpt, []string{}, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, []string{})
	}

	{
		const (
			key         = cfgKeyAllocExcludeCStates
			longOpt     = cfgKeyAllocExcludeCStates
			description = "Comma-separated list of allocation client states that will be excluded from notifications. \n" +
				"This takes priority over include-states. The set of client states is defined by the \n" +
				"Nomad API - see the AllocClientStatus* constants in the docs: https://godoc.org/github.com/hashicorp/nomad/api#pkg-constants."
		)

		flags.StringSlice(longOpt, []string{}, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, []string{})
	}

	rootCmd.AddCommand(cmd)

	return nil
}

func runDeployments(_ *cobra.Command, _ []string) {

	cfg, err := getConfig()
	if err != nil {
		log.Error().Err(err).Msg("unable to load nomad-toast config")
		os.Exit(sysexits.Software)
	}

	n, err := notifier.NewNotifier(cfg.Slack, cfg.UI, watcher.Deployments)
	if err != nil {
		log.Error().Err(err).Msg("unable to build new allocations notifier")
		os.Exit(sysexits.Software)
	}

	w, region, err := watcher.NewWatcher(cfg.Nomad, watcher.Allocations, n.MsgChan)
	if err != nil {
		log.Error().Err(err).Msg("unable to build new allocations watcher")
		os.Exit(sysexits.Software)
	}

	go n.Run(region)
	go w.Run()

	sigCh := make(chan os.Signal, 5)

	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	for {
		sig := <-sigCh
		log.Info().
			Str("signal", sig.String()).
			Msg("received OS signal to shutdown")
		os.Exit(sysexits.OK)
	}
}
