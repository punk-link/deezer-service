package main

import (
	"fmt"
	"main/services"

	"github.com/nats-io/nats.go"
	consulClient "github.com/punk-link/consul-client"
	envManager "github.com/punk-link/environment-variable-manager"
	httpclient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	platformContracts "github.com/punk-link/platform-contracts"
)

func main() {
	logger := logger.New()

	environmentName := getEnvironmentName()
	logger.LogInfo("Deezer API is running as '%s'", environmentName)

	consul, _ := getConsulClient("deezer-service", environmentName)
	natsSettingsValues, err := consul.Get("NatsSettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain Nats settings from Consul: '%s'", err.Error())
		return
	}
	natsSettings := natsSettingsValues.(map[string]interface{})

	natsConnection, err := nats.Connect(natsSettings["Url"].(string))
	if err != nil {
		logger.LogError(err, err.Error())
	}

	deezerService := services.NewDeezerService(logger, *httpclient.DefaultConfig(logger))
	queueProcessingService := platformContracts.NewQueueProcessingService(logger, natsConnection)
	queueProcessingService.Process(deezerService)
}

func getConsulClient(storageName string, environmentName string) (*consulClient.ConsulClient, error) {
	isExist, consulAddress := envManager.TryGetEnvironmentVariable("PNKL_CONSUL_ADDR")
	if !isExist {
		return nil, fmt.Errorf("can't find value of the '%s' environment variable", "PNKL_CONSUL_ADDR")
	}

	isExist, consulToken := envManager.TryGetEnvironmentVariable("PNKL_CONSUL_TOKEN")
	if !isExist {
		return nil, fmt.Errorf("can't find value of the '%s' environment variable", "PNKL_CONSUL_TOKEN")
	}

	consul, err := consulClient.New(&consulClient.ConsulConfig{
		Address:         consulAddress,
		EnvironmentName: environmentName,
		StorageName:     storageName,
		Token:           consulToken,
	})

	return consul, err
}

func getEnvironmentName() string {
	isExist, name := envManager.TryGetEnvironmentVariable("GO_ENVIRONMENT")
	if !isExist {
		return "Development"
	}

	return name
}
