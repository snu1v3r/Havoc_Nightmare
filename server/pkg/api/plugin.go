package api

import (
	"Havoc/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *ServerApi) pluginList(ctx *gin.Context) {
	logger.Debug("got request on /api/plugin/list")

	ctx.JSON(http.StatusOK, api.PluginList())
	return
}

func (api *ServerApi) pluginFile(ctx *gin.Context) {
	if !api.sanityCheck(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	logger.Debug("got request on /api/plugin/list")

}
