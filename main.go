package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/foreflight/terraform-provider-teamcity/teamcity"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: teamcity.Provider})
}
