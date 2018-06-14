package main

import (
	"github.com/Scalingo/terraform-provider-encrypted/encrypted"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: encrypted.Provider,
	})
}
