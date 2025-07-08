package main

import (
	"log"

	"github.com/nuea/backend-golang-test/cmd/http/di"
)

//go:generate go run github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal --parseDepth 1 -o internal/docs

// @title Backend Golang Test
// @version 1.0
// @description API for http gateway
//
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey DeviceID
// @in header
// @name X-Device-Id
//
// @query.collection.format multi
func main() {
	ctn, stop, err := di.InitContainer()
	if err != nil {
		log.Panicf("Unable to start service. Error: %s", err)
	}
	defer stop()
	ctn.Run()
}
