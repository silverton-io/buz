resource "aws_iam_role" "firehose" {
  name               = "${local.system_env_base}firehose"
  assume_role_policy = data.aws_iam_policy_document.firehose_assume_role.json
}

resource "aws_iam_role_policy" "firehose_bucket" {
  name   = "${local.system_env_base}firehose"
  role   = aws_iam_role.firehose.name
  policy = data.aws_iam_policy_document.firehose_bucket.json
}

# Lambda
resource "aws_iam_role" "lambda_role" {
  name               = "${local.system_env_base}lambda"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_policy.json
}

resource "aws_iam_policy" "lambda_policy" {
  name   = "${local.system_env_base}lambda"
  policy = data.aws_iam_policy_document.lambda_role_policy.json
}

resource "aws_iam_role_policy_attachment" "lambda_role_attachment" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}
