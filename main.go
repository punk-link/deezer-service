package main

import (
	"main/services"

	httpclient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"

	envManager "github.com/punk-link/environment-variable-manager"
	runtime "github.com/punk-link/streaming-platform-runtime"
	common "github.com/punk-link/streaming-platform-runtime/common"
	"github.com/punk-link/streaming-platform-runtime/startup"
)

func main() {
	logger := logger.New()
	envManager := envManager.New()

	environmentName := common.GetEnvironmentName(envManager)
	logger.LogInfo("%s is running as '%s'", SERVICE_NAME, environmentName)

	serviceOptions := runtime.NewServiceOptions(logger, envManager, environmentName, SERVICE_NAME)

	deezerService := services.NewDeezerService(logger, httpclient.DefaultConfig(logger))
	go startup.ProcessUrls(serviceOptions, deezerService)

	startup.RunServer(serviceOptions)
}

const SERVICE_NAME = "deezer-service"
