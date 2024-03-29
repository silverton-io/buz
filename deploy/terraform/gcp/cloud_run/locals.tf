locals {
  buz_config_var  = "BUZ_CONFIG_PATH"
  buz_config_dir  = "/etc/buz/"
  buz_config_path = "${local.buz_config_dir}config.yml"
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
  default_output               = "${local.system_env_base}events"
  default_subscription         = "${local.system_env_base}events"
  default_table                = var.default_bigquery_table
  deadletter_output            = "${local.system_env_base}invalid-events"
  deadletter_subscription      = "${local.system_env_base}invalid-events"
  deadletter_table             = var.deadletter_bigquery_table
}
