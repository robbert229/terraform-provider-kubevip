package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const kubevipVersion = "v0.6.3"

// kubevip is the provider implementation.
type kubevipProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kubevipProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *kubevipProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kubevip"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *kubevipProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *kubevipProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *kubevipProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewManifestDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *kubevipProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
