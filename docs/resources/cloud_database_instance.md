---
subcategory: Cloud Database
page_title: "Bizfly Cloud: bizflycloud_cloud_database_instance"
description: |-
    Provide a Bizfly Cloud Database Instances resource. This can be used to create, modify, and delete.
---

# Resource: bizflycloud_cloud_database_instance

Provides a Bizfly Cloud Database Instances resource. This can be used to create, modify, and delete.

## Example

Create a database instance

```hcl
# Normal create a database instance without secondary
resource "bizflycloud_cloud_database_instance" "terraform_patroni_postgres" {
  name = "patroni_postgres"
  autoscaling = {
    enable           = 0
    volume_limited   = 100
    volume_threshold = 90
  }
  datastore = {
    type       = "Postgres"
    name       = "patroni-13.11"
    version_id = "2c4d00d1-b4e1-4a29-a7a6-25c5bb78a70f"
  }

  availability_zone = "HN1"
  flavor_name       = "2c_4g"
  instance_type     = "enterprise"
  network_ids = [
    "0706f928-02f6-4c4a-935c-fc22fdccbe81"
  ]
  public_access = true
  volume_size   = 20
}
```

Restore a backup - Create a cloud database instance from a backup

```hcl
data "bizflycloud_cloud_database_backup" "backup_daily" {
  id = "7201f1dd-6185-46a2-b488-87f2a4fce70e"
}

# Restore a backup - Create a cloud database instance from a backup
resource "bizflycloud_cloud_database_instance" "terraform_patroni_postgres" {
  name = "patroni_postgres"
  autoscaling = {
    enable           = 0
    volume_limited   = 100
    volume_threshold = 90
  }
  backup_id = "7201f1dd-6185-46a2-b488-87f2a4fce70e"
  datastore = {
    type       = data.bizflycloud_cloud_database_backup.backup_daily.datastore.type
    name       = data.bizflycloud_cloud_database_backup.backup_daily.datastore.name
    version_id = data.bizflycloud_cloud_database_backup.backup_daily.datastore.version_id
  }
  availability_zone = "HN1"
  flavor_name       = "2c_4g"
  instance_type     = "enterprise"

  network_ids = [
    "0706f928-02f6-4c4a-935c-fc22fdccbe81"
  ]

  public_access = true
  volume_size   = 20
}

data "bizflycloud_cloud_database_node" "terraform_patroni_postgres_primary" {
  id = resource.bizflycloud_cloud_database_instance.terraform_patroni_postgres.nodes.0.id
}
```

Restore a backup - Create a cloud database instance from a backup and secondary
and create databases/ users when start

```hcl
data "bizflycloud_cloud_database_backup" "backup_daily" {
  id = "7201f1dd-6185-46a2-b488-87f2a4fce70e"
}

# Restore a backup - Create a cloud database instance from a backup and secondary
resource "bizflycloud_cloud_database_instance" "terraform_patroni_mongo" {
  name = "mongo"
  autoscaling = {
    enable           = 1
    volume_limited   = 100
    volume_threshold = 90
  }
  backup_id = data.bizflycloud_cloud_database_backup.backup_daily.id
  datastore = {
    type       = data.bizflycloud_cloud_database_backup.backup_daily.datastore.type
    name       = data.bizflycloud_cloud_database_backup.backup_daily.datastore.name
    version_id = data.bizflycloud_cloud_database_backup.backup_daily.datastore.version_id
  }
  availability_zone = "HN1"
  flavor_name       = "2c_4g"
  instance_type     = "enterprise"

  secondaries {
   availability_zone = "HN1"
   quantity          = 2
  }

  network_ids = [
    "0706f928-02f6-4c4a-935c-fc22fdccbe81"
  ]

  init_databases = ["games", "genshin", "impact"]
  users {
    username  = "yanfei"
    password  = "password"
    databases = ["games", "genshin", "impact"]
  }

  users {
    username  = "nahida"
    password  = "password"
    databases = ["games", "genshin", "impact"]
  }

  users {
    username  = "nilou"
    password  = "password"
    databases = ["games", "genshin", "impact"]
  }

  public_access = true
  volume_size   = 20
}

```

### Argument Reference

The following arguments are supported:

-   `autoscaling` - Information about autoscaling volume, ... for this database instance
-   `availability_zone` - (Required) The data center - that will allocate database instance primary node
-   `backup_id` - The backup id that want restore to new database instance
-   `datastore` - (Required) The datastore of database instance
-   `flavor_name` - (Required) The setting CPU/ RAM of this database node. This being applied for all database nodes
-   `init_databases` - The list databases name - that want create when master database node start
-   `instance_type` - (Required) The instance type of this database instance
-   `name` - (Required) Name of this database instance
-   `network_ids` - (Required)
-   `public_access` - Whether database instance can be connect from public internet
-   `secondaries` - Whether database instance with secondaries?
    -   `availability_zone` - The data center - that will allocate database instance members
-   `users` - The list users - that want create when master database node start
    -   `name` - The username of this user
    -   `password` - The password of this user
    -   `host` - (Optional) The source host a valid to this user access to database. This is valid with MariaDB/ MySQL
    -   `databases` - The list databases name that user can be access. If not define this value, user can access to all databases.
-   `volume_size` - (Required) The size of storage provision for this database instance, unit is `Gigabytes`
-   `configuration_group` - Define custom config for database nodes of this instance
    -   `id` - The ID of configuration group include custom configs that attach to all members of this database instance
    -   `apply_immediately` - Default: `true`. The node being restart after when attach configuration group if required

### Atrributes Reference

The following attributes are exported:

-   `created_at` - The time that init this database instance
-   `dns` - The DNS of this database instance (This current just have with MongoDB)
-   `id` - ID of this database instance
-   `nodes` - The list nodes - that are member of this database instance
-   `status` - The status of this database instance. Have some status like: `ACTIVE`, `RESIZE`, `ERROR`, `BUILD`
