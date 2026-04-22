# ── 미디어 버킷 ──────────────────────────────────────────────

resource "aws_s3_bucket" "media" {
    bucket = "yuno-media-${var.env}"
    tags = var.tags
}

resource "aws_s3_bucket_public_access_block" "media" {
    bucket = aws_s3_bucket.media.id
    block_public_acls = true
    block_public_policy = true
    ignore_public_acls = true
    restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "media_bucket" {
    bucket = aws_s3_bucket.media.id
    versioning_configuration {
        status = "Disabled"
    }
}

resource "aws_s3_bucket_cors_configuration" "media" {
    bucket = aws_s3_bucket.media.id
    cors_rule {
        allowed_headers = ["*"]
        allowed_methods = ["GET", "PUT", "POST", "HEAD"]
        allowed_origins = var.media_cors_origins
        expose_headers = ["ETag"]
        max_age_seconds = 3000
    }
}

# original/ 업로드 후 90일 → Glacier 이전 
resource "aws_s3_bucket_lifecycle_configuration" "media" {
    bucket = aws_s3_bucket.media.id
    rule {
        id = "archive-originals"
        status = "Enabled"
        filter {
            prefix = "original/"
        }
        transition {
            days = 90
            storage_class = "GLACIER"
        }
    }
}

# CloudFront에서 접근 가능하도록 정책 설정
resource "aws_s3_bucket_policy" "media" {
    bucket = aws_s3_bucket.media.id
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Principal = {
                    Service = "cloudfront.amazonaws.com"
                }
                Action = "s3:GetObject"
                Resource = "${aws_s3_bucket.media.arn}/*"
                Condition = {
                    StringEquals = {
                        "AWS:SourceArn" = var.media_distribution_arn
                    }
                }
            }
        ]
    })
}

# ── 프론트엔드 버킷 ──────────────────────────────────────────

resource "aws_s3_bucket" "frontend" {
    bucket = "yuno-frontend-${var.env}"
}

resource "aws_s3_bucket_public_access_block" "frontend" {
    bucket = aws_s3_bucket.frontend.id
    block_public_acls = true
    block_public_policy = true
    ignore_public_acls = true
    restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "frontend" {
    bucket = aws_s3_bucket.frontend.id
    versioning_configuration {
        status = "Disabled"
    }
}

resource "aws_s3_bucket_policy" "frontend" {
    bucket = aws_s3_bucket.frontend.id
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Principal = {
                    Service = "cloudfront.amazonaws.com"
                }
                Action = "s3:GetObject"
                Resource = "${aws_s3_bucket.frontend.arn}/*"
                Condition = {
                    StringEquals = {
                        "AWS:SourceArn" = var.frontend_distribution_arn
                    }
                }
            }
        ]
    })
}

# ── API빌드 결과물(Lambda Zip파일) 버킷 ────────────────────────────────────

resource "aws_s3_bucket" "lambda_zip_bucket" {
    bucket = "yuno-lambda-zip-${var.env}"
    tags = var.tags
}

resource "aws_s3_bucket_public_access_block" "lambda_zip_bucket" {
    bucket = aws_s3_bucket.lambda_zip_bucket.id
    block_public_acls = true
    block_public_policy = true
    ignore_public_acls = true
    restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "lambda_zip_bucket" {
    bucket = aws_s3_bucket.lambda_zip_bucket.id
    versioning_configuration {
        status = "Disabled"
    }
}
