terraform {
  required_providers {
    ory-network  = {
      source = "registry.terraform.io/ShipitSmarter/ory-network"
    }
  }
}

provider "ory-network" {
    api_key = ""
    host = ""
}

data "ory-network_identities" "example" {}


output "identities" {
  value = data.ory-network_identities.example
}