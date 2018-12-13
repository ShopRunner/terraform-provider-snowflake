package snowflake

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: CreateUser,
		Update: UpdateUser,
		Read:   ReadUser,
		Delete: DeleteUser,

		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "localhost",
			},

			"plaintext_password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				StateFunc: hashSum,
			},
			"password": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"plaintext_password"},
				Sensitive:     true,
				Deprecated:    "Please use plaintext_password instead",
			},
			"default_role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NULL",
			},
		},
	}
}

func CreateUser(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("CREATE USER \"%s@%s\"",
		d.Get("user").(string),
		d.Get("host").(string))

	var password string
	if v, ok := d.GetOk("plaintext_password"); ok {
		password = v.(string)
	} else {
		password = d.Get("password").(string)
	}

	if password != "" {
		stmtSQL = stmtSQL + fmt.Sprintf(" PASSWORD = \"%s\"", password)
	}

	stmtSQL = stmtSQL + fmt.Sprintf(" DEFAULT_ROLE = \"%s\"", d.Get("default_role").(string))

	log.Println("Executing statement:", stmtSQL)
	_, err := db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	user := fmt.Sprintf("%s@%s", d.Get("user").(string), d.Get("host").(string))
	d.SetId(user)

	return nil
}

func UpdateUser(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(*providerConfiguration)

	var newpw interface{}
	if d.HasChange("plaintext_password") {
		_, newpw = d.GetChange("plaintext_password")
	} else if d.HasChange("password") {
		_, newpw = d.GetChange("password")
	} else {
		newpw = nil
	}

	var newdefrole interface{}
	if d.HasChange("default_role") {
		_, newdefrole = d.GetChange("default_role")
	} else {
		newdefrole = nil
	}

	if newpw != nil || newdefrole != nil {
		stmtSQL := fmt.Sprintf("ALTER USER \"%s@%s\" SET ",
			d.Get("user").(string),
			d.Get("host").(string))

		if newpw != nil {
			stmtSQL = stmtSQL + fmt.Sprintf(" PASSWORD = \"%s\"",
				newpw.(string))
		}

		if newdefrole != nil {
			stmtSQL = stmtSQL + fmt.Sprintf(" DEFAULT_ROLE = \"%s\"",
				newdefrole.(string))
		}

		log.Println("Executing query:", stmtSQL)
		_, err := conf.DB.Exec(stmtSQL)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("SHOW USERS LIKE '%s@%s'",
		d.Get("user").(string),
		d.Get("host").(string))

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() && rows.Err() == nil {
		d.SetId("")
	}
	return rows.Err()
}

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	stmtSQL := fmt.Sprintf("DROP USER \"%s@%s\"",
		d.Get("user").(string),
		d.Get("host").(string))

	log.Println("Executing statement:", stmtSQL)

	_, err := db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}
