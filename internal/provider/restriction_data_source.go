package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

const restrictionDataSourceName = "restriction"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RestrictionDataSource{}

func NewRestrictionDataSource() datasource.DataSource {
	return &RestrictionDataSource{}
}

// RestrictionDataSource defines the remote path restriction implementation.
type RestrictionDataSource struct {
	client *radarr.Radarr
}

func (d *RestrictionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + restrictionDataSourceName
}

func (d *RestrictionDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->Single [Restriction](../resources/restriction).",
		Attributes: map[string]tfsdk.Attribute{
			"required": {
				MarkdownDescription: "Required.",
				Computed:            true,
				Type:                types.StringType,
			},
			"ignored": {
				MarkdownDescription: "Ignored.",
				Computed:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Restriction ID.",
				Required:            true,
				Type:                types.Int64Type,
			},
		},
	}, nil
}

func (d *RestrictionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.Radarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *radarr.Radarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *RestrictionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var restriction *Restriction

	resp.Diagnostics.Append(req.Config.Get(ctx, &restriction)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get remote path restriction current value
	response, err := d.client.GetRestrictionsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", restrictionDataSourceName, err))

		return
	}

	// Map response body to resource schema attribute
	value, err := findRestriction(restriction.ID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(tools.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", restrictionDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+restrictionDataSourceName)

	restriction.write(ctx, value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &restriction)...)
}

func findRestriction(id int64, restrictions []*radarr.Restriction) (*radarr.Restriction, error) {
	for _, m := range restrictions {
		if m.ID == id {
			return m, nil
		}
	}

	return nil, tools.ErrDataNotFoundError(restrictionDataSourceName, "id", strconv.Itoa(int(id)))
}