variable "env" {
    type = string
    description = "배포 환경"
}

variable "common_tags" {
    type = map(string)
    description = "태그"
    default = {}
}

variable "api_lambda_role_arn" {
    type = string
    description = "API Lambda Role ARN"
}

variable "resize_lambda_role_arn" {
    type = string
    description = "Resize Lambda Role ARN"
}

variable "resize_ecr_image_uri" {
    type = string
    description = "Resize Lambda ECR Image URI"
}

variable "media_bucket_arn" {
    type = string
    description = "Media Bucket ARN"
}

variable "media_bucket_id" {
    type = string
    description = "Media Bucket ID"
}

variable "api_lambda_zip_bucket_arn" {
    type = string
    description = "API Lambda Zip Bucket ARN"
}

variable "api_lambda_zip_bucket_name" {
    type = string
    description = "API Lambda Zip Bucket Name"
}

variable "api_lambda_zip_bucket_key" {
    type = string
    description = "API Lambda Zip Bucket Key"
}

variable "face_recognition_lambda_role_arn" {
    type        = string
    description = "Face Recognition Lambda Role ARN"
}

variable "face_recognition_ecr_image_uri" {
    type        = string
    description = "Face Recognition Lambda ECR Image URI"
}

variable "app_env_vars" {
    type        = map(string)
    description = "API Lambda 환경변수"
    sensitive   = true
    default     = {}
}