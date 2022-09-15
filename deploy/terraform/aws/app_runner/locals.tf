locals {
  buz_config_var      = "HONEYPOT_CONFIG_PATH"
  buz_config_path     = "/etc/buz"
  system_env_base     = "${var.system}-${var.env}-"
  artifact_repository = "${local.system_env_base}${var.repository_name}"
  buz_source_image    = "ghcr.io/silverton-io/buz:${var.buz_version}"
  service_name        = "${local.system_env_base}-collector"
  config              = "${local.system_env_base}config"
  schema_bucket       = "${local.system_env_base}${var.schema_bucket_name}"
  valid_bucket        = "${local.system_env_base}${var.valid_bucket_name}"
  invalid_bucket      = "${local.system_env_base}${var.invalid_bucket_name}"
  valid_topic         = "${local.system_env_base}${var.valid_firehose_name}"
  invalid_topic       = "${local.system_env_base}${var.invalid_firehose_name}"
}