package controllers

import (
	"binder_api/configuration"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AppController struct {
	Configuration *configuration.AppConfiguration
	StartTime     time.Time
}

func (controller AppController) GetAppRevision(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"commitSha": controller.Configuration.CommitRevision, "startTime": controller.StartTime})
}

func (controller AppController) RegisterAppEndpoints(router *gin.Engine) {
	router.GET("/app-revision", controller.GetAppRevision)
}

func ProvideAppController(config *configuration.AppConfiguration) *AppController {
	return &AppController{Configuration: config, StartTime: time.Now().UTC()}
}
