package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/blamelesshq/blameless-examples/slo/packages/models"
	"github.com/blamelesshq/blameless-examples/slo/packages/utils"
	"github.com/cheynewallace/tabby"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func sli() *cobra.Command {
	sli := &cobra.Command{
		Use:   "sli",
		Short: "SLI domain primary command",
		Long:  `SLI commands begin here. `,
	}

	sli.AddCommand(sliCreate())
	sli.AddCommand(sliGet())
	sli.AddCommand(ingest())

	return sli
}

type Types struct {
	Name string
	Id   int
}

func intPrompt(label string) int {
	validateInt := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return fmt.Errorf("unable to parse integer for %s", label)
		}
		return nil
	}

	intP := promptui.Prompt{
		Label:    label,
		Validate: validateInt,
	}
	result, err := intP.Run()
	if err != nil {
		log.Fatal(err)
	}

	id, err := strconv.Atoi(result)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func stringPrompt(label string) string {
	validateString := func(input string) error {
		if len(input) < 1 {
			return fmt.Errorf("must provide at least one character for %s", label)
		}
		return nil
	}
	stringP := promptui.Prompt{
		Label:    label,
		Validate: validateString,
	}
	result, err := stringP.Run()
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func sliCreate() *cobra.Command {
	sliTypes := []Types{
		{Name: "Availability", Id: 1},
		{Name: "Latency", Id: 2},
		{Name: "Throughput", Id: 3},
		{Name: "Saturation", Id: 4},
		{Name: "Durability", Id: 5},
		{Name: "Correctness", Id: 6},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .Id | white }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Id | white }})",
		Selected: "\U0001F336 {{ .Name | green | cyan }}",
		Details: `
		--------- SLI Types ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Id:" | faint }}	{{ .Id }}
		`,
	}

	searcher := func(input string, index int) bool {
		t := sliTypes[index]
		name := strings.Replace(strings.ToLower(t.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "Create a new SLI",
		Long:  `Create a new SLI`,
		Run: func(cmd *cobra.Command, args []string) {
			orgId := utils.IntPrompt("Org ID")
			name := utils.StringPrompt("Name")
			description := stringPrompt("Description")
			sliTypePrompt := promptui.Select{
				Label:     "Sli Type",
				Items:     sliTypes,
				Templates: templates,
				Size:      4,
				Searcher:  searcher,
			}
			sliType, _, err := sliTypePrompt.Run()
			if err != nil {
				log.Fatalf("Unable to parse selection: \n%+v", err)
			}
			serviceId := intPrompt("Service ID")

			sliBody := &models.SliBody{
				Name:         name,
				Description:  description,
				DataSourceId: 5,
				SliTypeId:    sliType,
				ServiceId:    serviceId,
			}

			metricPath := &models.MetricPath{}
			switch t := sliType + 1; t {
			case sliTypes[0].Id:
				goodRequest := utils.StringPrompt("Good Request Query")
				validRequest := utils.StringPrompt("Valid Request Query")
				availability := &models.AvailabilityStruct{
					GoodRequest:  goodRequest,
					ValidRequest: validRequest,
				}
				metricPath = &models.MetricPath{
					Availability: availability,
				}
			case sliTypes[1].Id:
				latencyReq := utils.StringPrompt("Latency Query")
				metricPath = &models.MetricPath{
					Latency: latencyReq,
				}
			case sliTypes[2].Id:
				throughputReq := utils.StringPrompt("Throughput Query")
				metricPath = &models.MetricPath{
					Throughput: throughputReq,
				}
			case sliTypes[3].Id:
				saturationReq := utils.StringPrompt("Saturation Query")
				metricPath = &models.MetricPath{
					Saturation: saturationReq,
				}
			case sliTypes[4].Id:
				durabilityReq := stringPrompt("Durability Query")
				metricPath = &models.MetricPath{
					Durability: durabilityReq,
				}
			case sliTypes[5].Id:
				correctnessReq := stringPrompt("Correctness Query")
				metricPath = &models.MetricPath{
					Correctness: correctnessReq,
				}
			}
			mp, err := json.Marshal(metricPath)
			if err != nil {
				log.Fatalf("error while marshaling metric path: %s", err)
			}
			sliBody.MetricPath = string(mp)
			postBody := &models.PostSliRequest{
				OrgId: orgId,
				Model: sliBody,
			}
			resp, err := models.PostSli(postBody)
			if err != nil {
				log.Fatalf("Unable to make regequest: \n%+v", err)
			}
			t := tabby.New()
			t.AddHeader("Org ID", "ID", "Name", "Description", "Data Source ID", "SLI Type ID", "Service ID", "User ID")
			t.AddLine(resp.Sli.Id,
				resp.Sli.Name,
				resp.Sli.Description,
				resp.Sli.DataSourceId,
				resp.Sli.SliTypeId,
				resp.Sli.ServiceId,
			)
			t.Print()
		},
	}

	return create
}

func sliGet() *cobra.Command {
	get := &cobra.Command{
		Use:   "get",
		Short: "get a SLI",
		Long:  `get a SLI`,
		Run: func(cmd *cobra.Command, args []string) {
			sliReq := &models.GetSliRequest{
				OrgId: intPrompt("Org ID"),
				Id:    intPrompt("SLI ID"),
			}
			resp, err := models.GetSli(sliReq)
			if err != nil {
				log.Fatalf("unable to complete request: \n%+v", err)
			}
			st, err := resp.Sli.GetSliType()
			if err != nil {
				log.Fatalf("unable to get SLI type: \n%+v", err)
			}
			t := tabby.New()
			t.AddHeader("Org ID", "ID", "Name", "Description", "Data Source ID", "SLI Type ID", "SLI Type", "Service ID", "User ID")
			t.AddLine(resp.Sli.OrgId,
				resp.Sli.Id,
				resp.Sli.Name,
				resp.Sli.Description,
				resp.Sli.DataSourceId,
				resp.Sli.SliTypeId,
				st.SliType.Name,
				resp.Sli.ServiceId,
			)
			t.Print()
		},
	}

	return get
}
