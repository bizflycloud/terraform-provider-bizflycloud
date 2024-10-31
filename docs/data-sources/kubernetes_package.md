---
subcategory: Cloud Kubernetes Engine
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_kubernetes_package"
description: |-
Provides a Bizfly Cloud Kubernetes Package
---

# Data Source: bizflycloud_kubernetes_package

Get information about Bizfly Cloud Kubernetes Package

## Example Usage

```hcl
# Get information of an kubernetes package
data "bizflycloud_kubernetes_package" "k8s_package_standard_1" {
  provision_type = "standard"
  name = "STANDARD - 1"
}
```

## Argument Reference

The following arguments are supported:

-   `provision_type` - (Required) The Kubernetes provision type of cluster: standard or everywhere

-   `name` - (Required) The name of Package

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of Kubernetes package
