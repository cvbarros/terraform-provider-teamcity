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