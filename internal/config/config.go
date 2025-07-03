package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type SystemConfig struct {
	HTTPPort    string `envconfig:"APP_HTTP_PORT" default:"8080"`
	GRPCPort    string `envconfig:"APP_GRPC_PORT" default:"8980"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"backend-golang-test"`
}

type AuthConfig struct {
	SecretKey            string        `envconfig:"AUTH_SECRET_KEY"`
	AccessTokenExpireTTL time.Duration `envconfig:"AUTH_ACCESS_TOKEN_EXPIRE_TTL" default:"5m"`
}

type MongoDBConfig struct {
	Host              string        `envconfig:"MONGODB_HOST"`
	User              string        `envconfig:"MONGODB_USER"`
	Password          string        `envconfig:"MONGODB_PASSWORD"`
	DatabaseName      string        `envconfig:"MONGODB_DATABASE_NAME" default:"-"`
	HeartbeatInterval time.Duration `envconfig:"MONGODB_HEARTBEAT_INTERVAL" default:"10s"`
	MaxPoolSize       uint64        `envconfig:"MONGODB_MAX_CONNECTION_POOL_SIZE" default:"20"`
	MinPoolSize       uint64        `envconfig:"MONGODB_MIN_CONNECTION_POOL_SIZE" default:"10"`
}

type BackendGolangTestGRPCConfig struct {
	GRPCTarget     string        `envconfig:"BACKEND_GOLANG_TEST_GRPC_TARGET" default:"localhost:8980"`
	RequestTimeout time.Duration `envconfig:"BACKEND_GOLANG_TEST_REQUEST_TIMEOUT" default:"10s"`
}

type AppConfig struct {
	System        SystemConfig
	MongoDB       MongoDBConfig
	BackendGoTest BackendGolangTestGRPCConfig
	Auth          AuthConfig
}

func (cfg *AppConfig) load() {
	envconfig.MustProcess("", &cfg.System)
	envconfig.MustProcess("", &cfg.Auth)
	envconfig.MustProcess("", &cfg.MongoDB)
	envconfig.MustProcess("", &cfg.BackendGoTest)
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
