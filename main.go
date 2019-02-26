package main

import (
	"github.com/ShopRunner/terraform-provider-snowflake/snowflake"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: snowflake.Provider})
}
