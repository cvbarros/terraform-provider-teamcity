# How to contribute #

We'd more than happy to accept all forms of contribution. Those may come in form of submitting a pull request, opening an issue, improving documentation or tests, or even helping the community that may ask questions about the project on issues. However, in order to keep work organized, it's good to follow some guidelines.

## Reporting issues ##

### Problems and/or bugs ###
Be descriptive about your issue. Specify which TeamCity, Golang and OS version you are running and provide helpful information through logs and/or debug information.

### Feature Requests ###
Feature requests are welcome, still, try to be as specific as possible about what the intended feature should be. Then, after proposing the feature, be open to discussion and patient regarding it's implementation.

## Developing ##
If you want to contribute in form of development, you may fork this repo, clone it locally then create a branch to work on. [This link](https://opensource.guide/how-to-contribute/#opening-a-pull-request) has helpful information to get started.

You'll first need [Go](https://golang.org/doc/install) installed on your machine (version 1.12+ is *required*).

> This project makes use of [Go Modules](https://github.com/golang/go/wiki/Modules). If you clone it inside `$GOPATH`, make sure you have `GO111MODULE=on` before using any `go` commands. The default is `auto`, which will disable modules support inside `$GOPATH`.

To build the provider, run `go build -o terraform-provider-teamcity`. This will generate the provider binary in the current directory.

Seen an issue that you want to tackle? Make sure to comment indicating that intention, so others know that someone is already willing to work on it.

Try to maintain the code conventions, or spot where they are inconsistent across the codebase, whenever possible. When submitting your work, if you would like to propose an enhancement to existing conventions, do it in a separate issue for discussion, other than the issued worked on.

### Running TeamCity Locally ###
TeamCity can be run with the official docker images from JetBrains. We leverage that by implementing full integration tests with every feature.
You are expected to fully run these tests and make sure your changes do not break existing functionality.

In order to run it locally you can first start a warmed-up TeamCity container, from the repository root folder:

```bash
$ ./integration_tests/start_teamcity.sh
```

This will start the needed dependencies to run all tests in your local machine as per [this configuration](./integration_tests/docker-compose.yml). After your environment is running, you may set the `TEAMCITY_ADDR`, `TEAMCITY_USER` and `TEAMCITY_PASSWORD` environment variables to the corresponding values (defaults to `admin`, `admin`, `http://127.0.0.1:8112`, respectively) and `TF_ACC` environment variable to `1` to enable Terraform acceptance tests:

```bash
$ export TEAMCITY_ADDR=http://localhost:8112
$ export TEAMCITY_USER=admin
$ export TEAMCITY_PASSWORD=admin
$ export TF_ACC=1
```

To run tests:

```bash
go test -v -timeout 180s ./...
```

If using an editor such as **Visual Studio Code** with `Go` integration, you can make your development easier by configuring your workspace settings:

```json
{
    "go.testEnvVars": {
        "TEAMCITY_ADDR":"http://localhost:8112",
        "TEAMCITY_USER":"admin",
        "TEAMCITY_PASSWORD":"admin",
        "TF_ACC":"1"
    },

    "go.testFlags": ["-v"],
    "go.buildFlags": ["-v"],
    "go.testTimeout": "180s"
}
```

####Using local [go-teamcity](https://github.com/cvbarros/go-teamcity) code

Add the following line to the go.mod file ```replace github.com/cvbarros/go-teamcity => <local path to repo>``` eg. ```replace github.com/cvbarros/go-teamcity => ~/Development/go-teamcity```

#### Running terraform examples

The terraform examples are in the /examples folder.  To apply an example to a local TeamCity server:
1. Build the terraform-provider-teamcity plugin: ```go build -o examples/<example folder>/terraform-provider-teamcity```
 eg. ```go build -o examples/groups/terraform-provider-teamcity```
2. Run teamcity server: ```./integration_tests/start_teamcity.sh```
3. Run example:
   1. cd examples/groups
   2. source .envrc
   3. ```terraform init```
   4. export TF_LOG=DEBUG
   5. ```terraform apply```
4. To undo changes made with terraform run ```terraform destroy```

If you have followed the steps but still need help to setup a development, help us by opening an issue so we can improve the documentation and the process for setting up the development environment. :+1:
