package main

import (
	"binder_api/controllers"
	"binder_api/db"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		fx.Provide(gin.Default),
		fx.Provide(db.ProvideDb),
		fx.Provide(controllers.ProvideAppController),
		fx.Provide(controllers.ProvideUserController),
		fx.Provide(controllers.ProvideControllers),
		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(controllers *controllers.Controllers, router *gin.Engine) {

	controllers.RegisterAllEndpoints(router)
	router.Run(":8080")
}
