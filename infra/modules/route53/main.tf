# 프론트엔드: yuno.app → CloudFront
resource "aws_route53_record" "frontend" {
  zone_id = var.hosted_zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = var.frontend_cf_domain
    zone_id                = "Z2FDTNDATAQYW2" # CloudFront 고정 hosted zone ID
    evaluate_target_health = false
  }
}

# 미디어: media.yuno.app → CloudFront
resource "aws_route53_record" "media" {
  zone_id = var.hosted_zone_id
  name    = "media.${var.domain_name}"
  type    = "A"

  alias {
    name                   = var.media_cf_domain
    zone_id                = "Z2FDTNDATAQYW2"
    evaluate_target_health = false
  }
}

# API: api.yuno.app → API Gateway 커스텀 도메인
resource "aws_route53_record" "api" {
  zone_id = var.hosted_zone_id
  name    = "api.${var.domain_name}"
  type    = "A"

  alias {
    name                   = var.api_gw_domain_target
    zone_id                = var.api_gw_domain_zone_id
    evaluate_target_health = false
  }
}
