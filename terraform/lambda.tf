/*resource "aws_lambda_function" "blog" {
  s3_bucket     = "${aws_s3_bucket.blog_lambda_deployments.id}"
  s3_key        = "${var.project}-deployment.zip"
  function_name = "${var.project}-lambda"
  handler       = "server"
  role          = "${aws_iam_role.lambda.arn}"
  runtime       = "go1.x"
}*/

resource "aws_cloudwatch_log_group" "blog" {
  name              = "/aws/lambda/pp-blog-api"
  retention_in_days = 14
}

resource "aws_s3_bucket" "blog_lambda_deployments" {
  bucket = "${var.project}-lambda-deployments-${var.region}"
  acl    = "private"
}

/*resource "aws_cloudwatch_log_group" "blog" {
  name              = "/aws/lambda/${aws_lambda_function.blog.function_name}"
  retention_in_days = 14
}

resource "aws_lambda_permission" "blog_with_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.blog.arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.blog.id}/*"
}*/

resource "aws_iam_role" "lambda" {
  name               = "${var.project}-lambda"
  assume_role_policy = "${data.aws_iam_policy_document.lambda_assume.json}"
}

data "aws_iam_policy_document" "lambda_assume" {
  statement {
    effect = "Allow"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"

      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
  }
}

data "aws_iam_policy_document" "lambda" {
  statement {
    effect = "Allow"

    actions = [
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:PutItem",
    ]

    resources = [
      "${aws_dynamodb_table.posts.arn}",
      "${aws_dynamodb_table.user.arn}",
      "${aws_dynamodb_table.comments.arn}",
      "${aws_dynamodb_table.reply.arn}",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "arn:aws:logs:*:*:*",
    ]
  }
}

resource "aws_iam_role_policy" "dynamo" {
  role   = "${aws_iam_role.lambda.id}"
  policy = "${data.aws_iam_policy_document.lambda.json}"
  name   = "${var.project}-lambda-dynamo"
}
