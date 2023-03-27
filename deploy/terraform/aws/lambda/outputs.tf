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

output "default_output" {
  value = local.default_output
}

output "deadletter_output" {
  value = local.deadletter_output
}

output "buz_function_url" {
  value = aws_lambda_function_url.buz.function_url
}

output "buz_cloudfront_url" {
    value = aws_cloudfront_distribution.buz.domain_name
}
