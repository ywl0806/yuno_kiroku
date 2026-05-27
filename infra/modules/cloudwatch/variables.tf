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

variable "video_queue_name" {
  type        = string
  description = "Video Processing Queue Name"
}

variable "video_task_scale_out_policy_arn" {
  type        = string
  description = "Video Task Scale Out Policy ARN"
}

variable "video_task_scale_in_policy_arn" {
  type        = string
  description = "Video Task Scale In Policy ARN"
}

variable "alert_email" {
  type        = string
  description = "알림 수신 이메일 주소"
}

variable "resize_dlq_name" {
  type        = string
  description = "Resize DLQ 이름"
}

variable "face_recognition_dlq_name" {
  type        = string
  description = "Face Recognition DLQ 이름"
}

variable "video_dlq_name" {
  type        = string
  description = "Video Processing DLQ 이름"
}

variable "api_log_group_name" {
  type        = string
  description = "API Lambda CloudWatch 로그 그룹 이름"
}