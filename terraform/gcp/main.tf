locals {
  activate_apis = [
    "artifactregistry.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com"
  ]
}

data "google_project" "project" {}

resource "google_project_service" "project_services" {
  for_each                   = toset(local.activate_apis)
  project                    = var.gcp_project
  service                    = each.value
  disable_on_destroy         = true
  disable_dependent_services = true
}

resource "google_storage_bucket" "buz_schemas" {
  name          = "buz-schemas-${data.google_project.project.number}"
  location      = "US"
  force_destroy = true
}

resource "google_pubsub_topic" "valid_topic" {
  name = "buz-valid"
}

resource "google_pubsub_topic" "invalid_topic" {
  name = "buz-invalid"
}

resource "google_secret_manager_secret" "buz_config" {
  secret_id = "buz-config"

  replication {
    user_managed {
      replicas {
        location = var.gcp_region
      }
    }
  }

  depends_on = [
    google_project_service.project_services
  ]
}

resource "google_secret_manager_secret_version" "buz_config" {
  secret      = google_secret_manager_secret.buz_config.id
  secret_data = file("gcp.yml")
}

resource "google_artifact_registry_repository" "buz_repository" {
  location      = var.gcp_region
  repository_id = "buz-repository"
  format        = "DOCKER"

  depends_on = [
    google_project_service.project_services
  ]
}

resource "google_project_iam_binding" "buz_config_secret_access" {
  project = var.gcp_project
  role    = "roles/secretmanager.secretAccessor"
  members = [
    "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
  ]

  # condition {
  #   title      = "buz-config-secret-access"
  #   expression = "resource.name == \"buz-config\""
  # }

  depends_on = [
    google_secret_manager_secret_version.buz_config
  ]
}

resource "google_cloud_run_service" "buz" {
  name     = "buz"
  location = var.gcp_region

  template {
    spec {
      timeout_seconds       = 300
      container_concurrency = 80

      volumes {
        name = "buz-config"
        secret {
          secret_name = google_secret_manager_secret.buz_config.secret_id
          items {
            key  = "latest"
            path = "config.yml"
          }
        }
      }

      containers {
        image = "${var.gcp_region}-docker.pkg.dev/${var.gcp_project}/buz-repository/buz@${local.buz_sha}"

        resources {
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }

        ports {
          container_port = 8080
        }

        env {
          name  = local.buz_config_path
          value = local.buz_config_path_value
        }

        volume_mounts {
          name       = "buz-config"
          mount_path = local.buz_config_path_value
        }
      }
    }
  }

  depends_on = [
    google_project_service.project_services,
    google_storage_bucket.buz_schemas,
    google_pubsub_topic.invalid_topic,
    google_pubsub_topic.valid_topic
  ]
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.buz.location
  project  = google_cloud_run_service.buz.project
  service  = google_cloud_run_service.buz.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

resource "google_cloud_run_domain_mapping" "buz" {
  location = "us-central1"
  name     = local.domain

  metadata {
    namespace = data.google_project.project.project_id
  }

  spec {
    route_name = google_cloud_run_service.buz.name
  }
}
