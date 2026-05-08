output "media_bucket_id" {
    value = aws_s3_bucket.media.id
}

output "frontend_bucket_id" {
    value = aws_s3_bucket.frontend.id
}

output "media_bucket_domain" {
    value = aws_s3_bucket.media.bucket_regional_domain_name
}

output "frontend_bucket_domain" {
    value = aws_s3_bucket.frontend.bucket_regional_domain_name
}

output "lambda_zip_bucket_arn" {
    value = aws_s3_bucket.lambda_zip_bucket.arn
}

output "lambda_zip_bucket_name" {
    value = aws_s3_bucket.lambda_zip_bucket.id
}

output "media_bucket_arn" {
    value = aws_s3_bucket.media.arn
}

output "media_bucket_name" {
    value = aws_s3_bucket.media.id
}