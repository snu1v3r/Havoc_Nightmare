package cmd

import (
	"fmt"
	"os"

	"Havoc/cmd/server"
	"Havoc/pkg/colors"

	"github.com/spf13/cobra"
)

var (
	HavocCli = &cobra.Command{
		Use:          "havoc",
		Short:        fmt.Sprintf("Havoc Framework [Version: %v] [CodeName: %v]", server.Version, server.CodeName),
		SilenceUsage: true,
		RunE:         serverRun,
	}

	flags server.TeamserverFlags
)

// init all flags
func init() {
	HavocCli.CompletionOptions.DisableDefaultCmd = true

	// server flags
	CobraServer.Flags().SortFlags = false
	CobraServer.Flags().StringVarP(&flags.Server.Profile, "profile", "", "", "set havoc teamserver profile")
	CobraServer.Flags().BoolVarP(&flags.Server.Debug, "debug", "", false, "enable debug mode")
	CobraServer.Flags().BoolVarP(&flags.Server.Default, "default", "d", false, "uses default profile (overwrites --profile)")

	// add commands to the teamserver cli
	HavocCli.Flags().SortFlags = false
	HavocCli.AddCommand(CobraServer)
}

func serverRun(cmd *cobra.Command, args []string) error {
	startMenu()

	if len(os.Args) <= 2 {
		err := cmd.Help()
		if err != nil {
			return err
		}
		os.Exit(0)
	}

	return nil
}

func startMenu() {
	fmt.Print(colors.Red("              _______           _______  _______ \n    │\\     /│(  ___  )│\\     /│(  ___  )(  ____ \\\n    │ )   ( ││ (   ) ││ )   ( ││ (   ) ││ (    \\/\n    │ (___) ││ (___) ││ │   │ ││ │   │ ││ │      \n    │  ___  ││  ___  │( (   ) )│ │   │ ││ │      \n    │ (   ) ││ (   ) │ \\ \\_/ / │ │   │ ││ │      \n    │ )   ( ││ )   ( │  \\   /  │ (___) ││ (____/\\\n    │/     \\││/     \\│   \\_/   (_______)(_______/"))
	fmt.Printf(" by %v\n\n", colors.BoldBlue("@C5pider"))
	fmt.Println("  	", colors.Red("pwn"), "and", colors.Blue("elevate"), "until it's done\n")
}
