output "media_distribution_domain" {
    value = aws_cloudfront_distribution.media.domain_name
}

output "media_signing_key_pair_id" {
    description = "CloudFront signed cookie용 public key ID (CLOUDFRONT_KEY_PAIR_ID 환경변수에 설정)"
    value       = local.enable_signed_cookies ? aws_cloudfront_public_key.media_signing[0].id : ""
}

output "media_distribution_id" {
    value = aws_cloudfront_distribution.media.id
}

output "media_distribution_arn" {
    value = aws_cloudfront_distribution.media.arn
}

output "frontend_distribution_domain" {
    value = aws_cloudfront_distribution.frontend.domain_name
}

output "frontend_distribution_id" {
    value = aws_cloudfront_distribution.frontend.id
}

output "frontend_distribution_arn" {
    value = aws_cloudfront_distribution.frontend.arn
}