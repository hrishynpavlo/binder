package main

import (
	"binder_api/configuration"
	"binder_api/controllers"
	"binder_api/db"
	"binder_api/logging"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {

	app := fx.New(
		fx.Provide(logging.ProviderLogger),
		fx.Provide(configuration.ProvideConfiguration),
		fx.Provide(gin.Default),
		fx.Provide(db.ProvideDb),
		fx.Provide(controllers.ProvideAppController),
		fx.Provide(controllers.ProvideUserController),
		fx.Provide(controllers.ProvideControllers),
		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(logger *zap.Logger, controllers *controllers.Controllers, router *gin.Engine) {

	controllers.RegisterAllEndpoints(router)
	logger.Debug("All endpoints registered")
	router.Run(":8080")
	logger.Info("Server successfully started")
}
