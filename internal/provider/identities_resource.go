package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ory/client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type x map[string]interface{}
type identityResourceModel struct {
	Id       types.String `tfsdk:"id"`
	SchemaId types.String `tfsdk:"schema_id"`
	Traits   types.Map    `tfsdk:"traits"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &identityResource{}
	_ resource.ResourceWithConfigure = &identityResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewIdentityResource() resource.Resource {
	return &identityResource{}
}

// identityResource is the resource implementation.
type identityResource struct {
	oryClient *client.IdentityApi
	context   *context.Context
}

// Metadata returns the resource type name.
func (r *identityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity"
}

// Schema defines the schema for the resource.
func (r *identityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
			},
			"schema_id": schema.StringAttribute{
				Required: true,
			},
			"traits": schema.MapAttribute{
				Required: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *identityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var identityResource identityResourceModel
	diags := req.Plan.Get(ctx, &identityResource)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	goObject := make(map[string]interface{})
	// read from the terraform data into the map
	if diags := identityResource.Traits.ElementsAs(ctx, &goObject, false); diags != nil {
		// error
	}
	body := client.CreateIdentityBody{
		SchemaId: identityResource.SchemaId.ValueString(),
		Traits:   goObject,
	}
	c := *r.oryClient
	auth := *r.context

	response, _, err := c.CreateIdentity(auth).CreateIdentityBody(body).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Identity",
			err.Error(),
		)
		return
	}

	identityResource.Id = types.StringValue(response.Id)
	diags = resp.State.Set(ctx, identityResource)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *identityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	c := *r.oryClient
	auth := *r.context

	goObject := make(map[string]interface{})
	// read from the terraform data into the map
	if diags := state.Traits.ElementsAs(ctx, &goObject, false); diags != nil {
		// error
	}
	identity, _, err := c.GetIdentity(auth, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Can't convert Traits", "Can't convert Traits")
	}
	state.Id = types.StringValue(identity.Id)
	state.SchemaId = types.StringValue(identity.SchemaId)
	state.Traits = types.Map{}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Identity",
			err.Error(),
		)
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *identityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *identityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *identityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.oryClient = &client.ApiClient.IdentityApi
	r.context = &client.Context
}
