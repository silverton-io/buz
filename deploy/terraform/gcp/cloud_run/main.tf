data "google_project" "project" {}

// NOTE
//  Uncomment this if you want to store tfstate in a
//  pre-existing GCS bucket
# terraform {
#   backend "gcs" {
#     bucket = "YOUR_TFSTATE_BUCKET"
#     prefix = "YOUR_TFSTATE_PREFIX"
#   }
# }

terraform {
  backend "gcs" {
    bucket = "silverton-tfstate"
    prefix = "buz/production/"
  }
}

module "template_files" {
  source = "hashicorp/dir/template"

  base_dir = "../../../../schemas"
}

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

resource "google_storage_bucket_object" "schemas" {
  for_each = module.template_files.files
  bucket   = google_storage_bucket.schemas.name
  name     = each.key
  source   = each.value.source_path
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
  secret_data = templatefile("${path.module}/config.yml.tftpl", {
    project          = var.gcp_project,
    dataset          = var.bigquery_dataset_name,
    system           = var.system,
    env              = var.env,
    port             = var.buz_service_container_port
    trackerDomain    = var.buz_domain,
    cookieDomain     = local.cookie_domain,
    schemaBucket     = local.schema_bucket,
    defaultOutput    = local.default_output,
    deadletterOutput = local.deadletter_output,
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
  name                       = local.service_name
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

        env {
          name  = "GOMEMLIMIT"
          value = floor(var.buz_service_gomemlimit_pct * var.buz_service_memory_limit)
        }

        volume_mounts {
          name       = local.config
          mount_path = local.buz_config_dir
        }
      }
    }
  }

  depends_on = [
    google_project_service.project_services,
    google_storage_bucket.schemas,
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

############################################################################
# Bigquery
############################################################################

resource "google_project_iam_member" "bigquery_viewer" {
  project = data.google_project.project.project_id
  role    = "roles/bigquery.metadataViewer"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "bigquery_editor" {
  project = data.google_project.project.project_id
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "buz" {
  dataset_id    = var.bigquery_dataset_name
  friendly_name = var.bigquery_dataset_name
  description   = "A dataset for Buz events"
  location      = var.bigquery_location
}
