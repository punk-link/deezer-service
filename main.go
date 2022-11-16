package main

import (
	"context"
	"errors"
	"fmt"
	"main/services"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	consulclient "github.com/punk-link/consul-client"
	httpclient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	contracts "github.com/punk-link/platform-contracts"
	common "github.com/punk-link/streaming-platform-runtime/common"
	processing "github.com/punk-link/streaming-platform-runtime/processing"
	"github.com/punk-link/streaming-platform-runtime/startup"
)

func main() {
	logger := logger.New()
	environmentName := common.GetEnvironmentName()
	logger.LogInfo("%s is running as '%s'", SERVICE_NAME, environmentName)
	consul := common.GetConsulClient(logger, SERVICE_NAME, environmentName)

	deezerService := services.NewDeezerService(logger, httpclient.DefaultConfig(logger))
	process(logger, consul, deezerService)
}

func process(logger logger.Logger, consul *consulclient.ConsulClient, service contracts.Platformer) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	hostSettingsValues, err := consul.Get("HostSettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain host settings from Consul: '%s'", err.Error())
	}
	hostSettings := hostSettingsValues.(map[string]any)

	app := startup.Configure(logger, consul, &startup.StartupOptions{
		EnvironmentName: "",
		GinMode:         hostSettings["Mode"].(string),
		ServiceName:     SERVICE_NAME,
	})
	app.Run()

	hostAddress := hostSettings["Address"]
	hostPort := hostSettings["Port"]
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", hostAddress, hostPort),
		Handler: app,
	}

	go func() {
		logger.LogInfo("Starting...")
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.LogError(err, "Listen error: %s\n", err.Error())
		}
	}()

	natsConnection := common.GetNatsConnection(logger, consul)
	queueProcessingService := processing.New(logger, natsConnection)
	go queueProcessingService.Process(ctx, &wg, service)

	wg.Wait()
	logger.LogInfo("Exiting...")
}

const SERVICE_NAME = "deezer-service"
