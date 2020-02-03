terraform {
  backend "s3" {
    bucket = "REPLACE_ME"
    key    = "terraform.tfstate"
    region = "us-east-2"
  }
}
