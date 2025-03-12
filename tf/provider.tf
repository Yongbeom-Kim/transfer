variable "project_id" {
  description = "The ID of the project in which to create resources."
  type        = string
}

variable "region" {
  description = "The region in which to create resources."
  type        = string
}

variable "environment" {
  description = "The environment for the resources (e.g., dev, prod)."
  type        = string
}

terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "6.24.0"
    }
    cockroach = {
      source = "cockroachdb/cockroach"
      version = "1.11.2"
    }
  }
}

provider "google" {
  project = var.project_id
  region = var.region
  default_labels = {
    managed_by = "opentofu"
    environment = var.environment
    project = var.project_id
  }
}


provider "cockroach" {
  # Configuration options
}