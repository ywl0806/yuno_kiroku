variable "env" {
    type = string
    description = "배포 환경"
}

variable "common_tags" {
    type = map(string)
    description = "태그"
    default = {}
}

variable "media_bucket_arn" {
    type = string
    description = "미디어 버킷 ARN"
}

variable "lambda_zip_bucket_arn" {
    type = string
    description = "Lambda Zip 버킷 ARN"
}

variable "ecr_repository_arns" {
    type = list(string)
    description = "ECR 레포지토리 ARN 목록"
}

variable "aws_region" {
    type = string
    description = "AWS 리전"
}

variable "aws_account_id" {
    type = string
    description = "AWS 계정 ID"
}