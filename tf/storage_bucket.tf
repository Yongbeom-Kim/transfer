variable "bucket_name" {
  description = "The name of the bucket to create."
  type        = string
}

resource "google_storage_bucket" "upload" {
  name          = var.bucket_name
  location      = "EU"
  force_destroy = true

  uniform_bucket_level_access = true
}