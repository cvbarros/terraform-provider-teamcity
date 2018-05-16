provider "teamcity" {
  address  = "http://192.168.99.100:8112"
  username = "admin"
  password = "admin"
}

resource "teamcity_project" "canary" {
  name = "Canary"
}

resource "teamcity_vcs_root_git" "canary_vcs" {
  name       = "Application"
  project_id = "${teamcity_project.canary.id}"

  url    = "https://github.com/kelseyhightower/nocode"
  branch = "refs/head/master"
}

resource "teamcity_buildconfiguration" "canary_pullrequest" {
  project_id          = "${teamcity_project.canary.id}"
  name                = "Pull Request"
  description         = "Inspection Build with \"Pull-Request\" hook"
  build_number_format = "2.0.%build.counter%"
  artifact_paths      = [""]

  options {
    status_widget         = false
    detect_hanging        = true
    allow_personal_builds = true
  }

  step {
    type = "powershell"
  }

  vcs_root {
    id             = "${teamcity_vcs_root_git.canary_vcs}"
    checkout_rules = ["+:*"]
  }

  env_params    = {}
  sys_params    = {}
  config_params = {}
  feature       = {}
}

resource "teamcity_buildconfiguration" "canary_build_release" {
  project_id          = "${teamcity_project.canary.id}"
  name                = "Build Release"
  description         = "Master build with \"BuildRelease\" hook"
  build_number_format = "2.0.%build.counter%"
  artifact_paths      = [""]

  options {
    status_widget         = false
    detect_hanging        = true
    allow_personal_builds = true
  }

  step {
    type = "powershell"
    file = "build.ps1"
    args = "-Target buildrelease -Verbosity %verbosity%"
  }

  vcs_root {
    id             = "${teamcity_vcs_root_git.canary_vcs}"
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_buildconfiguration" "canary_release_testing" {
  project_id          = "${teamcity_project.canary.id}"
  name                = "Release To Tesring"
  description         = "Perform a deployment to Testing environment"
  build_number_format = "2.0.%build.counter%"

  step {
    type = "powershell"
    file = "build.ps1"
    args = "-Target release -Verbosity %verbosity%"
  }
}

resource "teamcity_build_trigger" "canary_vcs_trigger" {
  build_config_id = "${teamcity_buildconfiguration.canary_pullrequest}"

  //schema.TypeString, validateFunc: validateTriggerType
  type = "vcs"

  //schema.TypeList
  rules = ["+:*"]
}

resource "teamcity_snapshot_dependency" "canary_release_testing" {
  build_config_id        = "${teamcity_buildconfiguration.canary_release_testing}"
  source_build_config_id = "${teamcity_buildconfiguration.canary_build_release}"
}
