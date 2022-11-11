package main

import (
	"context"
	"flag"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {

	var debug bool
	flag.BoolVar(&debug, "debug", false, "set this to true if you want to debug the code using delve")
	flag.Parse()

	providerserver.Serve(context.Background(), NewBuildOnAWSProvider, providerserver.ServeOpts{
		Debug:   debug,
		Address: "aws.amazon.com/terraform/buildonaws",
	})

}
