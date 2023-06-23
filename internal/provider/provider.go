// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ory/client-go"
)

// Ensure oryNetworkProvider satisfies various provider interfaces.
var _ provider.Provider = &oryNetworkProvider{}

// oryNetworkProvider defines the provider implementation.
type oryNetworkProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// oryNetworkProviderModel describes the provider data model.
type oryNetworkProviderModel struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"api_key"`
}

type oryProviderResponse struct {
	ApiClient *client.APIClient
	Context   context.Context
}

func (p *oryNetworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "orynetwork"
	resp.Version = p.version
}

func (p *oryNetworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Host of ORY Network",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Api Key to authenticate to Ory",
				Sensitive:           true,
				Required:            true,
			},
		},
	}
}

func (p *oryNetworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data oryNetworkProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: validate incoming model

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example oryClient configuration for data sources and resources

	configuration := client.NewConfiguration()
	configuration.Host = data.Host.ValueString()
	apiClient := client.NewAPIClient(configuration)
	auth := context.WithValue(context.Background(), client.ContextAccessToken, data.ApiKey.ValueString())

	response := oryProviderResponse{
		Context:   auth,
		ApiClient: apiClient,
	}
	resp.DataSourceData = response
	resp.ResourceData = response
}

func (p *oryNetworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *oryNetworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewIdentityDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &oryNetworkProvider{
			version: version,
		}
	}
}
