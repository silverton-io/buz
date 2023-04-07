output "gcp_project" {
  value = var.gcp_project
}

output "gcp_region" {
  value = var.gcp_region
}

output "buz_domain" {
  value = var.buz_domain
}

output "buz_version" {
  value = var.buz_version
}

output "schema_bucket" {
  value = local.schema_bucket
}

output "default_topic" {
  value = local.default_output
}

output "deadletter_topic" {
  value = local.deadletter_output
}

output "buz_service_id" {
  value = google_cloud_run_service.buz.id
}

output "buz_service_status" {
  value = google_cloud_run_service.buz.status
}

output "bigquery_dataset" {
  value = var.bigquery_dataset_name
}

output "default_table" {
  value = local.default_table
}

output "deadletter_table" {
  value = local.deadletter_table
}
