variable "aws_region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "aws_profile" {
  description = "AWS Profile"
  type        = string
  default     = "default"
}

variable "system" {
  description = "The name of the Buz implementation. \n\nExample: buz"
  type        = string
}

variable "env" {
  description = "The name of the Buz environment. \n\nExample: development/staging/production"
  type        = string
  default     = "dev"
}

variable "debug" {
  description = "The debug environment variable"
  type        = string
  default     = "1"
}

variable "buz_domain" {
  description = "The subdomain to map Buz to. \n\nExample: track.yourdomain.com"
  type        = string
}

variable "buz_image_repo" {
  description = "The Buz image repository"
  type        = string
  default     = "ghcr.io/silverton-io"
}

variable "buz_version" {
  description = "The version of Buz to run."
  type        = string
  default     = "v0.18.3"
}

variable "buz_lambda_memory_limit" {
  description = "The lambda memory limit"
  type        = number
  default     = 128
}

variable "buz_lambda_timeout" {
  description = "The lambda timeout"
  type        = number
  default     = 5
}

variable "buz_service_container_port" {
  description = "The service container port"
  type        = number
  default     = 8080
}

variable "schema_bucket_name" {
  description = "The name of the AWS bucket for schemas. \n\nPLEASE NOTE! Buckets are globally unique so you may need to be creative."
  type        = string
}

variable "events_bucket_name" {
  description = "The name of the AWS bucket for events. \n\nPLEASE NOTE! Buckets are globally unique so you may need to be creative."
  type        = string
}

variable "firehose_buffer_size" {
  description = "The size of the firehose buffer, in MiB"
  type        = number
  default     = 128 # Maximum
}

variable "firehose_buffer_interval" {
  description = "The buffer interval, in seconds"
  type        = number
  default     = 60 # 1 minute
}

variable "certificate_arn" {
  description = "(Optional) The ACM certificate arn to use for the pretty dns name"
}
