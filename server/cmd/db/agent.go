package db

import (
	"fmt"
	"os"

	"Havoc/pkg/db"
	"Havoc/pkg/logger"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (d *DatabaseCli) AgentList() error {
	var (
		list   = table.NewWriter()
		agents []db.Agent
		err    error
		pivots int
	)

	if agents, err = d.database.AgentList(); err != nil {
		return err
	}

	if len(agents) == 0 {
		logger.Warn("no agents has been registered")
		return err
	}

	list.SetOutputMirror(os.Stdout)
	list.AppendHeader(table.Row{"uuid", "type", "status", "parent", "note"})
	list.SetStyle(table.StyleColoredDark)

	for _, agent := range agents {
		list.AppendRow(table.Row{agent.Uuid, agent.Type, agent.Status, agent.Parent, agent.Note})

		if len(agent.Parent) != 0 {
			pivots++
		}
	}

	list.AppendFooter(table.Row{fmt.Sprintf("Agents: %v", len(agents)), "", "", fmt.Sprintf("Pivots: %v", pivots)})
	list.Render()

	return nil
}

func (d *DatabaseCli) AgentRemove(uuid string) error {
	var (
		err error
	)

	if err = d.database.AgentRemove(uuid); err != nil {
		logger.Error("failed to remove agent: %v", err)
		return err
	}

	logger.Info("agent %v successfully removed", uuid)

	return nil
}

func (d *DatabaseCli) AgentClear() error {
	var (
		agents  []db.Agent
		err     error
		removed int
	)

	if agents, err = d.database.AgentList(); err != nil {
		return err
	}

	if len(agents) == 0 {
		logger.Warn("no agents has been registered")
		return err
	}

	for _, agent := range agents {
		if err = d.database.AgentRemove(agent.Uuid); err != nil {
			logger.Error("failed to remove agent: %v", err)
			continue
		}

		removed++
	}

	logger.Info("cleared %v agents from the database", removed)

	return nil
}
