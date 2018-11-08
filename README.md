Terraform Provider TeamCity
==================
[![Build Status](https://travis-ci.org/cvbarros/terraform-provider-teamcity.svg?branch=master)](https://travis-ci.org/cvbarros/terraform-provider-teamcity)

Motivation:
https://cvbarros.io/2018/11/building-builds---teamcity-pipelines-as-code-using-terraform/

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/cvbarros/terraform-provider-teamcity`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/cvbarros
$ git clone git@github.com:cvbarros/terraform-provider-teamcity
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/cvbarros/terraform-provider-teamcity
$ go build -o $GOPATH/bin/terraform-provider-teamcity
```

Using the provider
----------------------

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

Developing
---------------------------

Please see [CONTRIBUTING](CONTRIBUTING.MD#developing).
