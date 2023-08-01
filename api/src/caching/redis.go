package caching

import (
	"binder_api/configuration"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ProvideRedis(logger *zap.Logger, config *configuration.AppConfiguration) *redis.Client {
	options, err := redis.ParseURL(config.RedisConnectionString)
	if err != nil {
		logger.Fatal("Error in redis options", zap.Error(err))
	}
	return redis.NewClient(options)
}
