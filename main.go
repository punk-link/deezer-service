package deezerservice

import (
	"main/infrastructure"

	"github.com/nats-io/nats.go"
)

type UpcContainer struct {
	Id  int
	Upc string
}

func main() {
	logger := infrastructure.NewLoggerWithoutInjection()

	environmentName := "Local"
	logger.LogInfo("Deezer API is running as '%s'", environmentName)

	natsConnection, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.LogError(err, err.Error())
	}

}

const PLATFORM_URL_REQUESTS_STREAM_SUBJECT = "PLATFORM-URL-REQUESTS.deezer"
const PLATFORM_URL_REQUESTS_CONSUMER_NAME = "deezer-platform-request-consumer"
