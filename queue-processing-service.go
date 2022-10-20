package deezerservice

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/punk-link/logger"
)

type QueueProcessingService struct {
	logger         *logger.Logger
	natsConnection *nats.Conn
}

func (t *QueueProcessingService) StartListen() {
	jetStreamContext, err := t.natsConnection.JetStream()
	subscription, err := t.getSubscription(err, jetStreamContext)
	if err != nil {
		t.logger.LogError(err, err.Error())
		return
	}

	deezerService := NewDeezerService(t.logger)
	for {
		containers := make([]UpcContainer, 10)
		messages, _ := subscription.Fetch(10)
		for i, message := range messages {
			message.Ack()

			var container UpcContainer
			_ = json.Unmarshal(message.Data, &container)

			containers[i] = container
		}

		deezerService.GetReleaseUrlsByUpc(containers)
	}
}

func (t *QueueProcessingService) getSubscription(err error, jetStreamContext nats.JetStreamContext) (*nats.Subscription, error) {
	if err != nil {
		return nil, err
	}

	return jetStreamContext.PullSubscribe(PLATFORM_URL_REQUESTS_STREAM_SUBJECT, PLATFORM_URL_REQUESTS_CONSUMER_NAME)
}
