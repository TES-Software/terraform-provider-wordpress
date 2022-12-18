package wordpress

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
  _ provider.Provider = &wordpressProvider{}
)

type wordpressProvider struct {}

func New() provider.Provider {
  return &wordpressProvider{}
}

func (p *wordpressProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *wordpressProvider) DataSources(_ context.Context) []func() datasource.DataSource {
  return []func() datasource.DataSource{}
}

func (p *wordpressProvider) Resources(_ context.Context) []func() resource.Resource {
  return []func() resource.Resource{}
}
