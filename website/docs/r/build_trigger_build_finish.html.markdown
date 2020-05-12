---
subcategory: "Build Configurations"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_build_trigger_build_finish"
description: |-
  Manages TeamCity build configuration triggers of "Finish Build" type.
---

# teamcity_build_trigger_build_finish

The Build Trigger Build Finish resource allows managing build configuration triggers of type "Finish Build".

## Example Usage

```hcl
resource "teamcity_project" "project" {
  name = "Go TeamCity SDK"
}

resource "teamcity_vcs_root_git" "vcs" {
  name       = "Application"
  project_id = teamcity_project.project.id
  url        = "https://github.com/cvbarros/go-teamcity"
  branch     = "refs/head/master"
}

resource "teamcity_build_config" "build_release" {
  project_id = teamcity_project.project.id
  name       = "Build Release"

  step {
    type = "cmd_line"
    file = "build.sh"
    args = "-t buildrelease"
  }

  vcs_root {
    id             = teamcity_vcs_root_git.vcs.id
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_build_config" "triggered_build" {
  project_id  = teamcity_project.project.id
  name        = "Triggered Build"
  description = "Build triggered when 'Build Release' is finished"

  step {
    type = "command_line"
    file = "build.sh"
    args = "-t release"
  }
}

resource "teamcity_build_trigger_build_finish" "buildrelease_finish_trigger" {
  build_config_id        = teamcity_build_config.build_release.id
  source_build_config_id = teamcity_build_config.triggered_build.id
  after_successful_only  = true
  branch_filter          = ["master", "feature"]
}
```

## Argument Reference

The following arguments are supported:

* `build_config_id` - (Required) ID of the build configuration which this trigger will be configured.

* `source_build_config_id` - (Required) ID of the build configuration that, when finished, will fire this trigger.

---

* `after_successful_only` - (Optional) If true, this trigger will fire only when `source_build_config_id` is successful. Defaults to `false`.

* `branch_filter` - (Optional) A list of branches that scope this trigger. Only finished builds in the given branches will fire the trigger.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id`- The auto-generated ID of the agent requirement.
