package plugin

import "errors"

func (system *System) ListenerStart(name, protocol string, options map[string]any) (map[string]string, error) {
	var (
		data map[string]string
		err  error
		ext  *Plugin
	)

	err = errors.New("protocol not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeListener {
			return true
		}

		if protocol == ext.Data["protocol"].(string) {
			data, err = ext.ListenerStart(name, options)
			return false
		}

		return true
	})

	return data, err
}

func (system *System) ListenerRestore(name, protocol, status string, config map[string]any) (map[string]string, error) {
	var (
		data map[string]string
		err  error
		ext  *Plugin
	)

	err = errors.New("protocol not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeListener {
			return true
		}

		if protocol == ext.Data["protocol"].(string) {
			data, err = ext.ListenerRestore(name, status, config)
			return false
		}

		return true
	})

	return data, err
}

func (system *System) ListenerRemove(name, protocol string) error {
	var (
		err error
		ext *Plugin
	)

	err = errors.New("protocol not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeListener {
			return true
		}

		if protocol == ext.Data["protocol"].(string) {
			err = ext.ListenerRemove(name)
			return false
		}

		return true
	})

	return err
}

func (system *System) ListenerRestart(name, protocol string) (string, error) {
	var (
		status string
		err    error
		ext    *Plugin
	)

	err = errors.New("protocol not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeListener {
			return true
		}

		if protocol == ext.Data["protocol"].(string) {
			status, err = ext.ListenerRestart(name)
			return false
		}

		return true
	})

	return status, err
}

func (system *System) ListenerStop(name, protocol string) (string, error) {
	var (
		status string
		err    error
		ext    *Plugin
	)

	err = errors.New("protocol not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeListener {
			return true
		}

		if protocol == ext.Data["protocol"].(string) {
			status, err = ext.ListenerStop(name)
			return false
		}

		return true
	})

	return status, err
}

func (system *System) ListenerEvent(name string, event map[string]any) (map[string]any, error) {
	var (
		err      error
		ext      *Plugin
		resp     map[string]any
		protocol string
	)

	err = errors.New("protocol not found")

	if protocol, err = system.havoc.ListenerProtocol(name); err != nil {
		return nil, err
	}

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeListener {
			return true
		}

		if protocol == ext.Data["protocol"].(string) {
			resp, err = ext.ListenerEvent(name, event)
			return false
		}

		return true
	})

	return resp, err
}

func (system *System) ListenerConfig(name string) (map[string]any, error) {
	var (
		data     map[string]any
		err      error
		ext      *Plugin
		protocol string
	)

	if protocol, err = system.havoc.ListenerProtocol(name); err != nil {
		return nil, err
	}

	err = errors.New("protocol not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeListener {
			return true
		}

		if protocol == ext.Data["protocol"].(string) {
			data, err = ext.ListenerConfig(name)
			return false
		}

		return true
	})

	return data, err
}

func (system *System) ListenerEdit(name string, config map[string]any) error {
	var (
		err      error
		ext      *Plugin
		protocol string
	)

	if protocol, err = system.havoc.ListenerProtocol(name); err != nil {
		return err
	}

	err = errors.New("protocol not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeListener {
			return true
		}

		if protocol == ext.Data["protocol"].(string) {
			err = ext.ListenerEdit(name, config)
			return false
		}

		return true
	})

	return err
}
