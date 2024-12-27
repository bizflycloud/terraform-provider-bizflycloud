---
subcategory: Cloud Simple Storage Bucket
page_title: "Bizfly Cloud: bizflycloud_loadbalancer"
description: |-
    Provides a Bizfly Cloud Simple Storage Bucket resource. This can be used to create, modify, and delete Simple Storage Buckets.
---

# Resource: bizflycloud_simple_storage_bucket

Provides a Bizfly Cloud Simple Storage Bucket resource. This can be used to create,
modify, and delete Simple Storage Bucket.

## Example Create Simple Storage Bucket with private access control  

```hcl
resource "bizflycloud_simple_storage_bucket" "bucket_example" {
    name = "newtest"
    location = "hn"
    acl = "private"
    default_storage_class = "COLD"
}
```

## Example Create a simple Storage Bucket with public access control

```hcl
# Create a new Simple Storage Bucket with only internal network
resource "bizflycloud_simple_storage_bucket" "bucket_example1" {
    name = "newtest"
    location = "hn"
    acl = "public-read"
    default_storage_class = "COLD"
}
```

## Argument Reference

The following arguments are supported:

-   `name` - (Required) The name of Simple Storage Bucket
-   `location` - (Required) The location of Simple Storage Bucket: `hn`
-   `acl` - (Required) - The type of acl: `public-read` or `private`
-   `default_storage_class` - (Required) The type of Simple Storage Bucket: `STANDARD` or `COLD`

## Attributes Reference

The following attributes are exported:

-   `name`- The name of the Simple storage bucket
-   `location` - The location of the Simple storage bucket
-   `description` - The description of Simple storage bucket
-   `default_storage_class` - The default storage type applied to data in the bucket
-   `num_objects` - Represents the total count of objects stored within the bucket
-   `created_at` - The time created simple storage
-   `size_kb` - The provisioning status of Simple Storage Bucket

