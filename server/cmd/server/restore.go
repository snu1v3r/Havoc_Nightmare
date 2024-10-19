package server

import (
	"Havoc/pkg/colors"
	"Havoc/pkg/db"
	"Havoc/pkg/logger"
)

var LogPrefix string

func init() {
	LogPrefix = colors.Blue("[database]") + " "
}

func (t *Teamserver) RestoreAgents() error {
	var (
		list []db.Agent
		err  error
	)

	// get a list of agents from the
	// database that aren't disabled
	list, err = t.database.AgentList()
	if err != nil {
		logger.DebugError(LogPrefix+"AgentList failed: %v", err)
		return err
	}

	for _, agent := range list {
		err = t.plugins.AgentRestore(agent.Type, agent.Uuid, agent.Parent, agent.Status, agent.Note, agent.Metadata)
		if err != nil {
			logger.Error(LogPrefix+"failed to restore agent %v (%v): %v", colors.Blue(agent.Uuid), agent.Type, err)
			continue
		}
	}

	return nil
}

func (t *Teamserver) RestoreListeners() error {
	var (
		list   []db.Listener
		data   map[string]string
		err    error
		path   string
		host   string
		port   string
		status string
	)

	// get a list of agents from the
	// database that aren't disabled
	list, err = t.database.ListenerList()
	if err != nil {
		logger.DebugError(LogPrefix+"ListenerList failed: %v", err)
		return err
	}

	// iterate over the queried listener
	// list and add each one to the server
	for _, listener := range list {
		logger.Debug("restoring listener %v", colors.Blue(listener.Name))

		// initialize directory
		path, err = t.ListenerDir(listener.Name)
		if err != nil {
			logger.DebugError(LogPrefix+"ListenerDir failed on %v: %v", listener.Name, err)
			return err
		}

		// restore the listener
		data, err = t.ListenerRestore(listener.Name, listener.Protocol, listener.Status, listener.Config)
		if err != nil {
			logger.Error(LogPrefix+"failed to restore listener %v (%v): %v", colors.Blue(listener.Name), listener.Protocol, err)
			continue
		}

		host, _ = data["host"]
		port, _ = data["port"]
		status, _ = data["status"]

		t.listener = append(t.listener, Listener{
			Name:     listener.Name,
			Protocol: listener.Protocol,
			Host:     host,
			Port:     port,
			Status:   status,
			Path:     path,
		})

		t.ListenerSetStatus(listener.Name, status)
	}

	return nil
}
