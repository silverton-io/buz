# Terraform GCP
Deploy Buz to Google Cloud via Terraform.

1. Create a file in this directory `terraform.tfvars`. Set the `gcp_project` and optionally set `gcp_region`.
```
gcp_project = "my-project-23456"
gcp_region  = "us-central1" #default value
```

2. In the command line authenticate GCP
```
gcloud auth application-default login
```

3. Create the artifact registry
```
terraform plan -target=google_artifact_registry_repository.buz_repository
```

4. After your registry has been created. Upload the appropriate buz image
```
gcloud auth configure-docker us-central1-docker.pkg.dev

docker pull ghcr.io/silverton-io/buz:v0.11.11@sha256:130ed9421579125e4f38089e4c2d1e07038fb26591a15082010a52b95f3a5dda # amd64

docker tag ghcr.io/silverton-io/buz:v0.11.11@sha256:130ed9421579125e4f38089e4c2d1e07038fb26591a15082010a52b95f3a5dda us-central1-docker.pkg.dev/{YOUR PROJECT NAME}/buz-repository/buz:v0.11.11

docker push us-central1-docker.pkg.dev/{YOUR PROJECT NAME}/buz-repository/buz:v0.11.11
```

5. Deploy the rest of the buz stack
```
terraform apply
```