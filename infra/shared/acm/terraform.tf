terraform {
  backend "s3" {
    bucket  = "terraform-state-bucket-yuno"
    key     = "shared/acm/terraform.tfstate"
    region  = "ap-northeast-1"
    encrypt = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.9"
    }
  }
}

# 기본 리전 (API Gateway 인증서용)
provider "aws" {
  region = "ap-northeast-1"
}

# CloudFront 인증서는 반드시 us-east-1
provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"
}

variable "domain_name" {
  type        = string
  description = "기본 도메인"
}
