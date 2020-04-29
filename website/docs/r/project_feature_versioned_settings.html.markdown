---
subcategory: "Projects"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_project_feature_versioned_settings"
description: |-
  Manages the Versioned Settings for a Project
---

# teamcity_project_feature_versioned_settings

Manages the Versioned Settings for a Project - which allows one/two-way syncing of Build Configurations to/from Source Control.

## Example Usage

```hcl
resource "teamcity_project" "example" {
  name = "Example"
}

resource "teamcity_vcs_root_git" "example" {
  name          = "application"
  project_id     = teamcity_project.example.id
  fetch_url      = "https://github.com/cvbarros/terraform-provider-teamcity"
  default_branch = "refs/head/master"
}

resource "teamcity_project_feature_versioned_settings" "example" {
  project_id     = teamcity_project.example.id
  vcs_root_id    = teamcity_vcs_root_git.example.id
  build_settings = "PREFER_VCS"
  format         = "kotlin"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) Specifies the Identifier for the Project where the Versioned Settings should be configured. Changing this forces a new resource to be created.

* `vcs_root_id` - (Required) Specifies the ID of the VCS Root which should be used to push/pull Versioned Settings from.

* `build_settings` - (Required) Specifies the Build Setting. Possible values are `ALWAYS_USE_CURRENT` (All the builds use current project settings from the TeamCity server. Settings changes in branches, history and personal builds are ignored.), `PREFER_CURRENT` (Builds use current project settings from the TeamCity server. Users can run a build with settings from VCS via the run custom build dialog.), `PREFER_VCS` (Builds in branches and history builds use settings from corresponding branch and revision in VCS. Developers can also change settings in **personal** builds.),

* `format` - (Required) Specifies the format used for the Versioned Settings. Possible values are `kotlin` and `xml`.

---

* `enabled` - (Optional) Should Versioned Settings Synchronization be enabled? Defaults to `true`.

* `show_changes` - (Optional) Should settings changes be shown in Builds? Defaults to `false`.

* `use_relative_ids` - (Optional) Should TeamCity generate Portable DSL Scripts? Defaults to `false`.

-> **Note:** `use_relative_ids` is only applicable when `format` is set to `kotlin`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the Versioned Settings Feature.

## Import

Project Versioned Settings can be imported using the ID of the Project, e.g.

```
$ terraform import teamcity_project_versioned_settings.example ProjectID
```
