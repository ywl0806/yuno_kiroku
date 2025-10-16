terraform {
  backend "s3" {
    bucket         = "terraform-state-bucket-yuno"
    key            = "terraform.tfstate"
    region         = "ap-northeast-1"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.9.0"
    }
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
      version = "1.14.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

provider "mongodbatlas" {
  public_key  = var.mongodb_atlas_public_key
  private_key = var.mongodb_atlas_private_key
}