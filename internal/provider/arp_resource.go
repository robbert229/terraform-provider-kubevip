package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	resourceARPManifest = "kubevip_arp_manifest"
)

func buildResourceARPManifest() *schema.Resource {
	return &schema.Resource{
		Create: resourceARPManifestCreate,
	}
}

func resourceARPManifestCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceARPManifestRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceARPManifestUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceARPManifestDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
