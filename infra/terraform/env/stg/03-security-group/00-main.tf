terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = var.region

  access_key = var.access_key
  secret_key = var.secret_key
}

terraform {
  backend "s3" {
    region                  = "ap-northeast-1"
    bucket                  = "furumaru-terraform"
    key                     = "stg/security-group.tfstate"
    profile                 = "default"
    shared_credentials_file = "~/.aws/credentials"
  }
}
