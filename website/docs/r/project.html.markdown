---
subcategory: "Projects"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_project"
description: |-
  Manages a Project
---

# teamcity_project

The Project resource allows managing a Projects. It is the base resource needed for provisioning Build Configurations, since they to be associated with a project that is not the `Root` project.

~> **WARNING:** Deleting a project resource will delete everything underneath it.

## Example Usage

```hcl
resource "teamcity_project" "parent" {
  name        = "Parent"
  description = "Parent Project, will be created under the 'Root' project"
}

resource "teamcity_project" "child" {
  name        = "Child"
  description = "Child Project, will be created under 'Parent' project"
  parent_id   = teamcity_project.parent.id

  config_params = {
    variable1 = "config_value1"
  }

  env_params = {
    variable1 = "env_value1"
  }

  sys_params = {
    variable1 = "system_value1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which the project will be created. TeamCity [automatically generates](https://confluence.jetbrains.com/display/TCD18/Identifier) a friendly `ID`  based on the name.

* `description` - (Optional) Description to be show under the project name.

* `parent_id` - (Optional) The ID of the Parent Project in the hierarchy which this project will be nested under. Leave it empty to create a top-level project under the `Root` project.

* `env_params` - (Optional) A map of parameters of type `Environment Variables`. Environment variables will be added to the environment of the processes launched by the build runner (without env. prefix).

* `config_params` - (Optional) A map of parameters of type `Configuration Parameters`. Configuration parameters are not passed into build, can be used in references only.

* `sys_params` - (Optional) A map of parameters of type `System Properties`. System properties will be passed into the build (without system. prefix), they are only supported by the build runners that understand the property notion.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the project.

## Import

Projects can be imported using their ID, e.g.

```
$ terraform import teamcity_project.example Parent_Child
```
