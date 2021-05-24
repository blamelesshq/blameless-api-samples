package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/blamelesshq/blameless-examples/slo/packages/config"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type PrometheusClient struct {
	client v1.API
}

type Prometheus interface {
	QueryRange(s string, start time.Time, end time.Time)
}

func client() *PrometheusClient {
	client, err := api.NewClient(api.Config{
		Address: fmt.Sprintf("%s:%i", config.Prometheus.Host, config.Prometheus.Port),

	})
	if err != nil {

	}

	v1Api := v1.NewAPI(client)
	return &PrometheusClient{
		client: v1Api,
	}
}

func (p *PrometheusClient) QueryRange(s string, start time.Time, end time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v1.Range{
		Start: ,
	}
	// p.client.QueryRange()
}