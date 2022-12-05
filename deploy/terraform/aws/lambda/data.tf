data "aws_caller_identity" "current" {}
data "aws_region" "current" {}


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

data "aws_iam_policy_document" "lambda_assume_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "lambda_role_policy" {
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
    effect = "Allow"
    actions = ["logs:CreateLogGroup"]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}"
    ]
  }

  statement {
    effect = "Allow"
    actions = ["logs:CreateLogStream", "logs:PutLogEvents"]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:log-group:/aws/lambda/${local.service_name}:*"
    ]
  }

}
