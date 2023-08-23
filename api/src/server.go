package main

import (
	"binder_api/caching"
	"binder_api/configuration"
	"binder_api/controllers"
	"binder_api/controllers/auth"
	"binder_api/db"
	"binder_api/logging"
	"binder_api/services"
	"binder_api/workers"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func ProvideGinEngine(logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	return router
}

func main() {

	app := fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(logging.ProviderLogger),
		fx.Provide(configuration.ProvideConfiguration),
		fx.Provide(ProvideGinEngine),
		fx.Provide(db.ProvideDb, db.ProvideUserRepository, db.ProvideFeedRepository),
		fx.Provide(caching.ProvideRedis),
		fx.Provide(services.ProvideFeedService),
		fx.Provide(auth.ProvideAuthService, controllers.ProvideAppController, controllers.ProvideUserController, controllers.ProvideFeedController, controllers.ProvideControllers),
		fx.Provide(workers.ProvideUserRegisteredChannel, workers.ProvideGeoMatcherWorker),
		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(logger *zap.Logger, controllers *controllers.Controllers, router *gin.Engine, geoWorker *workers.GeoWorker) {

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	controllers.RegisterAllEndpoints(router)
	logger.Debug("All endpoints registered")

	go geoWorker.StartGeoEnrichment()

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("ERROR on server start", zap.Error(err))
	}
}
