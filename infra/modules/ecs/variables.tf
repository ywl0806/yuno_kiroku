variable "env" {
  type        = string
  description = "배포 환경 (production, develop)"
}

variable "common_tags" {
  type = map(string)
  description = "태그"
  default = {}
}

variable "aws_region" {
  type = string
}

variable "ai_image_uri" {
  type        = string
  description = "ECS AI Task 컨테이너 이미지 URI"
}

variable "ecs_task_role_arn" {
  type        = string
  description = "ECS Task Role ARN"
}

variable "ecs_task_execution_role_arn" {
  type        = string
  description = "ECS Task Execution Role ARN"
}

variable "vpc_id" {
  type        = string
  description = "ECS Task가 실행될 VPC ID"
}

variable "face_recognition_queue_url" {
  type        = string
  description = "Face Recognition Queue URL"
}