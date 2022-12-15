package services

import (
	"fmt"
	"main/models"
	"net/http"
	"time"

	httpclient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	platformContracts "github.com/punk-link/platform-contracts"
)

type DeezerService struct {
	httpClient httpclient.HttpClient[models.UpcResponse]
	logger     logger.Logger
}

func NewDeezerService(logger logger.Logger, httpClientConfig *httpclient.HttpClientConfig) platformContracts.Platformer {
	httpClient := httpclient.New[models.UpcResponse](httpClientConfig)

	return &DeezerService{
		httpClient: httpClient,
		logger:     logger,
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
		response, err := t.makeRequest("GET", fmt.Sprintf("album/upc:%s", container.Upc))
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

func (t *DeezerService) makeRequest(method string, url string) (*models.UpcResponse, error) {
	request, err := getRequest(method, url)
	if err != nil {
		t.logger.LogWarn("can't build an http request: %s", err.Error())
		return nil, err
	}

	return t.httpClient.MakeRequest(request)
}

func getRequest(method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, "http://api.deezer.com/"+url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

const REQUEST_TIMEOUT_DURATION = time.Millisecond * 500
const RETRY_TIMEOUT_DURATION = time.Second * 5
