terraform {
  backend "s3" {
    bucket  = "terraform-state-bucket-yuno"
    key     = "environments/production/terraform.tfstate"
    region  = "ap-northeast-1"
    encrypt = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.40"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}
