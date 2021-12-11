---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_wan_ip"
sidebar_current: "docs-bizflycloud-datasource-wan-ip"
description: |-
Provides a Bizfly Cloud WAN IP resource. This can be used to create, modify, and delete WAN IP.
---

# bizflycloud\_wan_ip

Get information about Bizfly Cloud WAN IP Resource.

## Example Usage

```hcl
# Get information of WAN IP from datasource
data "bizflycloud_wan_ip" "wan_ip" {
  ip_address = "103.148.57.12"
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` (Required) - The IP address of the WAN IP.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the WAN IP.
* `name` - Name of the WAN IP.
* `availability_zone` - Availability zone of the WAN IP.
* `status` - Status of the WAN IP.
* `network_id` - Network ID of the WAN IP.
* `tenant_id` - Tenant ID of the WAN IP.
* `device_id` - Device ID of the WAN IP.
* `security_groups` - Security group IDs of the WAN IP.
* `description` - Description of the WAN IP.
* `bandwidth` - Bandwidth of the WAN IP.
* `billing_type` - Billing type of the WAN IP.
* `ip_address` - IP address of the WAN IP.
* `ip_version` - IP version of the WAN IP.
