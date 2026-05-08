locals {
    is_custom_domain = var.acm_certificate_arn != ""
}

# ── 미디어 CloudFront ──────────────────────────────────────────

resource "aws_cloudfront_origin_access_control" "media" {
    name = "yuno-media-oac-${var.env}"
    origin_access_control_origin_type = "s3"
    signing_behavior = "always"
    signing_protocol = "sigv4"
}

resource "aws_cloudfront_cache_policy" "media_long_ttl" {
    name = "yuno-media-long-ttl-${var.env}"
    min_ttl = 0
    default_ttl = 86400
    max_ttl = 31536000

    parameters_in_cache_key_and_forwarded_to_origin{
        cookies_config {
            cookie_behavior = "none"
        }
        headers_config {
            header_behavior = "none"
        }
        query_strings_config {
            query_string_behavior = "none"
        }
    }
}


resource "aws_cloudfront_distribution" "media" {
    origin {
        domain_name = var.media_bucket_domain
        origin_id = "S3-media-${var.env}"
        origin_access_control_id = aws_cloudfront_origin_access_control.media.id
    }

    enabled = true
    is_ipv6_enabled = true
    comment = "yuno-media-${var.env}"

    aliases = local.is_custom_domain && var.media_domain != "" ? [var.media_domain] : []

    default_cache_behavior {
        allowed_methods = ["GET", "HEAD"]
        cached_methods = ["GET", "HEAD"]
        target_origin_id = "S3-media-${var.env}"
        viewer_protocol_policy = "redirect-to-https"
        cache_policy_id = aws_cloudfront_cache_policy.media_long_ttl.id
    }
    
    # 커스텀 도메인 사용 시 ACM 인증서 사용
    dynamic "viewer_certificate" {
        for_each = local.is_custom_domain ? [1] : []
        content {
            acm_certificate_arn = var.acm_certificate_arn
            ssl_support_method = "sni-only"
            minimum_protocol_version = "TLSv1.2_2021"
        }
    }
    
    # 커스텀 도메인 사용 안 할 시 CloudFront 기본 인증서 사용
    dynamic "viewer_certificate" {
        for_each = local.is_custom_domain ? [] : [1]
        content {
            cloudfront_default_certificate = true
        }
    }

    # 한국, 일본만 접근 가능
    restrictions {
        geo_restriction {
            restriction_type = "whitelist"
            locations = ["KR", "JP"]
        }
    }

    tags = var.common_tags
}

# ── 프론트엔드 CloudFront ──────────────────────────────────────────

resource "aws_cloudfront_origin_access_control" "frontend" {
    name = "yuno-frontend-oac-${var.env}"
    origin_access_control_origin_type = "s3"
    signing_behavior = "always"
    signing_protocol = "sigv4"
}

resource "aws_cloudfront_distribution" "frontend" {
    origin {
        domain_name = var.frontend_bucket_domain
        origin_id = "S3-frontend-${var.env}"
        origin_access_control_id = aws_cloudfront_origin_access_control.frontend.id
    }

    enabled = true
    is_ipv6_enabled = true
    comment = "yuno-frontend-${var.env}"

    aliases = local.is_custom_domain && var.frontend_domain != "" ? [var.frontend_domain] : []

    default_cache_behavior {
        allowed_methods = ["GET", "HEAD", "OPTIONS"]
        cached_methods = ["GET", "HEAD"]
        target_origin_id = "S3-frontend-${var.env}"
        viewer_protocol_policy = "redirect-to-https"
        cache_policy_id = aws_cloudfront_cache_policy.media_long_ttl.id
    }
    
    dynamic "viewer_certificate" {
        for_each = local.is_custom_domain ? [1] : []
        content {
            acm_certificate_arn = var.acm_certificate_arn
            ssl_support_method = "sni-only"
            minimum_protocol_version = "TLSv1.2_2021"
        }
    }
    
    dynamic "viewer_certificate" {
        for_each = local.is_custom_domain ? [] : [1]
        content {
            cloudfront_default_certificate = true
        }
    }
    
    custom_error_response {
        error_caching_min_ttl = 60
        error_code         = 403
        response_code      = 200
        response_page_path = "/index.html"
    }

    custom_error_response {
        error_caching_min_ttl = 60
        error_code         = 404
        response_code      = 200
        response_page_path = "/index.html"
    }

    // 한국, 일본만 접근 가능
    restrictions {
        geo_restriction {
            restriction_type = "whitelist"
            locations = ["KR", "JP"]
        }
    }

    tags = var.common_tags
}