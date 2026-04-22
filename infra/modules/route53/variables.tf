variable "hosted_zone_id" {
  type        = string
  description = "Route53 hosted zone ID"
}

variable "domain_name" {
  type        = string
  description = "기본 도메인 (예: yuno.app)"
}

variable "frontend_cf_domain" {
  type        = string
  description = "프론트엔드 CloudFront distribution domain"
}

variable "media_cf_domain" {
  type        = string
  description = "미디어 CloudFront distribution domain"
}

variable "api_gw_domain_target" {
  type        = string
  description = "API Gateway 커스텀 도메인 타겟"
}

variable "api_gw_domain_zone_id" {
  type        = string
  description = "API Gateway 커스텀 도메인의 Route53 hosted zone ID"
}
