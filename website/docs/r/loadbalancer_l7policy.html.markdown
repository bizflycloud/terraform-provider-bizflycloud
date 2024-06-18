---
subcategory: Cloud Load Balancer
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_loadbalancer_l7policy"
description: |-
  Provides a Bizfly Cloud L7 policy of Load Balancer resource. This can be used to create, modify, and delete l7 policy of Load Balancer.
---

# Resource: bizflycloud_loadbalancer_l7policy

Provides a Bizfly Cloud L7 policy of Load Balancer resource. This can be used to create,
modify, and delete l7 policy of Load Balancer.

## Example Create L7 policy for Load Balancer

```hcl
# Create a new L7 policy for Load Balancer
resource "bizflycloud_loadbalancer_l7policy" "tf_policy" {
    name = "bizfly-l7-policy"
    action = "REDIRECT_TO_POOL"
    redirect_pool_id = "${bizflycloud_loadbalancer_pool.pool.id}"
    listener_id = "${bizflycloud_loadbalancer_listener.listener.id}"
    position = 1
    rules {
        invert = true
        type = "HOST_NAME"
        compare_type = "EQUAL_TO"
        value = "manage.bizflycloud.vn"
    }
    rules {
        invert = false
        type = "HEADER"
        compare_type = "EQUAL_TO"
        key = "X-Tenant-Id"
        value = "12345"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of l7 policy
* `action` - (Required) The action for l7 policy: `REDIRECT_TO_POOL`, `REJECT`, `REDIRECT_TO_URL`, `REDIRECT_PREFIX`
* `redirect_pool_id` - (Optional) The redirect pool id for l7 policy with action is `REDIRECT_TO_POOL`
* `redirect_prefix` - (Optional) The redirect prefix for l7 policy with action is `REDIRECT_PREFIX`
* `redirect_url` - (Optional) The redirect url for l7 policy with action is `REDIRECT_TO_URL`
* `listener_id` - (Required) The ID of listener to which l7 policy will apply
* `position` - (Optional) The position in list l7 policy of listener (Default: 1)
* `rules` - (Optional) The list ACLs of l7 policy
  - `invert` - (Required) The invert: `true`, `false`
  - `type` - (Required) The type: `HOST_NAME`, `PATH`, `HEADER`, `FILE_TYPE`
  - `compare_type` - (Required) The compare type: `EQUAL_TO`, `REGEX`, `CONTAINS`, `ENDS_WITH`, `STARTS_WITH`
  - `key` - (Optional) The key with rule type is `HEADER`
  - `value` - (Required) The value
## Attributes Reference

The following attributes are exported:
* `id` - The ID of l7 policy
* `name` - The name of l7 policy
* `action` - The action for l7 policy
* `redirect_pool_id` - The redirect pool id for l7 policy
* `redirect_prefix` - The redirect prefix for l7 policy 
* `redirect_url` - The redirect url for l7 policy
* `listener_id` - The ID of listener to which l7 policy is applied
* `position` - The position in list l7 policy of listener
* `rules` - The list ACLs of l7 policy
  - `id` - The ID of l7 policy rule
  - `invert` - The invert
  - `type` - The type
  - `compare_type` - The compare type
  - `key` - The key 
  - `value` - The value
  - `operating_status` - The operating status
  - `provisioning_status` - The provisioning status
  - `project_id` - The project id
  - `created_at` - The created at 
  - `updated_at` - The updated at

## Import

Bizfly Cloud loadbalancer l7 policy resource can be imported using the l7 policy id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_loadbalancer_l7policy.tf_policy l7policy-id
```