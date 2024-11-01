---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_ssh_key"
description: |-
    Provides a Bizfly Cloud SSH Key resource. This can be used to create and delete ssh key.
---

# Resource: bizflycloud_ssh_key

Provides a Bizfly Cloud SSH key resource. This can be used to create,
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

-   `name` - (Required) The name of the ssh key.
-   `public_key` - (Required) The public key of the ssh key.

## Attributes Reference

The following attributes are exported:

-   `name`- The name of the ssh key
-   `public_key` - The public key of the ssh key
-   `fingerprint` - The finger print of the ssh key

## Import

Bizfly Cloud ssh key resource can be imported using the ssh key name

```
$ terraform import bizflycloud_ssh_key.sshkey ssh-key-name
```
