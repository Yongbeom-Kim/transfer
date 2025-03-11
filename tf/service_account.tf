resource "google_project_organization_policy" "disable_service_account_key_creation" {
  project = var.project_id
  constraint = "iam.disableServiceAccountKeyCreation"

  boolean_policy {
    enforced = false
  }
}

resource "google_service_account" "backend" {
  account_id   = "backend-${var.environment}"
  display_name = "Backend Service Account"
}

resource "google_storage_bucket_iam_member" "backend_service_account_storage_admin" {
  bucket = var.bucket_name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_service_account_key" "backend_service_account_key" {
  service_account_id = google_service_account.backend.name

  depends_on = [google_project_organization_policy.disable_service_account_key_creation]
}

resource "local_file" "backend_service_account_key_file" {
  content  = base64decode(google_service_account_key.backend_service_account_key.private_key)
  filename = "${path.module}/../backend/service_account_key/${var.environment}.json"
}

