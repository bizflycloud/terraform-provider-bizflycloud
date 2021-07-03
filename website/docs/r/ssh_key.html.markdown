---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_ssh_key"
sidebar_current: "docs-bizflycloud-resource-ssh-key"
description: |-
  Provides a BizFly Cloud SSH Key resource. This can be used to create and delete ssh key.
---

# bizflycloud\_ssh_key

Provides a BizFly Cloud SSH key resource. This can be used to create,
and delete ssh key.
## Example Usage

```hcl
# Create a new ssh key
resource "bizflycloud_ssh_key" "sshkey" {
    name = "ssh-key-1"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDgXX0Kdd3XKojgj7maVd3PsPApzh9n2lT2CtgcJs8jw9i3mit5SZu02QFS772Pa9VdGeSjbqxtADLRpnuigW5ii0dHBQTgWqx593Cs7QKRhyRPb88u0TFCZynRwfMRnb6qngiKoWp5TtaHuIY+7kS8SyqNVIwoCYlr9a4ePX8rwydf9crhJocgKb2LgQkdW3TBE5QAvxbruYlj201jjXFeE5BtE4QER0QyY5MqW8MAgG98N3w95pKIffhHZ4TO4A3zgpWbNn1ROproZgV+9COzZ7WYuvPWqWdLAntd9b1/lLnDrDHXa/lrefJXJVamhz4i1cfIZ/p+aFWG0a7DpL5b"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ssh key.
* `public_key` - (Required) The public key of the ssh key.


## Attributes Reference

The following attributes are exported:

* `name`- The name of the ssh key
* `public_key` - The public key of the ssh key
* `fingerprint` - The finger print of the ssh key

## Import

BizFly Cloud server resource can be imported using the server id

You can obtain server id from:

- The UUID part in URL while managing the server in the web interface `https://hn.manage.bizflycloud.vn/iaas-cloud/servers/<ID>/details`

- API


```
$ terraform import bizflycloud_ssh_key.example 123e4567-e89b-12d3-a456-426614174000
```