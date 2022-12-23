package wordpress

import (
  "testing"

  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestAccResourceUser(t *testing.T) {
  resource.Test(t, resource.TestCase{
    ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error) {
      "wordpress": providerserver.NewProtocol6WithError(New()),
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
