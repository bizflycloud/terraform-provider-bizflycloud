---
subcategory: Cloud Simple Storage
page_title: "Bizfly Cloud: bizflycloud_simple_storage_bucket_acl"
description: |-
  Provides a Bizfly Cloud Simple Storage Bucket Acl resource. This can be change Simple Storage Bucket Acl
---


# Resource: bizflycloud_simple_storage_bucket_acl

Provides a Bizfly Cloud Simple Storage Bucket Acl resource. This can be change Simple Storage Bucket Acl

## Example Update Simple Storage Bucket Acl

```hcl
resource "bizflycloud_simple_storage_bucket_acl" "bucket_acl_example" {
    bucket_name = "newtest2"
    acl = "public-read"
}
```


## Argument Reference

The following arguments are supported:

-   `bucket_name` - (Required) The name of Simple Storage Bucket
-   `acl` - (Required) - The type of acl: `private` or `public-read`

## Attributes Reference

The following attributes are exported:

-   `message` - A string containing a status message about the ACL update, indicating the new ACL setting. For example: `"Bucket ACL đã được thay đổi thành: private"`.

-   `owner` - Information about the bucket's owner, including:
    - `id` - The unique identifier of the bucket's owner.
    - `display_name` - The display name of the bucket's owner.

-   `grants` - A list of grants specifying the permissions granted for the bucket. Each grant includes:
    - `permission` - The level of access granted. Possible values include:
        - `FULL_CONTROL` - Grants full control over the bucket.
        - `<none>` - Do not grant access to the bucket.
    - `grantee` - Details about the entity receiving the permission. Includes:
        - `type` - The type of grantee. Possible values include:
            - `CanonicalUser` - A specific user in Bizfly Cloud.
        - `id` - The unique identifier of the grantee (if applicable).
        - `display_name` - The display name of the grantee (if applicable).
        - `email` - The email address of the grantee (if available, otherwise `null`).
        - `uri` - A URI identifying the group (if the grantee type is `Group`, otherwise `null`).

## Import

Bizfly Cloud Simple Storage Bucket Acl resource can be imported using the server id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_simple_storage_bucket_acl.bucket_acl_example bucket_name
```