package main

import (
	"flag"
	"github.com/dmacvicar/terraform-provider-libvirt/libvirt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	defer libvirt.CleanupLibvirtConnections()
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: "registry.terraform.io/dmacvicar/libvirt",
		ProviderFunc: libvirt.Provider,
	}

	plugin.Serve(opts)
}
