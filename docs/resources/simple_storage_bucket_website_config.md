---
subcategory: Cloud Simple Storage
page_title: "Bizfly Cloud: bizflycloud_simple_storage_bucket_versioning"
description: |-
  Provides a Bizfly Cloud Simple Storage Bucket Website Config resource. This can be change Simple Storage Bucket Website Config
---

# Resource: bizflycloud_simple_storage_bucket_website_config

Provides a Bizfly Cloud Simple Storage Bucket Website Config resource. This can be change Simple Storage Bucket Website Config

## Example Update Simple Storage Bucket Website Config

```hcl
resource "bizflycloud_simple_storage_bucket_website_config" "bucket_website_config_example" {
    bucket_name = "newtest"
    index = "tttt.html"
    error = "okokokoefe"
}
```
    
## Argument Reference

The following arguments are supported:

-   `bucket_name` - (Required) The name of Simple Storage Bucket
-   `index` - (Required) The file index of Simple Storage Bucket
-   `error` - (Required) The error of Simple Storage Bucket

## Attributes Reference

The following attributes are exported:
-   `message` - This attribute helps verify that the changes have been applied and confirms the action.
-   `website_url` -  This attribute provides the publicly accessible web address where the bucket’s website can be viewed.
-   `index` - When users visit the website URL, this file will be loaded as the homepage. It is typically an HTML file like index.html
-   `error` - This file acts as a fallback to display error messages, guiding users when the requested content is not available.

## Import

Bizfly Cloud Simple Storage Bucket Website Config resource can be imported using the server id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_simple_storage_bucket_website_config.bucket_website_config_example bucket_name
```

