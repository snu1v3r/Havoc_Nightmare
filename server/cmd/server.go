package cmd

import (
	"os"
	"time"

	"Havoc/cmd/server"
	"Havoc/pkg/colors"
	"Havoc/pkg/logger"

	"github.com/spf13/cobra"
)

var CobraServer = &cobra.Command{
	Use:          "server",
	Short:        "team server command",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			DirPath, _  = os.Getwd()
			ServerTimer = time.Now()
			Server      *server.Teamserver
			err         error
		)

		if len(os.Args) <= 2 {
			err := cmd.Help()
			if err != nil {
				return err
			}
			os.Exit(0)
		}

		if Server = server.NewTeamserver(); Server == nil {
			logger.Error("failed to create server")
			return nil
		}

		Server.SetFlags(flags)

		startMenu()

		if flags.Server.Debug {
			logger.SetDebug(true)
			logger.Debug("debug mode enabled")
		}

		logger.Info("%v [Version: %v %v]", colors.BoldWhite("Havoc Framework"), server.Version, colors.BoldBlue(server.CodeName))

		if flags.Server.Default {
			err = Server.Profile(DirPath + "/data/havoc.toml")
			if err != nil {
				return nil
			}
		} else if flags.Server.Profile != "" {
			err = Server.Profile(flags.Server.Profile)
			if err != nil {
				return nil
			}
		} else {
			logger.Error("no profile specified")
			logger.Error("specify a profile with --profile or choose the standard profile with --default")
			os.Exit(1)
		}

		logger.Info("time: " + colors.Yellow(ServerTimer.Format("02/01/2006 15:04:05")))

		Server.Start()

		return nil
	},
}
