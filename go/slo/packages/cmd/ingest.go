package cmd

import (
	"github.com/spf13/cobra"
)

func ingest() *cobra.Command {
	ingest := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest domain primary command",
		Long:  `SLI ingest commands begin here. `,
	}

	// sli.AddCommand(sliCreate())
	// sli.AddCommand(sliGet())

	return ingest
}
