resource "snowflake_warehouse" "shoprunner_warehouse_terraform" {
      name              =   "shoprunner_terraform"
      warehouse_size    =   "SMALL"
      auto_resume       =   false
      auto_suspend      =   600
      comment           =   "shoprunner terraform development warehouse"
}

resource "snowflake_database" "shoprunner_database_terraform" {
      name              =   "shoprunner_terraform_db"
      comment           =   "shoprunner terraform development database"
}
