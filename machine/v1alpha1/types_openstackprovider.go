/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenstackProviderSpec is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an OpenStack Instance. It is used by the Openstack machine actuator to create a single machine instance.
// TODO(cglaubitz): We might consider to change this to OpenstackMachineProviderSpec
// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=2
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type OpenstackProviderSpec struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// The name of the secret containing the openstack credentials
	CloudsSecret *corev1.SecretReference `json:"cloudsSecret"`

	// The name of the cloud to use from the clouds secret
	CloudName string `json:"cloudName"`

	// The flavor reference for the flavor for your server instance.
	Flavor string `json:"flavor"`

	// The name of the image to use for your server instance.
	// If the RootVolume is specified, this will be ignored and use rootVolume directly.
	Image string `json:"image"`

	// The ssh key to inject in the instance
	// +optional
	KeyName string `json:"keyName,omitempty"`

	// The machine ssh username
	// +optional
	SshUserName string `json:"sshUserName,omitempty"`

	// A networks object. Required parameter when there are multiple networks defined for the tenant.
	// When you do not specify the networks parameter, the server attaches to the only network created for the current tenant.
	// +optional
	Networks []NetworkParam `json:"networks,omitempty"`

	// Create and assign additional ports to instances
	// +optional
	Ports []PortOpts `json:"ports,omitempty"`

	// +optional
	FloatingIP string `json:"floatingIP,omitempty"`

	// The availability zone from which to launch the server.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The names of the security groups to assign to the instance
	// +optional
	SecurityGroups []SecurityGroupParam `json:"securityGroups,omitempty"`

	// The name of the secret containing the user data (startup script in most cases)
	// +optional
	UserDataSecret *corev1.SecretReference `json:"userDataSecret,omitempty"`

	// Whether the server instance is created on a trunk port or not.
	// +optional
	Trunk bool `json:"trunk,omitempty"`

	// Machine tags
	// Requires Nova api 2.52 minimum!
	// +optional
	Tags []string `json:"tags,omitempty"`

	// Metadata mapping. Allows you to create a map of key value pairs to add to the server instance.
	// +optional
	ServerMetadata map[string]string `json:"serverMetadata,omitempty"`

	// Config Drive support
	// +optional
	ConfigDrive *bool `json:"configDrive,omitempty"`

	// The volume metadata to boot from
	// +optional
	RootVolume *RootVolume `json:"rootVolume,omitempty"`

	// The server group to assign the machine to.
	// +optional
	ServerGroupID string `json:"serverGroupID,omitempty"`

	// The server group to assign the machine to. A server group with that
	// name will be created if it does not exist. If both ServerGroupID and
	// ServerGroupName are non-empty, they must refer to the same OpenStack
	// resource.
	// +optional
	ServerGroupName string `json:"serverGroupName,omitempty"`

	// The subnet that a set of machines will get ingress/egress traffic from
	// +optional
	PrimarySubnet string `json:"primarySubnet,omitempty"`
}

type SecurityGroupParam struct {
	// Security Group UID
	// +optional
	UUID string `json:"uuid,omitempty"`
	// Security Group name
	// +optional
	Name string `json:"name,omitempty"`
	// Filters used to query security groups in openstack
	// +optional
	Filter SecurityGroupFilter `json:"filter,omitempty"`
}

type SecurityGroupFilter struct {
	// +optional
	ID string `json:"id,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	Description string `json:"description,omitempty"`
	// +optional
	TenantID string `json:"tenantId,omitempty"`
	// +optional
	ProjectID string `json:"projectId,omitempty"`
	// +optional
	Limit int `json:"limit,omitempty"`
	// +optional
	Marker string `json:"marker,omitempty"`
	// +optional
	SortKey string `json:"sortKey,omitempty"`
	// +optional
	SortDir string `json:"sortDir,omitempty"`
	// +optional
	Tags string `json:"tags,omitempty"`
	// +optional
	TagsAny string `json:"tagsAny,omitempty"`
	// +optional
	NotTags string `json:"notTags,omitempty"`
	// +optional
	NotTagsAny string `json:"notTagsAny,omitempty"`
}

type NetworkParam struct {
	// The UUID of the network. Required if you omit the port attribute.
	// +optional
	UUID string `json:"uuid,omitempty"`
	// A fixed IPv4 address for the NIC.
	// +optional
	FixedIp string `json:"fixedIp,omitempty"`
	// Filters for optional network query
	// +optional
	Filter Filter `json:"filter,omitempty"`
	// Subnet within a network to use
	// +optional
	Subnets []SubnetParam `json:"subnets,omitempty"`
	// NoAllowedAddressPairs disables creation of allowed address pairs for the network ports
	// +optional
	NoAllowedAddressPairs bool `json:"noAllowedAddressPairs,omitempty"`
	// PortTags allows users to specify a list of tags to add to ports created in a given network
	// +optional
	PortTags []string `json:"portTags,omitempty"`
	// VNICType sets the type of the OpenStack network port that will be created
	// +optional
	VNICType string `json:"vnicType,omitempty"`
	// Profile allows setting port binding-profile during creation
	// +optional
	Profile map[string]string `json:"profile,omitempty"`
	// PortSecurity optionally enables or disables security on ports managed by OpenStack
	// +optional
	PortSecurity *bool `json:"portSecurity,omitempty"`
}

type Filter struct {
	// +optional
	Status string `json:"status,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	Description string `json:"description,omitempty"`
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`
	// +optional
	TenantID string `json:"tenantId,omitempty"`
	// +optional
	ProjectID string `json:"projectId,omitempty"`
	// +optional
	Shared *bool `json:"shared,omitempty"`
	// +optional
	ID string `json:"id,omitempty"`
	// +optional
	Marker string `json:"marker,omitempty"`
	// +optional
	Limit int `json:"limit,omitempty"`
	// +optional
	SortKey string `json:"sortKey,omitempty"`
	// +optional
	SortDir string `json:"sortDir,omitempty"`
	// +optional
	Tags string `json:"tags,omitempty"`
	// +optional
	TagsAny string `json:"tagsAny,omitempty"`
	// +optional
	NotTags string `json:"notTags,omitempty"`
	// +optional
	NotTagsAny string `json:"notTagsAny,omitempty"`
}

type SubnetParam struct {
	// The UUID of the network. Required if you omit the port attribute.
	// +optional
	UUID string `json:"uuid,omitempty"`

	// Filters for optional network query
	// +optional
	Filter SubnetFilter `json:"filter,omitempty"`

	// PortTags are tags that are added to ports created on this subnet
	// +optional
	PortTags []string `json:"portTags,omitempty"`

	// PortSecurity optionally enables or disables security on ports managed by OpenStack
	// +optional
	PortSecurity *bool `json:"portSecurity,omitempty"`
}

type SubnetFilter struct {
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	Description string `json:"description,omitempty"`
	// +optional
	EnableDHCP *bool `json:"enableDhcp,omitempty"`
	// +optional
	NetworkID string `json:"networkId,omitempty"`
	// +optional
	TenantID string `json:"tenantId,omitempty"`
	// +optional
	ProjectID string `json:"projectId,omitempty"`
	// +optional
	IPVersion int `json:"ipVersion,omitempty"`
	// +optional
	GatewayIP string `json:"gateway_ip,omitempty"`
	// +optional
	CIDR string `json:"cidr,omitempty"`
	// +optional
	IPv6AddressMode string `json:"ipv6AddressMode,omitempty"`
	// +optional
	IPv6RAMode string `json:"ipv6RaMode,omitempty"`
	// +optional
	ID string `json:"id,omitempty"`
	// +optional
	SubnetPoolID string `json:"subnetpoolId,omitempty"`
	// +optional
	Limit int `json:"limit,omitempty"`
	// +optional
	Marker string `json:"marker,omitempty"`
	// +optional
	SortKey string `json:"sortKey,omitempty"`
	// +optional
	SortDir string `json:"sortDir,omitempty"`
	// +optional
	Tags string `json:"tags,omitempty"`
	// +optional
	TagsAny string `json:"tagsAny,omitempty"`
	// +optional
	NotTags string `json:"notTags,omitempty"`
	// +optional
	NotTagsAny string `json:"notTagsAny,omitempty"`
}

type PortOpts struct {
	// ID of the OpenStack network on which to create the port
	NetworkID string `json:"networkID" required:"true"`
	// Used to make the name of the port unique
	NameSuffix string `json:"nameSuffix" required:"true"`
	// +optional
	Description string `json:"description,omitempty"`
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`
	// +optional
	MACAddress string `json:"macAddress,omitempty"`
	// Specify pairs of subnet and/or IP address. These should be subnets of the network with the given NetworkID.
	// +optional
	FixedIPs []FixedIPs `json:"fixedIPs,omitempty"`
	// +optional
	TenantID string `json:"tenantID,omitempty"`
	// +optional
	ProjectID string `json:"projectID,omitempty"`
	// The names of the security groups to assign to the port
	// +optional
	SecurityGroups *[]string `json:"securityGroups,omitempty"`
	// +optional
	AllowedAddressPairs []AddressPair `json:"allowedAddressPairs,omitempty"`
	// +optional
	Tags []string `json:"tags,omitempty"`

	// The ID of the host where the port is allocated
	// +optional
	HostID string `json:"hostID,omitempty"`

	// The virtual network interface card (vNIC) type that is bound to the
	// neutron port.
	// +optional
	VNICType string `json:"vnicType,omitempty"`

	// A dictionary that enables the application running on the specified
	// host to pass and receive virtual network interface (VIF) port-specific
	// information to the plug-in.
	// +optional
	Profile map[string]string `json:"profile,omitempty"`

	// enable or disable security on a given port
	// incompatible with securityGroups and allowedAddressPairs
	// +optional
	PortSecurity *bool `json:"portSecurity,omitempty"`

	// Enables and disables trunk at port level. If not provided, openStackMachine.Spec.Trunk is inherited.
	// +optional
	Trunk *bool `json:"trunk,omitempty"`
}

type AddressPair struct {
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`
	// +optional
	MACAddress string `json:"macAddress,omitempty"`
}

type FixedIPs struct {
	SubnetID string `json:"subnetID"`
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`
}

type RootVolume struct {
	// +optional
	SourceType string `json:"sourceType,omitempty"`
	// +optional
	SourceUUID string `json:"sourceUUID,omitempty"`
	DeviceType string `json:"deviceType"`
	// +optional
	VolumeType string `json:"volumeType,omitempty"`
	// +optional
	Size int `json:"diskSize,omitempty"`
	// +optional
	Zone string `json:"availabilityZone,omitempty"`
}

// OpenstackClusterProviderSpec is the providerSpec for OpenStack in the cluster object
// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=2
type OpenstackClusterProviderSpec struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// NodeCIDR is the OpenStack Subnet to be created. Cluster actuator will create a
	// network, a subnet with NodeCIDR, and a router connected to this subnet.
	// If you leave this empty, no network will be created.
	// +optional
	NodeCIDR string `json:"nodeCidr,omitempty"`
	// DNSNameservers is the list of nameservers for OpenStack Subnet being created.
	// +optional
	DNSNameservers []string `json:"dnsNameservers,omitempty"`
	// ExternalNetworkID is the ID of an external OpenStack Network. This is necessary
	// to get public internet to the VMs.
	// +optional
	ExternalNetworkID string `json:"externalNetworkId,omitempty"`

	// ManagedSecurityGroups defines that kubernetes manages the OpenStack security groups
	// for now, that means that we'll create two security groups, one allowing SSH
	// and API access from everywhere, and another one that allows all traffic to/from
	// machines belonging to that group. In the future, we could make this more flexible.
	ManagedSecurityGroups bool `json:"managedSecurityGroups"`

	// Tags for all resources in cluster
	// +optional
	Tags []string `json:"tags,omitempty"`

	// Default: True. In case of server tag errors, set to False
	// +optional
	DisableServerTags bool `json:"disableServerTags,omitempty"`
}

// OpenstackClusterProviderStatus contains the status fields
// relevant to OpenStack in the cluster object.
// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=2
type OpenstackClusterProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Network contains all information about the created OpenStack Network.
	// It includes Subnets and Router.
	// +optional
	Network *Network `json:"network,omitempty"`

	// ControlPlaneSecurityGroups contains all the information about the OpenStack
	// Security Group that needs to be applied to control plane nodes.
	// TODO: Maybe instead of two properties, we add a property to the group?
	// +optional
	ControlPlaneSecurityGroup *SecurityGroup `json:"controlPlaneSecurityGroup,omitempty"`

	// GlobalSecurityGroup contains all the information about the OpenStack Security
	// Group that needs to be applied to all nodes, both control plane and worker nodes.
	// +optional
	GlobalSecurityGroup *SecurityGroup `json:"globalSecurityGroup,omitempty"`
}

// Network represents basic information about the associated OpenStack Neutron Network
type Network struct {
	Name string `json:"name"`
	ID   string `json:"id"`

	// +optional
	Subnet *Subnet `json:"subnet,omitempty"`
	// +optional
	Router *Router `json:"router,omitempty"`
}

// Subnet represents basic information about the associated OpenStack Neutron Subnet
type Subnet struct {
	Name string `json:"name"`
	ID   string `json:"id"`

	CIDR string `json:"cidr"`
}

// Router represents basic information about the associated OpenStack Neutron Router
type Router struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// SecurityGroup represents the basic information of the associated
// OpenStack Neutron Security Group.
type SecurityGroup struct {
	Name  string              `json:"name"`
	ID    string              `json:"id"`
	Rules []SecurityGroupRule `json:"rules"`
}

// SecurityGroupRule represent the basic information of the associated OpenStack
// Security Group Role.
type SecurityGroupRule struct {
	ID              string `json:"name"`
	Direction       string `json:"direction"`
	EtherType       string `json:"etherType"`
	SecurityGroupID string `json:"securityGroupID"`
	PortRangeMin    int    `json:"portRangeMin"`
	PortRangeMax    int    `json:"portRangeMax"`
	Protocol        string `json:"protocol"`
	RemoteGroupID   string `json:"remoteGroupID"`
	RemoteIPPrefix  string `json:"remoteIPPrefix"`
}

// Equal checks if two SecurityGroupRules are the same.
func (r SecurityGroupRule) Equal(x SecurityGroupRule) bool {
	return (r.Direction == x.Direction &&
		r.EtherType == x.EtherType &&
		r.PortRangeMin == x.PortRangeMin &&
		r.PortRangeMax == x.PortRangeMax &&
		r.Protocol == x.Protocol &&
		r.RemoteGroupID == x.RemoteGroupID &&
		r.RemoteIPPrefix == x.RemoteIPPrefix)

}
