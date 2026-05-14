locals {
  common_tags = {
    Project     = "yuno"
    Environment = var.env
    ManagedBy   = "terraform"
  }
}

# ── CloudWatch Log Groups ────────────────────────────────────

resource "aws_cloudwatch_log_group" "api" {
  name              = "/aws/lambda/yuno-api-${var.env}"
  retention_in_days = 30
  tags              = local.common_tags
}

resource "aws_cloudwatch_log_group" "resize" {
  name              = "/aws/lambda/yuno-resize-${var.env}"
  retention_in_days = 30
  tags              = local.common_tags
}

resource "aws_cloudwatch_log_group" "face_recognition" {
  name              = "/aws/lambda/yuno-face-recognition-${var.env}"
  retention_in_days = 30
  tags              = local.common_tags
}

# ── API Lambda ───────────────────────────────────────────────

resource "aws_lambda_function" "api" {
  function_name = "yuno-api-${var.env}"
  role          = var.api_lambda_role_arn
  package_type  = "Zip"
  depends_on = [aws_cloudwatch_log_group.api]
  s3_bucket = var.api_lambda_zip_bucket_name
  s3_key = var.api_lambda_zip_bucket_key
  architectures = ["arm64"]

  runtime = "provided.al2023"
  handler = "bootstrap"

  environment {
    variables = var.app_env_vars
  }

  tags = local.common_tags
}

# ── Resize Lambda ────────────────────────────────────────────

resource "aws_lambda_function" "resize" {
  function_name = "yuno-resize-${var.env}"
  role          = var.resize_lambda_role_arn
  package_type  = "Image"
  image_uri     = var.resize_ecr_image_uri
  architectures = ["arm64"]
  memory_size   = 2048
  timeout       = 600

  environment {
    variables = var.app_env_vars
  }

  depends_on = [aws_cloudwatch_log_group.resize]

  tags = local.common_tags

}

# ── S3 → Resize Lambda 이벤트 트리거 ────────────────────────

resource "aws_lambda_permission" "allow_s3_resize" {
  statement_id  = "AllowS3InvokeResize"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.resize.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = var.media_bucket_arn
}

resource "aws_s3_bucket_notification" "media_put" {
  bucket = var.media_bucket_id

  lambda_function {
    lambda_function_arn = aws_lambda_function.resize.arn
    events              = ["s3:ObjectCreated:Put"]
    filter_prefix       = "*/original/"
  }

  depends_on = [aws_lambda_permission.allow_s3_resize]
}
