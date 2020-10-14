Terraform Provider TeamCity
==================
[![Build Status](https://travis-ci.com/cvbarros/terraform-provider-teamcity.svg?branch=master)](https://travis-ci.com/cvbarros/terraform-provider-teamcity)

Motivation:

[Building Builds: TeamCity Pipelines as Code using Terraform](https://cvbarros.io/2018/11/building-builds---teamcity-pipelines-as-code-using-terraform/)

Installation
------------
## Terraform Registry - Terraform 0.13+ (Recommended)
Teamcity Provider [is listed on the Terraform Registry](https://registry.terraform.io/providers/cvbarros/teamcity/latest).  
The easiest to get started way is to follow [instructions for requiring
providers](https://www.terraform.io/docs/configuration/provider-requirements.html#requiring-providers).  

Example:
```hcl
terraform {
  required_providers {
    teamcity = {
      source  = "cvbarros/teamcity"
      version = ">= 1.0.1"
    }
  }
}

provider "teamcity" {
  # Configuration options
}
```

Once configured, running `terraform init` will download the provider automatically from the registry.

## Binaries (Recommended for Terraform up to 0.12.x)
Another way to install this provider is to use one of the binary distributions available as 
[Releases](https://github.com/cvbarros/terraform-provider-teamcity/releases). 
It contains pre-compiled bundles for Linux, MacOS and Windows.

Download and unpack the version for your OS/arch.  
In the example below, we use MacOS (darwin_amd64) and version `v1.0.1`:

```bash
VERSION=v1.0.1
wget https://github.com/cvbarros/terraform-provider-teamcity/releases/download/$VERSION/terraform-provider-teamcity_darwin_amd64_$VERSION.zip
tar -xvf terraform-provider-teamcity_darwin_amd64_$VERSION.zip 
```

Then, copy the output file to your `~/.terraform.d/plugins` directory. 
For Windows, use the directory: `%APPDATA%\terraform.d\plugins`

> **Note**: If you never installed any terraform providers before, you'll have to create that directory.

```bash
cp terraform-provider-teamcity_$VERSION ~/.terraform.d/plugins/
``` 

## Build from Source

### Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x (advised 0.13+)
-	[Go](https://golang.org/doc/install) 1.14+ (to build the provider plugin)


Clone the repository to the directory of your choice, then from the root of repository, run the target below.
It is advisable to checkout a tag corresponding to a version first, instead of `master`.

If the [plugins directory](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) doesn't exist, you have to create it first.

```bash
$ VERSION=v1.0.1 make install
```
This will build a binary for your platform and place it on the plugins directory. 

### Documentation

Documentation on available on Terraform Registry:
https://registry.terraform.io/providers/cvbarros/teamcity/latest/docs
 
Contributions to documentation are invaluable!  
If you'd like to help keep it error-free and up to date, please see [CONTRIBUTING](CONTRIBUTING.MD).

Example Configurations
----------------------
You can find several sample configurations under `/examples`.  
As a best effort basis, the samples will be following 
the same terraform version as the provider is built against.  

There is a `.terraform-version` file that indicates the 
minimum terraform version to which the samples are compatible. Use of [tfenv](https://github.com/tfutils/tfenv) is recommended, so you can run multiple
versions of terraform locally and switch based on the `.terraform-version` file.  

Please feel free to open an issue if you encounter any problems, or contribute to new sample configurations.

Developing
---------------------------

Please see [CONTRIBUTING](CONTRIBUTING.MD#developing).
