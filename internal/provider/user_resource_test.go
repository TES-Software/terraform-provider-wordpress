package provider

import (
  "os"
  "testing"

  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func testAccPreCheck(t *testing.T) {
  // Check if WordPress environment variables are set
  if v := os.Getenv("WP_ENDPOINT"); v == "" {
    t.Skip("WP_ENDPOINT must be set for acceptance tests")
  }
  if v := os.Getenv("WP_USER"); v == "" {
    t.Skip("WP_USER must be set for acceptance tests")
  }
  if v := os.Getenv("WP_PASSWORD"); v == "" {
    t.Skip("WP_PASSWORD must be set for acceptance tests")
  }
}

func TestAccResourceUser(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck: func() { testAccPreCheck(t) },
    ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error) {
      "wordpress": providerserver.NewProtocol6WithError(New("test")()),
    },
    Steps: []resource.TestStep{
      {
        Config: testResource,
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("wordpress_user.test", "username", "test1"),
          resource.TestCheckResourceAttr("wordpress_user.test", "email", "test1@example.com"),
        ),
      },
    },
  })
}

const testResource = `
resource "wordpress_user" "test" {
  email = "test1@example.com"
  username = "test1"
  password = "password"
}
`
