output "buz_version" {
  value = var.buz_version
}

output "schema_bucket" {
  value = local.schema_bucket
}

output "events_bucket" {
  value = local.events_bucket
}

output "valid_stream" {
  value = local.valid_stream
}

output "invalid_stream" {
  value = local.invalid_stream
}
