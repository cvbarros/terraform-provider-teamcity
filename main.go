package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/yext/terraform-provider-teamcity/teamcity"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: teamcity.Provider})
}
