package cmd

import (
	"fmt"
	"os"
	"time"
	"io/ioutil"


	"Havoc/cmd/server"
	"Havoc/pkg/colors"
	"Havoc/pkg/events"
	"Havoc/pkg/logger"
	"Havoc/pkg/logr"

	"github.com/spf13/cobra"
)

var CobraServer = &cobra.Command{
	Use:          "server",
	Short:        "teamserver command",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			ServerTimer = time.Now()
			LogrPath    = "/log/" + ServerTimer.Format("2006.01.02._15:04:05")
			Server      *server.Teamserver
			err	    error
		)

		if len(os.Args) <= 2 {
			err = cmd.Help()
			if err != nil {
				return err
			}
			os.Exit(0)
		}
		DirPath := flags.Server.Persist 
		_ = os.MkdirAll(DirPath + "/data", os.ModePerm)
		_ = os.MkdirAll(DirPath + "/loot", os.ModePerm)
		_ = os.MkdirAll(DirPath + "/log", os.ModePerm)
		_ = os.MkdirAll(DirPath + "/profiles", os.ModePerm)
		logger.LoggerInstance.STDERR = os.Stderr
		logger.SetStdOut(os.Stderr)
		
		Server = server.NewTeamserver(DirPath + DatabasePath)
		Server.SetServerFlags(flags)

		logr.LogrInstance = logr.NewLogr(DirPath, LogrPath)
		if logr.LogrInstance == nil {
			logger.Error("Failed to create Logr loot folder")
			return nil
		}

		logr.LogrInstance.LogrSendText = func(text string) {
			var pk = events.Teamserver.Logger(text)

			Server.EventAppend(pk)
			Server.EventBroadcast("", pk)
		}

		logr.LogrInstance.ServerStdOutInit()

		startMenu()

		if flags.Server.Debug {
			logger.SetDebug(true)
			logger.Debug("Debug mode enabled")
		}

		logger.ShowTime(flags.Server.Verbose)

		logger.Info(fmt.Sprintf("Havoc Framework [Version: %v] [CodeName: %v]", VersionNumber, VersionName))

		if flags.Server.Default {
			_, err = os.Stat(DirPath + "/profiles/havoc.yaotl")
			if err != nil {
				cwd, _ := os.Getwd()
				data, _ := ioutil.ReadFile(cwd + "/profiles/havoc.yaotl")
				_ = ioutil.WriteFile(DirPath + "/profiles/havoc.yaotl", data, 0644)
			}
			Server.SetProfile(DirPath + "/profiles/havoc.yaotl")
		} else if flags.Server.Profile != "" {
			fmt.Sprintf("profile")
			Server.SetProfile(flags.Server.Profile)
		} else {
			logger.Error("No profile specified. Specify a profile with --profile or choose the standard profile with --default")
			os.Exit(1)
		}

		if !Server.FindSystemPackages() {
			logger.Error("Please install needed packages. Refer to the Wiki for more help.")
			os.Exit(1)
		}

		logger.Info("Time: " + colors.Yellow(ServerTimer.Format("02/01/2006 15:04:05")))
		logger.Info("Teamserver logs saved under: " + colors.Blue(DirPath + LogrPath))

		// start teamserver
		Server.Start(DirPath)

		os.Exit(0)

		return nil
	},
}
