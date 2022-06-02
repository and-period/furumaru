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
    key                     = "stg/vpc.tfstate"
    profile                 = "default"
    shared_credentials_file = "~/.aws/credentials"
  }
}

module "vpc" {
  source = "./../../../modules/vpc"

  #####################################################################
  # Common
  #####################################################################
  tags = var.tags

  #####################################################################
  # VPC
  #####################################################################
  name                             = var.vpc_name
  cidr_block                       = "10.110.0.0/16"
  assign_generated_ipv6_cidr_block = false
  instance_tenancy                 = "default"
  enable_dns_support               = true
  enable_dns_hostnames             = true

  #####################################################################
  # Route Tables
  #####################################################################
  route_tables = [
    "furumaru-stg-rt-pub",
    "furumaru-stg-rt-pri-app",
    "furumaru-stg-rt-pri-db",
  ]

  #####################################################################
  # Subnets
  #####################################################################
  subnets = [
    // public subnet
    {
      name                    = "furumaru-stg-sub-pub-1a"
      availability_zone       = "ap-northeast-1a"
      cidr_block              = "10.110.0.0/24"
      map_public_ip_on_launch = true
      route_table_name        = "furumaru-stg-rt-pub"
    },
    {
      name                    = "furumaru-stg-sub-pub-1c"
      availability_zone       = "ap-northeast-1c"
      cidr_block              = "10.110.1.0/24"
      map_public_ip_on_launch = true
      route_table_name        = "furumaru-stg-rt-pub"
    },
    // private subnet (application)
    {
      name                    = "furumaru-stg-sub-pri-app-1a"
      availability_zone       = "ap-northeast-1a"
      cidr_block              = "10.110.10.0/24"
      map_public_ip_on_launch = false
      route_table_name        = "furumaru-stg-rt-pri-app"
    },
    {
      name                    = "furumaru-stg-sub-pri-app-1c"
      availability_zone       = "ap-northeast-1c"
      cidr_block              = "10.110.11.0/24"
      map_public_ip_on_launch = false
      route_table_name        = "furumaru-stg-rt-pri-app"
    },
    // private subnet (database)
    {
      name                    = "furumaru-stg-sub-pri-db-1a"
      availability_zone       = "ap-northeast-1a"
      cidr_block              = "10.110.20.0/24"
      map_public_ip_on_launch = false
      route_table_name        = "furumaru-stg-rt-pri-db"
    },
    {
      name                    = "furumaru-stg-sub-pri-db-1c"
      availability_zone       = "ap-northeast-1c"
      cidr_block              = "10.110.21.0/24"
      map_public_ip_on_launch = false
      route_table_name        = "furumaru-stg-rt-pri-db"
    },
  ]

  #####################################################################
  # Internet Gateway
  #####################################################################
  create_internet_gateway = true
  internet_gateway_name   = "furumaru-stg-igw"

  #####################################################################
  # NAT Gateway
  #####################################################################
  nat_gateways = [
    {
      name        = "furumaru-stg-ngw"
      subnet_name = "furumaru-stg-sub-pub-1a"
      eip_name    = "furumaru-stg-eip-ngw"
    },
  ]
}
