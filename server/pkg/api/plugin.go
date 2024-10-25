package api

import (
	"Havoc/pkg/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"Havoc/pkg/logger"
)

func (api *ServerApi) pluginList(ctx *gin.Context) {
	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	logger.Debug("got request on /api/plugin/list")

	ctx.JSON(http.StatusOK, api.PluginList())
	return
}

func (api *ServerApi) pluginResource(ctx *gin.Context) {
	var (
		err      = errors.New("invalid request")
		data     map[string]any
		body     []byte
		name     string
		resource string
	)

	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	logger.Debug("got request on /api/plugin/resource")

	// read from request the login data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("failed to read from server api login request: " + err.Error())
		goto ERROR
	}

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &data); err != nil {
		logger.DebugError("failed to unmarshal bytes to map: " + err.Error())
		err = errors.New("invalid request")
		goto ERROR
	}

	// retrieve the plugin name from the request
	if name, err = utils.MapKey[string](data, "name"); err != nil {
		logger.DebugError("failed to parse name from request: %v", err)
		goto ERROR
	}

	// retrieve the file resource name to retrieve from the plugin
	if resource, err = utils.MapKey[string](data, "resource"); err != nil {
		logger.DebugError("failed to parse file from request: %v", err)
		goto ERROR
	}

	// retrieve the resource file from the registered plugins
	if body, err = api.PluginResource(name, resource); err != nil {
		logger.DebugError("failed to retrieve plugin resource (%v : %v): %v", name, resource, err)
		goto ERROR
	}

	// prepare request headers and response
	ctx.Data(http.StatusOK, "application/octet-stream", body)

	return
ERROR:
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
