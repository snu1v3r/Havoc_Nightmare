package server

import (
	"Havoc/pkg/plugin"
	"fmt"
	"os"
	"path/filepath"
)

func (t *Teamserver) PluginList() []map[string]any {
	var (
		list    []map[string]any
		plugins []*plugin.Plugin
		err     error
	)

	if plugins, err = t.plugins.List(); err != nil {
		return nil
	}

	for _, entry := range plugins {
		list = append(list, map[string]any{
			"name":      entry.Name,
			"type":      entry.Type,
			"version":   entry.Version,
			"author":    entry.Author,
			"resources": entry.Resources,
			"data":      entry.Data,
		})
	}

	return list
}

func (t *Teamserver) PluginResource(name, resource string) ([]byte, error) {
	var (
		_plugin *plugin.Plugin
		data    []byte
		err     error
		found   bool
	)

	// get the plugin object by name
	if _plugin, err = t.plugins.Plugin(name); err != nil {
		return nil, err
	}

	// check if the resource is available and or valid
	for _, entry := range _plugin.Resources {
		if entry == resource {
			found = true
			break
		}
	}

	// if the specified resource is not available
	// under the plugin then abort
	if !found {
		return nil, fmt.Errorf("plugin %s does not contain %s resource", name, resource)
	}

	// now read the resource under the plugins file path
	if data, err = os.ReadFile(filepath.Dir(_plugin.Path) + "/" + resource); err != nil {
		return nil, err
	}

	return data, err
}
