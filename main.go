// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/kube-cloud/terraform-provider-vsphere/vsphere"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: vsphere.Provider,
	}

	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "kube-cloud/vsphere"
	}

	plugin.Serve(opts)
}
