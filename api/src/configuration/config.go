package configuration

import "github.com/kelseyhightower/envconfig"

type AppConfiguration struct {
	DbConnectionString string `envconfig:"BINDER_DB_CONNECTION_STRING" default:"postgres://binder_usr:binder_best_app@localhost:5432/binder_all?sslmode=disable"`
	CommitRevision     string `envconfig:"BINDER_COMMIT_REVISION" default:"local"`
}

func ProvideConfiguration() *AppConfiguration {
	var config AppConfiguration
	envconfig.MustProcess("BINDER", &config)
	return &config
}
