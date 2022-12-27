terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "us-east-1" # only us-east-1

  access_key = var.access_key
  secret_key = var.secret_key
}

terraform {
  backend "s3" {
    region                  = "ap-northeast-1"
    bucket                  = "furumaru-terraform"
    key                     = "shared/route53.tfstate"
    profile                 = "default"
    shared_credentials_file = "~/.aws/credentials"
  }
}
