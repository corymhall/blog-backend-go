/*data "aws_vpc" "default" {
  default = true
}

data "aws_subnet_ids" "default" {
  vpc_id = "${data.aws_vpc.default.id}"
}

resource "aws_lb" "blog" {
  name               = "${var.project}-lb"
  internal           = false
  load_balancer_type = "application"

  subnets = [
    "${data.aws_subnet_ids.default.ids}",
  ]
}


resource "aws_lb_listener" "blog" {
  load_balancer_arn = "${aws_lb.blog.arn}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-Ext-2018-06"
  certificate_arn   = "${data.aws_acm_certificate.blog.arn}"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.blog.arn}"
  }
}

resource "aws_lb_target_group" "blog" {
  name        = "${var.project}-lambda-tg"
  target_type = "lambda"
}

resource "aws_lb_target_group_attachment" "blog_lambda" {
  target_group_arn = "${aws_lb_target_group.blog.arn}"
  target_id        = "${aws_lambda_function.blog.arn}"
  depends_on       = ["aws_lambda_permission.blog_with_lb"]
}*/

