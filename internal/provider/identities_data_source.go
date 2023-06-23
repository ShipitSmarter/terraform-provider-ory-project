// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ory/client-go"
)

var (
	_ datasource.DataSource              = &identityDataSource{}
	_ datasource.DataSourceWithConfigure = &identityDataSource{}
)

type identitiesDataSourceModel struct {
	Identities []identitiesModel `tfsdk:"identities"`
}

type identitiesModel struct {
	ID types.String `tfsdk:"id"`
}

func (d *identityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(oryProviderResponse)
	if !ok {

		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.ApiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.oryClient = &client.ApiClient.IdentityApi
	d.context = &client.Context
}

// Ensure provider defined types fully satisfy framework interfaces.
func NewIdentityDataSource() datasource.DataSource {
	return &identityDataSource{}
}

type identityDataSource struct {
	oryClient *client.IdentityApi
	context   *context.Context
}

func (d *identityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identities"
}

func (d *identityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"identities": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *identityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state identitiesDataSourceModel

	c := *d.oryClient
	auth := *d.context

	identities, _, err := c.ListIdentities(auth).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Identities",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, identity := range identities {
		identityState := identitiesModel{
			ID: types.StringValue(identity.Id),
		}

		state.Identities = append(state.Identities, identityState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
