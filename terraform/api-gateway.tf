/*resource "aws_api_gateway_rest_api" "blog" {
  name               = "pp-blog-api"
  description        = "Pleasant Places Blog API"
  binary_media_types = ["**"]
}

resource "aws_api_gateway_deployment" "blog_deployment" {
  depends_on = ["aws_api_gateway_integration.blog"]

  rest_api_id = "${aws_api_gateway_rest_api.blog.id}"
  stage_name  = "production"
}

resource "aws_api_gateway_base_path_mapping" "blog" {
  api_id      = "${aws_api_gateway_rest_api.blog.id}"
  stage_name  = "${aws_api_gateway_deployment.blog_deployment.stage_name}"
  domain_name = "${aws_api_gateway_domain_name.blog.domain_name}"
}

resource "aws_api_gateway_method" "blog" {
  rest_api_id   = "${aws_api_gateway_rest_api.blog.id}"
  resource_id   = "${aws_api_gateway_resource.blog_proxy_resource.id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_resource" "blog_proxy_resource" {
  rest_api_id = "${aws_api_gateway_rest_api.blog.id}"
  parent_id   = "${aws_api_gateway_rest_api.blog.root_resource_id}"
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_domain_name" "blog" {
  certificate_arn = "REPLACE_ME"
  domain_name     = "REPLACE_ME"
}

resource "aws_api_gateway_integration" "blog" {
  rest_api_id             = "${aws_api_gateway_rest_api.blog.id}"
  resource_id             = "${aws_api_gateway_resource.blog_proxy_resource.id}"
  http_method             = "${aws_api_gateway_method.blog.http_method}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.blog.arn}/invocations"
}*/

