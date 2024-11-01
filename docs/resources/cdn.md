---
subcategory: Cloud CDN
page_title: "Bizfly Cloud: bizflycloud_cdn"
description: |-
    Provides a Bizfly Cloud CDN resource. This can be used to create, modify, and delete CDN.
---

# Resource: bizflycloud_cdn

Provides a Bizfly Cloud CDN resource. This can be used to create,
modify, and delete CDN.

## Example Usage

```hcl
# Create a new CDN resource
resource "bizflycloud_cdn" "domain_com" {
    domain = "cdn.domain.com"
    origin =  {
        upstream_addrs = "origin.domain.com"
        upstream_host = "origin.domain.com"
        upstream_proto = "https"
        name = "origin-domain"
    }
}
```

## Argument Reference

The following arguments are supported:

-   `domain` - (Required) Specifies the domain name of the CDN Endpoint.
-   `origin` - (Required) The origin of the CDN endpoint
    -   `name` - (Required) The name of the origin.
    -   `upstream_addrs` - (Required) A string that determines the hostname/IP address of the origin server. This string can be a domain name, Storage Account endpoint, Web App endpoint or IPv4 address.
    -   `upstream_host` - (Optional) The host header CDN provider will send along with content requests to origins.
    -   `upstream_proto` - (Optional) Origin protocol (http/https). Default value is http

## Attributes Reference

The following attributes are exported:

-   `domain_id`: Identifier for the CDN Endpoint. Example: 9805be94-ccf1-4551-b7d0-5e8fcd7804b6
-   `domain_cdn`: Domain name corresponding to the CDN Endpoint. Example: domain.cdn.vccloud.vn
