package plugin

import (
	"Havoc/pkg/logger"
	"Havoc/pkg/utils"
	"errors"
	"fmt"
	"plugin"
	"reflect"
	"sync"
)

const (
	TypeListener   = "listener"
	TypeAgent      = "agent"
	TypeManagement = "management"
)

type HavocInterface interface {
	LogInfo(fmt string, args ...any)
	LogError(fmt string, args ...any)
	LogWarn(fmt string, args ...any)
	LogDebug(fmt string, args ...any)
	LogDbgError(fmt string, args ...any)
	LogFatal(fmt string, args ...any)
	LogPanic(fmt string, args ...any)

	Version() map[string]string

	ListenerRegisterType(name string, listener map[string]any) error
	ListenerProtocol(name string) (string, error)

	AgentRegisterType(name string, agent map[string]any) error
}

type IPluginBasic interface {
	Register(havoc any) map[string]any
}

type IPluginListener interface {
	ListenerRegister() map[string]any
	ListenerStart(name string, options map[string]any) (map[string]string, error)
	ListenerRestore(name, status string, options map[string]any) (map[string]string, error)
	ListenerEdit(name string, config map[string]any) error
	ListenerStop(name string) (string, error)
	ListenerRestart(name string) (string, error)
	ListenerRemove(name string) error
	ListenerEvent(name string, event map[string]any) (map[string]any, error)
	ListenerConfig(name string) (map[string]any, error)
}

type IPluginAgent interface {
	AgentRegister() map[string]any
	AgentRestore(uuid, parent, status, note string, serialized []byte) error
	AgentRemove(uuid string) error
	AgentUpdate(uuid string) error
	AgentGenerate(config map[string]any) (string, []byte, map[string]any, error)
	AgentExecute(uuid string, data map[string]any, wait bool) (map[string]any, error)
	AgentProcess(context map[string]any, request []byte) ([]byte, error)
	AgentGet(uuid string) (map[string]any, error)
	AgentInterface(uuid string) (any, error)
}

type IPluginManagement interface{}

type Plugin struct {
	Name string
	Type string
	Data map[string]any

	// the loaded plugin interface that are callable
	IPluginBasic
	IPluginListener
	IPluginAgent
	IPluginManagement
}

type System struct {
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
func NewPluginSystem(havoc HavocInterface) *System {
	return &System{
		havoc: havoc,
	}
}

// RegisterPlugin is going to register a specified havoc plugin
func (system *System) RegisterPlugin(path string) (*Plugin, error) {
	var (
		err      error
		open     *plugin.Plugin
		lookup   plugin.Symbol
		inter    IPluginBasic
		register map[string]any
		name     string
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
	inter = lookup.(IPluginBasic)

	// try to register the plugin
	register = inter.Register(system.havoc)

	// check if the plugin already has been registered
	// to the manager to avoid loading the same plugin
	// over and over again

	if name, err = utils.MapKey[string](register, "name"); err != nil {
		logger.Error("failed to get plugin name: %v", err)
		return nil, err
	}

	// check if the plugin already has been registered yet
	if ext, err = system.Plugin(name); !IsPluginNotFound(err) {
		return nil, err
	}

	// if has been registered then return the
	// errAlreadyRegistered error and abort
	if ext != nil {
		return nil, errAlreadyRegistered
	}

	// add plugin to the internal sync
	// map and make it available
	if ext, err = system.AddPlugin(register, lookup); err != nil {
		return nil, err
	}

	return ext, nil
}

// AddPlugin the registered plugin to see if there
// wasn't given any faulty or lacking info and
// creates a havoc Plugin object
func (system *System) AddPlugin(register map[string]any, inter any) (*Plugin, error) {
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
	if err = system.CheckAndInsertInterface(ext, inter); err != nil {
		return nil, err
	}

	if err = system.interactPlugin(ext); err != nil {
		return nil, err
	}

	// add the ext to the sync map
	system.loaded.Store(ext.Name, ext)

	return ext, nil
}

func (system *System) Plugin(name string) (*Plugin, error) {
	var (
		val any
		ok  bool
	)

	if val, ok = system.loaded.Load(name); ok {
		return val.(*Plugin), nil
	}

	return nil, fmt.Errorf("plugin not found: %v", name)
}

// CheckAndInsertInterface
// this method checks if a specific plugin is exporting all the
// needed methods for the returned plugin type, if it does then
// cast the right interface to the plugin object
func (system *System) CheckAndInsertInterface(extension *Plugin, inter any) error {
	var (
		reflection = reflect.TypeOf(inter)
		ok         bool
	)

	switch extension.Type {
	case TypeAgent:

		// sanity check if the agent methods exist

		if _, ok = reflection.MethodByName("AgentRegister"); !ok {
			return fmt.Errorf("AgentRegister not found")
		}

		if _, ok = reflection.MethodByName("AgentGenerate"); !ok {
			return fmt.Errorf("AgentGenerate not found")
		}

		if _, ok = reflection.MethodByName("AgentRestore"); !ok {
			return fmt.Errorf("AgentRestore not found")
		}

		if _, ok = reflection.MethodByName("AgentExecute"); !ok {
			return fmt.Errorf("AgentExecute not found")
		}

		if _, ok = reflection.MethodByName("AgentProcess"); !ok {
			return fmt.Errorf("AgentProcess not found")
		}

		if _, ok = reflection.MethodByName("AgentGet"); !ok {
			return fmt.Errorf("AgentGet not found")
		}

		if _, ok = reflection.MethodByName("AgentInterface"); !ok {
			return fmt.Errorf("AgentInterface not found")
		}

		// cast the interface
		// we found everything we searched for
		extension.IPluginAgent = inter.(IPluginAgent)

		break

	case TypeListener:

		// sanity check if the listener methods exist

		if _, ok = reflection.MethodByName("ListenerRegister"); !ok {
			return fmt.Errorf("ListenerRegister not found")
		}

		if _, ok = reflection.MethodByName("ListenerStart"); !ok {
			return fmt.Errorf("ListenerStart not found")
		}

		if _, ok = reflection.MethodByName("ListenerEdit"); !ok {
			return fmt.Errorf("ListenerEdit not found")
		}

		if _, ok = reflection.MethodByName("ListenerRestart"); !ok {
			return fmt.Errorf("ListenerRestart not found")
		}

		if _, ok = reflection.MethodByName("ListenerRestore"); !ok {
			return fmt.Errorf("ListenerRestore not found")
		}

		if _, ok = reflection.MethodByName("ListenerRemove"); !ok {
			return fmt.Errorf("ListenerRemove not found")
		}

		if _, ok = reflection.MethodByName("ListenerStop"); !ok {
			return fmt.Errorf("ListenerStop not found")
		}

		if _, ok = reflection.MethodByName("ListenerEvent"); !ok {
			return fmt.Errorf("ListenerEvent not found")
		}

		if _, ok = reflection.MethodByName("ListenerConfig"); !ok {
			return fmt.Errorf("ListenerConfig not found")
		}

		// cast the interface
		// we found everything we searched for
		extension.IPluginListener = inter.(IPluginListener)

		break

	case TypeManagement:
		break

	default:
		return fmt.Errorf("invalid plugin type: \"%v\" not found", extension.Type)
	}

	return nil
}

// interactPlugin
// interact with plugin by calling the plugin
// register function for the type of plugin
func (system *System) interactPlugin(extension *Plugin) error {

	switch extension.Type {
	case TypeAgent:

		extension.Data = extension.AgentRegister()
		if err := system.havoc.AgentRegisterType(extension.Name, extension.Data); err != nil {
			return err
		}

		break

	case TypeListener:

		extension.Data = extension.ListenerRegister()
		if err := system.havoc.ListenerRegisterType(extension.Name, extension.Data); err != nil {
			return err
		}

		break

	case TypeManagement:
		break

	default:
		return fmt.Errorf("invalid plugin type: \"%v\" not found", extension.Type)
	}

	return nil
}
