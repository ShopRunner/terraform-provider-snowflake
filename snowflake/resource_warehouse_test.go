package snowflake

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestWarehouseSnowflakeDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeWarehouseConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_warehouse", "name", "shoprunner_terraform"),
					resource.TestCheckResourceAttr("snowflake_warehouse", "warehouse_size", "SMALL"),
				),
			},
		},
	})
}

var testSnowflakeWarehouseConfig = `
resource "snowflake_warehouse" "shoprunner_warehouse_terraform" {
      name              =   "shoprunner_terraform"
      warehouse_size    =   "SMALL"
}
`
