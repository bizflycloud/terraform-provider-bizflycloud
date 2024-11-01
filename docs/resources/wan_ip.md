---
subcategory: "Cloud Server"
page_title: "Bizfly Cloud: bizflycloud_wan_ip"
description: |-
    Provides a Bizfly Cloud WAN IP resource. This can be used to create, modify, and delete WAN IP.
---

# Resource: bizflycloud_wan_ip

Provides a Bizfly Cloud WAN IP resource. This can be used to create,
modify, and delete WAN IP.

## Example Usage

```hcl
# Create a new WAN IP and attach to a server
resource "bizflycloud_wan_ip" "test_wan_1" {
  name = "sapd-wan-ip-tf2"
  availability_zone = "HN1"
  attached_server = "61fe3c90-7db0-47ba-b034-06de66a0869b"
}
```

## Argument Reference

The following arguments are supported:

-   `name` - (Required) Name of the WAN IP.
-   `availability_zone` - (Required) Availability zone of the WAN IP.
-   `firewall_ids` - Firewall IDs of the WAN IP.

## Attributes Reference

The following attributes are exported:

-   `id` - ID of the WAN IP.
-   `name` - Name of the WAN IP.
-   `availability_zone` - Availability zone of the WAN IP.
-   `status` - Status of the WAN IP.
-   `network_id` - Network ID of the WAN IP.
-   `tenant_id` - Tenant ID of the WAN IP.
-   `server_id` - Server ID of the WAN IP.
-   `firewall_ids` - Firewall IDs of the WAN IP.
-   `description` - Description of the WAN IP.
-   `bandwidth` - Bandwidth of the WAN IP.
-   `billing_type` - Billing type of the WAN IP.
-   `ip_address` - IP address of the WAN IP.
-   `ip_version` - IP version of the WAN IP.

## Import

Bizfly Cloud wan ip resource can be imported using the wan ip id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_wan_ip.test_wan_1 wan-ip-id
```
