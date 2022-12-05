output "buz_version" {
  value = var.buz_version
}

# output "buz_url" {
#   value = aws_apigatewayv2_api.lambda.api_endpoint
# }

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


output "buz_function_url" {
  value = aws_lambda_function_url.buz.function_url
}