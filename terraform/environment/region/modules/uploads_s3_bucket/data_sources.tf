
data "aws_region" "current" {
  provider = aws.region
}

data "aws_caller_identity" "current" {
  provider = aws.region
}

data "aws_default_tags" "current" {
  provider = aws.region
}

data "aws_kms_alias" "source_default_key" {
  name     = "alias/aws/s3"
  provider = aws.region
}

data "aws_s3_bucket" "access_log" {
  bucket   = "s3-access-logs-${data.aws_default_tags.current.tags.application}-${data.aws_default_tags.current.tags.account-name}-${data.aws_region.current.name}"
  provider = aws.region
}