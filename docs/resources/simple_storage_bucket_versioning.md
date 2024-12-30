---
subcategory: Cloud Simple Storage
page_title: "Bizfly Cloud: bizflycloud_simple_storage_bucket_versioning"
description: |-
  Provides a Bizfly Cloud Simple Storage Bucket Versioning resource. This can be change Simple Storage Bucket Versioning
---

# Resource: bizflycloud_simple_storage_bucket_versioning

Provides a Bizfly Cloud Simple Storage Bucket Versioning resource. This can be change Simple Storage Bucket Versioning

## Example Update Simple Storage Bucket Versioning

```hcl
resource "bizflycloud_simple_storage_bucket_versioning" "bucket_versioning_example" {
    bucket_name = "newtest2"
    versioning = false
}
```


## Argument Reference

The following arguments are supported:

-   `bucket_name` - (Required) The name of Simple Storage Bucket
-   `versioning` - (Required) The versioning of Simple Storage Bucket: `true` or `false`

## Attributes Reference

The following attributes are exported:

-   `message` - A message confirming the result of the versioning update for the Simple Storage Bucket.
-   `status` - The current versioning status of the Simple Storage Bucket. Possible values are:
    - `Suspended`: Versioning is disabled for the bucket.
    - `Enabled`: Versioning is active for the bucket.

## Import

Bizfly Cloud Simple Storage Bucket Versioning resource can be imported using the server id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_simple_storage_bucket_versioning.bucket_versioning_example bucket_name
```

