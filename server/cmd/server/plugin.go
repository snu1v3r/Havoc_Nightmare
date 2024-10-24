package server

import "Havoc/pkg/plugin"

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

func (t *Teamserver) PluginFile(name, file string) ([]byte, error) {
	return nil, nil
}
