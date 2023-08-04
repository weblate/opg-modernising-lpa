resource "aws_s3_bucket_notification" "bucket_notification" {
  count  = var.enable_autoscan ? 1 : 0
  bucket = var.data_store_bucket.id

  lambda_function {
    id                  = "bucket-av-scan"
    lambda_function_arn = aws_lambda_function.lambda_function.arn
    events              = ["s3:ObjectCreated:Put"]
  }

  depends_on = [
    aws_lambda_permission.allow_bucket_to_run
  ]
  provider = aws.region
}