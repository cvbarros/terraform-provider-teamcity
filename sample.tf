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

  // How to set build counter?
  //build_counter = 1
  // 
  artifact_paths = [""]

  // schema.Resource, Computed: true
  options {
    status_widget         = false
    detect_hanging        = true
    allow_personal_builds = true
  }

  //schema.TypeSet
  step {
    type = "powershell"
  }

  // schema.TypeMap Default
  params {}

  // schema.TypeSet
  feature {}
}

resource "teamcity_build_trigger" "canary_vcs_trigger" {
  build_config_id = "${teamcity_buildconfiguration.canary_pullrequest}"

  //schema.TypeString, validateFunc: validateTriggerType
  type = "vcs"

  //schema.TypeList
  rules = ["+:*"]
}
