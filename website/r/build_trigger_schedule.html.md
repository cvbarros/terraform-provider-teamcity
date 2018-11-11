---
layout: teamcity
page_title: "TeamCity: teamcity_build_trigger_schedule resource"
sidebar_current: "docs-teamcity-build-trigger-schedule"
description: |-
  Manages TeamCity build configuration "Schedule" build triggers.
---

# teamcity\_build\_trigger\_schedule

The Build Trigger Schedule resource allows managing build configuration scheduled triggers.

## Example Usage

```hcl
resource "teamcity_project" "project" {
  name = "Go TeamCity SDK"
}

resource "teamcity_vcs_root_git" "project_vcs" {
  name       = "Application"
  project_id = "${teamcity_project.project.id}"

  url    = "https://github.com/cvbarros/go-teamcity-sdk"
  branch = "refs/head/master"
}

resource "teamcity_buildconfiguration" "triggered_build" {
  project_id          = "${teamcity_project.project.id}"
  name                = "Triggered Build"
  description         = "Build triggered on schedules"
  build_number_format = "0.0.%build.counter%"
  artifact_paths      = [""]

  step {
    type = "command_line"
    file = "build.sh"
    args = "-t buildrelease"
  }

  vcs_root {
    id             = "${teamcity_vcs_root_git.project_vcs}"
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_build_trigger_schedule" "schedule_trigger" {
  build_config_id = "${teamcity_buildconfiguration.triggered_build.id}"

  schedule = "daily"
  hour = 12
  minute = 37
}
```

## Argument Reference

The following arguments are supported:

* `build_config_id`: (Required) ID of the build configuration which this trigger will be configured.

* `schedule`: (Required) `daily` to fire once a day, or `weekly` to fire once a week. `cron` not supported yet.

* `hour`: (Required) Hour at which the trigger will fire.

* `minute`: (Optional) Minute at which the trigger will fire. Defaults to `0 (zero)`, which will be at full hour.

* `timezone`: (Optional) Which timezone the time configured corresponds to. TeamCity use `Quartz` library for Java to specify timezones, so a any string that is a valid timezone ID for the platform can be used. There is no validation on this value, only when trying to create the resource on the TeamCity API. Uses `SERVER` by default which is server-configured timezone.

* `weekday`: (Optional) When using `weekly`  trigger, specifies the full english day o the week which the trigger will fire, e.g: `"Monday"`, `"Wednesday"`, `"Sunday"`.

* `with_pending_changes_only`: (Optional) If true, when this trigger will only fire if the build has VCS pending changes. Defaults to `false`.

* `rules`: (Optional) A list of rules that specifies the changes to be considered. Only used if `with_pending_changes_only` is set.

* `queue_optimization`: (Optional) If true, a queued build can be replaced with an already started build or more recent one. Defaults to `true`.

* `enforce_clean_checkout`: (Optional) If true, all files in the checkout directory will be deleted before the build. Defaults to `false`.

* `enforce_clean_checkout_dependencies`: (Optional) If true, server will peform a clean checkout also for dependencies. Defaults to `false`.
 
* `on_all_compatible_agents`: (Optional) If true, when this trigger fires, the build will be ran on all compatible agents. Defaults to `false`.

* `with_pending_changes_only`: (Optional) If true, when this trigger fires, the build will be ran on all compatible agents. Defaults to `false`.

* `only_if_watched_changes`: (Optional) If set, this trigger will only fire if the source build has changed since the last trigger. If set, `watched_build_config_id` must be set. Defaults to `false`. 

* `watched_build_config_id`: (Optional) Determines the ID of build configuration to watch.

* `revision`: (Optional) Configures which revision to consider from the artifacts produced by the source build. `lastSuccessful` uses artifacts produced by the last successful build. `lastPinned`, artifacts from the last pinned build for the source build configuration. `lastFinished` collects artifacts from the last finished build, successful or not. `buildTag` uses the VCS Tag specified in `watched_branch`.

* `watched_branch`: (Optional) Specifies which tag to use when `revision` is `buildTag`.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id`: The auto-generated ID of the agent requirement.
