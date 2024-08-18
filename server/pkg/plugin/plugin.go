package plugin

import (
	"Havoc/pkg/utils"
	"errors"
	"fmt"
	"plugin"
	"reflect"
	"sync"
)

const (
	PluginTypeListener   = "listener"
	PluginTypeAgent      = "agent"
	PluginTypeManagement = "management"
)

type PluginListener struct {
	Name    string
	Type    string
	Options map[string]any
}

type HavocInterface interface {
	LogInfo(fmt string, args ...any)
	LogError(fmt string, args ...any)
	LogWarn(fmt string, args ...any)
	LogFatal(fmt string, args ...any)
	LogPanic(fmt string, args ...any)
	LogDebug(fmt string, args ...any)
	LogDebugError(fmt string, args ...any)

	Version() map[string]string

	ListenerRegister(name string, listener map[string]any) error
	ListenerProtocol(name string) (string, error)

	AgentRegisterType(name string, agent map[string]any) error
}

type BasicInterface interface {
	Register(havoc any) map[string]any
}

type ListenerInterface interface {
	ListenerRegister() map[string]any
	ListenerStart(name string, options map[string]any) (map[string]string, error)
	ListenerEdit(name string, config map[string]any) error
	ListenerStop(name string) (string, error)
	ListenerRestart(name string) (string, error)
	ListenerRemove(name string) error
	ListenerEvent(name string, event map[string]any) (map[string]any, error)
	ListenerConfig(name string) (map[string]any, error)
}

type AgentInterface interface {
	AgentRegister() map[string]any
	AgentGenerate(ctx map[string]any, config map[string]any) (string, []byte, map[string]any, error)
	AgentExecute(uuid string, data map[string]any, wait bool) (map[string]any, error)
	AgentProcess(ctx map[string]any, request []byte) ([]byte, error)
	AgentGet(uuid string) (map[string]any, error)
}

type ManagementInterface interface{}

type Plugin struct {
	Name string
	Type string
	Data map[string]any

	// the loaded plugin interface that are callable
	BasicInterface
	ListenerInterface
	AgentInterface
	ManagementInterface
}

type PluginSystem struct {
	// havoc teamserver instance that is
	// going to be passed to every loaded plugin
	havoc HavocInterface

	// loaded havoc plugins
	// those plugins have been
	// registered and are safe to use
	loaded sync.Map
}

// NewPluginSystem
// create a new plugin system instance
func NewPluginSystem(havoc HavocInterface) *PluginSystem {
	var ps = new(PluginSystem)

	ps.havoc = havoc

	return ps
}

// RegisterPlugin is going to register a specified havoc plugin
func (s *PluginSystem) RegisterPlugin(path string) (*Plugin, error) {
	var (
		err      error
		open     *plugin.Plugin
		lookup   plugin.Symbol
		inter    BasicInterface
		register map[string]any
		ext      *Plugin
		ok       bool
	)

	// try to open plugin
	if open, err = plugin.Open(path); err != nil {
		return nil, err
	}

	// try to get the exported plugin object
	// to register the plugin and retrieve metadata about it
	if lookup, err = open.Lookup("Plugin"); err != nil {
		return nil, err
	}

	// reflect the method and
	// check if it's a valid interface
	if _, ok = reflect.TypeOf(lookup).MethodByName("Register"); !ok {
		return nil, errors.New("method \"Register\" not found inside of plugin")
	}

	// cast the looked up symbol to
	// the HavocPlugin interface
	inter = lookup.(BasicInterface)

	// try to register the plugin
	register = inter.Register(s.havoc)

	// add plugin to the internal sync
	// map and make it available
	if ext, err = s.AddPlugin(register, lookup); err != nil {
		return nil, err
	}

	return ext, nil
}

// AddPlugin the registered plugin to see if there
// wasn't given any faulty or lacking info and
// creates a havoc Plugin object
func (s *PluginSystem) AddPlugin(register map[string]any, inter any) (*Plugin, error) {
	var (
		ext = new(Plugin)
		err error
	)

	if len(register) == 0 {
		return nil, errors.New("register is empty")
	}

	// get the name of the registered plugin
	switch register["name"].(type) {
	case string:
		// check if the name is empty.
		// if yes then return an error
		if len(register["name"].(string)) == 0 {
			return nil, errors.New("register.name is empty")
		}

		// set the name of the extension/plugin
		ext.Name = register["name"].(string)
		break

	default:
		return nil, errors.New("register.name is not a string")
	}

	// get the type of the registered plugin
	switch register["type"].(type) {
	case string:
		// check if the type is empty.
		// if yes then return an error
		if len(register["type"].(string)) == 0 {
			return nil, errors.New("register.type is empty")
		}

		// set the name of the extension/plugin
		ext.Type = register["type"].(string)
	default:
		return nil, errors.New("register.type is not a string")
	}

	// sanity check interface and insert it into the plugin
	if err = s.CheckAndInsertInterface(ext, inter); err != nil {
		return nil, err
	}

	if err = s.interactPlugin(ext); err != nil {
		return nil, err
	}

	// add the ext to the sync map
	s.loaded.Store(utils.GenerateID(32), ext)

	return ext, nil
}

// CheckAndInsertInterface
// this method checks if a specific plugin is exporting all the
// needed methods for the returned plugin type, if it does then
// cast the right interface to the plugin object
func (s *PluginSystem) CheckAndInsertInterface(extension *Plugin, inter any) error {
	var (
		reflection = reflect.TypeOf(inter)
		ok         bool
	)

	switch extension.Type {
	case PluginTypeAgent:

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("AgentRegister"); !ok {
			return fmt.Errorf("AgentRegister not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("AgentGenerate"); !ok {
			return fmt.Errorf("AgentGenerate not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("AgentExecute"); !ok {
			return fmt.Errorf("AgentExecute not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("AgentProcess"); !ok {
			return fmt.Errorf("AgentProcess not found")
		}

		// cast the interface
		// we found everything we searched for
		extension.AgentInterface = inter.(AgentInterface)

		break

	case PluginTypeListener:

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("ListenerRegister"); !ok {
			return fmt.Errorf("ListenerRegister not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("ListenerStart"); !ok {
			return fmt.Errorf("ListenerStart not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("ListenerEdit"); !ok {
			return fmt.Errorf("ListenerEdit not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("ListenerRestart"); !ok {
			return fmt.Errorf("ListenerRestart not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("ListenerRemove"); !ok {
			return fmt.Errorf("ListenerRemove not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("ListenerStop"); !ok {
			return fmt.Errorf("ListenerStop not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("ListenerEvent"); !ok {
			return fmt.Errorf("ListenerEvent not found")
		}

		// sanity check if the method exist
		if _, ok = reflection.MethodByName("ListenerConfig"); !ok {
			return fmt.Errorf("ListenerConfig not found")
		}

		// cast the interface
		// we found everything we searched for
		extension.ListenerInterface = inter.(ListenerInterface)

		break

	case PluginTypeManagement:
		break

	default:
		return fmt.Errorf("invalid plugin type: \"%v\" not found", extension.Type)
	}

	return nil
}

// interactPlugin
// interact with plugin by calling the plugin
// register function for the type of plugin
func (s *PluginSystem) interactPlugin(extension *Plugin) error {

	switch extension.Type {
	case PluginTypeAgent:

		if err := s.havoc.AgentRegisterType(extension.Name, extension.AgentRegister()); err != nil {
			return err
		}

		break

	case PluginTypeListener:

		if err := s.havoc.ListenerRegister(extension.Name, extension.ListenerRegister()); err != nil {
			return err
		}

		break

	case PluginTypeManagement:
		break

	default:
		return fmt.Errorf("invalid plugin type: \"%v\" not found", extension.Type)
	}

	return nil
}
