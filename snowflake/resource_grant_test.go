package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestDatabaseSnowflakeGrant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeGrantConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_grant", "user", "shoprunner_terraform_user"),
					resource.TestCheckResourceAttr(
						"snowflake_grant", "host", "192.168.0.1"),
					resource.TestCheckResourceAttr(
						"snowflake_grant", "database", "db"),
					resource.TestCheckResourceAttrSet(
						"snowflake_grant", "privileges"),
					resource.TestCheckResourceAttr(
						"snowflake_grant", "grant", "false"),
				),
			},
		},
	})
}

var testSnowflakeGrantConfig = `
resource "snowflake_grant" "shoprunner_grant_terraform" {
      user       =   "shoprunner_terraform_user"
	  host    	 =   "192.168.0.1"
	  database   =   "db"
	  privileges =   ["all"]
}
`
