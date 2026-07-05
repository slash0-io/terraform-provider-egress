package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefaultFeedURL is the public feed. Overridable via the feed_url provider
// attribute or the EGRESS_FEED_URL environment variable (attribute wins).
const DefaultFeedURL = "https://egresshq.github.io/feed/v1"

var _ provider.Provider = (*egressProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider { return &egressProvider{version: version} }
}

type egressProvider struct{ version string }

type providerModel struct {
	FeedURL types.String `tfsdk:"feed_url"`
}

func (p *egressProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "egress"
	resp.Version = p.version
}

func (p *egressProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data sources for third-party service IP ranges (Stripe, GitHub, Datadog, ...), " +
			"backed by a versioned public feed built exclusively from each vendor's official publication.",
		Attributes: map[string]schema.Attribute{
			"feed_url": schema.StringAttribute{
				Optional: true,
				Description: "Feed base URL (the directory containing index.json). Supports http(s):// and file://. " +
					"Defaults to the EGRESS_FEED_URL environment variable, then the public feed.",
			},
		},
	}
}

func (p *egressProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	feedURL := DefaultFeedURL
	if v := os.Getenv("EGRESS_FEED_URL"); v != "" {
		feedURL = v
	}
	if !cfg.FeedURL.IsNull() && cfg.FeedURL.ValueString() != "" {
		feedURL = cfg.FeedURL.ValueString()
	}
	client := NewFeedClient(feedURL, "terraform-provider-egress/"+p.version)
	resp.DataSourceData = client
}

func (p *egressProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewRangesDataSource, NewServicesDataSource}
}

func (p *egressProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}
