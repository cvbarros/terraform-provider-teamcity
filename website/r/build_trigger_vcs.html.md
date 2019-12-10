---
layout: teamcity
page_title: "TeamCity: teamcity_build_trigger_vcs resource"
sidebar_current: "docs-teamcity-build-trigger-vcs"
description: |-
  Manages TeamCity build configuration "VCS" build triggers.
---

# teamcity\_build\_trigger\_vcs

The Build Trigger VCS resource allows managing build configuration triggers of type "VCS", that will fire builds when VCS changes are detected.

## Example Usage

```hcl
resource "teamcity_project" "project" {
  name = "Go TeamCity SDK"
}

resource "teamcity_vcs_root_git" "vcs" {
  name       = "Application"
  project_id = "${teamcity_project.project.id}"

  url    = "https://github.com/cvbarros/go-teamcity"
  branch = "refs/head/master"
}

resource "teamcity_build_config" "triggered_build" {
  project_id          = teamcity_project.project.id
  name                = "Triggered Build"

  step {
    type = "command_line"
    file = "build.sh"
    args = "-t release"
  }
}

resource "teamcity_build_trigger_vcs" "vcs_trigger" {
    build_config_id = teamcity_build_config.triggered_build.id

    rules = ["-:*.md"]
    branch_filter = ["master"]
}
```

## Argument Reference

The following arguments are supported:

* `build_config_id`: (Required) ID of the build configuration which this trigger will be configured.

* `rules`: (Required) A list of rules: +|-:[Ant-like wildcard] that can make this trigger fire.

* `branch_filter`: (Optional) A list of branches. Only changes in the scoped branches will fire this trigger.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id`: The auto-generated ID of the agent requirement.
