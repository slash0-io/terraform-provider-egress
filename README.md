# terraform-provider-egress

Terraform data sources for the published IP ranges of third-party services — Stripe, GitHub, Datadog, Okta, Cloudflare, and 25+ more — so `0.0.0.0/0` never has to appear in your egress rules again.

```hcl
data "egress_ranges" "stripe_api" {
  service = "stripe"
  purpose = "api"
}

resource "aws_security_group_rule" "stripe_egress" {
  type              = "egress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = data.egress_ranges.stripe_api.ipv4_cidrs
  security_group_id = aws_security_group.app.id
  description       = "Stripe API"
}
```

## How it works

The provider reads a **versioned public feed** ([egresshq/feed](https://github.com/egresshq/feed)) rebuilt continuously from each vendor's *official* publication — never from third-party aggregators. Every service document carries a provenance chain: the upstream URL, retrieval timestamp, and the SHA-256 of the upstream body it was derived from.

The catalog records two things most sources don't model:

- **`direction`** — `egress` ranges are what your workloads connect *to* (SG egress rules); `ingress` ranges are what the service connects *from* (webhook/agent sources, for SG ingress rules). A third of published ranges are the latter, and putting them in egress rules is a silent no-op.
- **`classification`** — `dedicated` (vendor-owned space, safe to pin), `mixed`, or `cdn-shared` (pinning would allowlist an entire CDN). The feed also publishes a **non-publishers list**: services whose vendors state that IP pinning is unsupported, with their recommended alternative.

## Data sources

| Data source | Purpose |
|---|---|
| `egress_ranges` | Current ranges for one service purpose (`stripe`/`api`, `github`/`hooks`, `datadog`/`agents`, …) |
| `egress_services` | The full catalog: slugs, purposes, directions, classifications |

Provider configuration: `feed_url` (optional) — defaults to the `EGRESS_FEED_URL` environment variable, then the public feed. `file://` URLs are supported for air-gapped or vendored feeds.

## Staying current

Terraform data sources refresh only at `plan`/`apply` time. If you apply infrequently, pair the provider with scheduled applies — or don't manage the drift yourself: the hosted tier keeps AWS-native managed prefix lists continuously updated and shared into your account, with staged rollouts and change notifications. *(In development.)*

## Development

```sh
go build -o bin/terraform-provider-egress .
go test ./...
```

To run [examples/basic](examples/basic) against a local feed: build a feed with the [generator](https://github.com/egresshq/feed), point a `dev_overrides` CLI config at `bin/`, and set `EGRESS_FEED_URL=file:///path/to/dist/v1`.

The feed schema in `internal/feedschema` is a vendored copy of the canonical definition in [egresshq/feed](https://github.com/egresshq/feed); schema v1 is frozen (additive changes only).

## License

[MPL-2.0](LICENSE)
