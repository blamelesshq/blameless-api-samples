package main

import (
	"github.com/blamelesshq/blameless-examples/slo/packages/config"
)

func main() {
	// Initialize config build-out
	config.NewConfig()

	// app := &cli.App{
	// 	Name: "create_sli",
	// 	HelpName: "h",
	// 	Usage: "Programatically create SLI",
	// 	Action: func(c *cli.Context) error {

	// 	}
	// }
	// First get auth token from Blameless' oauth service (Token information should be provided by AE/SE)

	// Next get the Blameless SLI ID or manually hard-code the id into the

	// **IF** you want to do programattic SLI creation enable the flag within the config.yaml
	// This will programmatically build an SLI through command prompt variables

	// Check the SLI to see if you've already completed the backfill
	// Backfilling is used to provided historical data so Blameless can accurrately provide
	// provide you a 28 day rolling window. Once exapnded to calendar year, this will also
	// be necessary. If backfill is not complete perform and update the SLI record to
	// set it completed.

	// We will schedule a regular ingest period to fetch data from our prometheus endpoint.
	// Its assumed since we don't store the metric query that you'd somehow store this information
	// with a foreign ID mapping back to the SLI ID provided from Blameless. This can be
	// stored in some simple file.
}
