package models

import (
	"encoding/json"

	"github.com/blamelesshq/blameless-examples/slo/packages/clients"
	"github.com/blamelesshq/blameless-examples/slo/packages/config"
)

type SliRawDataBody struct {
	SliId        int `json:"sliId" binding:"required"`
	Latency      int `json:"latency,omitempty"`
	ValidRequest int `json:"validRequest,omitempty"`
	GoodRequest  int `json:"goodRequest,omitempty"`
	Throughput   int `json:"throughput,omitempty"`  // Ensure supported in Blameless frontend product before using
	Correctness  int `json:"correctness,omitempty"` // Ensure supported in Blameless frontend product before using
	Saturation   int `json:"saturation,omitempty"`  // Ensure supported in Blameless frontend product before using
	Durability   int `json:"durability,omitempty"`  // Ensure supported in Blameless frontend product before using
	Start        int `json:"start" binding:"required"`
	End          int `json:"end" binding:"required"`
}

type PostManyRequest struct {
	OrgId   int              `json:"orgId"`
	SliType string           `json:"sliType"`
	RawData []SliRawDataBody `json:"rawData"`
}

type PostManyResponse struct {
	SliRawData *[]SliRawDataBody `json:"sliRawData"`
}

type SliRawData interface {
	PostMany(c *clients.BlamelessClient, data *[]SliRawDataBody) (*PostManyResponse, error)
}

func PostMany(c *clients.BlamelessClient, data []SliRawDataBody) (*PostManyResponse, error) {
	payload := &PostManyRequest{
		OrgId:   config.Environment().Blameless.OrgId,
		SliType: "latency",
		RawData: data,
	}
	postBody, err := json.Marshal(payload)
	if err != nil {
		return &PostManyResponse{}, err
	}
	resp, err := c.Post(c.SloTimeseriesService, "SliRawDataPostMany", postBody)
	if err != nil {
		return &PostManyResponse{}, err
	}
	var resultBody *PostManyResponse
	if err := json.Unmarshal(resp, &resultBody); err != nil {
		return &PostManyResponse{}, err
	}

	return resultBody, nil
}
