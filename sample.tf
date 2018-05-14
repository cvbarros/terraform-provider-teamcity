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

  // schema.TypeSet Elem: schema.Resource { Schema: ... }
  trigger {
    //schema.TypeString, validateFunc: validateTriggerType
    type = "vcs"

    //schema.TypeList
    rules = ["+:*"]
  }

  // schema.TypeMap Default
  params {}

  // schema.TypeSet
  feature {}
}
