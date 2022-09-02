variable "gcp_project" {
  description = "GCloud Project ID"
  type        = string
  default     = "silvertonio"
}

variable "gcp_region" {
  description = "GCloud Region"
  type        = string
  default     = "us-central1"
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
  description = "The version of Buz to run. \n\nExample: v0.11.11"
  type        = string
}

variable "valid_topic_name" {
  description = "The name of the Pub/Sub topic for valid events"
  type        = string
  default     = "valid-events"
}

variable "invalid_topic_name" {
  description = "The name of the Pub/Sub topic for invalid events"
  type        = string
  default     = "invalid-events"
}
