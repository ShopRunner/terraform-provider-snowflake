package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestDatabaseSnowflakeDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeDatabaseConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_database", "name", "shoprunner_terraform_db"),
					resource.TestCheckResourceAttr("snowflake_database", "comment", "A test comment"),
				),
			},
		},
	})
}

var testSnowflakeDatabaseConfig = `
resource "snowflake_database" "shoprunner_database_terraform" {
      name       =   "shoprunner_terraform_db"
      comment    =   "A test comment"
}
`
