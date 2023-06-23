terraform {
  required_providers {
    hashicups = {
      source = "registry.terraform.io/ShipitSmarter/ory-network"
    }
  }
}

provider "hashicups" {}

data "hashicups_coffees" "example" {}