variable "env" {
    type = string
    description = "배포 환경"
}

variable "common_tags" {
    type = map(string)
    description = "태그"
    default = {}
}

variable "media_bucket_id" {
    type = string
    description = "미디어 버킷 ID"
}

variable "media_bucket_domain" {
    type = string
    description = "미디어 버킷 도메인"
}

variable "frontend_bucket_id" {
    type = string
    description = "프론트엔드 버킷 ID"
}

variable "frontend_bucket_domain" {
    type = string
    description = "프론트엔드 버킷 도메인"
}

variable "acm_certificate_arn" {
    type = string
    description = "ACM 인증서 ARN"
    default = ""
}

variable "media_domain" {
    type = string
    description = "미디어 CDN 도메인"
    default = ""
}

variable "cloudfront_public_key_pem" {
    type        = string
    description = "CloudFront signed cookie용 RSA 공개키 PEM (미설정 시 signed cookie 비활성화)"
    default     = ""
    sensitive   = true
}

variable "frontend_domain" {
    type = string
    description = "프론트엔드 도메인"
    default = ""
}