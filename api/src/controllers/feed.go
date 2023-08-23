package controllers

import (
	"binder_api/controllers/auth"
	"binder_api/db"
	"binder_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FeedController struct {
	repo        *db.FeedRepository
	feed        services.FeedProvider
	authService *auth.AuthService
}

func ProvideFeedController(repository *db.FeedRepository, feed services.FeedProvider, auth *auth.AuthService) *FeedController {
	return &FeedController{repo: repository, feed: feed, authService: auth}
}

func (controller FeedController) RegisterFeedEndpoints(router *gin.Engine) {
	router.GET("/api/feed", controller.authService.AuthMiddleware, controller.getUserFeed)
}

func (controller FeedController) getUserFeed(c *gin.Context) {
	id := c.GetInt64("userId")

	feed := controller.feed.GetOrAdd(id)

	c.JSON(http.StatusOK, feed)
	return
}
