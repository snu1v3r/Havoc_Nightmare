package plugin

import "errors"

func (s *PluginSystem) AgentGenerate(ctx map[string]any, config map[string]any) (string, []byte, map[string]any, error) {
    var (
        err  error
        ext  *Plugin
        bin  []byte
        cfg  map[string]any
        name string
    )

    err = errors.New("agent not found")

    s.loaded.Range(func(key, value any) bool {
        ext = value.(*Plugin)

        if ext.Type != PluginTypeAgent {
            return true
        }

        if ctx["name"] == ext.AgentRegister()["name"] {
            name, bin, cfg, err = ext.AgentGenerate(ctx, config)
            return false
        }

        return true
    })

    return name, bin, cfg, err
}

func (s *PluginSystem) AgentExecute(plugin, uuid string, data map[string]any, wait bool) (map[string]any, error) {
    var (
        resp map[string]any
        err  error
        ext  *Plugin
    )

    err = errors.New("agent not found")

    s.loaded.Range(func(key, value any) bool {
        ext = value.(*Plugin)

        if ext.Type != PluginTypeAgent {
            return true
        }

        if plugin == ext.AgentRegister()["name"] {
            resp, err = ext.AgentExecute(uuid, data, wait)
            return false
        }

        return true
    })

    return resp, err
}

func (s *PluginSystem) AgentProcess(ctx map[string]any, request []byte) ([]byte, error) {
    var (
        err error
        ext *Plugin
        res []byte
    )

    err = errors.New("agent not found")

    s.loaded.Range(func(key, value any) bool {
        ext = value.(*Plugin)

        if ext.Type != PluginTypeAgent {
            return true
        }

        if ctx["name"] == ext.AgentRegister()["name"] {
            res, err = ext.AgentProcess(ctx, request)
            return false
        }

        return true
    })

    return res, err
}

func (s *PluginSystem) AgentGet(plugin, uuid string) (map[string]any, error) {
    var (
        err error
        ext *Plugin
        ctx map[string]any
    )

    err = errors.New("agent not found")

    s.loaded.Range(func(key, value any) bool {
        ext = value.(*Plugin)

        if ext.Type != PluginTypeAgent {
            return true
        }

        if plugin == ext.AgentRegister()["name"] {
            ctx, err = ext.AgentGet(uuid)
            return false
        }

        return true
    })

    return ctx, err
}
