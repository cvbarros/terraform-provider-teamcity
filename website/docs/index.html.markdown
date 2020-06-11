---
layout: teamcity
page_title: "Provider: TeamCity"
sidebar_current: "docs-teamcity-index"
description: |-
  TeamCity provider is used to manage TeamCity resources. The provider must be configured with the proper credentials before it can be used.

# TeamCity Provider

The [TeamCity](https://www.jetbrains.com/teamcity/) provider is used to interact with the
resources supported by TeamCity. The provider needs to be configured
with the proper credentials before it can be used.

~> Important Interacting with TeamCity from Terraform causes any sensitive parameters and variables to be persisted in both Terraform's state file and in any generated plan files. Whenever dealing with sensitive parameters in projects and/or build configurations, files should be treated as sensitive and protected accordingly

This provider is meant primarily for enabling _pipelines as code_ for TeamCity.

## Provider Arguments

The provider configuration block accepts the following arguments. In general, it's better to set them via indicated environment variables to keep the configuration safe.

* `address` - (Required) Address of TeamCity server. This is a URL with a scheme, a hostname and port but no path. May be set via the `TEAMCITY_ADDR` environment variable.

---

If using Token Authentication - the following fields can be specified:

* `token` - (Required) The API Token which should be used to authenticate to TeamCity. The user must have broad permissions to create resources and manage projects at the appropriate hierarchy path. Refer to TeamCity documentation on how to apply proper roles and permissions for the user. It is recommended to be set via `TEAMCITY_TOKEN` environment variable.

---

If using a Username/Password for authentication - the following fields can be specified:

* `username` - (Required) Username that will be used to authenticate to TeamCity. The user must have broad permissions to create resources and manage projects at the appropriate hierarchy path. Refer to TeamCity documentation on how to apply proper roles and permissions for the user. It is recommended to be set via `TEAMCITY_USER` environment variable.

* `password` - (Required) Matching password for the user to authenticate to TeamCity. It is recommended to be set via `TEAMCITY_PASSWORD` environment variable.

## Example Usage

```hcl
provider "teamcity" {
  # It is strongly recommended to configure this provider through the
  # environment variables described above, so that each user can have
  # separate credentials set in the environment.

  address = "http://127.0.0.1:8112"
}

# Creates a project "Terraformed project" in the Root level hierarchy
resource "teamcity_project" "project" {
  name = "Terraformed project"
}
```
