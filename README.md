Terraform Provider TeamCity
==================
[![Build Status](https://travis-ci.org/cvbarros/terraform-provider-teamcity.svg?branch=master)](https://travis-ci.org/cvbarros/terraform-provider-teamcity)

Motivation:

[Building Builds: TeamCity Pipelines as Code using Terraform](https://cvbarros.io/2018/11/building-builds---teamcity-pipelines-as-code-using-terraform/)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x (advised 0.12+)
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/cvbarros/terraform-provider-teamcity`

```sh
$ mkdir -p $GOPATH/src/github.com/cvbarros/terraform-provider-teamcity; cd $GOPATH/src/github.com/cvbarros
$ git clone git@github.com:cvbarros/terraform-provider-teamcity
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/cvbarros/terraform-provider-teamcity
$ go build -o $GOPATH/bin/terraform-provider-teamcity
```

Using the provider
----------------------

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory, run `terraform init` to initialize it.

Example Configurations
----------------------
You can find several sample configurations under `/examples`. As a best effort basis, the samples will be following 
the same terraform version as the provider is built against.  

There is a `.terraform-version` file that indicates the 
minimum terraform version to which the samples are compatible. Use of [tfenv](https://github.com/tfutils/tfenv) is recommended so you can run multiple
versions of terraform locally and switch based on the `.terraform-version` file.  

Please feel free to open an issue if you encounter any problems, or contribute to new sample configurations.

Developing
---------------------------

Please see [CONTRIBUTING](CONTRIBUTING.MD#developing).
