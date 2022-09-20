data "aws_caller_identity" "current" {}


module "template_files" {
  source   = "hashicorp/dir/template"
  base_dir = "../../../../schemas/"
}

resource "aws_kinesis_firehose_delivery_stream" "buz_valid" {
  name        = local.valid_stream
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn        = aws_iam_role.firehose_role.arn
    bucket_arn      = aws_s3_bucket.buz_events.arn
    buffer_size     = var.firehose_buffer_size
    buffer_interval = var.firehose_buffer_interval

    prefix              = "valid/${local.s3_dynamic_prefix}/"
    error_output_prefix = "err"


    dynamic_partitioning_configuration {
      enabled = true
    }

    processing_configuration {
      enabled = true

      processors {
        type = "MetadataExtraction"
        parameters {
          parameter_name  = "JsonParsingEngine"
          parameter_value = "JQ-1.6"
        }
        parameters {
          parameter_name  = "MetadataExtractionQuery"
          parameter_value = "{vendor:.event.vendor,namespace:.event.namespace,version:.event.version}"
        }
      }
    }
  }
}

resource "aws_kinesis_firehose_delivery_stream" "buz_invalid" {
  name        = local.invalid_stream
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn        = aws_iam_role.firehose_role.arn
    bucket_arn      = aws_s3_bucket.buz_events.arn
    buffer_size     = var.firehose_buffer_size
    buffer_interval = var.firehose_buffer_interval

    prefix              = "invalid/${local.s3_dynamic_prefix}/"
    error_output_prefix = "err/invalid/"

    dynamic_partitioning_configuration {
      enabled = true
    }

    processing_configuration {
      enabled = true

      processors {
        type = "MetadataExtraction"
        parameters {
          parameter_name  = "JsonParsingEngine"
          parameter_value = "JQ-1.6"
        }
        parameters {
          parameter_name  = "MetadataExtractionQuery"
          parameter_value = "{vendor:.event.vendor,namespace:.event.namespace,version:.event.version}"
        }
      }
    }
  }
}

resource "aws_s3_bucket" "buz_events" {
  bucket = local.events_bucket
}

resource "aws_s3_object" "schemas" {
  for_each = module.template_files.files
  bucket   = aws_s3_bucket.buz_schemas.bucket
  key      = each.key
  source   = each.value.source_path
  etag     = filemd5(each.value.source_path)
}

resource "aws_s3_bucket" "buz_schemas" {
  bucket = local.schema_bucket
}

resource "aws_s3_bucket_acl" "events_acl" {
  bucket = aws_s3_bucket.buz_events.id
  acl    = "private"
}

resource "aws_s3_bucket_acl" "schemas_acl" {
  bucket = aws_s3_bucket.buz_schemas.id
  acl    = "private"
}

resource "aws_ecr_repository" "buz_repository" {
  name                 = local.artifact_repository
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

data "aws_ecr_image" "buz_image" {
  repository_name = aws_ecr_repository.buz_repository.name
  image_tag       = var.buz_version
  depends_on = [
    null_resource.build_and_push_image
  ]
}

resource "null_resource" "configure_docker" {
  triggers = {
    build_number = var.buz_version
  }
  provisioner "local-exec" {
    command = "aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com"
  }
  depends_on = [
    aws_ecr_repository.buz_repository
  ]
}

resource "local_file" "config" {
  filename = "config.yml.build"
  content = templatefile("config.tftpl", {
    system        = var.system,
    env           = var.env,
    mode          = "debug",
    port          = "8080",
    trackerDomain = var.buz_domain,
    cookieDomain  = local.cookie_domain,
    schemaBucket  = local.schema_bucket,
    validStream   = local.valid_stream,
    invalidStream = local.invalid_stream,
  })
}

resource "local_file" "dockerfile" {
  filename = "Dockerfile.build"
  content = templatefile("Dockerfile.tftpl", {
    sourceImage = local.buz_source_image
  })
}

resource "null_resource" "build_and_push_image" {
  triggers = {
    build_number = timestamp()
  }
  provisioner "local-exec" {
    command = "docker buildx build --push --platform=linux/amd64 -f Dockerfile.build -t ${aws_ecr_repository.buz_repository.repository_url}:${var.buz_version} ."
  }
  depends_on = [
    null_resource.configure_docker
  ]
}

resource "null_resource" "config_cleanup" {
  provisioner "local-exec" {
    command = "rm ${local_file.dockerfile.filename} || true && rm ${local_file.config.filename} || true"
  }
  depends_on = [
    null_resource.build_and_push_image
  ]
}

resource "aws_apprunner_auto_scaling_configuration_version" "buz" {
  auto_scaling_configuration_name = local.service_name

  max_concurrency = var.buz_service_container_concurrency
  min_size        = var.buz_service_container_min_count
  max_size        = var.buz_service_container_max_count
}

resource "aws_apprunner_observability_configuration" "buz" {
  observability_configuration_name = local.service_name

  trace_configuration {
    vendor = "AWSXRAY"
  }
}

resource "aws_apprunner_service" "buz" {
  service_name = local.service_name

  observability_configuration {
    observability_configuration_arn = aws_apprunner_observability_configuration.buz.arn
    observability_enabled           = true
  }
  auto_scaling_configuration_arn = aws_apprunner_auto_scaling_configuration_version.buz.arn

  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.apprunner_service_role.arn
    }
    image_repository {
      image_configuration {
        port = var.buz_service_container_port
        runtime_environment_variables = {
          (local.buz_config_var) : local.buz_config_path
        }
      }
      image_identifier      = "${aws_ecr_repository.buz_repository.repository_url}@${data.aws_ecr_image.buz_image.image_digest}"
      image_repository_type = "ECR"
    }
    auto_deployments_enabled = false
  }

  instance_configuration {
    cpu               = var.buz_service_cpu_limit
    memory            = var.buz_service_memory_limit
    instance_role_arn = aws_iam_role.apprunner_instance_role.arn
  }

  depends_on = [
    null_resource.build_and_push_image,
    aws_iam_role.apprunner_service_role
  ]
}

# resource "aws_apprunner_custom_domain_association" "buz" {
#   domain_name = var.buz_domain
#   service_arn = aws_apprunner_service.buz.arn
# }
