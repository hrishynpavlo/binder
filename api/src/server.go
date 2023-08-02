package main

import (
	"binder_api/caching"
	"binder_api/configuration"
	"binder_api/controllers"
	"binder_api/db"
	"binder_api/logging"
	"binder_api/workers"

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
		fx.Provide(caching.ProvideRedis),
		fx.Provide(controllers.ProvideAppController, controllers.ProvideUserController, controllers.ProvideControllers),
		fx.Provide(workers.ProvideMatcherWorker, workers.ProvideUserRegisteredChannel, workers.ProvideGeoMatcherWorker),
		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(logger *zap.Logger, controllers *controllers.Controllers, router *gin.Engine, matcher *workers.MatcherWorker, geoMatcher *workers.GeoMatcherWorker) {

	controllers.RegisterAllEndpoints(router)
	logger.Debug("All endpoints registered")

	go geoMatcher.StartWorker()

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("ERROR on server start", zap.Error(err))
	}
}
