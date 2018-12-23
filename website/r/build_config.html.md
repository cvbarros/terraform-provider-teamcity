---
layout: teamcity
page_title: "TeamCity: teamcity_build_config resource"
sidebar_current: "docs-teamcity-build-config"
description: |-
  Manages TeamCity build configurations
---

# teamcity\_build\_config

The Build Configuration resource allows managing TeamCity build configurations.

## Example Usage

```hcl
resource "teamcity_project" "project" {
  name = "My Project"
}

resource "teamcity_build_config" "build" {
  name = "BuildRelease"
  description = "Build 'My Project'"
  project_id = "${teamcity_project.project.id}"

  settings {
    build_number_format = "1.2.%build.counter%"
  }

	step {
		type = "cmd_line"
		name = "build_script"
		file = "./build.sh"
		args = "default_target --verbose"
	}

  step {
    type = octopus.push.package
    host          = "https://prod.octoserver.com"
    api_key       = "API-ABCD123"
    package_paths = "*"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name`: (Required) Specifies the name which the build configuration will be have. TeamCity [automatically generates](https://confluence.jetbrains.com/display/TCD18/Identifier) a friendly `ID`  based on name. If duplicate names are found within a same project, TeamCity will append a number to the end of the `ID`. It is better to avoid duplicating build configuration names in the scope of the same project.

* `description`: (Optional) Description for this build configuration.

* `project_id`: (Required) ID of the project under which this build configuration will be created.

* `env_params`: (Optional) A map of parameters of type `Environment Variables`. Environment variables will be added to the environment of the processes launched by the build runner (without env. prefix).

* `config_params`: (Optional) A map of parameters of type `Configuration Parameters`. Configuration parameters are not passed into build, can be used in references only.

* `sys_params`: (Optional) A map of parameters of type `System Properties`. System properties will be passed into the build (without system. prefix), they are only supported by the build runners that understand the property notion.

The `vcs_root` block is used to manage attaching VCS Roots to this build configuration. For every attachment, use a `vcs_root` block to configure it:

* `id`: (Required) ID of the VCS Root to attach.

* `checkout_rules`: (Optional) A list of strings specifying set of rules in the form of +|-:VCSPath\[\=\>AgentPath\]. Used to add/exclude which files should be checked out from this VCS Root

The `step` block supports the attributes below. You may use one step block per step you want the build configuration to perform. The order configured will be the order they will be performed.

* `type`: (Required) Specify `cmd_line` for command line runner, `powershell` for powershell runner or `octopus.push.package` for the Octopus Push Package runner.

* `name`: (Optional) A named reference for this step. If not specified, TeamCity will generate it based on runner.

* `file`: (Optional) If calling an external script, this is the file name to run. Do not use this with `code`.

* `code`: (Optional) Inline script code to call. Do not use this with `file`.

* `args`: (Optional) Arguments to pass to external script specified in `file`.

If you selected the `octopus.push.package` runner, here are the available options:
* `host`: (Required) Octopus web portal URL.
* `api_key`: (Required) Octopus API key.
* `package_paths`: (Required) Package path patterns.
* `force_push`: (Optional).  Defaults to `true`.
* `publish_artifacts`: (Optional). Defaults to `true`.
* `additional_command_line_arguments`: (Optional). Defaults to `''` (none).

The `settings` block supports:

* `configuration_type`: (Optional) Build Configuration Type. Use `"REGULAR"`, `"DEPLOYMENT"` or `"COMPOSITE"`. Defaults to `"REGULAR"`

* `build_number_format`: (Optional) Build Number Format. The format may include '%build.counter%' as a placeholder for the build counter value, for example, `"1.%build.counter%"`.

* `build_counter`: (Optional) Build Counter. Must be at least `0` (zero). Defaults to `0` (zero).

* `allow_personal_builds`: (Optional) If true, it allows triggering builds manually from UI in "Run...".

* `artifact_paths`: (Optional) A list of paths in the form of [+:]source [ => target] to include and -:source [ => target] to exclude files or directories to publish as build artifacts. Ant-style wildcards are supported, e.g. use **/* => target_directory, -: **/folder1 => target_directory to publish all files except for folder1 into the target_directory.

* `detect_hanging`: (Optional) If true, enables hanging builds detection. Defaults to `true`.

* `status_widget`: (Optional) If true, enables hanging builds detection. Defaults to `false`.

* `concurrent_limit`: (Optional) Limit the number of simultaneously running builds. Must be at least `0` (zero). Defaults to `0` (zero), which means unlimited.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id`: The auto-generated ID of the build configuration.

## Import
Build Configurations can be imported using their ID, e.g.

```
$ terraform import teamcity_build_config.teamcity_build_config MyProject_BuildRelease
```
