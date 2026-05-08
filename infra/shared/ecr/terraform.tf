terraform {
  backend "s3" {
    bucket  = "terraform-state-bucket-yuno"
    key     = "shared/ecr/terraform.tfstate"
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
