package plugin

import "errors"

func (system *System) AgentGenerate(implant string, config map[string]any) (string, []byte, map[string]any, error) {
	var (
		err  error
		ext  *Plugin
		bin  []byte
		cfg  map[string]any
		name string
	)

	err = errors.New("agent not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeAgent {
			return true
		}

		if implant == ext.Data["name"] {
			name, bin, cfg, err = ext.AgentGenerate(config)
			return false
		}

		return true
	})

	return name, bin, cfg, err
}

func (system *System) AgentRemove(plugin, uuid string) error {
	var (
		err error
		ext *Plugin
	)

	err = errors.New("agent not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeAgent {
			return true
		}

		if plugin == ext.Data["name"] {
			err = ext.AgentRemove(uuid)
			return false
		}

		return true
	})

	return err
}

func (system *System) AgentRestore(plugin, uuid, parent, status, note string, serialized []byte) error {
	var (
		err error
		ext *Plugin
	)

	err = errors.New("agent not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeAgent {
			return true
		}

		if plugin == ext.Data["name"] {
			err = ext.AgentRestore(uuid, parent, status, note, serialized)
			return false
		}

		return true
	})

	return err
}

func (system *System) AgentUpdate(plugin, uuid string) error {
	var (
		err error
		ext *Plugin
	)

	err = errors.New("agent not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeAgent {
			return true
		}

		if plugin == ext.Data["name"] {
			err = ext.AgentUpdate(uuid)
			return false
		}

		return true
	})

	return err
}

func (system *System) AgentExecute(plugin, uuid string, data map[string]any, wait bool) (map[string]any, error) {
	var (
		resp map[string]any
		err  error
		ext  *Plugin
	)

	err = errors.New("agent not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeAgent {
			return true
		}

		if plugin == ext.Data["name"] {
			resp, err = ext.AgentExecute(uuid, data, wait)
			return false
		}

		return true
	})

	return resp, err
}

func (system *System) AgentProcess(implant string, context map[string]any, request []byte) ([]byte, error) {
	var (
		err error
		ext *Plugin
		res []byte
	)

	err = errors.New("agent not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeAgent {
			return true
		}

		if implant == ext.Data["name"] {
			res, err = ext.AgentProcess(context, request)
			return false
		}

		return true
	})

	return res, err
}

func (system *System) AgentGet(plugin, uuid string) (map[string]any, error) {
	var (
		err error
		ext *Plugin
		ctx map[string]any
	)

	err = errors.New("agent not found")

	system.loaded.Range(func(key, value any) bool {
		ext = value.(*Plugin)

		if ext.Type != TypeAgent {
			return true
		}

		if plugin == ext.Data["name"] {
			ctx, err = ext.AgentGet(uuid)
			return false
		}

		return true
	})

	return ctx, err
}
