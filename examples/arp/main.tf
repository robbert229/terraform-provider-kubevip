terraform {
  required_providers {
    kubevip = {
      source = "robbert229/kubevip"
    }
  }
}

provider "kubevip" {

}


# deploy arp manifest

data "kubevip_pod_manifest" "manifest" {
  interface       = "ens192"
  address         = "192.168.0.40"
  controlplane    = true
  services        = true
  leader_election = true
  arp             = true
}
