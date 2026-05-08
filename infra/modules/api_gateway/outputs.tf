output "api_endpoint" {
  value       = aws_apigatewayv2_api.main.api_endpoint
  description = "API Gateway endpoint URL"
}

output "api_id" {
  value = aws_apigatewayv2_api.main.id
}

output "execution_arn" {
  value = aws_apigatewayv2_api.main.execution_arn
}

output "custom_domain_target" {
  description = "Route53 ALIAS 레코드 타겟 (커스텀 도메인 미사용 시 빈 문자열)"
  value       = var.api_domain != "" ? aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].target_domain_name : ""
}

output "custom_domain_zone_id" {
  description = "Route53 ALIAS 레코드용 hosted zone ID"
  value       = var.api_domain != "" ? aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].hosted_zone_id : ""
}
