# This configuration sample shows how to create build chains with artifact dependency link
provider "teamcity" {
  address  = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password
}

resource "teamcity_project" "project" {
  name = "Samples - Artifact Dependency Project"
}

resource "teamcity_build_config" "dependency_source_build" {
  project_id  = teamcity_project.project.id
  name        = "Source Build"
  description = "Source build that produces artifacts"
  settings {
    artifact_paths = ["+:*.json => /config/*.json"]
  }
}

resource "teamcity_build_config" "dependency_dependent_build" {
  project_id  = teamcity_project.project.id
  name        = "Dependent Build"
  description = "Build that has artifact dependency on other build"
}

resource "teamcity_artifact_dependency" "artifact_dependency" {
  #Required
  build_config_id = teamcity_build_config.dependency_dependent_build.id
  #Required
  source_build_config_id = teamcity_build_config.dependency_source_build.id

  #Controls from which source build it should get the artifacts
  #It can be "lastSuccessful", "lastPinned", "lastFinished", "sameChainOrLastFinished", "buildNumber" or "buildTag"
  #'buildNumber' uses artifacts produced by the source build that has a specific build number
  #'buildTag' uses the artifacts produced by the source build that have a specific VCS tag
  #Default: "lastSuccessful"
  dependency_revision = "lastFinished"

  #Defines the rules to grab artifacts from the source build's outputs. Required.
  path_rules = ["+:*"]

  #Revision value will mean different things depending on 'dependency_revision' attribute
  #'buildNumber' it means which build number to refer to
  #'buildTag', controls which VCS tag this dependency will look for
  #Otherwise, this attribute is ignored and not set into state
  #Required if using 'buildNumber' or 'buildTag' for 'dependency_revision'
  revision = "latest"

  #Instructs to clean the destination paths (configured with path_rules) before downloading artifacts from the source build
  #Optional
  #Default: false
  clean_destination = true
}
