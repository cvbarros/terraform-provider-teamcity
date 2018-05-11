provider "teamcity" {
  address = "192."
}

resource "teamcity_project" "canary" {}

resource "teamcity_vcs_root_git" "canary_vcs" {
  name       = "Application"
  project_id = "${teamcity_project.canary.id}"

  url    = "https://github.com/kelseyhightower/nocode"
  branch = "refs/head/master"
}

resource "teamcity_buildconfiguration" "canary_pullrequest" {
  project_id = "${teamcity_project.canary.id}"
}
