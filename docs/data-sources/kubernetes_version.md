---
subcategory: Cloud Kubernetes Engine
page_title: "Bizfly Cloud: bizflycloud_kubernetes_version"
description: |-
    Provides a Bizfly Cloud Kubernetes Version
---

# Data Source: bizflycloud_kubernetes_version

Get ìnformation about Bizfly Cloud Kubernetes Version

## Example Usage

```hcl
# Get information of an kubernetes version
data "bizflycloud_kubernetes_version" "test_k8s_version" {
  version = "v1.17.9"
}
```

## Argument Reference

The following arguments are supported:

-   `version` - (Required) The Kubernetes version

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of Kubernetes version
-   `version` - The Kubernetes version
