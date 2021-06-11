package ingest

import (
	"log"
	"strconv"
	"time"

	"github.com/blamelesshq/blameless-examples/slo/packages/clients"
	"github.com/blamelesshq/blameless-examples/slo/packages/config"
	"github.com/blamelesshq/blameless-examples/slo/packages/models"
)

type Ingest interface {
	Backfill(p *clients.PrometheusClient, sliType string, sli *models.SliBody) error
	Regular(p *clients.PrometheusClient, sliType string, sli *models.SliBody) (*models.PostManyResponse, error)
}

func buildModel(id int, tuples []clients.Values, sliType *models.SliTypeBody) ([]models.SliRawDataBody, error) {
	rawDatas := make([]models.SliRawDataBody, len(tuples))
	for i, t := range tuples {
		model := models.SliRawDataBody{
			SliId: id,
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
	return rawDatas, nil
}

// This is expensive to do in a linear programmatic fashion, you should use a distribute queue system for this
func Backfill(p *clients.PrometheusClient, query string, sli *models.SliBody) error {
	bClient := clients.NewBlamelessClient()
	resp, err := sli.GetSliType()
	if err != nil {
		return err
	}

	start := time.Now().AddDate(0, 0, -28).Truncate(24 * time.Hour)
	var now time.Time
	// We support 28 day rolling windows so at max you should only backfill 56 days
	for d := 1; d <= config.Environment().Ingest.Backfill; d++ {
		now = start
		for h := 1; h <= 24; h++ {
			from := now.Add(time.Hour * time.Duration(h-1))
			to := from.Add(time.Hour * time.Duration(h))
			log.Printf("BACKFILLING SLI (ID): %d | QUERY %s | FROM: %s | TO: %s", sli.Id, query, from, to)
			tuples, err := p.QueryRange(query, from, to)
			if err != nil {
				return err
			}
			rawDatas, err := buildModel(sli.Id, tuples, resp.SliType)
			if err != nil {
				return err
			}
			models.PostMany(bClient, rawDatas)
		}
		now = now.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	}

	return nil
}

func Regular(p *clients.PrometheusClient, query string, sli *models.SliBody) (*models.PostManyResponse, error) {
	now := time.Now()
	start := now.Add(time.Duration(-config.Environment().Ingest.Period/60) * time.Minute)
	tuples, err := p.QueryRange(query, start, now)
	if err != nil {
		return nil, err
	}
	resp, err := sli.GetSliType()
	if err != nil {
		return nil, err
	}
	rawDatas, err := buildModel(sli.Id, tuples, resp.SliType)
	if err != nil {
		return nil, err
	}
	results, err := models.PostMany(clients.NewBlamelessClient(), rawDatas)
	if err != nil {
		return nil, err
	}
	return results, nil
}
