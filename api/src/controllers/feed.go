package controllers

import (
	"binder_api/db"
	"binder_api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedController struct {
	repo *db.FeedRepository
	feed services.FeedProvider
}

func ProvideFeedController(repository *db.FeedRepository, feed services.FeedProvider) *FeedController {
	return &FeedController{repo: repository, feed: feed}
}

func (controller FeedController) RegisterFeedEndpoints(router *gin.Engine) {
	router.GET("/api/user/:id/feed", controller.getUserFeed)
}

func (controller FeedController) getUserFeed(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)

	feed := controller.feed.GetOrAdd(id)

	c.JSON(http.StatusOK, feed)
	return
}
