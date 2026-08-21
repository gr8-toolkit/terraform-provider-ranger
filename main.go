package main

import (
  "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
  "github.com/gr8-toolkit/terraform-provider-ranger/ranger"
)

func main() {
  plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: ranger.Provider,
  })
}
