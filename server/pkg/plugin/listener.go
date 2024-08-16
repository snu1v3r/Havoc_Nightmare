package plugin

import "errors"

func (s *PluginSystem) ListenerStart(name, protocol string, options map[string]any) (map[string]string, error) {
	var (
		data map[string]string
		err  error
		ext  *Plugin
	)

	err = errors.New("protocol not found")

	s.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != PluginTypeListener {
			return true
		}

		if protocol == ext.ListenerRegister()["protocol"].(string) {
			data, err = ext.ListenerStart(name, options)
			return false
		}

		return true
	})

	return data, err
}

func (s *PluginSystem) ListenerRemove(name, protocol string) error {
	var (
		err error
		ext *Plugin
	)

	err = errors.New("protocol not found")

	s.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != PluginTypeListener {
			return true
		}

		if protocol == ext.ListenerRegister()["protocol"].(string) {
			err = ext.ListenerRemove(name)
			return false
		}

		return true
	})

	return err
}

func (s *PluginSystem) ListenerRestart(name, protocol string) (string, error) {
	var (
		status string
		err    error
		ext    *Plugin
	)

	err = errors.New("protocol not found")

	s.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != PluginTypeListener {
			return true
		}

		if protocol == ext.ListenerRegister()["protocol"].(string) {
			status, err = ext.ListenerRestart(name)
			return false
		}

		return true
	})

	return status, err
}

func (s *PluginSystem) ListenerStop(name, protocol string) (string, error) {
	var (
		status string
		err    error
		ext    *Plugin
	)

	err = errors.New("protocol not found")

	s.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != PluginTypeListener {
			return true
		}

		if protocol == ext.ListenerRegister()["protocol"].(string) {
			status, err = ext.ListenerStop(name)
			return false
		}

		return true
	})

	return status, err
}

func (s *PluginSystem) ListenerEvent(name string, event map[string]any) (map[string]any, error) {
	var (
		err      error
		ext      *Plugin
		resp     map[string]any
		protocol string
	)

	err = errors.New("protocol not found")

	if protocol, err = s.havoc.ListenerProtocol(name); err != nil {
		return nil, err
	}

	s.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != PluginTypeListener {
			return true
		}

		if protocol == ext.ListenerRegister()["protocol"].(string) {
			resp, err = ext.ListenerEvent(name, event)
			return false
		}

		return true
	})

	return resp, err
}

func (s *PluginSystem) ListenerConfig(name string) (map[string]any, error) {
	var (
		data     map[string]any
		err      error
		ext      *Plugin
		protocol string
	)

	if protocol, err = s.havoc.ListenerProtocol(name); err != nil {
		return nil, err
	}

	err = errors.New("protocol not found")

	s.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != PluginTypeListener {
			return true
		}

		if protocol == ext.ListenerRegister()["protocol"].(string) {
			data, err = ext.ListenerConfig(name)
			return false
		}

		return true
	})

	return data, err
}

func (s *PluginSystem) ListenerEdit(name string, config map[string]any) error {
	var (
		err      error
		ext      *Plugin
		protocol string
	)

	if protocol, err = s.havoc.ListenerProtocol(name); err != nil {
		return err
	}

	err = errors.New("protocol not found")

	s.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != PluginTypeListener {
			return true
		}

		if protocol == ext.ListenerRegister()["protocol"].(string) {
			err = ext.ListenerEdit(name, config)
			return false
		}

		return true
	})

	return err
}
