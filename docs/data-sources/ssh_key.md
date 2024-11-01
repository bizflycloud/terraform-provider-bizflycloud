---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_ssh_key"
description: |-
    Provides a Bizfly Cloud SSH Key datasource. This can be used to read SSH Key.
---

# Data Source: bizflycloud_ssh_key

Get information about Bizfly Cloud SSH Key resource.

## Example Usage

```hcl
# Get information of SSH Key
data "bizflycloud_ssh_key" "ssh_key" {
  name = "ssh_key_name"
}
```

## Argument Reference

The following arguments are supported:

-   `name` - The name of SSH Key

## Attributes Reference

The following attributes are exported:

-   `name`- The name of the SSH key.
-   `public_key`- The public key of the SSH key
-   `fingerprint`- The finger
