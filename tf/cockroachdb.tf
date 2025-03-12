variable "cockroach_cluster_name" {
  type = string
  description = "The name of the CockroachDB cluster"
}

variable "cockroach_db_name" {
  type = string
  description = "The name of the CockroachDB database"
}

variable "cockroach_admin_user_name" {
  type = string
  description = "The name of the CockroachDB admin user"
}

variable "cockroach_admin_user_password" {
  type = string
  sensitive = true
  description = "The password of the CockroachDB admin user"
}

resource "cockroach_cluster" "db" {
  name = var.cockroach_cluster_name
  cloud_provider = "AWS"
  plan = "BASIC"
  serverless = {}
  regions = [
    {
        name = "ap-southeast-1"
    }
  ]
}

resource "cockroach_database" "db" {
  name = var.cockroach_db_name
  cluster_id = cockroach_cluster.db.id
}

resource "cockroach_sql_user" "admin_user" {
  name = var.cockroach_admin_user_name
  password = var.cockroach_admin_user_password
  cluster_id = cockroach_cluster.db.id
}

output "cockroach_db_url" {
  value = cockroach_cluster.db.regions[0].sql_dns
}
