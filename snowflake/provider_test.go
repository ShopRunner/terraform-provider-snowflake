package snowflake

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

var testSnowflakeProviders map[string]terraform.ResourceProvider
var testSnowflakeProvider *schema.Provider

func init() {
	testSnowflakeProvider = Provider().(*schema.Provider)
	testSnowflakeProviders = map[string]terraform.ResourceProvider{
		"snowflakedb": testSnowflakeProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
