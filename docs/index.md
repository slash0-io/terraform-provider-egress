# egress Provider

Data sources for the published IP ranges of third-party services (Stripe, GitHub, Datadog, Okta, …), backed by a versioned public feed built exclusively from each vendor's official publication.

## Example Usage

```terraform
provider "egress" {}

data "egress_ranges" "stripe_api" {
  service = "stripe"
  purpose = "api"
}
```

## Schema

### Optional

- `feed_url` (String) Feed base URL (the directory containing `index.json`). Supports `http(s)://` and `file://`. Defaults to the `EGRESS_FEED_URL` environment variable, then the public feed.
