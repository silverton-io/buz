
data "aws_iam_policy_document" "firehose_assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["firehose.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "firehose_bucket" {
  statement {
    effect = "Allow"

    actions = [
      "s3:AbortMultipartUpload",
      "s3:GetBucketLocation",
      "s3:GetObject",
      "s3:ListBucket",
      "s3:ListBucketMultipartUploads",
      "s3:PutObject",
    ]

    resources = [
      aws_s3_bucket.events.arn,
      "${aws_s3_bucket.events.arn}/*",
    ]
  }
}

resource "aws_iam_role" "firehose" {
  name               = "${local.system_env_base}firehose"
  assume_role_policy = data.aws_iam_policy_document.firehose_assume_role.json
}

resource "aws_iam_role_policy" "firehose_bucket" {
  name   = "${local.system_env_base}firehose"
  role   = aws_iam_role.firehose.name
  policy = data.aws_iam_policy_document.firehose_bucket.json
}

# App Runner Service 
resource "aws_iam_role" "apprunner" {
  name               = "${local.system_env_base}apprunner"
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

resource "aws_iam_role_policy_attachment" "apprunner_attachment" {
  role       = aws_iam_role.apprunner.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess"
}

# App Runner Instance 
resource "aws_iam_role" "apprunner_instance_role" {
  name               = "${local.system_env_base}apprunner-instance"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.apprunner_instance_assume_policy.json
}

resource "aws_iam_policy" "apprunner_policy" {
  name   = "${local.system_env_base}apprunner"
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

}
