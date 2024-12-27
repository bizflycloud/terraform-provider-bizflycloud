
# Resource: bizflycloud_simple_storage_access_key

Provides a Bizfly Cloud Simple Storage Access Key resource. This can be used to create,
modify, and delete Simple Storage Bucket access keys.

## Example Create Simple Storage Access Key with private access control

```hcl
resource "bizflycloud_simple_storage_access_key" "access_key_example" {
    subuser_id = "subuser-id-example"
    access_key = "your-access-key"
    secret_key = "your-secret-key"
}
```

## Argument Reference

The following arguments are supported:

- `subuser_id` - (Required) The unique identifier of the subuser associated with this access key.
- `access_key` - (Required) The primary key used for authenticating and accessing the resources in the Simple Storage Bucket. This key must be unique and comply with security guidelines.
- `secret_key` - (Required) A secret string paired with the `access_key` to ensure secure authentication. Keep it confidential and do not expose it publicly.

## Attributes Reference

The following attributes are exported:

- `subuser_id` - The identifier linking the access key to the correct subuser in the Bizfly Cloud system.
- `access_key` - The access key associated with the subuser.
