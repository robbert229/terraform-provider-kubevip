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

resource "kubevip_arp_manifest" "manifest" {
  interface       = "ens192"
  address         = "192.168.0.40"
  controlplan     = true
  services        = true
  leader-election = true
}
