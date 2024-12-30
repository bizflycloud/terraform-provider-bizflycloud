---
subcategory: Cloud Simple Storage
page_title: "Bizfly Cloud: bizflycloud_simple_storage_bucket"
description: |-
    Provides a Bizfly Cloud Simple Storage Bucket resource. This can be used to create, modify, and delete Simple Storage Bucket.
---

# Resource: bizflycloud_simple_storage_bucket

Provides a Bizfly Cloud Simple Storage Bucket resource. This can be used to create,
modify, and delete Simple Storage Bucket.

## Example Create Simple Storage Bucket with COLD type storage  

```hcl
resource "bizflycloud_simple_storage_bucket" "bucket_example" {
    name = "newtest"
    location = "hn"
    default_storage_class = "COLD"
}
```

## Example Create a simple Storage Bucket with STANDARD type storage

```hcl
# Create a new Simple Storage Bucket with only internal network
resource "bizflycloud_simple_storage_bucket" "bucket_example1" {
    name = "newtest"
    location = "hn"
    default_storage_class = "STANDARD"
}
```

## Argument Reference

The following arguments are supported:

-   `name` - (Required) The name of Simple Storage Bucket
-   `location` - (Required) The location of Simple Storage Bucket: `hn` or `hcm`
-   `acl` - (Required) - The type of acl: `private`
-   `default_storage_class` - (Required) The type of Simple Storage Bucket: `STANDARD` or `COLD`

## Attributes Reference

The following attributes are exported:

-   `name`- The name of the Simple Storage Bucket
-   `location` - The location of the Simple Storage Bucket
-   `description` - The description of Simple Storage Bucket
-   `default_storage_class` - The default storage type applied to data in the bucket
-   `num_objects` - Represents the total count of objects stored within the bucket
-   `created_at` - The time created simple storage
-   `size_kb` - The provisioning status of Simple Storage Bucket

## Import

Bizfly Cloud Simple Storage Bucket resource can be imported using the server id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_simple_storage_bucket.bucket_example bucket_name
```