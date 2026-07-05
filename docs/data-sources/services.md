# egress_services (Data Source)

The feed catalog: every service with published, pinnable IP ranges.

## Example Usage

```terraform
data "egress_services" "catalog" {}

output "catalog" {
  value = [for s in data.egress_services.catalog.services : "${s.slug} (${s.classification})"]
}
```

## Schema

### Read-Only

- `services` (List of Object)
  - `slug` (String)
  - `name` (String)
  - `category` (String)
  - `classification` (String) `dedicated` | `mixed` | `cdn-shared`
  - `purposes` (List of Object)
    - `key` (String)
    - `direction` (String) `egress` | `ingress` | `both`
- `sync_token` (String)
- `generated_at` (String)
