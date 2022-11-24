package services

import (
	"fmt"
	"main/models"
	"time"

	httpclient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	platformContracts "github.com/punk-link/platform-contracts"
)

type DeezerService struct {
	httpClientConfig *httpclient.HttpClientConfig
	logger           logger.Logger
}

func NewDeezerService(logger logger.Logger, httpClientConfig *httpclient.HttpClientConfig) *DeezerService {
	return &DeezerService{
		httpClientConfig: httpClientConfig,
		logger:           logger,
	}
}

func (t *DeezerService) GetBatchSize() int {
	return 10
}

func (t *DeezerService) GetPlatformName() string {
	return platformContracts.Deezer
}

func (t *DeezerService) GetReleaseUrlsByUpc(upcContainers []platformContracts.UpcContainer) []platformContracts.UrlResultContainer {
	results := make([]platformContracts.UrlResultContainer, 0)
	attemptNumber := 3
	for _, container := range upcContainers {
	REQUEST:
		var response models.UpcResponse
		err := makeRequest(t.logger, t.httpClientConfig, "GET", fmt.Sprintf("album/upc:%s", container.Upc), &response)
		if err != nil {
			continue
		}

		if response.Error.Code == 4 || response.Error.Code == 100 || response.Error.Code == 700 {
			if attemptNumber == 0 {
				continue
			}

			attemptNumber--
			time.Sleep(RETRY_TIMEOUT_DURATION)
			goto REQUEST
		}

		if response.Error.Code != 0 {
			continue
		}

		results = append(results, platformContracts.UrlResultContainer{
			Id:           container.Id,
			PlatformName: platformContracts.Deezer,
			Upc:          container.Upc,
			Url:          response.Url,
		})

		attemptNumber = 3
		time.Sleep(REQUEST_TIMEOUT_DURATION)
	}

	return results
}

const REQUEST_TIMEOUT_DURATION = time.Millisecond * 500
const RETRY_TIMEOUT_DURATION = time.Second * 5
