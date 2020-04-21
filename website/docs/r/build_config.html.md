---
subcategory: "Build Configurations"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_build_config"
description: |-
  Manages TeamCity build configurations
---

# teamcity_build_config

The Build Configuration resource allows managing TeamCity build configurations.

## Example Usage

```hcl
resource "teamcity_project" "project" {
  name = "My Project"
}

resource "teamcity_build_config" "build" {
  name        = "BuildRelease"
  description = "Build 'My Project'"
  project_id  = teamcity_project.project.id

  settings {
    build_number_format = "1.2.%build.counter%"
  }

  step {
    type = "cmd_line"
    name = "build_script"
    file = "./build.sh"
    args = "default_target --verbose"
  }
}
```

## Build Configuration Templates

Build Configurations can be managed with the same resource for both regular and templates.

To manage it as a template, specify the `is_template` attribute as `true`.

To associate templates with a given build configuration, use the `templates` list attribute.

```hcl
resource "teamcity_project" "project" {
  name = "My Project"
}

resource "teamcity_build_config" "build" {
  name        = "MainBuildConfiguration"
  description = "Build 'My Project'"
  project_id  = teamcity_project.project.id

  templates = [teamcity_build_config.template1.id, teamcity_build_config.template2.id]
}

resource "teamcity_build_config" "template1" {
  name = "Build Config Template 1"
  # Description is not supported for Build Configuration Templates! https://youtrack.jetbrains.com/issue/TW-63617
  project_id  = teamcity_project.project.id
  is_template = true
}

resource "teamcity_build_config" "template2" {
  name        = "Build Config Template 2"
  project_id  = teamcity_project.project.id
  is_template = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which the build configuration will be have. TeamCity [automatically generates](https://confluence.jetbrains.com/display/TCD18/Identifier) a friendly `ID`  based on name. If duplicate names are found within a same project, TeamCity will append a number to the end of the `ID`. It is better to avoid duplicating build configuration names in the scope of the same project.

* `project_id` - (Required) ID of the project under which this build configuration will be created.

---

* `description`: (Optional) Description for this build configuration.

~> **Note:** Descriptions cannot be specified for Templates - [see this YouTrack issue for more information](https://youtrack.jetbrains.com/issue/TW-63617.)

* `config_params` - (Optional) A map of parameters of type `Configuration Parameters`. Configuration parameters are not passed into build, can be used in references only.

* `env_params` - (Optional) A map of parameters of type `Environment Variables`. Environment variables will be added to the environment of the processes launched by the build runner (without env. prefix).

* `is_template` - (Optional) If true, the build configuration will be managed as a template. Defaults to `false`.

* `settings` - (Optional) One or more `settings` blocks as defined below.

* `step` - (Optional) One or more `step` blocks as defined below, used as Build Steps in the Build Configuration.

* `sys_params` - (Optional) A map of parameters of type `System Properties`. System properties will be passed into the build (without system. prefix), they are only supported by the build runners that understand the property notion.

* `templates` - (Optional) A list of Build Configuration Template IDs to associate to this build configuration.

* `vcs_root` - (Optional) One or more `vcs_root` blocks as defined below, used to manage attaching VCS Roots to this build configuration.

---

The `settings` block supports the following arguments:

* `allow_personal_builds` - (Optional) If true, it allows triggering builds manually from UI in "Run...".

* `artifact_paths` - (Optional) A list of paths in the form of [+:]source [ => target] to include and -:source [ => target] to exclude files or directories to publish as build artifacts. Ant-style wildcards are supported, e.g. use **/* => target_directory, -: **/folder1 => target_directory to publish all files except for folder1 into the target_directory.

* `build_counter` - (Optional) Build Counter. Must be at least `0` (zero). Defaults to `0` (zero).

* `build_number_format` - (Optional) Build Number Format. The format may include '%build.counter%' as a placeholder for the build counter value, for example, `"1.%build.counter%"`.

* `concurrent_limit` - (Optional) Limit the number of simultaneously running builds. Must be at least `0` (zero). Defaults to `0` (zero), which means unlimited.

* `configuration_type` - (Optional) Build Configuration Type. Use `"REGULAR"`, `"DEPLOYMENT"` or `"COMPOSITE"`. Defaults to `"REGULAR"`

* `detect_hanging` - (Optional) If true, enables hanging builds detection. Defaults to `true`.

* `status_widget` - (Optional) If true, enables hanging builds detection. Defaults to `false`.

---

The `step` block supports the following arguments:

* `type` - (Required) Specify `cmd_line` for command line runner or `powershell` for powershell runner.

* `name` - (Optional) A named reference for this step. If not specified, TeamCity will generate it based on runner.

* `file` - (Optional) If calling an external script, this is the file name to run. Do not use this with `code`.

* `code` - (Optional) Inline script code to call. Do not use this with `file`.

* `args` - (Optional) Arguments to pass to external script specified in `file`.

---

The `vcs_root` block supports the following arguments:

* `id` - (Required) The ID of the VCS Root to attach.

* `checkout_rules` - (Optional) A list of strings specifying set of rules in the form of `+|-:VCSPath\[\=\>AgentPath\]`. Used to add/exclude which files should be checked out from this VCS Root

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the build configuration.

## Import

Build Configurations can be imported using their ID, e.g.

```
$ terraform import teamcity_build_config.example MyProject_BuildRelease
```
