data "aws_route53_zone" "main" {
  name = "REPLACE_ME"
}

data "aws_acm_certificate" "blog" {
  domain   = "REPLACE_ME"
  statuses = ["ISSUED"]
}

/*resource "aws_route53_record" "blog_api" {
  zone_id = "${data.aws_route53_zone.main.zone_id}"
  name    = "${aws_api_gateway_domain_name.blog.domain_name}"
  type    = "A"

  alias {
    name                   = "${aws_api_gateway_domain_name.blog.cloudfront_domain_name}"
    zone_id                = "${aws_api_gateway_domain_name.blog.cloudfront_zone_id}"
    evaluate_target_health = false
  }
}*/

