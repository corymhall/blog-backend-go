resource "aws_s3_bucket" "tfstate" {
  bucket = "${var.project}-terraform-tfstate-${var.region}"
  acl    = "private"

  tags {
    Name        = "${var.project}-terraform-tfstate-${var.region}"
    Provisioner = "Terraform"
  }
}
