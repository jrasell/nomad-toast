package main

import (
	"fmt"
	"os"

	"github.com/jrasell/nomad-toast/cmd/nomad-toast/allocations"
	"github.com/jrasell/nomad-toast/cmd/nomad-toast/deployments"
	"github.com/jrasell/nomad-toast/pkg/buildconsts"
	"github.com/jrasell/nomad-toast/pkg/config"
	"github.com/jrasell/nomad-toast/pkg/logger"
	"github.com/pkg/errors"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "nomad-toast",
		Short: "\nnomad-toast in an open source tool for receiving notifications based on HashiCorp Nomad events.\n" +
			"It is designed to increase observability throughout an organisation and provide insights within\n" +
			"chatops style environments.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			buildconsts.SetProgramName(cmd.Use)

			if err := logger.Setup(config.GetLogConfig()); err != nil {
				return errors.Wrap(err, "error configuring logging")
			}

			return nil
		},
		Version: buildconsts.GetVersion(),
	}

	config.RegisterLogConfig(rootCmd)
	config.RegisterNomadConfig(rootCmd)
	config.RegisterSlackConfig(rootCmd)
	config.RegisterUIConfig(rootCmd)

	if err := registerCommands(rootCmd); err != nil {
		fmt.Println("error registering commands:", err)
		os.Exit(sysexits.Software)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(sysexits.Software)
	}
}

func registerCommands(rootCmd *cobra.Command) error {
	if err := deployments.RegisterCommand(rootCmd); err != nil {
		return err
	}

	if err := allocations.RegisterCommand(rootCmd); err != nil {
		return err
	}

	return nil
}
