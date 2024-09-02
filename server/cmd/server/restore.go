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
		logger.DebugError("ListenerList failed: %v", err)
		return err
	}

	// iterate over the queried listener
	// list and add each one to the server
	for _, listener := range list {
		// initialize directory
		path, err = t.ListenerDir(listener.Name)
		if err != nil {
			logger.DebugError("ListenerDir failed on %v: %v", listener.Name, err)
			return err
		}

		// restore the listener
		data, err = t.ListenerRestore(listener.Name, listener.Protocol, listener.Status, listener.Config)
		if err != nil {
			logger.DebugError("ListenerRestore failed on %v: %v", listener.Name, err)
			return err
		}

		host, _ = data["host"]
		port, _ = data["port"]
		status, _ = data["status"]

		t.listener = append(t.listener, Handler{
			Name: listener.Name,
			Data: map[string]any{
				"protocol":    listener.Protocol,
				"host":        host,
				"port":        port,
				"status":      status,
				"config.path": path,
			},
		})
	}

	return nil
}
