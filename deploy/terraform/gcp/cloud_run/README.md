# Terraform GCP
Deploy Buz to Google Cloud via Terraform.

## Prerequisites

You will need `terraform` and `gcloud`.

## Spinning up Buz

**Auth Gcloud**

```
gcloud auth application-default login
```


**Apply Terraform**

```
terraform apply
```

**[Optional] - Create and populate terraform.tfvars**

If you don't want to pass terraform variables interactively, you can optionally create a `terraform.tfvars` file in this directory and populate it:

```
gcp_project = "my-project-23456"
gcp_region  = "us-central1"
system      = "buz"
buz_domain  = "track.yourdomain.com"
buz_version = "v0.15.1"
```
