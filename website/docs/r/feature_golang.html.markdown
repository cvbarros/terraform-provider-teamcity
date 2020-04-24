---
subcategory: "Build Configurations"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_feature_golang"
description: |-
  Manages an Golang Build Feature for a Build Configuration
---

# teamcity_feature_golang

Manages an Golang Build Feature for a Build Configuration

## Example Usage

```hcl
resource "teamcity_project" "example" {
  name = "Example Project"
}

resource "teamcity_build_config" "example" {
  name        = "Example Build"
  project_id  = teamcity_project.example.id
}

resource "teamcity_feature_golang" "example" {
  build_config_id = teamcity_build_config.example.id
}
```

## Argument Reference

The following arguments are supported:

* `build_config_id` - (Required) Specifies the ID of the Build Configuration for which a Golang Build Feature should be configured.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the Golang Build Feature.

## Import

Golang Build Features can be imported using their ID, e.g.

```
$ terraform import teamcity_feature_golang.example "ProjectID|golang"
```

-> **Note:** This is a Terraform specific ID comprised of "ProjectID|FeatureID" - where featureID is likely `golang`
