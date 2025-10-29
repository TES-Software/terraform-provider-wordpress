package provider

import (
	"context"
  "os"

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

func New(version string) func() provider.Provider {
  return func() provider.Provider {
    return &wordpressProvider{
      version: version,
    }
  }
}

func (p *wordpressProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp provider.MetadataResponse) {
  resp.TypeName = "wordpress"
  resp.Version = p.version
}

func (p *wordpressProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
  var data wordpressProviderModel

  user := os.Getenv("WP_USER")
  endpoint := os.Getenv("WP_ENDPOINT")
  password := os.Getenv("WP_PASSWORD")

  resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

  if data.Endpoint.ValueString() != "" {
    endpoint = data.Endpoint.ValueString()
  }

  if data.User.ValueString() != "" {
    user = data.User.ValueString()
  }

  if data.Password.ValueString() != "" {
    password = data.Password.ValueString()
  }

  if endpoint == "" {
    resp.Diagnostics.AddError(
      "Missing Endpoint Configuration",
      "While configuring the provider, the Endpoint was not found in "+
      "the WP_ENDPOINT environment variable or provider "+
      "configuration block endpoint attribute.",
    )
  }

  if user == "" {
    resp.Diagnostics.AddError(
      "Missing User Configuration",
      "While configuring the provider, the User was not found in "+
      "the WP_USER environment variable or provider "+
      "configuration block user attribute.",
    )
  }

  if password == "" {
    resp.Diagnostics.AddError(
      "Missing Password Configuration",
      "While configuring the provider, the Password was not found in "+
      "the WP_PASSWORD environment variable or provider "+
      "configuration block password attribute.",
    )
  }
  clientOptions := &wcl.Options{
    BaseAPIURL: endpoint,
    Username: user,
    Password: password,
  }
  resp.ResourceData = clientOptions
}

func (p *wordpressProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
  resp.Schema = schema.Schema{
    Attributes: map[string]schema.Attribute{
      "password": schema.StringAttribute{
        Sensitive: true,
        Optional: true,
      },
      "user": schema.StringAttribute{
        Optional: true,
      },
      "endpoint": schema.StringAttribute{
        Optional: true,
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
