locals {
  buz_config_var  = "HONEYPOT_CONFIG_PATH"
  buz_config_path = "/etc/buz"
  activate_apis = [
    "artifactregistry.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com"
  ]
  domain_parts    = split(".", var.buz_domain)
  cookie_domain   = join(".", slice(local.domain_parts, 1, length(local.domain_parts))) # Assumes Buz is running on a subdomain and the cookie should be on root
  system_env_base = "${var.system}-${var.env}-"
  # Image
  artifact_registry_location   = "${var.gcp_region}-docker.pkg.dev"
  artifact_registry_root       = "${local.artifact_registry_location}/${var.gcp_project}"
  artifact_registry_repository = "${local.system_env_base}repository"
  buz_source_image             = "ghcr.io/silverton-io/buz:${var.buz_version}"
  buz_image                    = "${local.artifact_registry_root}/${local.artifact_registry_repository}/buz:${var.buz_version}"
  # Config
  config = "${local.system_env_base}config"
  # Schema Bucket
  schema_bucket = "${local.system_env_base}${var.schema_bucket_name}"
  # Pub/Sub
  invalid_topic = "${local.system_env_base}${var.invalid_topic_name}"
  valid_topic   = "${local.system_env_base}${var.valid_topic_name}"
}
