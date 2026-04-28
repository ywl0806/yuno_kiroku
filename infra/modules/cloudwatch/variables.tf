variable "env" {
  type = string
  description = "배포 환경"
}

variable "common_tags" {
  type = map(string)
  description = "태그"
  default = {}
}

variable "ai_task_scale_out_policy_arn" {
  type = string
  description = "AI Task Scale Out ARN"
}

variable "ai_task_scale_in_policy_arn" {
  type = string
  description = "AI Task Scale In ARN"
}

variable "face_recognition_queue_name" {
  type = string
  description = "Face Recognition Queue Name"
}