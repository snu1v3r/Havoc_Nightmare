package db

import (
	"Havoc/cmd/server"
	"Havoc/pkg/db"
	"Havoc/pkg/logger"
)

type Flags struct {
	ClearAll bool

	Agent struct {
		List   bool
		Remove string
		Clear  bool
	}

	Listener struct {
		List   bool
		Remove string
		Clear  bool
	}
}

type DatabaseCli struct {
	server   *server.Teamserver
	database *db.Database
}

func NewDatabaseCli() *DatabaseCli {
	var (
		teamserver *server.Teamserver
		database   *db.Database
		err        error
	)

	if teamserver = server.NewTeamserver(); teamserver == nil {
		logger.Error("failed to create server")
		return nil
	}

	if database, err = db.NewDatabase(teamserver.ConfigPath() + server.DatabasePath); err != nil {
		logger.Error("failed to open database: %v", err)
		return nil
	}

	return &DatabaseCli{
		server:   teamserver,
		database: database,
	}
}
