package configuration

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type AppConfiguration struct {
	DbConnectionString string `envconfig:"BINDER_DB_CONNECTION_STRING" default:"postgres://binder_app:binder_best_app@localhost:5432/binder_all?sslmode=disable"`
	CommitRevision     string `envconfig:"BINDER_COMMIT_REVISION" default:"local"`
}

func ProvideConfiguration(logger *zap.Logger) *AppConfiguration {
	var config AppConfiguration
	if err := envconfig.Process("BINDER", &config); err != nil {
		logger.Fatal("Reading of environment configuration failed")
		panic("EXIT")
	}
	return &config
}
