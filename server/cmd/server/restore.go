package server

import (
	"Havoc/pkg/db"
	"Havoc/pkg/logger"
)

func (t *Teamserver) RestoreAgents() error {
	var (
		list []db.Agent
		err  error
	)

	// get a list of agents from the
	// database that aren't disabled
	list, err = t.database.AgentLists()
	if err != nil {
		logger.DebugError("AgentLists failed: %v", err)
		return err
	}

	for _, agent := range list {
		err = t.plugins.AgentRestore(agent.Type, agent.Uuid, agent.Parent, agent.Status, agent.Note, agent.Metadata)
		if err != nil {
			logger.DebugError("AgentRestore failed: %v", err)
			return err
		}
	}

	return nil
}

func (t *Teamserver) RestoreListeners() error {
	return nil
}
