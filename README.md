# kubevip

This is a terraform provider that generates [kube-vip](https://kube-vip.io/) manifests. It relies on using the same source code that kube-vip's command line interface uses, but is wrapped in a terraform provider. This makes it easier to manage lots of clusters with `kube-vip`. It's currently used in production environments, though I haven't tested many use-cases outside of the happy path.

## Example

In the following example we are generating a kube-vip daemonset manifest.

```terraform
data "kubevip_manifest" "k8sc1" {
  type            = "daemonset"
  interface       = "eth0"
  address         = "192.168.103.9"
  controlplane    = true
  arp             = true
  leader_election = true
  in_cluster      = true
  taint           = true
}
```
