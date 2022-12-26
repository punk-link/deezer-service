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

	appSecrets := common.GetAppSecrets(envManager, logger, SECRET_ENGINE_NAME, SERVICE_NAME)
	serviceOptions := runtime.NewServiceOptions(logger, appSecrets, environmentName, SERVICE_NAME)

	deezerService := services.NewDeezerService(logger, httpclient.DefaultConfig(logger))
	go startup.ProcessUrls(serviceOptions, deezerService)

	startup.RunServer(serviceOptions)
}

const SECRET_ENGINE_NAME = "secrets"
const SERVICE_NAME = "deezer-service"
