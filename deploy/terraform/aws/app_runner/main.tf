data "aws_caller_identity" "current" {}

resource "aws_kinesis_firehose_delivery_stream" "buz_valid" {
  name        = local.valid_topic
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose_role.arn
    bucket_arn = aws_s3_bucket.buz_valid.arn

    prefix = local.s3_dynamic_prefix

    processing_configuration {
      enabled = "true"

      processors {
        type = "MetadataExtraction"
        parameters {
          parameter_name  = "JsonParsingEngine"
          parameter_value = "{namespace:.event.namespace}"
        }
        parameters {
          parameter_name  = "JsonParsingEngine"
          parameter_value = "{version:.event.version}"
        }
      }
    }
  }
}

resource "aws_kinesis_firehose_delivery_stream" "buz_invalid" {
  name        = local.invalid_topic
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose_role.arn
    bucket_arn = aws_s3_bucket.buz_invalid.arn

    prefix = local.s3_dynamic_prefix

    processing_configuration {
      enabled = "true"

      processors {
        type = "MetadataExtraction"
        parameters {
          parameter_name  = "JsonParsingEngine"
          parameter_value = "{namespace:.event.namespace}"
        }
        parameters {
          parameter_name  = "JsonParsingEngine"
          parameter_value = "{version:.event.version}"
        }
      }
    }
  }
}

resource "aws_s3_bucket" "buz_valid" {
  bucket = local.valid_bucket
}

resource "aws_s3_bucket" "buz_invalid" {
  bucket = local.invalid_bucket
}

resource "aws_s3_bucket" "buz_schemas" {
  bucket = local.schema_bucket
}

resource "aws_s3_bucket_acl" "valid_acl" {
  bucket = aws_s3_bucket.buz_valid.id
  acl    = "private"
}

resource "aws_s3_bucket_acl" "invalid_acl" {
  bucket = aws_s3_bucket.buz_invalid.id
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
  image_tag       = "latest"
}

resource "aws_secretsmanager_secret" "buz_config" {
  name = local.config
}

resource "aws_secretsmanager_secret_version" "buz_config" {
  secret_id     = aws_secretsmanager_secret.buz_config.id
  secret_string = templatefile("config.tftpl", {})
}

resource "null_resource" "configure_docker" {
  triggers = {
    build_number = timestamp()
  }
  provisioner "local-exec" {
    command = "aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com"
  }
  depends_on = [
    aws_ecr_repository.buz_repository
  ]
}

resource "null_resource" "pull_and_push_image" {
  triggers = {
    build_number = timestamp()
  }
  provisioner "local-exec" {
    command = "docker pull ${local.buz_source_image} && docker tag ${local.buz_source_image} ${aws_ecr_repository.buz_repository.arn}:latest && docker push ${aws_ecr_repository.buz_repository.arn}:latest"
  }
  depends_on = [
    null_resource.configure_docker
  ]
}

resource "aws_apprunner_auto_scaling_configuration_version" "buz" {
  auto_scaling_configuration_name = local.service_name

  max_concurrency = var.buz_service_container_concurrency
  min_size        = 1
  max_size        = 5
}

resource "aws_apprunner_service" "buz" {
  service_name = local.service_name

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
    null_resource.pull_and_push_image,
    aws_iam_role.apprunner_service_role
  ]
}
