variable "env" {
  type        = string
  description = "배포 환경"
}

variable "common_tags" {
  type        = map(string)
  description = "버킷 태그"
  default     = {}
}
