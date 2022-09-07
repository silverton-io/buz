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

output "valid_topic" {
  value = local.valid_topic
}

output "invalid_topic" {
  value = local.invalid_topic
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

output "bigquery_valid_events_table" {
  value = var.bigquery_valid_events_table_name
}

output "bigquery_invalid_events_table" {
  value = var.bigquery_invalid_events_table_name
}
