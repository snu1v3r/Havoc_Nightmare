package cmd

import (
	"Havoc/pkg/logger"
	"fmt"
	"os"
	"time"

	dbcmd "Havoc/cmd/db"
	"Havoc/cmd/server"
	"Havoc/pkg/colors"
	"github.com/spf13/cobra"
)

var (
	serverFlags   server.Flags
	databaseFlags dbcmd.Flags

	HavocCmd = &cobra.Command{
		Use:          "havoc",
		Short:        fmt.Sprintf("Havoc Framework [Version: %v] [CodeName: %v]", server.Version, server.CodeName),
		SilenceUsage: true,
		RunE:         havocRun,
	}

	cliServer = &cobra.Command{
		Use:          "server",
		Short:        "team server command",
		SilenceUsage: true,
		RunE:         serverRun,
	}

	cliDatabase = &cobra.Command{
		Use:          "database",
		Short:        "interact with the havoc database",
		SilenceUsage: true,
		RunE:         databaseRun,
	}
)

// init all flags
func init() {
	HavocCmd.CompletionOptions.DisableDefaultCmd = true

	// server command flags
	cliServer.Flags().SortFlags = false
	cliServer.Flags().StringVarP(&serverFlags.Profile, "profile", "", "", "set havoc teamserver profile")
	cliServer.Flags().BoolVarP(&serverFlags.Debug, "debug", "", false, "enable debug mode")
	cliServer.Flags().BoolVarP(&serverFlags.Default, "default", "d", false, "uses default profile (overwrites --profile)")

	// database command flags
	cliDatabase.Flags().SortFlags = false

	cliDatabase.Flags().BoolVarP(&databaseFlags.ClearAll, "clear-all", "", false, "clear all data from the database")

	cliDatabase.Flags().BoolVarP(&databaseFlags.Agent.List, "agent-list", "", false, "list all registered agents from the database")
	cliDatabase.Flags().StringVarP(&databaseFlags.Agent.Remove, "agent-remove", "", "", "remove an agent uuid from the database")
	cliDatabase.Flags().BoolVarP(&databaseFlags.Agent.Clear, "agent-clear", "", false, "clear all registered agents from the database")

	cliDatabase.Flags().BoolVarP(&databaseFlags.Listener.List, "listener-list", "", false, "list all available listeners from the database")
	cliDatabase.Flags().StringVarP(&databaseFlags.Listener.Remove, "listener-remove", "", "", "remove a listener from the database")
	cliDatabase.Flags().BoolVarP(&databaseFlags.Listener.Clear, "listener-clear", "", false, "clear all listeners from the database")

	// add commands to the teamserver cli
	HavocCmd.Flags().SortFlags = false
	HavocCmd.AddCommand(cliServer)
	HavocCmd.AddCommand(cliDatabase)
}

func havocRun(cmd *cobra.Command, args []string) error {
	if len(os.Args) <= 2 {
		err := cmd.Help()
		if err != nil {
			return err
		}
		os.Exit(0)
	}

	return nil
}

func serverRun(cmd *cobra.Command, args []string) error {
	var (
		DirPath, _  = os.Getwd()
		ServerTimer = time.Now()
		Server      *server.Teamserver
		err         error
	)

	if len(os.Args) <= 2 {
		err = cmd.Help()
		if err != nil {
			return err
		}
		os.Exit(0)
	}

	if Server = server.NewTeamserver(); Server == nil {
		logger.Error("failed to create server")
		return nil
	}

	Server.SetFlags(serverFlags)

	if serverFlags.Debug {
		logger.SetDebug(true)
		logger.Debug("debug mode enabled")
	}

	logger.Info("%v [Version: %v %v]", colors.BoldWhite("Havoc Framework"), server.Version, colors.BoldBlue(server.CodeName))

	if serverFlags.Default {
		err = Server.Profile(DirPath + "/data/havoc.toml")
		if err != nil {
			return nil
		}
	} else if serverFlags.Profile != "" {
		err = Server.Profile(serverFlags.Profile)
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
}

func databaseRun(cmd *cobra.Command, args []string) error {
	var (
		err      error
		database *dbcmd.DatabaseCli
	)

	if len(os.Args) <= 2 {
		err = cmd.Help()
		if err != nil {
			return err
		}
		return nil
	}

	if database = dbcmd.NewDatabaseCli(); database == nil {
		return nil
	}

	if databaseFlags.ClearAll {
		_ = database.AgentClear()
		_ = database.ListenerClear()
	} else if databaseFlags.Agent.List {
		_ = database.AgentList()
	} else if databaseFlags.Agent.Clear {
		_ = database.AgentClear()
	} else if len(databaseFlags.Agent.Remove) != 0 {
		_ = database.AgentRemove(databaseFlags.Agent.Remove)
	} else if databaseFlags.Listener.List {
		_ = database.ListenerList()
	} else if databaseFlags.Listener.Clear {
		_ = database.ListenerClear()
	} else if len(databaseFlags.Listener.Remove) != 0 {
		_ = database.ListenerRemove(databaseFlags.Listener.Remove)
	}

	return nil
}
