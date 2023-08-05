package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedController struct {
}

func ProvideFeedController() *FeedController {
	return &FeedController{}
}

func (controller FeedController) RegisterFeedEndpoints(router *gin.Engine) {
	router.GET("/api/user/:id/feed", controller.getUserFeed)
}

func (controller FeedController) getUserFeed(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)

	c.JSON(http.StatusOK, id)
}
