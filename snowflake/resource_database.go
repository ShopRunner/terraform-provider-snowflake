package snowflake

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const unknownDatabaseErrCode = 1049
const (
	dbNameAttr    = "name"
	dbCommentAttr = "comment"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: createDatabase,
		Update: updateDatabase,
		Read:   readDatabase,
		Delete: deleteDatabase,

		Schema: map[string]*schema.Schema{
			dbNameAttr: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Identifier for the Snowflake database ",
			},
			dbCommentAttr: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    false,
				Description: "Specifies a comment for the database.",
			},
		},
	}
}

func createDatabase(d *schema.ResourceData, meta interface{}) error {
	dbName := d.Get(whNameAttr).(string)
	db := meta.(*providerConfiguration).DB
	b := bytes.NewBufferString("CREATE  DATABASE IF NOT EXISTS ")
	fmt.Fprint(b, dbName)
	fmt.Fprintf(b, " ")
	// Wrap string values in quotes
	for _, attr := range []string{whCommentAttr} {
		fmt.Fprintf(b, " %s='%v' ", attr, d.Get(attr))
	}

	sql := b.String()
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error creating database sql(%s) \n %q: {{err}}", sql, dbName), err)
	}
	d.SetId(dbName)
	return readDatabase(d, meta)
}

func updateDatabase(d *schema.ResourceData, meta interface{}) error {
	dbName := d.Get(whNameAttr).(string)
	db := meta.(*providerConfiguration).DB
	b := bytes.NewBufferString("ALTER DATABASE IF EXISTS ")
	fmt.Fprint(b, dbName)
	fmt.Fprintf(b, " SET ")
	// Wrap string values in quotes
	for _, attr := range []string{dbCommentAttr} {
		fmt.Fprintf(b, " %s='%v' ", attr, d.Get(attr))
	}

	sql := b.String()
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error altering database %q: {{err}}", dbName), err)
	}
	d.SetId(dbName)
	return readDatabase(d, meta)
}

func readDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	databaseName := d.Id()
	stmtSQL := fmt.Sprintf("show databases like '%s'", databaseName)

	fmt.Printf(" Read Database Executing query: %s \n", stmtSQL)
	log.Println("Executing query:", stmtSQL)

	var createdOn, name, isDefault, isCurrent, origin, owner, comment, options, retentionTime sql.NullString

	err := db.QueryRow(stmtSQL).Scan(
		&createdOn, &name, &isDefault, &isCurrent, &origin, &owner, &comment, &options, &retentionTime,
	)

	if err != nil {
		return fmt.Errorf("Error during show databases like: %s", err)
	}

	d.Set(dbNameAttr, name)
	d.Set(dbCommentAttr, comment)
	return nil
}

func deleteDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	dbName := d.Get(whNameAttr).(string)
	sql := fmt.Sprintf("DROP DATABASE  %s ", dbName)
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error dropping database %q: {{err}}", dbName), err)
	}
	return nil
}
