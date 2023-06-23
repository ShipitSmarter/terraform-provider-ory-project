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

data "orynetwork_identity" "example" {
    schema_id = "preset://email"
    traits = {
        "email" = "test@test.com"
    }
}