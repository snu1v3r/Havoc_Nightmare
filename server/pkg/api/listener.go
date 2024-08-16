package api

import (
	"Havoc/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (api *ServerApi) listenerStart(ctx *gin.Context) {
	var (
		body     []byte
		err      error
		listener map[string]any
		name     string
		protocol string
		options  map[string]any
	)

	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// read from request the login data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("Failed to read from server api login request: " + err.Error())
		goto ERROR
	}

	logger.Debug("got request on /api/listener/start:" + fmt.Sprintf("%s", string(body)))

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &listener); err != nil {
		logger.DebugError("Failed to unmarshal bytes to map: " + err.Error())
		return
	}

	if val, ok := listener["name"]; ok {
		// get name from listener start request
		switch val.(type) {
		case string:
			name = val.(string)
		default:
			logger.DebugError("Failed retrieve name: invalid type")
			goto ERROR
		}
	} else {
		logger.DebugError("failed retrieve listener name: not found")
		goto ERROR
	}

	if val, ok := listener["protocol"]; ok {
		// get protocol from listener start request
		switch val.(type) {
		case string:
			protocol = val.(string)
		default:
			logger.DebugError("Failed retrieve protocol: invalid type")
			goto ERROR
		}
	} else {
		logger.DebugError("failed retrieve listener protocol: not found")
		goto ERROR
	}

	if val, ok := listener["data"]; ok {
		// get options from listener start request
		switch val.(type) {
		case map[string]any:
			options = val.(map[string]any)
		default:
			logger.DebugError("Failed retrieve listener data: invalid type")
			goto ERROR
		}
	}

	if err = api.teamserver.ListenerStart(name, protocol, options); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
	return

ERROR:
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func (api *ServerApi) listenerStop(ctx *gin.Context) {
	var (
		body     []byte
		err      error
		name     string
		listener map[string]any
	)

	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// read from request the login data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("Failed to read from server api login request: " + err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	logger.Debug("got request on /api/listener/stop:" + fmt.Sprintf("%s", string(body)))

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &listener); err != nil {
		logger.DebugError("Failed to unmarshal bytes to map: " + err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if val, ok := listener["name"]; ok {
		// get name from listener start request
		switch val.(type) {
		case string:
			name = val.(string)
		default:
			logger.DebugError("failed retrieve name: invalid type")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		logger.DebugError("failed retrieve listener name: not found")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = api.teamserver.ListenerStop(name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
	return
}

func (api *ServerApi) listenerRestart(ctx *gin.Context) {
	var (
		body     []byte
		err      error
		name     string
		listener map[string]any
	)

	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// read from request the login data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("Failed to read from server api login request: " + err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	logger.Debug("got request on /api/listener/restart: " + fmt.Sprintf("%s", string(body)))

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &listener); err != nil {
		logger.DebugError("Failed to unmarshal bytes to map: " + err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if val, ok := listener["name"]; ok {
		// get name from listener start request
		switch val.(type) {
		case string:
			name = val.(string)
		default:
			logger.DebugError("failed retrieve name: invalid type")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		logger.DebugError("failed retrieve listener name: not found")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = api.teamserver.ListenerRestart(name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
	return
}

func (api *ServerApi) listenerRemove(ctx *gin.Context) {
	var (
		body     []byte
		err      error
		name     string
		listener map[string]any
	)

	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// read from request the login data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("Failed to read from server api login request: " + err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	logger.Debug("got request on /api/listener/remove: " + fmt.Sprintf("%s", string(body)))

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &listener); err != nil {
		logger.DebugError("Failed to unmarshal bytes to map: " + err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if val, ok := listener["name"]; ok {
		// get name from listener start request
		switch val.(type) {
		case string:
			name = val.(string)
		default:
			logger.DebugError("failed retrieve name: invalid type")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		logger.DebugError("failed retrieve listener name: not found")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = api.teamserver.ListenerRemove(name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
	return
}

func (api *ServerApi) listenerEdit(ctx *gin.Context) {
	var (
		body   []byte
		err    error
		data   map[string]any
		config map[string]any
		name   string
	)

	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// read from request the event data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("Failed to read from server api login request: " + err.Error())
		goto ERROR
	}

	logger.Debug("got request on /api/listener/edit: " + fmt.Sprintf("%s", string(body)))

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &data); err != nil {
		logger.DebugError("Failed to unmarshal bytes to map: " + err.Error())
		return
	}

	if val, ok := data["name"]; ok {
		// get protocol from listener start request
		switch val.(type) {
		case string:
			name = val.(string)
		default:
			logger.DebugError("failed retrieve listener name: invalid type")
			goto ERROR
		}
	} else {
		logger.DebugError("failed retrieve listener name: not found")
		goto ERROR
	}

	if val, ok := data["data"]; ok {
		// get protocol from listener start request
		switch val.(type) {
		case map[string]any:
			config = val.(map[string]any)
		default:
			logger.DebugError("failed retrieve listener data: invalid type")
			goto ERROR
		}
	} else {
		logger.DebugError("failed retrieve listener data: not found")
		goto ERROR
	}

	// process listener event
	if err = api.teamserver.ListenerEdit(name, config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
	return

ERROR:
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func (api *ServerApi) listenerEvent(ctx *gin.Context) {
	var (
		body     []byte
		err      error
		event    map[string]any
		protocol string
	)

	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// read from request the event data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("Failed to read from server api login request: " + err.Error())
		goto ERROR
	}

	logger.Debug("got request on /api/listener/event:" + fmt.Sprintf("%s", string(body)))

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &event); err != nil {
		logger.DebugError("Failed to unmarshal bytes to map: " + err.Error())
		return
	}

	if val, ok := event["protocol"]; ok {
		// get protocol from listener start request
		switch val.(type) {
		case string:
			protocol = val.(string)
		default:
			logger.DebugError("Failed retrieve protocol: invalid type")
			goto ERROR
		}
	} else {
		logger.DebugError("failed retrieve listener protocol: not found")
		goto ERROR
	}

	// process listener event
	if event, err = api.teamserver.ListenerEvent(protocol, event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, event)
	return

ERROR:
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func (api *ServerApi) listenerConfig(ctx *gin.Context) {
	var (
		body   []byte
		err    error
		data   map[string]any
		config map[string]any
		name   string
	)

	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// read from request the event data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("Failed to read from server api login request: " + err.Error())
		goto ERROR
	}

	logger.Debug("got request on /api/listener/config:" + fmt.Sprintf("%s", string(body)))

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &data); err != nil {
		logger.DebugError("Failed to unmarshal bytes to map: " + err.Error())
		return
	}

	if val, ok := data["name"]; ok {
		// get protocol from listener start request
		switch val.(type) {
		case string:
			name = val.(string)
		default:
			logger.DebugError("failed retrieve listener name: invalid type")
			goto ERROR
		}
	} else {
		logger.DebugError("failed retrieve listener name: not found")
		goto ERROR
	}

	// process listener event
	if config, err = api.teamserver.ListenerConfig(name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, config)
	return

ERROR:
	ctx.AbortWithStatus(http.StatusInternalServerError)
}
