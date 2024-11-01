---
page_title: "Provider: Bizfly Cloud"
description: |-
    The Bizfly Cloud provider is used to interact with the resources supported by Bizfly Cloud. The provider needs to be configured with the proper credentials before it can be used.
---

Bizfly Cloud Provider

The Bizfly Cloud provider is used to interact with the
resources supported by Bizfly Cloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Bizfly Cloud Provider
provider "bizflycloud" {
    auth_method = "password"
    region_name = "HaNoi"
    email = "email@domain.com"
    password = "password"
    project_id = "project_id"
}

# Create a database server
resource "bizflycloud_server" "db" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

-   `auth_method` - (Required) This is Authentication method. Value can be password or application_credential. Alternatively, this can also be specified
    using environment variables ordered by precedence:
    -   `BIZFLYCLOUD_AUTH_METHOD`
-   `api_endpoint` - (Optional) This can be used to override the base URL for
    Bizfly Cloud API requests (Defaults to the value of the `BIZFLYCLOUD_API_ENDPOINT`
    environment variable or `https://manage.bizflycloud.vn` if unset).
-   `email` - (Optional) This is your email to authenticate with Bizfly Cloud. Alternatively, this can also be specified using environment
    variables ordered by precedence:

    -   `BIZFLYCLOUD_EMAIL`

-   `password` - (Optional) This is your password to authenticate with Bizfly Cloud. Alternatively, this can also be specified using environment
    variables ordered by precedence:

    -   `BIZFLYCLOUD_PASSWORD`

-   `application_credential_id` - (Optional) This is your application credential ID authenticate with Bizfly Cloud. Alternatively, this can also be specified using environment
    variables ordered by precedence:

    -   `BIZFLYCLOUD_APPLICATION_CREDENTIAL_ID`

-   `application_credential_secret` - (Optional) This is your application credential secret authenticate with Bizfly Cloud. Alternatively, this can also be specified using environment
    variables ordered by precedence:

    -   `BIZFLYCLOUD_APPLICATION_CREDENTIAL_SECRET`

-   `region_name` - (Required) This is the region of resource you are working. Alternatively, this can also be specified using environment variables ordered by precedence:

    -   `BIZFLYCLOUD_REGION_NAME`

-   `project_id` - (Optional) This is the project ID of resource you are working. Alternatively, this can also be specified using environment variables ordered by precedence:
    -   `BIZFLYCLOUD_PROJECT_ID`
