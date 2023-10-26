terraform {
  required_providers {
    kubevip = {
      source = "registery.terraform.io/robbert229/kubevip"
    }
  }
}

provider "kubevip" {

}


# deploy arp manifest

data "kubevip_manifest" "manifest" {
  type            = "daemonset"
  interface       = "ens192"
  address         = "192.168.0.40"
  controlplane    = true
  services        = true
  leader_election = true
  arp             = true
  in_cluster      = true
}

output "manifest" {
  value = data.kubevip_manifest.manifest.raw_yaml
}
