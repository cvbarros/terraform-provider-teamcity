---
subcategory: "Build Configurations"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_snapshot_dependency"
description: |-
  Manages TeamCity snapshot dependencies
---

# teamcity_snapshot_dependency

The Snapshot Dependency resource allows managing build dependencies of "Snapshot" type.

## Example Usage

```hcl
resource "teamcity_project" "project" {
  name = "Project"
}

resource "teamcity_build_config" "source" {
  name       = "source"
  project_id = teamcity_project.project.id
}

resource "teamcity_build_config" "dependant" {
  name       = "dependant"
  project_id = teamcity_project.project.id
}

resource "teamcity_snapshot_dependency" "dependency" {
  source_build_config_id = teamcity_build_config.source.id
  build_config_id        = teamcity_build_config.dependant.id
}
```

## Argument Reference

The following arguments are supported:

* `build_config_id` - (Required) The ID of build configuration this dependency will be created.

* `source_build_config_id` - (Required) The ID of build configuration this dependency relates to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the dependency.
