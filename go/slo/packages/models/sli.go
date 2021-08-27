package models

import (
	"encoding/json"

	"github.com/blamelesshq/blameless-examples/slo/packages/clients"
)

var BlamelessSourceId = 5

var Types = &SliTypes{
	Latency:      "Latency",
	Availability: "Availability",
	Saturation:   "Saturation",
	Throughput:   "Throughput",
	Correctness:  "Correctness",
	Durability:   "Durability",
}

type SliTypes struct {
	Latency      string
	Availability string
	Saturation   string
	Throughput   string
	Correctness  string
	Durability   string
}

type AvailabilityStruct struct {
	GoodRequest  string `json:"goodRequest,omitempty"`
	ValidRequest string `json:"validRequest,omitempty"`
}

type MetricPath struct {
	Latency      string              `json:"latency,omitempty"`
	Availability *AvailabilityStruct `json:"availability,omitempty"`
	Throughput   string              `json:"throughput,omitempty"`
	Saturation   string              `json:"saturation,omitempty"`
	Correctness  string              `json:"correctness,omitempty"`
	Durability   string              `json:"durability,omitempty"`
}

type SliBody struct {
	OrgId        int    `json:"orgId,omitempty"`
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DataSourceId int    `json:"dataSourceId"` // Datasource ID 5 === Blameless API
	SliTypeId    int    `json:"sliTypeId"`
	ServiceId    int    `json:"serviceId"`
	UserId       int    `json:"userId,omitempty"`
	Checkpoint   int    `json:"checkpoint,omitempty"`
	MetricPath   string `json:"metricPath,omitempty"`
}

type SliTypeRequest struct {
	Id int `json:"id" binding:"required"`
}

type SliTypeBody struct {
	Id   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type SliTypeResponse struct {
	SliType *SliTypeBody `json:"sliType"`
}

type PostSliRequest struct {
	OrgId int      `json:"orgId" binding:"required"`
	Model *SliBody `json:"model" binding:"required"`
}

type GetSliRequest struct {
	OrgId int `json:"orgId" binding:"required"`
	Id    int `json:"id" binding:"required"`
}

type SliResponse struct {
	Sli *SliBody `json:"sli"`
}

type Sli interface {
	GetSli(req *GetSliRequest) (*SliResponse, error)
	PostSli(req *PostSliRequest) (*SliResponse, error)
	Update()
}

func GetSli(req *GetSliRequest) (*SliResponse, error) {
	c := clients.NewBlamelessClient()
	payload, err := json.Marshal(&req)
	if err != nil {
		return &SliResponse{}, err
	}
	resp, err := c.Post(c.SloService, "GetSLI", payload)
	if err != nil {
		return &SliResponse{}, err
	}
	var resultBody *SliResponse
	if err := json.Unmarshal(resp, &resultBody); err != nil {
		return &SliResponse{}, err
	}
	return resultBody, nil
}

func PostSli(req *PostSliRequest) (*SliResponse, error) {
	c := clients.NewBlamelessClient()
	payload, err := json.Marshal(&req)
	if err != nil {
		return &SliResponse{}, err
	}
	resp, err := c.Post(c.SloService, "CreateSLI", payload)
	if err != nil {
		return &SliResponse{}, err
	}
	var resultBody *SliResponse
	if err := json.Unmarshal(resp, &resultBody); err != nil {
		return &SliResponse{}, err
	}
	return resultBody, nil
}

func (s *SliBody) GetSliType() (*SliTypeResponse, error) {
	c := clients.NewBlamelessClient()
	request := &SliTypeRequest{
		Id: s.SliTypeId,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return &SliTypeResponse{}, err
	}
	resp, err := c.Post(c.SloService, "GetSliType", payload)
	if err != nil {
		return &SliTypeResponse{}, err
	}
	var bodyResponse *SliTypeResponse
	if err := json.Unmarshal(resp, &bodyResponse); err != nil {
		return &SliTypeResponse{}, err
	}
	return bodyResponse, nil
}
