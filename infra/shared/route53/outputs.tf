output "hosted_zone_id" {
  value = aws_route53_zone.main.zone_id
}

output "name_servers" {
  description = "도메인 등록기관에 입력할 NS 레코드"
  value       = aws_route53_zone.main.name_servers
}
