
data "google_project" "project" {}

resource "google_project_service" "project_services" {
  for_each                   = toset(local.activate_apis)
  project                    = var.gcp_project
  service                    = each.value
  disable_on_destroy         = true
  disable_dependent_services = true
}

resource "google_storage_bucket" "schemas" {
  name          = local.schema_bucket
  location      = var.schema_bucket_location
  force_destroy = true
}

resource "google_pubsub_topic" "valid_topic" {
  name = local.valid_topic
}

resource "google_pubsub_topic" "invalid_topic" {
  name = local.invalid_topic
}

resource "google_secret_manager_secret" "buz_config" {
  secret_id = local.config

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
  secret = google_secret_manager_secret.buz_config.id
  secret_data = templatefile("config.tftpl", {
    project       = var.gcp_project,
    system        = var.system,
    env           = var.env,
    port          = var.buz_service_container_port
    trackerDomain = var.buz_domain,
    cookieDomain  = local.cookie_domain,
    schemaBucket  = local.schema_bucket,
    validTopic    = local.valid_topic,
    invalidTopic  = local.invalid_topic,
  })
}

resource "google_artifact_registry_repository" "buz_repository" {
  location      = var.gcp_region
  repository_id = local.artifact_repository
  format        = "DOCKER"

  depends_on = [
    google_project_service.project_services
  ]
}

resource "null_resource" "configure_docker" {
  triggers = {
    build_number = timestamp()
  }
  provisioner "local-exec" {
    command = "gcloud auth configure-docker ${local.artifact_registry_location}"
  }
  depends_on = [
    google_artifact_registry_repository.buz_repository
  ]
}

resource "null_resource" "pull_and_push_image" {
  triggers = {
    build_number = timestamp()
  }
  provisioner "local-exec" {
    command = "docker pull ${local.buz_source_image} --platform=linux/amd64 && docker tag ${local.buz_source_image} ${local.buz_image} && docker push ${local.buz_image}"
  }
  depends_on = [
    google_artifact_registry_repository.buz_repository,
    null_resource.configure_docker
  ]
}

resource "google_project_iam_binding" "buz_config_secret_access" {
  project = var.gcp_project
  role    = "roles/secretmanager.secretAccessor"
  members = [
    "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
  ]

  depends_on = [
    google_secret_manager_secret_version.buz_config
  ]
}

resource "google_cloud_run_service" "buz" {
  name                       = var.system
  location                   = var.gcp_region
  autogenerate_revision_name = true

  template {
    spec {
      timeout_seconds       = var.buz_service_timeout_seconds
      container_concurrency = var.buz_service_container_concurrency

      volumes {
        name = local.config
        secret {
          secret_name = google_secret_manager_secret.buz_config.secret_id
          items {
            key  = "latest"
            path = "config.yml"
          }
        }
      }

      containers {
        image = local.buz_image

        resources {
          limits = {
            cpu    = var.buz_service_cpu_limit
            memory = var.buz_service_memory_limit
          }
        }

        ports {
          container_port = var.buz_service_container_port
        }

        env {
          name  = local.buz_config_var
          value = local.buz_config_path
        }

        volume_mounts {
          name       = local.config
          mount_path = local.buz_config_path
        }
      }
    }
  }

  depends_on = [
    google_project_service.project_services,
    google_storage_bucket.schemas,
    google_pubsub_topic.invalid_topic,
    google_pubsub_topic.valid_topic,
    null_resource.pull_and_push_image,
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
  location = var.gcp_region
  name     = var.buz_domain

  metadata {
    namespace = data.google_project.project.project_id
  }

  spec {
    route_name = google_cloud_run_service.buz.name
  }
}
