package ingest

import (
	"strconv"
	"time"

	"github.com/blamelesshq/blameless-examples/slo/packages/clients"
	"github.com/blamelesshq/blameless-examples/slo/packages/config"
	"github.com/blamelesshq/blameless-examples/slo/packages/models"
)

type Ingest interface {
	Backfill(b *clients.BlamelessClient, p *clients.PrometheusClient, query string, sli *models.SliBody) error
	Regular(b *clients.BlamelessClient, p *clients.PrometheusClient, query string, sli *models.SliBody) (*models.PostManyResponse, error)
}

func IngestData() {

}

// This is expensive to do in a linear programmatic fashion, you should use a distribute queue system for this
func Backfill(b *clients.BlamelessClient, p *clients.PrometheusClient, query string, sli *models.SliBody) error {
	// We support 28 day rolling windows so at max you should only backfill 56 days
	for i := 1; i <= 28; i++ {

	}
	return nil
}

func Regular(b *clients.BlamelessClient, p *clients.PrometheusClient, query string, sli *models.SliBody) (*models.PostManyResponse, error) {
	now := time.Now()
	start := now.Add(time.Duration(-config.Environment().Ingest.Period/60) * time.Minute)
	tuples, err := p.QueryRange(query, start, now)
	if err != nil {
		return nil, err
	}
	rawDatas := make([]models.SliRawDataBody, len(tuples))
	sliType, err := sli.GetSliType(b)
	if err != nil {
		return nil, err
	}
	for i, t := range tuples {
		model := models.SliRawDataBody{
			SliId: sli.Id,
			Start: t.Time,
			End:   t.Time * config.Environment().Ingest.Step,
		}
		value, err := strconv.Atoi(t.Value)
		if err != nil {
			return nil, err
		}

		switch name := sliType.Name; name {
		case models.Types.Latency:
			model.Latency = value
		case models.Types.Throughput:
			model.Throughput = value
		case models.Types.Saturation:
			model.Saturation = value
		case models.Types.Durability:
			model.Durability = value
		case models.Types.Correctness:
			model.Correctness = value
		}
		rawDatas[i] = model
	}
	results, err := models.PostMany(b, rawDatas)
	if err != nil {
		return nil, err
	}
	return results, nil
}
