package cmd

import (
	"log"

	"github.com/blamelesshq/blameless-examples/slo/packages/models"
	"github.com/blamelesshq/blameless-examples/slo/packages/utils"
	"github.com/spf13/cobra"
)

func ingest() *cobra.Command {
	ingest := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest domain primary command",
		Long:  `SLI ingest commands begin here. `,
	}

	return ingest
}

func ingestSli() *cobra.Command {
	ingest := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest scheduling for an SLI",
		Long:  `Ingest an SLIs backfill or regular ingest period`,
		Run: func(cmd *cobra.Command, args []string) {
			orgId := utils.IntPrompt("Org ID")
			sliId := utils.IntPrompt("SLI ID")
			backfill := utils.BooleanPrompt("Backfill ?")

			resp, err := models.GetSli(&models.GetSliRequest{
				OrgId: orgId,
				Id:    sliId,
			})

			if err != nil {
				log.Fatalf("unable to fetch SLI: %+v", err)
			}

			sliType, err := resp.Sli.GetSliType()
			if err != nil {
				log.Fatalf("unable to get SLI type: %+v", err)
			}
		},
	}
	return ingest
}
