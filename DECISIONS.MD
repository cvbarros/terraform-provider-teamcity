# 2018-05-15
When implementing `steps` for build configurations, I decided to make them an _embedded resource_ rather than a separate resource. Example:

```hcl
resource "teamcity_build_configuration" "build_release" {
   // ...

   step {
      type = "powershell"
      name = "build release"
      file = "build.ps1"
      args = "-Target buildrelease -Verbosity detailed"
   }

   step {
      type = "commandline"
      name = "notify"
      file = "notify.cmd"
   }
}
```

#2018-05-22
Build Features will be implemented as a separate resource per "Build Feature Type". So for each kind of feature, there will be a separate resource that can be linked to the build configuration. Example:

```hcl
resource "teamcity_build_configuration" "pull_request" {
   // ...
}

resource "teamcity_feature_commit_status_publisher" "github_commit_update" {
   build_config_id = "${teamcity_build_configuration.pull_request.id}
   publisher = "github"
   
   github {
      auth_type = "password"
      host = "https://api.github.com"
      username = "username"
      password = "${var.github_password}"
   }
}
```