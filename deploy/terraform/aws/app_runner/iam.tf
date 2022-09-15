resource "aws_iam_role" "firehose_role" {
  name = local.firehose_role

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "firehose.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role" "app_runner_role" {
  name = local.app_runner_role

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = [
            "build.apprunner.amazonaws.com",
            "tasks.apprunner.amazonaws.com"
          ]
        }
      },
      {
        Action = [
          "firehose:Put*",
        ]
        Effect = "Allow"
        Resource = [
          aws_kinesis_firehose_delivery_stream.buz_valid.name,
          aws_kinesis_firehose_delivery_stream.buz_invalid.name,
        ]
      },
      {
        Action = [
          "s3:Get*",
        ]
        Effect = "Allow"
        Resource = [
          aws_s3_bucket.buz_schemas.id,
          "${aws_s3_bucket.buz_schemas.id}/*",
        ]
      },
      {
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Effect = "Allow"
        Resource = [
          aws_secretsmanager_secret.buz_config.arn
        ]
      },
      {
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchCheckLayerAvailability",
          "ecr:BatchGetImage",
          "ecr:DescribeImages",
          "ecr:GetAuthorizationToken"
        ],
        Effect = "Allow",
        Resource = [
          aws_ecr_repository.buz_repository.arn
        ]
      }
    ]
  })
}