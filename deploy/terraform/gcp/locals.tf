locals {
  buz_config_var  = "BUZ_CONFIG_PATH"
  buz_config_path = "/etc/buz/config.yml"
  activate_apis = [
    "artifactregistry.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com"
  ]
  domain_parts                 = split(".", var.buz_domain)
  cookie_domain                = join(".", slice(local.domain_parts, 1, length(local.domain_parts))) # Assumes Buz is running on a subdomain and the cookie should be on root
  system_env_base              = "${var.system}-${var.env}-"
  artifact_repository          = "${local.system_env_base}repository"
  artifact_registry_location   = "${var.gcp_region}-docker.pkg.dev"
  artifact_registry_root       = "${local.artifact_registry_location}/${var.gcp_project}"
  artifact_registry_repository = "${local.system_env_base}repository"
  buz_source_image             = "ghcr.io/silverton-io/buz:${var.buz_version}"
  buz_image                    = "${local.artifact_registry_root}/${local.artifact_registry_repository}/buz:${var.buz_version}"
  service_name                 = "${local.system_env_base}collector"
  config                       = "${local.system_env_base}config"
  schema_bucket                = "${local.system_env_base}${var.schema_bucket_name}"
  invalid_topic                = "${local.system_env_base}invalid-events"
  valid_topic                  = "${local.system_env_base}events"
  valid_events_subscription    = "${local.system_env_base}events"
  invalid_events_subscription  = "${local.system_env_base}invalid-events"
  events_table_fqn             = "${var.gcp_project}.${var.bigquery_dataset_name}.${var.bigquery_valid_events_table_name}"
  invalid_events_table_fqn     = "${var.gcp_project}.${var.bigquery_dataset_name}.${var.bigquery_invalid_events_table_name}"
}
