
# Resource: bizflycloud_simple_storage_bucket_cors

Provides a Bizfly Cloud Simple Storage Bucket Cors resource. This can be change Simple Storage Bucket Cors

## Example Update Simple Storage Bucket Cors

```hcl
resource "bizflycloud_simple_storage_bucket_cors" "example" {
  bucket_name = "newtest"

  rules {
    allowed_origin  = "http://ahoho.com"
    allowed_methods = ["PUT"]
    allowed_headers = ["Content-Type"]
    max_age_seconds = 6400
  }

  rules {
    allowed_origin  = "http://another-origin.com"
    allowed_methods = ["POST"]
    allowed_headers = ["Authorization"]
    max_age_seconds = 7200
  }
}
```

## Argument Reference

The following arguments are supported:

-   `bucket_name` - (Required) The name of the Simple Storage Bucket for which the CORS configuration will be applied.

-   `rules` - (Required) A block defining the CORS rules to be applied. Each `rules` block supports the following arguments:
    -   `allowed_origin` - (Required) Specifies the origin that is allowed to make requests to the bucket. Must be a valid URI.
    -   `allowed_methods` - (Required) A list of HTTP methods that are allowed for the specified origin. Examples include `GET`, `PUT`, `POST`, etc.
    -   `allowed_headers` - (Optional) A list of headers that are allowed in requests from the specified origin.
    -   `max_age_seconds` - (Optional) The maximum amount of time, in seconds, that the results of a preflight request can be cached by the browser.

## Attributes Reference

The following attributes are exported:

-   `bucket_name` - The name of the bucket for which the CORS configuration is applied.

-   `rules` - A list of rules defining the applied CORS configuration. Each rule contains:
    -   `allowed_origin` - The origin that is allowed to make requests to the bucket.
    -   `allowed_methods` - A list of HTTP methods that are allowed for the specified origin.
    -   `allowed_headers` - A list of headers allowed in requests from the specified origin.
    -   `max_age_seconds` - The maximum time, in seconds, that the results of a preflight request can be cached by the browser.
