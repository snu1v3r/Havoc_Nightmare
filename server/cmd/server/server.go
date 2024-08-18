package server

import "errors"

func (t *Teamserver) ServerAgentRegister(uuid, _type string, data map[string]any) error {
	var agent = &Agent{
		uuid:   uuid,
		plugin: _type,
	}

	// check if the given uuid already exists
	if _, ok := t.agents.Load(uuid); ok {
		return errors.New("agent with given uuid already exists")
	}

	// store the agent data with the given uuid
	t.agents.Store(uuid, agent)

	t.LogDebug("agent: %v", agent)

	t.UserBroadcast(false, t.EventCreate(EventAgentInitialize, map[string]any{
		"uuid": uuid,
		"type": _type,
		"meta": data,
	}))

	return nil
}

func (t *Teamserver) ServerAgentExist(uuid string) bool {
	_, found := t.agents.Load(uuid)

	return found
}

func (t *Teamserver) ServerAgent(uuid string) (map[string]any, error) {
	var (
		agent *Agent
		val   any
		ok    bool
	)

	// check if the given uuid already exists
	if val, ok = t.agents.Load(uuid); ok {
		agent = val.(*Agent)
	} else {
		return nil, errors.New("agent with given uuid doesn't exists")
	}

	return t.plugins.AgentGet(agent.plugin, uuid)
}

func (t *Teamserver) ServerAgentType(uuid string) (string, error) {
	var (
		agent any
		ok    bool
	)

	// check if the given uuid already exists
	if agent, ok = t.agents.Load(uuid); !ok {
		return "", errors.New("agent with given uuid doesn't exists")
	}

	return agent.(*Agent).plugin, nil
}

func (t *Teamserver) ServerAgentNote(uuid string) (string, error) {
	var (
		agent any
		ok    bool
	)

	// check if the given uuid already exists
	if agent, ok = t.agents.Load(uuid); !ok {
		return "", errors.New("agent with given uuid doesn't exists")
	}

	return agent.(*Agent).note, nil
}

func (t *Teamserver) ServerAgentDelete(uuid string) error {
	// check if the given uuid exists
	if _, ok := t.agents.Load(uuid); !ok {
		return errors.New("agent with given uuid doesn't exists")
	}

	t.agents.Delete(uuid)

	return nil
}

func (t *Teamserver) ServerAgentList() []string {
	var list []string

	t.agents.Range(func(k, v interface{}) bool {
		agent := v.(*Agent)
		list = append(list, agent.uuid)
		return true
	})

	return list
}

// ServerAgentStatus updates the agent object status
// and broadcasts the event to every connected client
func (t *Teamserver) ServerAgentStatus(uuid string, status string) error {
	var agent *Agent

	// check if the given uuid already exists
	if val, ok := t.agents.Load(uuid); !ok {
		return errors.New("agent with given uuid doesn't exists")
	} else {
		agent = val.(*Agent)
	}

	agent.status = status

	t.UserBroadcast(false, t.EventCreate(EventAgentStatus, map[string]any{
		"uuid":   uuid,
		"status": status,
	}))

	return nil
}
