package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

const indexerResourceName = "indexer"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexerResource{}
var _ resource.ResourceWithImportState = &IndexerResource{}

var (
	indexerIntSliceFields = []string{"categories", "multiLanguages", "requiredFlags", "codecs", "mediums"}
	indexerBoolFields     = []string{"allowZeroSize", "removeYear", "rankedOnly"}
	indexerIntFields      = []string{"delay", "minimumSeeders", "seedTime"}
	indexerStringFields   = []string{"additionalParameters", "apiKey", "apiPath", "baseUrl", "captchaToken", "cookie", "passkey", "username", "aPIUser", "aPIKey"}
	indexerFloatFields    = []string{"seedRatio"}
)

func NewIndexerResource() resource.Resource {
	return &IndexerResource{}
}

// IndexerResource defines the indexer implementation.
type IndexerResource struct {
	client *radarr.Radarr
}

// Indexer describes the indexer data model.
type Indexer struct {
	Categories              types.Set     `tfsdk:"categories"`
	Mediums                 types.Set     `tfsdk:"mediums"`
	Codecs                  types.Set     `tfsdk:"codecs"`
	RequiredFlags           types.Set     `tfsdk:"required_flags"`
	Tags                    types.Set     `tfsdk:"tags"`
	MultiLanguages          types.Set     `tfsdk:"multi_languages"`
	Cookie                  types.String  `tfsdk:"cookie"`
	APIKey                  types.String  `tfsdk:"api_key"`
	ConfigContract          types.String  `tfsdk:"config_contract"`
	Implementation          types.String  `tfsdk:"implementation"`
	Protocol                types.String  `tfsdk:"protocol"`
	Username                types.String  `tfsdk:"username"`
	Passkey                 types.String  `tfsdk:"passkey"`
	BaseURL                 types.String  `tfsdk:"base_url"`
	CaptchaToken            types.String  `tfsdk:"captcha_token"`
	AdditionalParameters    types.String  `tfsdk:"additional_parameters"`
	APIPath                 types.String  `tfsdk:"api_path"`
	APIUser                 types.String  `tfsdk:"api_user"`
	Name                    types.String  `tfsdk:"name"`
	Priority                types.Int64   `tfsdk:"priority"`
	SeedTime                types.Int64   `tfsdk:"seed_time"`
	MinimumSeeders          types.Int64   `tfsdk:"minimum_seeders"`
	DownloadClientID        types.Int64   `tfsdk:"download_client_id"`
	Delay                   types.Int64   `tfsdk:"delay"`
	ID                      types.Int64   `tfsdk:"id"`
	SeedRatio               types.Float64 `tfsdk:"seed_ratio"`
	AllowZeroSize           types.Bool    `tfsdk:"allow_zero_size"`
	RankedOnly              types.Bool    `tfsdk:"ranked_only"`
	EnableRss               types.Bool    `tfsdk:"enable_rss"`
	EnableAutomaticSearch   types.Bool    `tfsdk:"enable_automatic_search"`
	EnableInteractiveSearch types.Bool    `tfsdk:"enable_interactive_search"`
	RemoveYear              types.Bool    `tfsdk:"remove_year"`
}

func (r *IndexerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerResourceName
}

func (r *IndexerResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/radarr/settings#indexers) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"enable_automatic_search": {
				MarkdownDescription: "Enable automatic search flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_interactive_search": {
				MarkdownDescription: "Enable interactive search flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_rss": {
				MarkdownDescription: "Enable RSS flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"priority": {
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"download_client_id": {
				MarkdownDescription: "Download client ID.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"config_contract": {
				MarkdownDescription: "Indexer configuration template.",
				Required:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "Indexer implementation name.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Indexer name.",
				Required:            true,
				Type:                types.StringType,
			},
			"protocol": {
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					tools.StringMatch([]string{"usenet", "torrent"}),
				},
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Indexer ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"allow_zero_size": {
				MarkdownDescription: "Allow zero size files.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"remove_year": {
				MarkdownDescription: "Remove year.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"ranked_only": {
				MarkdownDescription: "Allow ranked only.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"delay": {
				MarkdownDescription: "Delay before grabbing.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"minimum_seeders": {
				MarkdownDescription: "Minimum seeders.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"seed_time": {
				MarkdownDescription: "Seed time.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"seed_ratio": {
				MarkdownDescription: "Seed ratio.",
				Optional:            true,
				Computed:            true,
				Type:                types.Float64Type,
			},
			"additional_parameters": {
				MarkdownDescription: "Additional parameters.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_user": {
				MarkdownDescription: "API User.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_path": {
				MarkdownDescription: "API path.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"base_url": {
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"captcha_token": {
				MarkdownDescription: "Captcha token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"cookie": {
				MarkdownDescription: "Cookie.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"passkey": {
				MarkdownDescription: "Passkey.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"categories": {
				MarkdownDescription: "Series list.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"multi_languages": {
				MarkdownDescription: "Language list.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"required_flags": {
				MarkdownDescription: "Required flags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"codecs": {
				MarkdownDescription: "Codecs.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"mediums": {
				MarkdownDescription: "Mediumd.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
		},
	}, nil
}

func (r *IndexerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.Radarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *radarr.Radarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IndexerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Indexer
	request := indexer.read(ctx)

	response, err := r.client.AddIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *Indexer

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Indexer current value
	response, err := r.client.GetIndexerContext(ctx, indexer.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Indexer
	request := indexer.read(ctx)

	response, err := r.client.UpdateIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer Indexer

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Indexer current value
	err := r.client.DeleteIndexerContext(ctx, indexer.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+indexerResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *Indexer) write(ctx context.Context, indexer *radarr.IndexerOutput) {
	i.EnableAutomaticSearch = types.BoolValue(indexer.EnableAutomaticSearch)
	i.EnableInteractiveSearch = types.BoolValue(indexer.EnableInteractiveSearch)
	i.EnableRss = types.BoolValue(indexer.EnableRss)
	i.Priority = types.Int64Value(indexer.Priority)
	i.DownloadClientID = types.Int64Value(indexer.DownloadClientID)
	i.ID = types.Int64Value(indexer.ID)
	i.ConfigContract = types.StringValue(indexer.ConfigContract)
	i.Implementation = types.StringValue(indexer.Implementation)
	i.Name = types.StringValue(indexer.Name)
	i.Protocol = types.StringValue(indexer.Protocol)
	i.Tags = types.SetValueMust(types.Int64Type, nil)
	i.MultiLanguages = types.SetValueMust(types.Int64Type, nil)
	i.RequiredFlags = types.SetValueMust(types.Int64Type, nil)
	i.Codecs = types.SetValueMust(types.Int64Type, nil)
	i.Mediums = types.SetValueMust(types.Int64Type, nil)
	i.Categories = types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, indexer.Tags, i.Tags.Type(ctx), &i.Tags)
	i.writeFields(ctx, indexer.Fields)
}

func (i *Indexer) writeFields(ctx context.Context, fields []*starr.FieldOutput) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(indexerStringFields, f.Name) {
			tools.WriteStringField(f, i)

			continue
		}

		if slices.Contains(indexerBoolFields, f.Name) {
			tools.WriteBoolField(f, i)

			continue
		}

		if slices.Contains(indexerIntFields, f.Name) {
			tools.WriteIntField(f, i)

			continue
		}

		if slices.Contains(indexerFloatFields, f.Name) {
			tools.WriteFloatField(f, i)

			continue
		}

		if slices.Contains(indexerIntSliceFields, f.Name) {
			tools.WriteIntSliceField(ctx, f, i)
		}
	}
}

func (i *Indexer) read(ctx context.Context) *radarr.IndexerInput {
	var tags []int

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	return &radarr.IndexerInput{
		EnableAutomaticSearch:   i.EnableAutomaticSearch.ValueBool(),
		EnableInteractiveSearch: i.EnableInteractiveSearch.ValueBool(),
		EnableRss:               i.EnableRss.ValueBool(),
		Priority:                i.Priority.ValueInt64(),
		DownloadClientID:        i.DownloadClientID.ValueInt64(),
		ID:                      i.ID.ValueInt64(),
		ConfigContract:          i.ConfigContract.ValueString(),
		Implementation:          i.Implementation.ValueString(),
		Name:                    i.Name.ValueString(),
		Protocol:                i.Protocol.ValueString(),
		Tags:                    tags,
		Fields:                  i.readFields(ctx),
	}
}

func (i *Indexer) readFields(ctx context.Context) []*starr.FieldInput {
	var output []*starr.FieldInput

	for _, b := range indexerBoolFields {
		if field := tools.ReadBoolField(b, i); field != nil {
			output = append(output, field)
		}
	}

	for _, n := range indexerIntFields {
		if field := tools.ReadIntField(n, i); field != nil {
			output = append(output, field)
		}
	}

	for _, f := range indexerFloatFields {
		if field := tools.ReadFloatField(f, i); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range indexerStringFields {
		if field := tools.ReadStringField(s, i); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range indexerIntSliceFields {
		if field := tools.ReadIntSliceField(ctx, s, i); field != nil {
			output = append(output, field)
		}
	}

	return output
}
