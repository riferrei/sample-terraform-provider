terraform {
  required_version = ">= 1.3.4"
  required_providers {
    buildonaws = {
      source  = "aws.amazon.com/terraform/buildonaws"
      version = "1.0"
    }
  }
}
