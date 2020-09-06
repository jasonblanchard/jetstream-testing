terraform {
  required_providers {
    jetstream = {
      source = "nats-io/jetstream"
      version = "0.0.18"
    }
  }
}

provider "jetstream" {
  servers = "localhost:4222"
}
