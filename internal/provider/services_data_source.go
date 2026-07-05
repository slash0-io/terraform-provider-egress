package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = (*servicesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*servicesDataSource)(nil)
)

func NewServicesDataSource() datasource.DataSource { return &servicesDataSource{} }

type servicesDataSource struct{ client *FeedClient }

type servicesModel struct {
	Services    types.List   `tfsdk:"services"`
	SyncToken   types.String `tfsdk:"sync_token"`
	GeneratedAt types.String `tfsdk:"generated_at"`
}

type serviceEntryModel struct {
	Slug           types.String       `tfsdk:"slug"`
	Name           types.String       `tfsdk:"name"`
	Category       types.String       `tfsdk:"category"`
	Classification types.String       `tfsdk:"classification"`
	Purposes       []purposeEntryModel `tfsdk:"purposes"`
}

type purposeEntryModel struct {
	Key       types.String `tfsdk:"key"`
	Direction types.String `tfsdk:"direction"`
}

var purposeObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":       types.StringType,
	"direction": types.StringType,
}}

var serviceObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"slug":           types.StringType,
	"name":           types.StringType,
	"category":       types.StringType,
	"classification": types.StringType,
	"purposes":       types.ListType{ElemType: purposeObjType},
}}

func (d *servicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (d *servicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The feed catalog: every service with published, pinnable IP ranges.",
		Attributes: map[string]schema.Attribute{
			"services": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"slug":           schema.StringAttribute{Computed: true},
						"name":           schema.StringAttribute{Computed: true},
						"category":       schema.StringAttribute{Computed: true},
						"classification": schema.StringAttribute{Computed: true},
						"purposes": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key":       schema.StringAttribute{Computed: true},
									"direction": schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
			"sync_token":   schema.StringAttribute{Computed: true},
			"generated_at": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *servicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*FeedClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *servicesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	idx, err := d.client.Index(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Reading feed index failed", err.Error())
		return
	}

	entries := make([]serviceEntryModel, 0, len(idx.Services))
	for _, s := range idx.Services {
		e := serviceEntryModel{
			Slug:           types.StringValue(s.Slug),
			Name:           types.StringValue(s.Name),
			Category:       types.StringValue(s.Category),
			Classification: types.StringValue(s.Classification),
		}
		for _, p := range s.Purposes {
			e.Purposes = append(e.Purposes, purposeEntryModel{
				Key:       types.StringValue(p.Key),
				Direction: types.StringValue(p.Direction),
			})
		}
		entries = append(entries, e)
	}

	list, diags := types.ListValueFrom(ctx, serviceObjType, entries)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state := servicesModel{
		Services:    list,
		SyncToken:   types.StringValue(idx.SyncToken),
		GeneratedAt: types.StringValue(idx.GeneratedAt),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func mustStringList(ctx context.Context, diags *diag.Diagnostics, ss []string) types.List {
	if ss == nil {
		ss = []string{}
	}
	l, d := types.ListValueFrom(ctx, types.StringType, ss)
	diags.Append(d...)
	return l
}
