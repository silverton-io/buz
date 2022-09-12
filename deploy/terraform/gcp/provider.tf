terraform {
  required_version = ">= 1.0.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.33.0"
    }
  }
}

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
}
