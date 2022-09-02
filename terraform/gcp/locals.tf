locals {
  buz_config_var  = "HONEYPOT_CONFIG_PATH"
  buz_config_path = "/etc/buz"
  activate_apis = [
    "artifactregistry.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com"
  ]
  domain_parts                 = split(".", var.buz_domain)
  cookie_domain                = join(".", slice(local.domain_parts, 1, length(local.domain_parts))) # Assumes Buz is running on a subdomain and the cookie should be on root
  artifact_registry_root       = "${var.gcp_region}-docker.pkg.dev/${var.gcp_project}"
  artifact_registry_repository = "${var.system}-repository"
  buz_source_image             = "ghcr.io/silverton-io/buz:${var.buz_version}"
  buz_image                    = "${local.artifact_registry_root}/${local.artifact_registry_repository}/buz:${var.buz_version}"
}
