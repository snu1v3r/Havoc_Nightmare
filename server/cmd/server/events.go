package server

import (
	"Havoc/pkg/logger"
	"encoding/json"
	"time"
)

const (
	EventUserLogin   = "user::login"
	EventUserLogout  = "user::logout"
	EventUserMessage = "user::message"
)

const (
	EventListenerRegister = "listener::register"
	EventListenerStart    = "listener::start"
	EventListenerEdit     = "listener::edit"
	EventListenerStop     = "listener::stop"
	EventListenerRemove   = "listener::remove"
	EventListenerLog      = "listener::log"
	EventListenerStatus   = "listener::status"
)

const (
	EventAgentRegister   = "agent::register"
	EventAgentInitialize = "agent::initialize"
	EventAgentCallback   = "agent::callback"
	EventAgentConsole    = "agent::console"
	EventAgentHeartbeat  = "agent::heartbeat"
	EventAgentNote       = "agent::note"
	EventAgentRemove     = "agent::remove"
	EventAgentStatus     = "agent::status"
	EventAgentBuildLog   = "agent::build-log"
)

func (t *Teamserver) EventCreate(event string, data any) map[string]any {
	return map[string]any{
		"time": time.Now().Format("02/01/2006 15:04:05"),
		"type": event,
		"data": data,
	}
}

func (t *Teamserver) EventCreateJson(event string, data []byte) map[string]any {
	var dict map[string]any

	if err := json.Unmarshal(data, &dict); err != nil {
		logger.DebugError("json.Unmarshal failed: %v", err)
		return nil
	}

	return map[string]any{
		"time": time.Now().Format("02/01/2006 15:04:05"),
		"type": event,
		"data": dict,
	}
}
