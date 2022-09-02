locals {
  buz_config_path       = "HONEYPOT_CONFIG_PATH"
  buz_config_path_value = "/etc/buz"
  activate_apis = [
    "artifactregistry.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com"
  ]
  domainParts  = split(".", var.buz_domain)
  cookieDomain = join(".", slice(local.domainParts, 1, length(local.domainParts))) # Assumes Buz is running on a subdomain and the cookie should be on root
}
