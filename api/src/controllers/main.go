package controllers

import "github.com/gin-gonic/gin"

type Controllers struct {
	AppController  *AppController
	UserController *UserController
}

func ProvideControllers(app *AppController, user *UserController) *Controllers {
	return &Controllers{AppController: app, UserController: user}
}

func (controllers Controllers) RegisterAllEndpoints(router *gin.Engine) {
	controllers.AppController.RegisterAppEndpoints(router)
	controllers.UserController.RegisterUserEndpoints(router)
}
