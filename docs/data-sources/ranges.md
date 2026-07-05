# egress_ranges (Data Source)

Current published IP ranges for one service purpose. Ranges refresh on every `terraform plan`/`apply`.

## Example Usage

```terraform
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
}
```

## Schema

### Required

- `service` (String) Service slug, e.g. `stripe`. Enumerate with `egress_services`.

### Optional

- `purpose` (String) Purpose key, e.g. `api` or `webhooks`. May be omitted only when the service publishes exactly one purpose.

### Read-Only

- `ipv4_cidrs` (List of String) Sorted IPv4 CIDRs.
- `ipv6_cidrs` (List of String) Sorted IPv6 CIDRs.
- `cidrs` (List of String) `ipv4_cidrs` followed by `ipv6_cidrs`.
- `direction` (String) `egress` = ranges you connect to (SG egress rules); `ingress` = ranges the service connects from (webhook sources — SG ingress rules).
- `classification` (String) `dedicated` | `mixed` | `cdn-shared`.
- `name` (String) Human-readable service name.
- `sync_token` (String) Feed sync token at generation time.
- `generated_at` (String) Feed generation timestamp (RFC 3339).
