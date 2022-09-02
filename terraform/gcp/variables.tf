variable "gcp_project" {
  description = "GCloud Project ID"
  type        = string
}

variable "gcp_region" {
  description = "GCloud Region"
  type        = string
  default     = "us-central1"
}

variable "system" {
  description = "The name of the Buz implementation"
  type        = string
  default     = "buzzy"
}

variable "env" {
  description = "The name of the Buz environment"
  type        = string
  default     = "dev"
}

variable "buz_domain" {
  description = "The domain or subdomain to map Buz to"
  type        = string
}

variable "buz_image_sha" {
  description = "The Buz image SHA"
  type        = string
  default     = "sha256:1083c0333c284dfa16dd7cc621f90b8a1197fe4d9905237e41f2f1a495481d92"
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
