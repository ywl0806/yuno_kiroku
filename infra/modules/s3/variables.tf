variable "env" {
  type        = string
  description = "배포 환경"
}

variable "media_cors_origins" {
  type        = list(string)
  description = "미디어 버킷 CORS 허용 origin 목록"
  default     = ["*"]
}

variable "common_tags" {
  type        = map(string)
  description = "버킷 태그"
  default     = {}
}

variable "frontend_distribution_arn" {
  type        = string
  description = "프론트엔드 CloudFront 배포 ARN"
  default     = ""
}

variable "media_distribution_arn" {
  type        = string
  description = "미디어 CloudFront 배포 ARN"
  default     = ""
}