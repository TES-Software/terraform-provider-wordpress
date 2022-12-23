package main

import (
	"context"
  "log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/yyamanoi1222/terraform-provider-wordpress/wordpress"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
  err :=providerserver.Serve(context.Background(), wordpress.New, providerserver.ServeOpts{
    Address: "registry.terrafomr.io/yyamanoi1222/wordpress",
  })

  if err != nil {
    log.Fatal(err.Error())
  }
}
