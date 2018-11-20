resource "snowflake_warehouse" "warehouse_terraform" {
      name              =   "dev_wh"
      warehouse_size    =   "SMALL"
      auto_resume       =   false
      auto_suspend      =   600
      comment           =   "terraform development warehouse"
}

resource "snowflake_database" "database_terraform" {
      name              =   "dev_db"
      comment           =   "terraform development database"
}
