package server

import (
	"Havoc/pkg/logger"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
)

func (t *Teamserver) AgentInitialize(uuid, plugin, status, note string, data map[string]any) error {
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

	logger.Debug("agent: %v", agent)

	t.UserBroadcast(false, t.EventCreate(EventAgentInitialize, map[string]any{
		"uuid": uuid,
		"type": plugin,
		"meta": data,
	}))

	return nil
}

func (t *Teamserver) AgentDbInsert(uuid, _type, parent, status, note string, metadata any) error {
	var (
		serialize bytes.Buffer
		err       error
	)

	gob.Register(map[string]any{})

	// encode the metadata object into a serialized
	// bytes buffer to be inserted into the database
	if err = gob.NewEncoder(&serialize).Encode(metadata); err != nil {
		return err
	}

	// insert the agent into the database
	err = t.database.AgentInsert(uuid, _type, parent, status, note, serialize.Bytes())
	if err != nil {
		return err
	}

	return err
}

func (t *Teamserver) AgentDbUpdate(uuid string, parent, status, note string, metadata any) error {
	var (
		serialize bytes.Buffer
		err       error
	)

	gob.Register(map[string]any{})

	// encode the metadata object into a serialized
	// bytes buffer to be inserted into the database
	if err = gob.NewEncoder(&serialize).Encode(metadata); err != nil {
		return err
	}

	// insert the agent into the database
	err = t.database.AgentUpdate(uuid, parent, status, note, serialize.Bytes())
	if err != nil {
		return err
	}

	return err
}

func (t *Teamserver) AgentDbDisable(uuid string) error {
	var err error

	// insert the agent into the database
	err = t.database.AgentDisable(uuid)
	if err != nil {
		return err
	}

	return err
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

func (t *Teamserver) AgentDelete(uuid string) error {
	// check if the given uuid exists
	if _, ok := t.agents.Load(uuid); !ok {
		return errors.New("agent with given uuid doesn't exists")
	}

	// TODO: also interact with the
	//       plugin to delete the agent

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

		// we honestly don't care much about it failed to update
		// continue and broadcast the status update
		_ = t.plugins.AgentUpdate(agent.plugin, uuid)

		t.UserBroadcast(false, t.EventCreate(EventAgentNote, map[string]any{
			"uuid": uuid,
			"note": note,
		}))

		return nil
	}

	return errors.New("agent by uuid not found")
}

// AgentSetStatus updates the agent object status
// and broadcasts the event to every connected client
func (t *Teamserver) AgentSetStatus(uuid string, status string) error {
	var agent *Agent

	// check if the given uuid already exists
	if val, ok := t.agents.Load(uuid); !ok {
		return errors.New("agent with given uuid doesn't exists")
	} else {
		agent = val.(*Agent)
	}

	agent.status = status

	// we honestly don't care much about it failed to update
	// continue and broadcast the status update
	_ = t.plugins.AgentUpdate(agent.plugin, uuid)

	t.UserBroadcast(false, t.EventCreate(EventAgentStatus, map[string]any{
		"uuid":   uuid,
		"status": status,
	}))

	return nil
}

func (t *Teamserver) AgentStatus(uuid string) (string, error) {
	var (
		agent any
		ok    bool
	)

	// check if the given uuid already exists
	if agent, ok = t.agents.Load(uuid); !ok {
		return "", errors.New("agent with given uuid doesn't exists")
	}

	return agent.(*Agent).status, nil
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

//
// Plugin interfaces and methods
//

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

func (t *Teamserver) AgentRestore(uuid, parent, status, note string, serialized []byte) error {
	var (
		agent *Agent
		value any
		ok    bool
	)

	// load stored agent by uuid from map
	if value, ok = t.agents.Load(uuid); ok {
		agent = value.(*Agent)
		return t.plugins.AgentRestore(agent.plugin, uuid, parent, status, note, serialized)
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
