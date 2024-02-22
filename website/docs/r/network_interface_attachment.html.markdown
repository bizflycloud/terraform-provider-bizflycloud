---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_network_interface_attachment"
sidebar_current: "docs-bizflycloud_network_interface_attachment"
description: |-
Provides a Bizfly Cloud Network Interface Attachment resource. This can be used to create, modify, and delete the 
  attachment of network interface to server.
---

# bizflycloud\_network\_interface\_attachment

Provides a Bizfly Cloud Network Interface Attachment resource. This can be used to create, modify, and delete the 
  attachment of network interface to server.

## Example Usage

```hcl
resource "bizflycloud_network_interface_attachment" "test_attachment" {
  server_id            = bizflycloud_server.tf_server1.id
  network_interface_id = bizflycloud_network_interface.test_lan1.id
  firewall_ids         = [
    "9ff7fbdc-1461-4713-b4f8-ba8a22699bb4",
    "0b1d6a90-24a0-4f86-b75f-fb07290e44dd",
    "348b3be8-0610-4397-a781-10be10998b90"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `server_id` - (Required) The ID of the server to attach the network interface to.
* `network_interface_id` - (Required) The ID of the network interface to attach to the server.
* `firewall_ids` - (Optional) A list of the firewall IDs of the network interface


## Attribute Reference

In addition to all arguments above, the following attributes are exported:
* `id` - The ID of the network interface.
* `server_id` - The ID of the server to attach the network interface to.
* `network_interface_id` - The ID of the network interface to attach to the server.
* `firewall_ids` - A list of the firewall IDs of the network interface


## Import

Bizfly Cloud network interface attachment resource can be imported using the network interface id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_network_interface_attachment.test_attachment network-interface-id
```