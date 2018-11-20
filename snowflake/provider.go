package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// DefaultSnowFlakeRegion mentions SnowFlake AWS Account Region
const DefaultSnowFlakeRegion = "us-east-1"

type providerConfiguration struct {
	DB            *sql.DB
	ServerVersion *version.Version
}

// Provider blah foo bar
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SF_ACCOUNT", nil),
				Description: "Name of Snowflake Account string to connect to",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "" {
						errors = append(errors, fmt.Errorf("Account must not be an empty string"))
					}
					return
				},
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SF_USER", nil),
				Description: "Snowflake user name to connect as ",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "" {
						errors = append(errors, fmt.Errorf("Username must not be an empty string"))
					}
					return
				},
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password to be used to connect to Snowflake Server",
				DefaultFunc: schema.EnvDefaultFunc("SF_PASSWORD", nil),
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Snowflake AWS region that is configured with account",
				DefaultFunc: schema.EnvDefaultFunc("SF_REGION", DefaultSnowFlakeRegion),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"snowflake_warehouse": resourceWarehouse(),
			"snowflake_database":  resourceDatabase(),
			//"snowflake_user":      resourceUser(),
			//"snowflake_grant":     resourceGrant(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	var username = d.Get("username").(string)
	var password = d.Get("password").(string)
	var account = d.Get("account").(string)
	var region = d.Get("region").(string)

	// database/sql is the thread-safe by default, so we can
	// safely re-use the same handle between multiple parallel
	// operations.

	var dataSourceName string

	if region == "us-west-2" {
		dataSourceName = fmt.Sprintf("%s:%s@%s", username, password, account)
	} else {
		dataSourceName = fmt.Sprintf("%s:%s@%s.%s", username, password, account, region)
	}

	db, err := sql.Open("snowflake", dataSourceName)

	ver, err := serverVersion(db)
	if err != nil {
		return nil, err
	}

	return &providerConfiguration{
		DB:            db,
		ServerVersion: ver,
	}, nil
}

var identQuoteReplacer = strings.NewReplacer("`", "``")

func quoteIdentifier(in string) string {
	return fmt.Sprintf("`%s`", identQuoteReplacer.Replace(in))
}

func serverVersion(db *sql.DB) (*version.Version, error) {
	rows, err := db.Query("SELECT  CURRENT_VERSION()")
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, fmt.Errorf("SELECT  CURRENT_VERSION() returned an empty set")
	}

	var versionString string
	rows.Scan(&versionString)
	return version.NewVersion(versionString)
}
