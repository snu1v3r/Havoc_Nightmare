package server

import (
	"Havoc/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
)

func (t *Teamserver) AgentRegister(name string, agent map[string]any) error {
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

func (t *Teamserver) AgentNote(uuid, note string) error {
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

func (t *Teamserver) AgentCallback(uuid string, _type string, data map[string]any) {
	var save = false

	if _type == "console" {
		save = true
	}

	t.UserBroadcast(save, t.EventCreate(EventAgentCallback, map[string]any{
		"uuid": uuid,
		"type": _type,
		"data": data,
	}))
}
