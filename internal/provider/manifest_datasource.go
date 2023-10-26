package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kube-vip/kube-vip/pkg/kubevip"
)

var (
	_ datasource.DataSource = &manifestDataSource{}
)

func NewManifestDataSource() datasource.DataSource {
	return &manifestDataSource{}
}

type manifestDataSource struct{}

// Metadata returns the data source type name.
func (d *manifestDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manifest"
}

type manifestTypeValidator struct {
}

// Description implements validator.String.
func (manifestTypeValidator) Description(context.Context) string {
	return "Validates that the manifest type is either `pod` or `daemonset`"
}

// MarkdownDescription implements validator.String.
func (manifestTypeValidator) MarkdownDescription(context.Context) string {
	return "Validates that the manifest type is either `pod` or `daemonset`"
}

// ValidateString implements validator.String.
func (manifestTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	str := req.ConfigValue.ValueString()
	if str == "pod" || str == "daemonset" {
		return
	}

	res.Diagnostics.AddAttributeError(req.Path, "Invalid type given, valid values are `pod` and `daemonset`", "")
}

var _ validator.String = manifestTypeValidator{}

// Schema defines the schema for the data source.
func (d *manifestDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "`kubevip_manifest` generates a static pod, or daemonSet manifest for kubevip.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Required:    true,
				Description: "either `pod` or `daemonset`",
				Validators:  []validator.String{manifestTypeValidator{}},
			},

			"interface": schema.StringAttribute{
				Description: "Name of the interface to bind to",
				Required:    true,
			},
			"controlplane": schema.BoolAttribute{
				Description: "Enable HA for control plane",
				Optional:    true,
			},
			"address": schema.StringAttribute{
				Description: "an address (IP or DNS name) to use as a VIP",
				Required:    true,
			},
			"services": schema.BoolAttribute{
				Description: "Enable Kubernetes services",
				Optional:    true,
			},
			"leader_election": schema.BoolAttribute{
				Description: "Use the Kubernetes leader election mechanism for clustering",
				Optional:    true,
			},
			"arp": schema.BoolAttribute{
				Description: "Enable Arp for VIP changes",
				Optional:    true,
			},
			"in_cluster": schema.BoolAttribute{
				Description: "Use the incluster token to authenticate to Kubernetes",
				Optional:    true,
			},
			"taint": schema.BoolAttribute{
				Description: "Taint the manifest for only running on control planes",
				Optional:    true,
			},
			"raw_yaml": schema.StringAttribute{
				Description: "The resulting yaml",
				Computed:    true,
			},
		},
	}
}

type daemonsetManifestDataSourceModel struct {
	Interface          types.String `tfsdk:"interface"`
	Address            types.String `tfsdk:"address"`
	EnableControlPlane types.Bool   `tfsdk:"controlplane"`
	EnableServices     types.Bool   `tfsdk:"services"`
	LeaderElection     types.Bool   `tfsdk:"leader_election"`
	ARP                types.Bool   `tfsdk:"arp"`
	InCluster          types.Bool   `tfsdk:"in_cluster"`
	Taint              types.Bool   `tfsdk:"taint"`
	RawYAML            types.String `tfsdk:"raw_yaml"`
	Type               types.String `tfsdk:"type"`
}

// Read refreshes the Terraform state with the latest data.
func (d *manifestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state daemonsetManifestDataSourceModel

	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	{
		hasher := sha256.New()
		hasher.Write([]byte(state.Address.ValueString()))
		id := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		_ = id
		// state.ID = types.StringValue(id)
	}

	initConfig := &kubevip.Config{
		Interface:          state.Interface.ValueString(),
		Address:            state.Address.ValueString(),
		EnableControlPlane: state.EnableControlPlane.ValueBool(),
		EnableServices:     state.EnableServices.ValueBool(),
		LeaderElection: kubevip.LeaderElection{
			EnableLeaderElection: state.LeaderElection.ValueBool(),
			RetryPeriod:          1,
			RenewDeadline:        3,
			LeaseDuration:        5,
			LeaseName:            "plndr-cp-lock",
		},
		EnableARP:         state.ARP.ValueBool(),
		Port:              6443,
		VIPCIDR:           "32",
		Namespace:         "kube-system",
		ServicesLeaseName: "plndr-svcs-lock",
	}
	initLoadBalancer := kubevip.LoadBalancer{}

	initConfig.LoadBalancers = append(initConfig.LoadBalancers, initLoadBalancer)
	// TODO - A load of text detailing what's actually happening
	if err := kubevip.ParseEnvironment(initConfig); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Environment From Config",
			err.Error(),
		)
		return
	}

	// The control plane has a requirement for a VIP being specified
	if initConfig.EnableControlPlane && (initConfig.VIP == "" && initConfig.Address == "" && !initConfig.DDNS) {
		resp.Diagnostics.AddError(
			"Invalid Address",
			"no address is specified for kube-vip to expose services on",
		)
		return
	}

	var cfg string

	switch state.Type.ValueString() {
	case "daemonset":
		cfg = kubevip.GenerateDaemonsetManifestFromConfig(
			initConfig,
			kubevipVersion,
			state.InCluster.ValueBool(),
			state.Taint.ValueBool(),
		)
	case "pod":
		cfg = kubevip.GeneratePodManifestFromConfig(
			initConfig,
			kubevipVersion,
			state.InCluster.ValueBool(),
		)
	}

	state.RawYAML = types.StringValue(cfg)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
