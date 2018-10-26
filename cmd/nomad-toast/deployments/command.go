package deployments

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jrasell/nomad-toast/pkg/notifier"
	"github.com/jrasell/nomad-toast/pkg/watcher"
	"github.com/rs/zerolog/log"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

// RegisterCommand is used to register the deployments nomad-toast command.
func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "Run the nomad-toast deployments watcher and notifier",
		Run: func(cmd *cobra.Command, args []string) {
			runDeployments(cmd, args)
		},
	}

	rootCmd.AddCommand(cmd)

	return nil
}

func runDeployments(_ *cobra.Command, _ []string) {

	cfg, err := getConfig()
	if err != nil {
		log.Error().Err(err).Msg("unable to load deployment config")
		os.Exit(sysexits.Software)
	}

	n, err := notifier.NewNotifier(cfg.Slack)
	if err != nil {
		log.Error().Err(err).Msg("unable to build new deployments notifier")
		os.Exit(sysexits.Software)
	}

	w, err := watcher.NewWatcher(cfg.Nomad, watcher.Deployments, n.MsgChan)
	if err != nil {
		log.Error().Err(err).Msg("unable to build new deployments watcher")
		os.Exit(sysexits.Software)
	}

	go n.Run()
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
