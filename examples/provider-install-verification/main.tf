terraform {
  required_providers {
    orynetwork = {
      source = "registry.terraform.io/ShipitSmarter/ory-network"
    }
  }
}

provider "orynetwork" {
    api_key = ""
    host = ""
}

data "orynetwork_identities" "example" {}


output "identities" {
  value = data.orynetwork_identities.example
}
