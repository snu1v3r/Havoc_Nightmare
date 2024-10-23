package plugin

import (
	"errors"
	"strings"
)

var (
	errAlreadyRegistered = errors.New("plugin already has been registered")
	errNotFound          = errors.New("plugin not found")
)

func IsAlreadyRegistered(err error) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(err.Error(), errAlreadyRegistered.Error())
}

func IsPluginNotFound(err error) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(err.Error(), errNotFound.Error())
}
