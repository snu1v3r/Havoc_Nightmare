package cmd

import (
	"fmt"
	"os"
	"time"

	"Havoc/cmd/db"
	"Havoc/cmd/plugin"
	"Havoc/cmd/server"
	"Havoc/pkg/colors"
	"Havoc/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	serverFlags   server.Flags
	databaseFlags db.Flags
	pluginFlags   plugin.Flags

	HavocCmd = &cobra.Command{
		Use:          "havoc",
		Short:        fmt.Sprintf("Havoc Framework [%v %v]", server.Version, server.CodeName),
		SilenceUsage: true,
		Run:          havocRun,
	}

	cliServer = &cobra.Command{
		Use:          "server",
		Short:        "teamserver command",
		SilenceUsage: true,
		Run:          serverRun,
	}

	cliDatabase = &cobra.Command{
		Use:          "database",
		Short:        "database manager command",
		SilenceUsage: true,
		Run:          databaseRun,
	}

	cliPlugin = &cobra.Command{
		Use:          "plugin",
		Short:        "plugin manager command",
		SilenceUsage: true,
		Run:          pluginRun,
	}
)

// init all flags
func init() {
	HavocCmd.CompletionOptions.DisableDefaultCmd = true

	// hide the help subcommand
	HavocCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// server command flags
	cliServer.Flags().SortFlags = false
	cliServer.Flags().StringVarP(&serverFlags.Profile, "profile", "", "", "set havoc teamserver profile")
	cliServer.Flags().BoolVarP(&serverFlags.Debug, "debug", "", false, "enable debug mode")
	cliServer.Flags().BoolVarP(&serverFlags.Default, "default", "d", false, "uses default profile (overwrites --profile)")

	// plugin command flags
	cliPlugin.Flags().SortFlags = false
	cliPlugin.Flags().StringVarP(&pluginFlags.Register, "register", "", "", "register plugin path")
	cliPlugin.Flags().StringVarP(&pluginFlags.UnRegister, "unregister", "", "", "unregister specified plugin")
	cliPlugin.Flags().BoolVarP(&pluginFlags.UnRegisterAll, "unregister-all", "", false, "unregister all registered plugin")
	cliPlugin.Flags().BoolVarP(&pluginFlags.List, "list", "", false, "list registered plugins")

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
	HavocCmd.AddCommand(cliPlugin)
	HavocCmd.AddCommand(cliDatabase)
}

func havocRun(cmd *cobra.Command, _ []string) {
	if len(os.Args) <= 2 {
		err := cmd.Help()
		if err != nil {
			logger.Error("error: %v", err)
		}
		os.Exit(0)
	}
}

func serverRun(cmd *cobra.Command, _ []string) {
	var (
		path, _    = os.Getwd()
		timer      = time.Now()
		teamserver *server.Teamserver
		err        error
	)

	if len(os.Args) <= 2 {
		err = cmd.Help()
		if err != nil {
			logger.Error("error: %v", err)
		}
		os.Exit(0)
	}

	//
	if teamserver = server.NewTeamserver(); teamserver == nil {
		logger.Error("failed to create server")
		return
	}

	teamserver.SetFlags(serverFlags)

	if serverFlags.Debug {
		logger.SetDebug(true)
		logger.Debug("debug mode enabled")
	}

	logger.Info("%v [%v %v]", colors.BoldWhite("Havoc Framework"), server.Version, colors.BoldBlue(server.CodeName))

	if serverFlags.Default {
		if err = teamserver.Profile(path + "/data/havoc.toml"); err != nil {
			logger.Error("failed to parse default profile: %v", err)
			return
		}
	} else if serverFlags.Profile != "" {
		err = teamserver.Profile(serverFlags.Profile)
		if err != nil {
			logger.Error("failed to parse profile %v: %v", serverFlags.Profile, err)
			return
		}
	} else {
		logger.Error("no profile specified")
		logger.Error("specify a profile with --profile or choose the standard profile with --default")
		os.Exit(1)
	}

	logger.Info("time: " + colors.Yellow(timer.Format("02/01/2006 15:04:05")))

	teamserver.Start()

	return
}

func databaseRun(cmd *cobra.Command, _ []string) {
	var (
		err      error
		database *db.DatabaseCli
	)

	if len(os.Args) <= 2 {
		err = cmd.Help()
		if err != nil {
			logger.Error("error: %v", err)
		}
		return
	}

	if database = db.NewDatabaseCli(); database == nil {
		logger.Error("failed to create database")
		return
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

	return
}

func pluginRun(cmd *cobra.Command, _ []string) {
	var (
		err     error
		manager *plugin.Plugin
	)

	if len(os.Args) <= 2 {
		err = cmd.Help()
		if err != nil {
			logger.Error("error: %v", err)
		}
		return
	}

	if manager = plugin.NewPlugin(); manager == nil {
		logger.Error("failed to create plugin manager")
		return
	}

	if len(pluginFlags.Register) != 0 {
		if err = manager.Register(pluginFlags.Register); err != nil {
			logger.Error("failed to register plugin: %v", err)
			return
		}
	} else if len(pluginFlags.UnRegister) != 0 {
		if err = manager.UnRegister(pluginFlags.UnRegister); err != nil {
			logger.Error("failed to unregister plugin: %v", err)
			return
		}
	} else if pluginFlags.UnRegisterAll {
		if err = manager.UnRegisterAll(); err != nil {
			logger.Error("failed to unregister plugin: %v", err)
			return
		}
	} else if pluginFlags.List {
		if err = manager.List(); err != nil {
			logger.Error("failed to list plugins: %v", err)
			return
		}
	}
}
