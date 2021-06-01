package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/blamelesshq/blameless-examples/slo/packages/config"
	"github.com/go-resty/resty/v2"
)

const apiPrefix = "/api/v1/services"
const sloService = "SLOServiceCrud"
const sloTimeseriesService = "SLOTimeSeriesServiceCrud"

type BlamelessClient struct {
	SloService           string
	SloTimeseriesService string
}

type Blameless interface {
	NewBlamelessClient() *BlamelessClient
	Post(service string, method string, body json.RawMessage) (json.RawMessage, error)
}

func NewBlamelessClient() *BlamelessClient {
	return &BlamelessClient{
		SloService:           sloService,
		SloTimeseriesService: sloTimeseriesService,
	}
}

func (c *BlamelessClient) Post(service string, method string, body json.RawMessage) (json.RawMessage, error) {
	client := resty.New()
	client.SetRetryCount(3).SetRetryWaitTime(5 * time.Second)
	client.SetHostURL(fmt.Sprintf("%s:%d", config.Environment().Blameless.Host, config.Environment().Blameless.Port))
	client.SetHeader("Accept", "application/json")

	client.SetAuthScheme("Bearer")
	client.SetAuthToken(config.Environment().Blameless.AuthToken)

	if c.SloService != service && c.SloTimeseriesService != service {
		return json.RawMessage{}, fmt.Errorf("service name %s is not a valid identifier", service)
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("%s/%s/%s", apiPrefix, service, method))
	if err != nil {
		return json.RawMessage{}, fmt.Errorf("unable to perform request \n%+v", err)
	}
	log.Println(resp.String())
	return resp.Body(), nil
}
