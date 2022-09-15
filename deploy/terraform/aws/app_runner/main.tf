resource "aws_kinesis_firehose_delivery_stream" "buz_valid" {
  name        = local.valid_topic
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose_role.arn
    bucket_arn = aws_s3_bucket.buz_valid.arn
  }
}

resource "aws_kinesis_firehose_delivery_stream" "buz_invalid" {
  name        = local.invalid_topic
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose_role.arn
    bucket_arn = aws_s3_bucket.buz_invalid.arn
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
    image_repository {
      image_configuration {
        port = var.buz_service_container_port
      }
      image_identifier      = "${aws_ecr_repository.buz_repository.repository_url}@${data.aws_ecr_image.buz_image.image_digest}"
      image_repository_type = "ECR"
    }
    auto_deployments_enabled = false
  }

  instance_configuration {
    cpu               = var.buz_service_cpu_limit
    memory            = var.buz_service_memory_limit
    instance_role_arn = aws_iam_role.app_runner_role.arn
  }
}
