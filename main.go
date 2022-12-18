package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/yyamanoi1222/terraform-provider-wordpress/wordpress"
)

func main() {
  providerserver.Serve(context.Background(), wordpress.New, providerserver.ServeOpts{
    Address: "wordpress",
  })
}
