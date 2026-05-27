variable "env" {
  type        = string
  description = "배포 환경"
}

variable "common_tags" {
  type        = map(string)
  description = "버킷 태그"
  default     = {}
}

variable "media_bucket_arn" {
  type        = string
  description = "미디어 버킷 ARN (resize 큐 정책 condition용)"
}
