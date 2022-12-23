package wordpress

import (
  "context"
  "fmt"
  "strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
  "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
  "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
  wcl "github.com/sogko/go-wordpress"
)

var (
  _ resource.Resource = &UserResource{}
  _ resource.ResourceWithConfigure = &UserResource{}
  _ resource.ResourceWithImportState = &UserResource{}
)

type UserResource struct {
  client *wcl.Client
}

type UserResourceModel struct {
  ID types.String `tfsdk:"id"`
  Name types.String `tfsdk:"name"`
  Nickname types.String `tfsdk:"nickname"`
  Username types.String `tfsdk:"username"`
  FirstName types.String `tfsdk:"first_name"`
  LastName types.String `tfsdk:"last_name"`
  Description types.String `tfsdk:"description"`
  Email types.String `tfsdk:"email"`
  Password types.String `tfsdk:"password"`
  Roles types.List `tfsdk:"roles"`
}

func NewUserResource() resource.Resource {
  return &UserResource{}
}

func (u *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
  if req.ProviderData == nil {
    return
  }

  client, ok := req.ProviderData.(*wcl.Client)

  if !ok {
    resp.Diagnostics.AddError(
      "Unexpected Resource Configure Type",
      fmt.Sprintf("Expected Client, got: %T. Please report this issue to the developer", req.ProviderData),
    )
    return
  }

  u.client = client
}

func (u *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
  resp.Schema = schema.Schema{
    Attributes: map[string]schema.Attribute{
      "id": schema.StringAttribute{
        MarkdownDescription: "ID of the user",
        Computed: true,
        PlanModifiers: []planmodifier.String{
          stringplanmodifier.UseStateForUnknown(),
        },
      },
      "name": schema.StringAttribute{
        MarkdownDescription: "Name of the user",
        Optional: true,
        Computed: true,
      },
      "nickname": schema.StringAttribute{
        MarkdownDescription: "Nickname of the user",
        Optional: true,
        Computed: true,
      },
      "username": schema.StringAttribute{
        MarkdownDescription: "Username of the user",
        Required:            true,
        PlanModifiers: []planmodifier.String{
          stringplanmodifier.RequiresReplace(),
        },
      },
      "email": schema.StringAttribute{
        MarkdownDescription: "Email of the user",
        Required:            true,
      },
      "password": schema.StringAttribute{
        MarkdownDescription: "Password of the user",
        Required:            true,
        Sensitive: true,
      },
      "roles": schema.ListAttribute{
        ElementType: types.StringType,
        Optional: true,
        Computed: true,
      },
      "first_name": schema.StringAttribute{
        Optional: true,
      },
      "last_name": schema.StringAttribute{
        Optional: true,
      },
      "description": schema.StringAttribute{
        Optional: true,
      },
    },
    MarkdownDescription: "Manages a user",
  }
}

func (u *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
  resp.TypeName = "wordpress_user"
}

func (u *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
  var data UserResourceModel

  resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

  if resp.Diagnostics.HasError() {
    return
  }

  var roles []string = []string{}

  for _, r := range data.Roles.Elements() {
    roles = append(roles, r.(types.String).ValueString())
  }

  requ := &wcl.User{
    Email: data.Email.ValueString(),
    Name: data.Name.ValueString(),
    Username: data.Username.ValueString(),
    Password: data.Password.ValueString(),
    Roles: roles,
  }

  if !data.Nickname.IsNull() {
    requ.Nickname = data.Nickname.ValueString()
  }

  if !data.Name.IsNull() {
    requ.Name = data.Name.ValueString()
  }

  wu, _, _, err := u.client.Users().Create(requ)

  if err != nil {
    resp.Diagnostics.AddError(
      "Error creating user",
      "Could not create uesr, unexpected error: "+err.Error(),
    )
    return
  }

  setUserModel(&data, wu)
  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }
}


func (u *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
  var data UserResourceModel

  resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

  iid, _ := strconv.Atoi(data.ID.ValueString())
  u.client.Users().Delete(int(iid), nil)
}

func (u *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
  var data UserResourceModel
  resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

  iid, _ := strconv.Atoi(data.ID.ValueString())
  wu, _, _, err := u.client.Users().Get(int(iid), "context=edit")

  if err != nil {
    resp.Diagnostics.AddError(
      "Error Reading User",
      "Could not read User ID "+data.ID.ValueString()+": "+err.Error(),
    )
    return
  }

  setUserModel(&data, wu)
  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }
}

func (u *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
  var data UserResourceModel

  resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

  if resp.Diagnostics.HasError() {
    return
  }

  iid, _ := strconv.Atoi(data.ID.ValueString())

  wu, _, _, err := u.client.Users().Update(int(iid), &wcl.User{
    Email: data.Email.ValueString(),
    Name: data.Name.ValueString(),
    Password: data.Password.ValueString(),
  })

  if err != nil {
    resp.Diagnostics.AddError(
      "Error updating user",
      "Could not update uesr, unexpected error: "+err.Error(),
    )
    return
  }
  setUserModel(&data, wu)
  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }
}

func (u *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
  resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func setUserModel(data *UserResourceModel, wu *wcl.User) {
  sid := strconv.Itoa(wu.ID)

  if sid != "" {
    data.ID = types.StringValue(sid)
  }

  if wu.Email != "" {
    data.Email = types.StringValue(wu.Email)
  }

  if wu.Name != "" {
    data.Name = types.StringValue(wu.Name)
  }

  if wu.Nickname != "" {
    data.Nickname = types.StringValue(wu.Nickname)
  }

  if wu.FirstName != "" {
    data.FirstName = types.StringValue(wu.FirstName)
  }

  if wu.LastName != "" {
    data.LastName = types.StringValue(wu.LastName)
  }

  if wu.Description != "" {
    data.Description = types.StringValue(wu.Description)
  }

  if wu.Username != "" {
    data.Username = types.StringValue(wu.Username)
  }

  if wu.Roles != nil {
    var roles []attr.Value
    for _, r := range wu.Roles {
      roles = append(roles, types.StringValue(r))
    }
    data.Roles, _ = types.ListValue(types.StringType, roles)
  }
}
