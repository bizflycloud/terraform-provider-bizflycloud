---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_dns"
sidebar_current: "docs-bizflycloud-resource-dns"
description: |-
  Provides a Bizfly Cloud DNS resource. This can be used to create, modify, and delete DNS.
---

# bizflycloud\_dns

Provides a Bizfly Cloud DNS resource. This can be used to create,
modify, and delete DNS.

## Example Usage

```hcl
# Create a new DNS zone
resource "bizflycloud_dns" "dns_zone" {
    name = "abc.xyz"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of DNS.
* `description` - (Optional) The description of DNS.
* `required` - (Optional) The required of DNS: true or false.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of DNS
* `name`- The name of DNS
* `description` - The description of DNS
* `required` - The required of DNS: true or false.
* `active` - The active of the DNS: true or false.
* `deleted` - Number deleted of DNS.
* `nameserver` - List name server of the DNS.
* `tenant_id` - The tenant ID of DNS.
* `ttl` - The time to live of DNS.
* `record_set` - The list record of DNS
  * `id` - The id of record.
  * `name` - The name of record.
  * `ttl` - The time to live of record.
  * `type` - The type of record.
* `create_at` - The created time.
* `updated_at` - The updated time.
