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

variable "buz_version" {
  description = "The version of Buz to run. \n\nExample: v0.11.11"
  type        = string
}

variable "buz_service_container_concurrency" {
  description = "The service container concurrency"
  type        = number
  default     = 200
}

variable "buz_service_cpu_limit" {
  description = "The service cpu limit"
  type        = string
  default     = "1"
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

variable "valid_firehose_name" {
  description = "The name of the Kinesis Firehose for valid events"
  type        = string
  default     = "valid-events"
}

variable "invalid_firehose_name" {
  description = "The name of the Kinesis Firehose for invalid events"
  type        = string
  default     = "invalid-events"
}
