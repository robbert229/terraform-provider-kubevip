package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kube-vip/kube-vip/pkg/kubevip"
)

func dataDaemonSetManifest() *schema.Resource {
	return &schema.Resource{
		Description: "`kubevip_daemonset_manifest` generates a static pod manifest for kubevip.",
		ReadContext: dataDaemonSetManifestRead,
		Schema: map[string]*schema.Schema{
			"interface": {
				Description: "Name of the interface to bind to",
				Required:    true,
				Type:        schema.TypeString,
			},
			"controlplane": {
				Description: "Enable HA for control plane",
				Optional:    true,
				Default:     false,
				Type:        schema.TypeBool,
			},
			"address": {
				Description: "an address (IP or DNS name) to use as a VIP",
				Type:        schema.TypeString,
				Required:    true,
			},
			"services": {
				Description: "Enable Kubernetes services",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"leader_election": {
				Description: "Use the Kubernetes leader election mechanism for clustering",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"arp": {
				Description: "Enable Arp for VIP changes",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"taint": {
				Description: "Taint the manifest for only running on control planes",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"in_cluster": {
				Description: "Use the incluster token to authenticate to Kubernetes",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},

			"raw_yaml": {
				Description: "The resulting yaml",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataDaemonSetManifestRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {

	initConfig := kubevip.Config{
		Interface:          d.Get("interface").(string),
		Address:            d.Get("address").(string),
		EnableControlPlane: d.Get("controlplane").(bool),
		EnableServices:     d.Get("services").(bool),
		LeaderElection: kubevip.LeaderElection{
			EnableLeaderElection: d.Get("leader_election").(bool),
		},
		EnableARP: d.Get("arp").(bool),
	}
	initLoadBalancer := kubevip.LoadBalancer{}

	initConfig.LoadBalancers = append(initConfig.LoadBalancers, initLoadBalancer)
	// TODO - A load of text detailing what's actually happening
	if err := kubevip.ParseEnvironment(&initConfig); err != nil {
		return diag.FromErr(fmt.Errorf("error parsing environment from config: %w", err))
	}

	// The control plane has a requirement for a VIP being specified
	if initConfig.EnableControlPlane && (initConfig.VIP == "" && initConfig.Address == "" && !initConfig.DDNS) {
		return diag.FromErr(fmt.Errorf("no address is specified for kube-vip to expose services on"))
	}

	inCluster := d.Get("in_cluster").(bool)
	taint := d.Get("taint").(bool)

	cfg := kubevip.GenerateDaemonsetManifestFromConfig(&initConfig, version, inCluster, taint)

	address := d.Get("address").(string)

	d.SetId(address)

	err := d.Set("raw_yaml", cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
