package plugin

import (
	"Havoc/cmd/server"
	"Havoc/pkg/colors"
	"Havoc/pkg/db"
	"Havoc/pkg/logger"
	"Havoc/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

type Flags struct {
	Register      string
	UnRegister    string
	UnRegisterAll bool
	List          bool
}

type Plugin struct {
	server   *server.Teamserver
	database *db.Database
}

func NewPlugin() *Plugin {
	var (
		teamserver *server.Teamserver
		database   *db.Database
		err        error
	)

	// to ensure that the database and config path is accessible/created
	if teamserver = server.NewTeamserver(); teamserver == nil {
		logger.Error("failed to create server")
		return nil
	}

	// create an object to interact with the database
	if database, err = db.NewDatabase(teamserver.ConfigPath() + server.DatabasePath); err != nil {
		logger.Error("failed to open database: %v", err)
		return nil
	}

	return &Plugin{
		server:   teamserver,
		database: database,
	}
}

// Register
// will register the specified plugin path to the team server
// database to be loaded everytime on server startup
func (p *Plugin) Register(path string) error {
	var (
		err    error
		data   []byte
		object []map[string]any
		count  int
		list   = table.NewWriter()
	)

	if data, err = os.ReadFile(path + "/plugins.json"); err != nil {
		logger.Error("failed to read plugins.json from %v: %v", path, err)
		return err
	}

	if err = json.Unmarshal(data, &object); err != nil {
		logger.Error("failed to unmarshal plugins.json from %v: %v", path, err)
		return err
	}

	list.SetOutputMirror(os.Stdout)
	list.AppendHeader(table.Row{"name", "type", "version", "path"})
	list.SetStyle(table.StyleColoredDark)

	for _, entry := range object {
		var (
			name    string
			_type   string
			version string
			_path   string
		)

		if name, err = utils.MapKey[string](entry, "name"); err != nil {
			logger.Error("failed to get plugin name: %v", err)
			continue
		}

		if _type, err = utils.MapKey[string](entry, "type"); err != nil {
			logger.Error("failed to get plugin type (%v): %v", name, err)
			continue
		}

		if version, err = utils.MapKey[string](entry, "version"); err != nil {
			logger.Error("failed to get plugin version (%v): %v", name, err)
			continue
		}

		if _path, err = utils.MapKey[string](entry, "path"); err != nil {
			logger.Error("failed to get plugin path (%v): %v", name, err)
			continue
		}

		// TODO: validate if name, _type and _path is valid !!
		// TODO: maybe resolve the absolute path of path and save that into the database

		err = p.database.PluginRegister(name, _type, version, path+"/"+_path)
		if err != nil {
			return err
		}

		list.AppendRow(table.Row{name, _type, version, path + "/" + _path})

		count++
	}

	list.AppendFooter(table.Row{fmt.Sprintf("Plugins registered: %v", count), "", ""})
	list.Render()

	return err
}

// UnRegister
// will unregister the specified plugin from the team server database
func (p *Plugin) UnRegister(name string) error {
	var (
		data *db.Plugin
		err  error
	)

	if data, err = p.database.Plugin(name); err != nil {
		return err
	}

	err = p.database.PluginUnRegister(name)
	if err != nil {
		return err
	}

	logger.Info("plugin %v (%v) has been unregistered", colors.Blue(data.Name), data.Type)

	return nil
}

// UnRegisterAll
// will unregister all registered plugins from the database
func (p *Plugin) UnRegisterAll() error {
	var (
		err     error
		plugins []db.Plugin
	)

	if plugins, err = p.database.PluginList(); err != nil {
		return err
	}

	for _, plugin := range plugins {
		if err = p.UnRegister(plugin.Name); err != nil {
			logger.Error("failed to unregister plugin %v: %v", colors.Blue(plugin.Name), err)
			continue
		}
	}

	return nil
}

// List
// will list the registered plugins from the
// database and display the name, type and path
func (p *Plugin) List() error {
	var (
		err     error
		plugins []db.Plugin
		list    = table.NewWriter()
	)

	if plugins, err = p.database.PluginList(); err != nil {
		return err
	}

	if len(plugins) == 0 {
		logger.Warn("no plugins has been registered")
		return nil
	}

	list.SetOutputMirror(os.Stdout)
	list.AppendHeader(table.Row{"name", "type", "version", "path"})
	list.SetStyle(table.StyleColoredDark)

	for _, plugin := range plugins {
		list.AppendRow(table.Row{plugin.Name, plugin.Type, plugin.Version, plugin.Path})
	}

	list.AppendFooter(table.Row{fmt.Sprintf("Plugins: %v", len(plugins)), "", ""})
	list.Render()

	return err
}
