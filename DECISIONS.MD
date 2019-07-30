# 2019-07-30
Project has been converted to use Go Modules. This follows on other major providers and Terraform project itself.

# 2018-09-22
Breaking change on `teamcity_build_trigger` resource. I've decided to create separate resource types for different types of build triggers, since the configuration complexity of these resources will greatly increase if kept on a single "build trigger" resource type.

Thus, three new resources will be created:
- `teamcity_build_trigger_vcs`: Corresponding to Build Triggers of type `VCS`
- `teamcity_build_trigger_build_finish`: Corresponding to Build Triggers of type `Finish Build`
- `teamcity_build_trigger_schedule`: Corresponding to Build Triggers of type `Schedule`


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

# 2018-05-22
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
