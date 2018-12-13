Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-$PROVIDER_NAME`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-$PROVIDER_NAME
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-$PROVIDER_NAME
$ make build
```

Using the provider
----------------------
## Fill in for each provider
```
go get github.com/snowflakedb/gosnowflake
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-$PROVIDER_NAME
...
```


Provider Configuration
--------------------------

Firstly, to use the provider you will need to create a user within Snowflake that can execute the resource requests made by Terraform

```sh
$ export SF_USER
$ export SF_PASSWORD
$ export SF_REGION
$ export SF_ACCOUNT
```

### Snowflake Warehouse Management
```
resource "snowflake_warehouse" "warehouse_terraform" {
      name              =   "dev_wh"
      warehouse_size    =   "SMALL"
      auto_resume       =   false
      auto_suspend      =   600
      comment           =   "terraform development warehouse"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------ | ------ |
| `name` | Name of the Snowflake warehouse | String | TRUE |
| `max_concurrency_level` | Max concurrent SQL statements that can run on warehouse | String | FALSE |
| `statement_queued_timeout_in_seconds` | Time, in seconds, an SQL statement can be queued before being cancelled | String | FALSE |
| `statement_timeout_in_seconds` | Time, in seconds, after which an SQL statement will be terminated | String | FALSE |
| `warehouse_size` | Size of the warehouse | String | FALSE |
| `max_cluster_count` | Min number of warehouses | String | FALSE |
| `min_cluster_count` | Max number of warehouses | String | FALSE |
| `auto_resume` | Should warehouse should auto resume | Boolean | FALSE |
| `auto_suspend` | Should warehouse should auto suspend | Boolean | FALSE |
| `initially_suspended` | Should warehouse start off suspended  | Boolean | FALSE |
| `resource_monitor` | Name of resource monitor assigned to warehouse | Boolean | FALSE |
| `comment` | Additional comments | String | FALSE |

### Snowflake Database Management
```
resource "snowflake_database" "database_terraform" {
      name              =   "dev_db"
      comment           =   "terraform development database"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------   | ------ |
| `name` | Name of the Snowflake database | String | TRUE |
| `comment` | Additional comments | String | FALSE |

### Snowflake User Management
```
resource "snowflake_user" "tf_test_user" {
  user = "terraform.test"
  host = "mydomain.org"
  plaintext_password = "12345QWERTYqwerty"
  default_role = "READONLY"
}
```

##### Properties
| Property | Description | Type | Required |
| ------ | ------ | ------ | ------ |
| `user` | The username of the user | String | TRUE |
| `host` | Host/TLD associated with the user. The default for this is localhost. This has a direct effect on the Username | String | FALSE |
| `plaintext_password` | Password of the user. Ensure that passwords conform to the complexity requirements by Snowflake | String | TRUE |
| `default_role` | Default role the user assumes. Defaults to `null` | String | FALSE |
