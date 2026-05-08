terraform {
  backend "s3" {
    bucket  = "terraform-state-bucket-yuno"
    key     = "shared/route53/terraform.tfstate"
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

provider "aws" {
  region = "ap-northeast-1"
}

variable "domain_name" {
  type        = string
  description = "구입한 도메인 (예: yuno.app)"
}
