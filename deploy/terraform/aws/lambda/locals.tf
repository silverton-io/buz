locals {
  domain_parts               = split(".", var.buz_domain)
  cookie_domain              = join(".", slice(local.domain_parts, 1, length(local.domain_parts))) # Assumes Buz is running on a subdomain and the cookie should be on root
  buz_debug_var              = "DEBUG"
  buz_config_var             = "BUZ_CONFIG_PATH"
  buz_config_path            = "/etc/buz/config.yml"
  system_env_base            = "${var.system}-${var.env}-"
  artifact_repository        = "${local.system_env_base}img"
  image                      = "buz:${var.buz_version}"
  buz_source_image           = "${var.buz_image_repo}/${local.image}"
  service_name               = "${local.system_env_base}collector"
  policy_name                = "${local.system_env_base}policy"
  config                     = "${local.system_env_base}config"
  schema_bucket              = "${local.system_env_base}${var.schema_bucket_name}"
  events_bucket              = "${local.system_env_base}${var.events_bucket_name}"
  default_output             = "buz_events"
  deadletter_output          = "buz_invalid_events"
  metadata_extraction_params = "{isValid:.isValid,vendor:.vendor,namespace:.namespace,version:.version}"
  s3_dynamic_prefix          = "isValid=!{partitionKeyFromQuery:isValid}/vendor=!{partitionKeyFromQuery:vendor}/namespace=!{partitionKeyFromQuery:namespace}/version=!{partitionKeyFromQuery:version}/year=!{timestamp:yyyy}/month=!{timestamp:MM}/day=!{timestamp:dd}/"
}
