package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type SystemConfig struct {
	GRPCPort    string `envconfig:"APP_GRPC_PORT" default:"8980"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"backend-golang-test"`
}

type AppConfig struct {
	System SystemConfig
}

func (cfg *AppConfig) load() {
	envconfig.MustProcess("", &cfg.System)
}

func ProvideCofig() *AppConfig {
	env, ok := os.LookupEnv("ENV")
	if ok && env != "" {
		_, b, _, _ := runtime.Caller(0)
		basePath := filepath.Dir(b)
		err := godotenv.Load(fmt.Sprintf("%v/../../.env.%v", basePath, env))
		if err != nil {
			err = godotenv.Load()
			if err != nil {
				panic(err)
			}
		}
	}
	cfg := &AppConfig{}
	cfg.load()
	return cfg
}
