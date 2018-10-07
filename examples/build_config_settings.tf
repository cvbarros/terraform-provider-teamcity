resource "teamcity_project" "nocode" {
  name = "No Code"
}


resource "teamcity_buildconfiguration" "nocode_config" {
  name = "SimpleConfig"
  description = "Configuration to showcase build configuration settings"
  project_id = "${teamcity_project.nocode.id}"

  #These settings appear on "General Settings" for build configurations in TeamCity's UI
  settings {
    #Type of build configuration: "regular" (default), "composite" or "deployment"
    configuration_type = "regular"

    #The format may include '%build.counter%' as a placeholder for the build counter value, for example, 1.%build.counter%.
    #It may also contain a reference to any other available parameter, for example, %build.vcs.number.VCSRootName%.
    #Note: The maximum length of a build number after all substitutions is 256 characters.
    build_number_format = "%build.counter%"

    #Positive int
    build_counter = 1

    #Set to false to disable personal builds. Default: true
    allow_personal_builds = true

    #Paths in the form of [+:]source [ => target] to include and -:source [ => target] to exclude files or directories to publish as build artifacts.
    artifact_paths = ["+:*.json => /config/*.json"]

    #Enable hanging builds detection. Default: true
    detect_hanging = true

    #Enable build status to be queried externally. Default: false
    status_widget = false

    #Int 0->unlimited. Defaults to '0', which means unlimited.
    concurrent_limit = 10
  }
}

