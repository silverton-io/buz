resource "aws_iam_role" "firehose_role" {
  name = "${local.service_name}FirehoseRole"

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

# App Runner Service 
resource "aws_iam_role" "apprunner_service_role" {
  name               = "${local.service_name}AppRunnerECRAccessRole"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.apprunner_service_assume_policy.json
}

data "aws_iam_policy_document" "apprunner_service_assume_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["build.apprunner.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "apprunner_service_role_attachment" {
  role       = aws_iam_role.apprunner_service_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess"
}

# App Runner Instance 
resource "aws_iam_role" "apprunner_instance_role" {
  name               = "${local.service_name}AppRunnerInstanceRole"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.apprunner_instance_assume_policy.json
}

resource "aws_iam_policy" "apprunner_policy" {
  name   = "${local.service_name}Apprunner"
  policy = data.aws_iam_policy_document.apprunner_instance_role_policy.json
}

resource "aws_iam_role_policy_attachment" "apprunner_instance_role_attachment" {
  role       = aws_iam_role.apprunner_instance_role.name
  policy_arn = aws_iam_policy.apprunner_policy.arn
}

data "aws_iam_policy_document" "apprunner_instance_assume_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["tasks.apprunner.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "apprunner_instance_role_policy" {
  statement {
    actions = ["firehose:Put*"]
    effect  = "Allow"
    resources = [
      aws_kinesis_firehose_delivery_stream.buz_valid.arn,
      aws_kinesis_firehose_delivery_stream.buz_invalid.arn
    ]
  }

  statement {
    actions = ["s3:Get*"]
    effect  = "Allow"
    resources = [
      aws_s3_bucket.buz_schemas.arn,
      "${aws_s3_bucket.buz_schemas.arn}/*",
    ]
  }

  statement {
    actions = ["secretsmanager:GetSecretValue"]
    effect  = "Allow"
    resources = [
      aws_secretsmanager_secret.buz_config.arn
    ]

  }
}