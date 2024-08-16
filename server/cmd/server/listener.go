package server

import (
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
			return p.Data["protocol"].(string), nil
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

func (t *Teamserver) ListenerInitDir(name string) (string, error) {
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
		protocol = t.listener[i].Data["protocol"].(string)

		if t.listener[i].Name == name {
			if err = t.plugins.ListenerRemove(name, protocol); err != nil {
				return err
			}

			t.UserBroadcast(true, t.EventCreate(EventListenerRemove, map[string]string{
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
		protocol = listener.Data["protocol"].(string)

		if listener.Name == name {
			if status, err = t.plugins.ListenerRestart(name, protocol); err != nil {
				return err
			}

			t.ListenerStatus(name, status)
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

			t.UserBroadcast(true, t.EventCreate(EventListenerStop, map[string]string{
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
		found  = false
	)

	for _, prot := range t.protocols {
		if val, ok := prot.Data["protocol"]; ok {
			if val.(string) == protocol {
				// check if we want to start a listener with the same
				// name by checking if a config has been specified
				if t.ListenerExists(name) && options != nil {
					return errors.New("listener already exists")
				}

				found = true

				if !t.ListenerExists(name) {
					if path, err = t.ListenerInitDir(name); err != nil {
						return errors.New("failed to create listener config path: " + err.Error())
					}
				}

				if data, err = t.plugins.ListenerStart(name, protocol, options); err != nil {
					return err
				}

				if t.ListenerExists(name) {
					status, ok = data["status"]

					t.UserBroadcast(true, t.EventCreate(EventListenerStart, map[string]string{
						"name":   name,
						"status": status,
					}))
				} else {
					host, ok = data["host"]
					port, ok = data["port"]
					status, ok = data["status"]

					t.listener = append(t.listener, Handler{
						Name: name,
						Data: map[string]any{
							"protocol":    protocol,
							"host":        host,
							"port":        port,
							"status":      status,
							"config.path": path,
						},
					})

					t.UserBroadcast(true, t.EventCreate(EventListenerStart, map[string]string{
						"name":     name,
						"protocol": protocol,
						"host":     host,
						"port":     port,
						"status":   status,
					}))
				}

				break
			}
		}
	}

	if !found {
		err = fmt.Errorf("listener protocol \"%v\" has not been registered by the server", protocol)
	}

	return err
}

func (t *Teamserver) ListenerEvent(name string, event map[string]any) (map[string]any, error) {
	var (
		err      error
		resp     map[string]any
		protocol string
	)

	err = errors.New("protocol not found")

	if protocol, err = t.ListenerProtocol(name); err != nil {
		return nil, err
	}

	for _, p := range t.protocols {
		if val, ok := p.Data["protocol"]; ok {
			if val.(string) == protocol {
				if resp, err = t.plugins.ListenerEvent(name, event); err != nil {
					return nil, err
				}

				break
			}
		}
	}

	return resp, err
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

	return err
}

func (t *Teamserver) ListenerLog(name string, format string, args ...any) {
	t.UserBroadcast(true, t.EventCreate(EventListenerLog, map[string]any{
		"name": name,
		"log":  fmt.Sprintf(format, args...),
	}))
}

func (t *Teamserver) ListenerStatus(name string, status string) {
	t.UserBroadcast(true, t.EventCreate(EventListenerStatus, map[string]any{
		"name":   name,
		"status": status,
	}))
}

func (t *Teamserver) ListenerConfigPath(name string) string {
	return t.configPath + "/listeners/" + name
}
