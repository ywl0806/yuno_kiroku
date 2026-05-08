output "media_distribution_domain" {
    value = aws_cloudfront_distribution.media.domain_name
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