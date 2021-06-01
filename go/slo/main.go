package main

import (
	"github.com/blamelesshq/blameless-examples/slo/packages/cmd"
	"github.com/blamelesshq/blameless-examples/slo/packages/config"
)

func main() {
	config.Environment()
	cmd.Execute()
}
