data "terraform_remote_state" "route53" {
  backend = "s3"
  config = {
    bucket = "terraform-state-bucket-yuno"
    key    = "shared/route53/terraform.tfstate"
    region = "ap-northeast-1"
  }
}

# ── CloudFront용 ACM (us-east-1 필수) ──────────────────────────

resource "aws_acm_certificate" "cloudfront" {
  provider                  = aws.us_east_1
  domain_name               = var.domain_name
  subject_alternative_names = ["*.${var.domain_name}"]
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "cloudfront_validation" {
  for_each = {
    for dvo in aws_acm_certificate.cloudfront.domain_validation_options : dvo.resource_record_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }...
  }

  zone_id = data.terraform_remote_state.route53.outputs.hosted_zone_id
  name    = each.value[0].name
  type    = each.value[0].type
  records = [each.value[0].record]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "cloudfront" {
  provider                = aws.us_east_1
  certificate_arn         = aws_acm_certificate.cloudfront.arn
  validation_record_fqdns = [for r in aws_route53_record.cloudfront_validation : r.fqdn]
}

# ── API Gateway용 ACM (ap-northeast-1) ─────────────────────────

resource "aws_acm_certificate" "api_gw" {
  domain_name       = "api.${var.domain_name}"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "api_gw_validation" {
  for_each = {
    for dvo in aws_acm_certificate.api_gw.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  zone_id = data.terraform_remote_state.route53.outputs.hosted_zone_id
  name    = each.value.name
  type    = each.value.type
  records = [each.value.record]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "api_gw" {
  certificate_arn         = aws_acm_certificate.api_gw.arn
  validation_record_fqdns = [for r in aws_route53_record.api_gw_validation : r.fqdn]
}
