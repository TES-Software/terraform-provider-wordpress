package wordpress

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
  wcl "github.com/sogko/go-wordpress"
)

var (
  _ provider.Provider = &wordpressProvider{}
)

type wordpressProvider struct {
  version string
}

type wordpressProviderModel struct {
  Endpoint types.String `tfsdk:"endpoint"`
  User types.String `tfsdk:"user"`
  Password types.String `tfsdk:"password"`
}

func New() provider.Provider {
  return &wordpressProvider{
    version: "0.0.1",
  }
}

func (p *wordpressProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp provider.MetadataResponse) {
  resp.TypeName = "wordpress"
  resp.Version = p.version
}

func (p *wordpressProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
  var data wordpressProviderModel

  resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

  endpoint := data.Endpoint.ValueString()
  user := data.User.ValueString()
  password := data.Password.ValueString()

  client := wcl.NewClient(&wcl.Options{
    BaseAPIURL: endpoint,
    Username: user,
    Password: password,
  })
  resp.ResourceData = client
}

func (p *wordpressProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
  resp.Schema = schema.Schema{
    Attributes: map[string]schema.Attribute{
      "password": schema.StringAttribute{
        Required: true,
      },
      "user": schema.StringAttribute{
        Required: true,
      },
      "endpoint": schema.StringAttribute{
        Required: true,
      },
    },
  }
}

func (p *wordpressProvider) DataSources(_ context.Context) []func() datasource.DataSource {
  return []func() datasource.DataSource{}
}

func (p *wordpressProvider) Resources(_ context.Context) []func() resource.Resource {
  return []func() resource.Resource{
    NewUserResource,
  }
}
