package server

import (
	"Havoc/pkg/logger"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ListenerRegister
// register a listener to the server and
// notify clients to a newly available handler.
func (t *Teamserver) ListenerRegister(name string, listener map[string]any) error {
	var (
		protocol Handler
		data     []byte
		err      error
	)

	protocol = Handler{
		Name: name,
		Data: listener,
	}

	// add the listener to the
	// available protocols lists
	t.protocols = append(t.protocols, protocol)

	// convert the object to a json
	if data, err = json.Marshal(protocol); err != nil {
		return err
	}

	// broadcast to all connected clients
	// and any future clients
	t.UserBroadcast(true, t.EventCreateJson(EventListenerRegister, data))

	return nil
}

func (t *Teamserver) ListenerProtocol(name string) (string, error) {
	for _, p := range t.listener {
		if p.Name == name {
			return p.Protocol, nil
		}
	}

	return "", errors.New("listener protocol not found")
}

func (t *Teamserver) ListenerExists(name string) bool {
	for _, listener := range t.listener {
		if listener.Name == name {
			return true
		}
	}

	return false
}

func (t *Teamserver) ListenerProtocolExists(protocol string) bool {
	for _, listener := range t.protocols {
		if listener.Data["protocol"].(string) == protocol {
			return true
		}
	}

	return false
}

func (t *Teamserver) ListenerDir(name string) (string, error) {
	var (
		path string
		err  error
	)

	path = t.ConfigPath() + "/listeners/" + name

	// create agents directory
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return path, err
	}

	return path, nil
}

func (t *Teamserver) ListenerRemove(name string) error {
	var (
		protocol string
		err      error
	)

	if !t.ListenerExists(name) {
		return errors.New("listener not found")
	}

	for i := 0; i < len(t.listener); i++ {
		protocol = t.listener[i].Protocol

		if t.listener[i].Name == name {
			if err = t.plugins.ListenerRemove(name, protocol); err != nil {
				return err
			}

			if err = t.database.ListenerRemove(name); err != nil {
				return err
			}

			t.UserBroadcast(false, t.EventCreate(EventListenerRemove, map[string]string{
				"name": name,
			}))

			// remove the current index from the listener entry
			t.listener = append(t.listener[:i], t.listener[i+1:]...)

			break
		}
	}

	return nil
}

func (t *Teamserver) ListenerRestart(name string) error {
	var (
		protocol string
		status   string
		err      error
	)

	if !t.ListenerExists(name) {
		return errors.New("listener not found")
	}

	for _, listener := range t.listener {
		protocol = listener.Protocol

		if listener.Name == name {
			if status, err = t.plugins.ListenerRestart(name, protocol); err != nil {
				return err
			}

			t.ListenerSetStatus(name, status, true)
			break
		}
	}

	return nil
}

func (t *Teamserver) ListenerStop(name string) error {
	var (
		protocol string
		status   string
		err      error
	)

	if !t.ListenerExists(name) {
		return errors.New("listener not found")
	}

	if protocol, err = t.ListenerProtocol(name); err != nil {
		return err
	}

	for _, listener := range t.listener {
		if listener.Name == name {
			if status, err = t.plugins.ListenerStop(name, protocol); err != nil {
				return err
			}

			t.ListenerSetStatus(name, status, false)
			t.UserBroadcast(false, t.EventCreate(EventListenerStop, map[string]string{
				"name":   name,
				"status": status,
			}))
			break
		}
	}

	return nil
}

func (t *Teamserver) ListenerStart(name, protocol string, options map[string]any) error {
	var (
		err    error
		data   map[string]string
		host   string
		port   string
		status string
		path   string
	)

	// check if specified protocol exists and has been registered to the server
	if t.ListenerProtocolExists(protocol) {
		// check if we want to start a listener with the same
		// name by checking if a config has been specified
		if t.ListenerExists(name) && options != nil {
			return errors.New("listener already exists")
		}

		if !t.ListenerExists(name) {
			if path, err = t.ListenerDir(name); err != nil {
				return errors.New("failed to create listener config path: " + err.Error())
			}
		}

		if data, err = t.plugins.ListenerStart(name, protocol, options); err != nil {
			return err
		}

		if t.ListenerExists(name) {
			status, _ = data["status"]

			t.ListenerSetStatus(name, status, false)
			t.UserBroadcast(false, t.EventCreate(EventListenerStart, map[string]string{
				"name":   name,
				"status": status,
			}))
		} else {
			host, _ = data["host"]
			port, _ = data["port"]
			status, _ = data["status"]

			t.listener = append(t.listener, Listener{
				Name:     name,
				Protocol: protocol,
				Host:     host,
				Port:     port,
				Status:   status,
				Path:     path,
			})

			err = t.DatabaseListenerInsert(name, protocol, status, options)
			if err != nil {
				// as long as we managed to start the listener but failed to
				// insert it into the database we should still try to continue
				logger.DebugError("DatabaseListenerInsert error on %v: %v", name, err)
			}

			t.UserBroadcast(false, t.EventCreate(EventListenerStart, map[string]string{
				"name":     name,
				"protocol": protocol,
				"host":     host,
				"port":     port,
				"status":   status,
			}))
		}
	} else {
		err = fmt.Errorf("listener protocol \"%v\" has not been registered by the server", protocol)
	}

	return err
}

func (t *Teamserver) ListenerRestore(name, protocol, status string, config []byte) (map[string]string, error) {
	var (
		err     error
		options = make(map[string]any)
	)

	err = errors.New("protocol not found")

	if !t.ListenerProtocolExists(protocol) {
		return nil, errors.New("listener protocol not found")
	}

	err = gob.NewDecoder(bytes.NewBuffer(config)).Decode(&options)
	if err != nil {
		logger.DebugError("gob.NewDecoder error: %v", err)
		return nil, err
	}

	return t.plugins.ListenerRestore(name, protocol, status, options)
}

func (t *Teamserver) ListenerEvent(name string, event map[string]any) (map[string]any, error) {
	var (
		err      error
		protocol string
	)

	err = errors.New("protocol not found")

	if protocol, err = t.ListenerProtocol(name); err != nil {
		return nil, err
	}

	if !t.ListenerProtocolExists(protocol) {
		return nil, errors.New("listener protocol not found")
	}

	return t.plugins.ListenerEvent(name, event)
}

func (t *Teamserver) ListenerConfig(name string) (map[string]any, error) {
	var (
		err  error
		resp map[string]any
	)

	err = errors.New("protocol not found")

	if resp, err = t.plugins.ListenerConfig(name); err != nil {
		return nil, err
	}

	return resp, err
}

func (t *Teamserver) ListenerEdit(name string, config map[string]any) error {
	var err error

	err = errors.New("protocol not found")

	if err = t.plugins.ListenerEdit(name, config); err != nil {
		return err
	}

	// TODO: edit the database entry as well

	return err
}

func (t *Teamserver) ListenerList() []map[string]string {
	var list []map[string]string

	for _, listener := range t.listener {
		list = append(list, map[string]string{
			"name":     listener.Name,
			"protocol": listener.Protocol,
			"host":     listener.Host,
			"port":     listener.Port,
			"status":   listener.Status,
		})
	}

	return list
}

func (t *Teamserver) ListenerLog(name string, format string, args ...any) {
	t.UserBroadcast(false, t.EventCreate(EventListenerLog, map[string]any{
		"name": name,
		"log":  fmt.Sprintf(format, args...),
	}))
}

// ListenerSetStatus
// sets the listener status, updates the listener
// database entry and broadcasts the status change
func (t *Teamserver) ListenerSetStatus(name string, status string, event bool) {
	var (
		err    error
		config map[string]any
	)

	for i := range t.listener {
		if t.listener[i].Name == name {
			t.listener[i].Status = status
			break
		}
	}

	// let's try to update the database listener status entry
	config, err = t.ListenerConfig(name)
	if err != nil {
		logger.DebugError("ListenerConfig error on %v: %v", name, err)
	} else {
		err = t.DatabaseListenerUpdate(name, status, config)
		if err != nil {
			logger.DebugError("DatabaseListenerUpdate error on %v: %v", name, err)
		}
	}

	if event {
		t.UserBroadcast(false, t.EventCreate(EventListenerStatus, map[string]any{
			"name":   name,
			"status": status,
		}))
	}
}

func (t *Teamserver) ListenerConfigPath(name string) string {
	return t.configPath + "/listeners/" + name
}

func (t *Teamserver) DatabaseListenerInsert(name, protocol, status string, config map[string]any) error {
	var (
		serialize bytes.Buffer
		err       error
	)

	// encode the config map into a serialized bytes
	// buffer to be inserted into the database
	if err = gob.NewEncoder(&serialize).Encode(config); err != nil {
		return err
	}

	// insert the listener into the database
	err = t.database.ListenerInsert(name, protocol, status, serialize.Bytes())
	if err != nil {
		return err
	}

	return err
}

func (t *Teamserver) DatabaseListenerUpdate(name, status string, config map[string]any) error {
	var (
		serialize bytes.Buffer
		err       error
	)

	// encode the config map into a serialized bytes
	// buffer to be inserted into the database
	if err = gob.NewEncoder(&serialize).Encode(config); err != nil {
		return err
	}

	// insert the listener into the database
	err = t.database.ListenerUpdate(name, status, serialize.Bytes())
	if err != nil {
		return err
	}

	return err
}

func (t *Teamserver) DatabaseListenerRemove(name string) error {
	return t.database.ListenerRemove(name)
}
