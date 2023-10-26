package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Schema is a struct to represent the configuration of the provider
type Schema struct {
}

// providerSchema returns the provider schema
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
