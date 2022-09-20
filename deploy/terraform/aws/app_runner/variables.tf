variable "aws_region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "system" {
  description = "The name of the Buz implementation. \n\nExample: buz"
  type        = string
}

variable "env" {
  description = "The name of the Buz environment"
  type        = string
  default     = "dev"
}

variable "buz_domain" {
  description = "The subdomain to map Buz to. \n\nExample: track.yourdomain.com"
  type        = string
}

variable "buz_version" {
  description = "The version of Buz to run. \n\nExample: v0.11.14"
  type        = string
}

variable "buz_service_container_concurrency" {
  description = "The service container concurrency"
  type        = number
  default     = 200
}

variable "buz_service_container_min_count" {
  description = "The minimum number of buz instances to run"
  type = number
  default = 1 # The minimum since app runner doesn't go to 0. Ouch.
}

variable "buz_service_container_max_count" {
  description = "The maximum number of buz instances to run"
  type = number
  default = 25 # The actual maximum w/out bumping limits
}

variable "buz_service_cpu_limit" {
  description = "The service cpu limit"
  type        = number
  default     = 1024
}

variable "buz_service_memory_limit" {
  description = "The service memory limit"
  type        = string
  default     = "2048"
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