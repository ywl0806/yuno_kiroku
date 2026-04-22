output "cloudfront_certificate_arn" {
  description = "CloudFront용 ACM 인증서 ARN (us-east-1)"
  value       = aws_acm_certificate_validation.cloudfront.certificate_arn
}

output "api_gw_certificate_arn" {
  description = "API Gateway용 ACM 인증서 ARN (ap-northeast-1)"
  value       = aws_acm_certificate_validation.api_gw.certificate_arn
}
