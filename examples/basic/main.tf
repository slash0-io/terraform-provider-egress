terraform {
  required_providers {
    egress = {
      source = "egresshq/egress"
    }
  }
}

# feed_url is resolved from (in order): this block, EGRESS_FEED_URL, the
# public feed default.
provider "egress" {}

# Stripe's API endpoints — the ranges your workloads connect out to.
data "egress_ranges" "stripe_api" {
  service = "stripe"
  purpose = "api"
}

# GitHub webhook sources — direction is "ingress": these belong in ingress
# rules on whatever receives your webhooks, not in egress rules.
data "egress_ranges" "github_hooks" {
  service = "github"
  purpose = "hooks"
}

data "egress_services" "catalog" {}

output "stripe_api_ipv4" {
  value = data.egress_ranges.stripe_api.ipv4_cidrs
}

output "stripe_direction" {
  value = data.egress_ranges.stripe_api.direction
}

output "github_hooks_cidrs" {
  value = data.egress_ranges.github_hooks.cidrs
}

output "catalog_size" {
  value = length(data.egress_services.catalog.services)
}

# The real-world shape: an egress security group rule fed by the data source.
# (Commented out so this example plans without AWS credentials.)
#
# resource "aws_security_group_rule" "stripe_egress" {
#   type              = "egress"
#   from_port         = 443
#   to_port           = 443
#   protocol          = "tcp"
#   cidr_blocks       = data.egress_ranges.stripe_api.ipv4_cidrs
#   security_group_id = aws_security_group.app.id
#   description       = "Stripe API (managed by egress feed ${data.egress_ranges.stripe_api.sync_token})"
# }
