terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}


provider "aws" {
  region = "eu-west-2"
}

variable "cluster_name" {
    type = string
}

module "cluster" {
  source = "./cluster"
  cluster_name = var.cluster_name
}

module "producer" {
  source = "./producer"
  bootstrap_server_url = module.cluster.bootstrap_server_url
}

module "consumers" {
  source = "./consumers"
  bootstrap_server_url = module.cluster.bootstrap_server_url
}