package controllers

import "github.com/gin-gonic/gin"

type Controllers struct {
	AppController  *AppController
	UserController *UserController
	FeedController *FeedController
}

func ProvideControllers(app *AppController, user *UserController, feed *FeedController) *Controllers {
	return &Controllers{AppController: app, UserController: user, FeedController: feed}
}

func (controllers Controllers) RegisterAllEndpoints(router *gin.Engine) {
	controllers.AppController.RegisterAppEndpoints(router)
	controllers.UserController.RegisterUserEndpoints(router)
	controllers.FeedController.RegisterFeedEndpoints(router)
}
