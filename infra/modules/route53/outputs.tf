output "frontend_fqdn" {
  value = aws_route53_record.frontend.fqdn
}

output "api_fqdn" {
  value = aws_route53_record.api.fqdn
}

output "media_fqdn" {
  value = aws_route53_record.media.fqdn
}
