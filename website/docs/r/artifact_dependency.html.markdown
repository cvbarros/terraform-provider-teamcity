---
subcategory: "Artifact Dependencies"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_artifact_dependency"
description: |-
  Manages TeamCity artifact dependencies
---

# teamcity_artifact_dependency

The Artifact Dependency resource allows managing build dependencies of "Artifact" type.

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

resource "teamcity_artifact_dependency" "dependency" {
  source_build_config_id = teamcity_build_config.source.id
  build_config_id        = teamcity_build_config.dependant.id
  path_rules             = ["+:**/* => target_dir", "-:**/folder1 => target_dir"]
}
```

## Argument Reference

The following arguments are supported:

* `source_build_config_id` - (Required) The ID of build configuration this dependency relates to.

* `build_config_id` - (Required) The ID of build configuration this dependency will be created.

* `dependency_revision` - (Optional) Configures which revision to consider from the artifact produced by the source build. `lastSuccessful` uses artifacts produced by the last successful build. `lastPinned`, artifacts from the last pinned build for the source build configuration. `lastFinished` collects artifacts from the last finished build, successful or not. `sameChainOrLastFinished` uses artifacts produced by source build that was triggered within the same build chain. `buildNumber` uses artifacts from the source build with specific build number. `buildTag` is the same as `buildNumber`, but considers VCS Tag instead.

* `revision` - (Optional) If using `buildNumber`, this is the parameter for which specific build number to consider. In case of `buildTag`, this refers to the tag name. Required in these cases.

* `path_rules` - (Required) A list of rules to match files that will have to be dowloaded from the source build that output artifacts. They can be specified in the format [+:|-:]SourcePath[!ArchivePath][=>DestinationPath].

* `clean_destination` - (Optional) If true, this will clean destination paths before downloading artifacts.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the dependency.
