
variable "domain_name" {
  type        = string
  description = "구입한 도메인 (예: yuno.app). 설정하면 frontend/api/media 서브도메인 자동 구성."
  default     = ""
}

variable "database_url" {
  type      = string
  sensitive = true
}

variable "auth_secret_key" {
  type      = string
  sensitive = true
}

variable "line_channel_id" {
  type = string
}

variable "line_channel_secret" {
  type      = string
  sensitive = true
}

variable "kakao_client_id" {
  type = string
}

variable "kakao_client_secret" {
  type      = string
  sensitive = true
}

variable "alert_email" {
  type = string
}

variable "cloudfront_public_key_pem" {  
  type      = string
  sensitive = true
}

variable "cloudfront_private_key" {
  type      = string
  sensitive = true
}