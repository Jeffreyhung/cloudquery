terraform {
  backend "s3" {
    bucket = "cq-plugins-source-aws-tf"
    key    = "iot"
    region = "us-east-1"
  }
}
