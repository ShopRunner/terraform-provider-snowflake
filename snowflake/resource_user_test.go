package snowflake

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestUserSnowflakeDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_user", "name", "shoprunner_terraform"),
					resource.TestCheckResourceAttr("snowflake_user", "user", "tf-test"),
				),
			},
		},
	})
}

var testSnowflakeUserConfig = `
resource "snowflake_user" "shoprunner_terraform" {
  user = "tf-test"
}
`
