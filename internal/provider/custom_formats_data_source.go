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

const customFormatsDataSourceName = "custom_formats"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatsDataSource{}

func NewCustomFormatsDataSource() datasource.DataSource {
	return &CustomFormatsDataSource{}
}

// CustomFormatsDataSource defines the download clients implementation.
type CustomFormatsDataSource struct {
	client *radarr.Radarr
}

// CustomFormats describes the download clients data model.
type CustomFormats struct {
	CustomFormats types.Set    `tfsdk:"custom_formats"`
	ID            types.String `tfsdk:"id"`
}

func (d *CustomFormatsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatsDataSourceName
}

func (d *CustomFormatsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Custom Formats](../resources/custom_format).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"custom_formats": {
				MarkdownDescription: "Download Client list..",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"include_custom_format_when_renaming": {
						MarkdownDescription: "Include custom format when renaming flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"name": {
						MarkdownDescription: "Custom Format name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"id": {
						MarkdownDescription: "Custom Format ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"specifications": {
						MarkdownDescription: "Specifications.",
						Computed:            true,
						Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
							"negate": {
								MarkdownDescription: "Negate flag.",
								Computed:            true,
								Type:                types.BoolType,
							},
							"required": {
								MarkdownDescription: "Computed flag.",
								Computed:            true,
								Type:                types.BoolType,
							},
							"name": {
								MarkdownDescription: "Specification name.",
								Computed:            true,
								Type:                types.StringType,
							},
							"implementation": {
								MarkdownDescription: "Implementation.",
								Computed:            true,
								Type:                types.StringType,
							},
							// Field values
							"value": {
								MarkdownDescription: "Value.",
								Computed:            true,
								Type:                types.StringType,
							},
							"min": {
								MarkdownDescription: "Min.",
								Computed:            true,
								Type:                types.Int64Type,
							},
							"max": {
								MarkdownDescription: "Max.",
								Computed:            true,
								Type:                types.Int64Type,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (d *CustomFormatsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CustomFormatsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormats

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get download clients current value
	response, err := d.client.GetCustomFormatsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", customFormatsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatsDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]CustomFormat, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.CustomFormats.Type(context.Background()), &data.CustomFormats)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}