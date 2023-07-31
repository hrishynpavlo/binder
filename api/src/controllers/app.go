package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppController struct {
}

func (controller AppController) GetAppRevision(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"commitSha": "local"})
}

func (controller AppController) RegisterAppEndpoints(router *gin.Engine) {
	router.GET("/app-revision", controller.GetAppRevision)
}

func ProvideAppController() *AppController {
	return &AppController{}
}
