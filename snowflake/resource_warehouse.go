package snowflake

import (
	"database/sql"
	"fmt"
	"log"

	"bytes"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	whNameAttr                            = "name"
	whMaxConcurrencyLevelAttr             = "max_concurrency_level"
	whStatementQueuedTimeOutInSecondsAttr = "statement_queued_timeout_in_seconds"
	whStatementTimeoutInSecondsAttr       = "statement_timeout_in_seconds"
	whCommentAttr                         = "comment"

	//Properties
	whSizeAttr           = "warehouse_size"
	whMaxClusterCount    = "max_cluster_count"
	whMinClusterCount    = "min_cluster_count"
	whAutoSuspend        = "auto_suspend"
	whAutoResume         = "auto_resume"
	whInitiallySuspended = "initially_suspended"
	whResourceMonitor    = "resource_monitor"
	//whUUIDAttr								= "UUID"
)

func resourceWarehouse() *schema.Resource {
	return &schema.Resource{
		Create: createWarehouse,
		Update: updateWarehouse,
		Read:   readWarehouse,
		Delete: deleteWarehouse,

		Schema: map[string]*schema.Schema{
			whNameAttr: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Identifier for the Snowflake warehouse;must be unique for your account ",
			},
			whMaxConcurrencyLevelAttr: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				ForceNew:    false,
				Description: "Specifies the maximum number of SQL statements (queries, DDL, DML, etc.) a warehouse cluster can execute concurrently.  ",
			},
			whStatementQueuedTimeOutInSecondsAttr: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     120,
				ForceNew:    false,
				Description: "Specifies the time, in seconds, a SQL statement (query, DDL, DML, etc.) can be queued on a warehouse before it is canceled by the system.",
			},
			whStatementTimeoutInSecondsAttr: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				ForceNew:    false,
				Description: "Specifies the time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system.",
			},
			whCommentAttr: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    false,
				Description: "Specifies a comment for the warehouse.",
			},
			whSizeAttr: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "XSMALL",
				ForceNew:    false,
				Description: "Specifies the size of virtual warehouse to create.",
			},
			whMaxClusterCount: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				ForceNew:    false,
				Description: "Specifies the maximum number of server clusters for the warehouse.",
			},
			whMinClusterCount: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				ForceNew:    false,
				Description: "Specifies the minimum number of server clusters for a multi-cluster warehouse. ",
			},
			whAutoSuspend: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				ForceNew:    false,
				Description: "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
			},
			whAutoResume: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    false,
				Description: "Specifies whether to automatically resume a warehouse when it is accessed by a SQL statement, ",
			},
			whInitiallySuspended: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    false,
				Description: "Specifies whether the warehouse is created initially in suspended state.",
			},
			whResourceMonitor: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    false,
				Description: "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
			},
		},
	}
}

func createWarehouse(d *schema.ResourceData, meta interface{}) error {
	whName := d.Get(whNameAttr).(string)
	db := meta.(*providerConfiguration).DB
	b := bytes.NewBufferString("CREATE  WAREHOUSE IF NOT EXISTS ")
	fmt.Fprint(b, whName)
	fmt.Fprintf(b, " WITH ")
	for _, attr := range []string{whMaxClusterCount, whMinClusterCount, whAutoSuspend, whAutoResume, whInitiallySuspended} {
		fmt.Fprintf(b, " %s=%v ", attr, d.Get(attr))
	}
	// Wrap string values in quotes
	for _, attr := range []string{whSizeAttr, whCommentAttr} {
		fmt.Fprintf(b, " %s='%v' ", attr, d.Get(attr))
	}

	sql := b.String()
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error creating warehouse sql(%s) \n %q: {{err}}", sql, whName), err)
	}
	d.SetId(whName)
	return readWarehouse(d, meta)
}

func updateWarehouse(d *schema.ResourceData, meta interface{}) error {
	whName := d.Get(whNameAttr).(string)
	db := meta.(*providerConfiguration).DB
	b := bytes.NewBufferString("ALTER WAREHOUSE IF EXISTS ")
	fmt.Fprint(b, whName)
	fmt.Fprintf(b, " SET ")
	for _, attr := range []string{whMaxClusterCount, whMinClusterCount, whAutoSuspend, whAutoResume} {
		fmt.Fprintf(b, " %s=%v ", attr, d.Get(attr))
	}
	// Wrap string values in quotes
	for _, attr := range []string{whSizeAttr, whCommentAttr} {
		fmt.Fprintf(b, " %s='%v' ", attr, d.Get(attr))
	}

	sql := b.String()
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error Altering warehouse %q: {{err}}", whName), err)
	}
	d.SetId(whName)
	return readWarehouse(d, meta)
}

func readWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB

	warehouseName := d.Id()
	stmtSQL := fmt.Sprintf("show warehouses like '%s'", warehouseName)

	fmt.Printf(" Read Warehouse Executing query: %s \n", stmtSQL)
	log.Println("Executing query:", stmtSQL)
	var name, state, instanceType, size, minClusterCount, maxClusterCount, startedClusters, running, queued sql.NullString
	var isDefault, isCurrent, autoSuspend, autoResume, available, provisioning, quiescing, other sql.NullString
	var createdOn, resumedOn, updatedOn, owner, comment, resourceMonitor sql.NullString
	var actives, pendings, failed, suspended, uuid, scalingPolicy sql.NullString

	err := db.QueryRow(stmtSQL).Scan(
		&name, &state, &instanceType, &size, &minClusterCount, &maxClusterCount, &startedClusters, &running, &queued,
		&isDefault, &isCurrent, &autoSuspend, &autoResume, &available, &provisioning, &quiescing, &other,
		&createdOn, &resumedOn, &updatedOn, &owner, &comment, &resourceMonitor,
		&actives, &pendings, &failed, &suspended, &uuid, &scalingPolicy,
	)
	if err != nil {
		//if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		//	if mysqlErr.Number == unknownDatabaseErrCode {
		//		d.SetId("")
		//		return nil
		//	}
		//}
		return fmt.Errorf("Error during show create warehouse: %s", err)
	}

	d.Set(whNameAttr, name)
	d.Set("state", state)
	d.Set("instance", instanceType)
	d.Set(whSizeAttr, size)
	d.Set("min_cluster_count", minClusterCount)
	d.Set("max_cluster_count", maxClusterCount)
	d.Set("started_clusters", startedClusters)
	d.Set("running", running)
	d.Set("queued", queued)
	d.Set("is_default", isDefault)
	d.Set("is_current", isCurrent)
	d.Set("auto_suspend", autoSuspend)
	d.Set("auto_resume", autoResume)
	d.Set("available", available)
	d.Set("provisioning", provisioning)
	d.Set("quiescing", quiescing)
	d.Set("other", other)
	d.Set("create_on", createdOn)
	d.Set("resumed_on", resumedOn)
	d.Set("updated_on", updatedOn)
	d.Set("owner", owner)
	d.Set("comment", comment)
	d.Set("resource_monitor", resourceMonitor)
	d.Set("actives", actives)
	d.Set("pendings", pendings)
	d.Set("failed", failed)
	d.Set("suspended", suspended)
	d.Set("id", uuid)
	d.Set("scaling_policy", scalingPolicy)
	return nil
}

func deleteWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*providerConfiguration).DB
	whName := d.Get(whNameAttr).(string)
	sql := fmt.Sprintf("DROP WAREHOUSE  %s ", whName)
	if _, err := db.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error droping warehouse %q: {{err}}", whName), err)
	}
	return nil
}
