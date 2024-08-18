package server

import (
	"Havoc/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
)

func (t *Teamserver) AgentInitialize(uuid, plugin string, data map[string]any) error {
	var agent = &Agent{
		uuid:   uuid,
		plugin: plugin,
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
		"type": plugin,
		"meta": data,
	}))

	return nil
}

func (t *Teamserver) AgentExist(uuid string) bool {
	_, found := t.agents.Load(uuid)

	return found
}

func (t *Teamserver) AgentData(uuid string) (map[string]any, error) {
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

func (t *Teamserver) AgentType(uuid string) (string, error) {
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

func (t *Teamserver) AgentNote(uuid string) (string, error) {
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

func (t *Teamserver) AgentDelete(uuid string) error {
	// check if the given uuid exists
	if _, ok := t.agents.Load(uuid); !ok {
		return errors.New("agent with given uuid doesn't exists")
	}

	t.agents.Delete(uuid)

	return nil
}

func (t *Teamserver) AgentList() []string {
	var list []string

	t.agents.Range(func(k, v interface{}) bool {
		agent := v.(*Agent)
		list = append(list, agent.uuid)
		return true
	})

	return list
}

// AgentStatus updates the agent object status
// and broadcasts the event to every connected client
func (t *Teamserver) AgentStatus(uuid string, status string) error {
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

func (t *Teamserver) AgentRegisterType(name string, agent map[string]any) error {
	var (
		payload Handler
		data    []byte
		err     error
	)

	payload = Handler{
		Name: name,
		Data: agent,
	}

	// add the payload to the
	// available payloads lists
	t.payloads = append(t.payloads, payload)

	// convert the object to a json
	if data, err = json.Marshal(payload); err != nil {
		return err
	}

	// broadcast to all connected clients
	// and any future clients
	t.UserBroadcast(true, t.EventCreateJson(EventAgentRegister, data))

	return nil
}

func (t *Teamserver) AgentProcess(ctx map[string]any, request []byte) ([]byte, error) {
	for _, agent := range t.payloads {
		if agent.Data["name"].(string) == ctx["name"].(string) {
			return t.plugins.AgentProcess(ctx, request)
		}
	}

	return nil, errors.New("agent to process request not found")
}

func (t *Teamserver) AgentGenerate(ctx map[string]any, config map[string]any) (string, []byte, map[string]any, error) {
	for _, agent := range t.payloads {
		if agent.Data["name"].(string) == ctx["name"].(string) {
			return t.plugins.AgentGenerate(ctx, config)
		}
	}

	return "", nil, nil, errors.New("agent to generate not found")
}

func (t *Teamserver) AgentSetNote(uuid, note string) error {
	var (
		agent *Agent
		value any
		ok    bool
	)

	// load stored agent by uuid from map
	if value, ok = t.agents.Load(uuid); ok {
		agent = value.(*Agent)
		agent.note = note

		t.UserBroadcast(false, t.EventCreate(EventAgentNote, map[string]any{
			"uuid": uuid,
			"note": note,
		}))

		return nil
	}

	return errors.New("agent by uuid not found")
}

func (t *Teamserver) AgentExecute(uuid string, data map[string]any, wait bool) (map[string]any, error) {
	var (
		agent *Agent
		value any
		ok    bool
	)

	// load stored agent by uuid from map
	if value, ok = t.agents.Load(uuid); ok {
		agent = value.(*Agent)
		return t.plugins.AgentExecute(agent.plugin, uuid, data, wait)
	}

	return nil, errors.New("agent by uuid not found")
}

func (t *Teamserver) AgentBuildLog(context map[string]any, format string, args ...any) {
	var (
		user string
		err  error
	)

	user = context["user"].(string)

	//
	// send the client the build log message
	//
	err = t.UserSend(user, t.EventCreate(EventAgentBuildLog, map[string]any{
		"log": fmt.Sprintf(format, args...),
	}))

	if err != nil {
		logger.DebugError("failed to send build log message to %v: %v", user, err)
		return
	}
}

func (t *Teamserver) AgentCallback(uuid string, data map[string]any) {
	t.UserBroadcast(false, t.EventCreate(EventAgentCallback, map[string]any{
		"uuid": uuid,
		"data": data,
	}))
}

func (t *Teamserver) AgentConsole(uuid string, data map[string]any) {
	t.UserBroadcast(false, t.EventCreate(EventAgentConsole, map[string]any{
		"uuid": uuid,
		"data": data,
	}))
}

func (t *Teamserver) AgentHeartbeat(uuid, time string) {
	t.UserBroadcast(false, t.EventCreate(EventAgentHeartbeat, map[string]any{
		"uuid": uuid,
		"time": time,
	}))
}
