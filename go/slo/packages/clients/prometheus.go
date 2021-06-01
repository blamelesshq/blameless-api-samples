package clients

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/blamelesshq/blameless-examples/slo/packages/config"
	"github.com/go-resty/resty/v2"
)

var prometheusOncer sync.Once
var pClient *resty.Client

type PrometheusClient struct {
	client *resty.Client
}

type Prometheus interface {
	NewPrometheusClient() *PrometheusClient
	QueryRange(query string, start time.Time, end time.Time)
}

type Values struct {
	Time  int
	Value string
}

type QueryRangeResponse struct {
	Status string
	Data   struct {
		Result []struct {
			Values []json.RawMessage `json:"values"`
		}
	}
}

func formatTime(t time.Time) string {
	return strconv.FormatFloat(float64(t.Unix())+float64(t.Nanosecond())/1e9, 'f', -1, 64)
}

func NewPrometheusClient() *PrometheusClient {
	prometheusOncer.Do(func() {
		pClient := resty.New()
		pClient.SetRetryCount(3).SetRetryWaitTime(5 * time.Second)
		pClient.SetHostURL(fmt.Sprintf("%s:%d", config.Environment().Prometheus.Host, config.Environment().Prometheus.Port))
		pClient.SetHeader("Accept", "application/json")

		// If you need to add auth
		// // BASIC
		// client.SetBasicAuth("user", "password")
		// // Auth Token
		// client.SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
	})

	return &PrometheusClient{
		client: pClient,
	}
}

func (v *Values) UnmarshalJSON(b []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	v.Time = int(tmp[0].(float64))
	v.Value = tmp[1].(string)
	return nil
}

func (p *PrometheusClient) QueryRange(query string, start time.Time, end time.Time) ([]Values, error) {
	step := time.Duration(config.Environment().Ingest.Step) * time.Second

	resp, err := p.client.R().SetQueryParams(map[string]string{
		"query": query,
		"start": formatTime(start),
		"end":   formatTime(end),
		"step":  step.String(),
	}).Get("/api/v1/query_range")

	if err != nil {
		return []Values{}, fmt.Errorf("error querying Prometheus instance %s\nError: %v", fmt.Sprintf("%s:%d", config.Environment().Prometheus.Host, config.Environment().Prometheus.Port), err)
	}

	if resp.StatusCode() != 200 {
		return []Values{}, fmt.Errorf("unsuccessful query")
	}

	var results *QueryRangeResponse
	if err := json.Unmarshal(resp.Body(), &results); err != nil {
		return []Values{}, fmt.Errorf("unable to successfully unmarshall: \n%v", err)
	}

	dataTuples := results.Data.Result[0].Values
	tuples := make([]Values, len(dataTuples))
	for i := 0; i < len(dataTuples); i++ {
		var v Values
		if err := json.Unmarshal(dataTuples[i], &v); err != nil {
			return []Values{}, err
		}
		tuples[i] = v
	}

	return tuples, nil
}
