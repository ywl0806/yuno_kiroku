variable "env" {
  type        = string
  description = "배포 환경 (production, develop)"
}

variable "aws_region" {
  type = string
}

variable "aws_account_id" {
  type = string
}

variable "api_lambda_invoke_arn" {
  type = string
}

variable "api_lambda_function_name" {
  type = string
}

variable "api_domain" {
  type        = string
  description = "API 커스텀 도메인 (예: api.yuno.app). 빈 문자열이면 커스텀 도메인 미사용."
  default     = ""
}

variable "acm_certificate_arn" {
  type        = string
  description = "API Gateway 커스텀 도메인용 ACM 인증서 ARN (ap-northeast-1)"
  default     = ""
}
