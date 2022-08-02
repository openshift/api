/*
Copyright 2022.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// IngressNodeFirewallICMPRule define ingress node firewall rule for ICMP and ICMPv6 protocols
type IngressNodeFirewallICMPRule struct {
	// imcpType defines ICMP Type Numbers (RFC 792).
	// if configured, this field matches against the ICMP/ICMPv6 header otherwise its ignored.
	// +optional
	// +kubebuilder:validation:Maximum:=255
	// +kubebuilder:validation:Minimum:=0
	ICMPType uint8 `json:"icmpType,omitempty"`

	// icmpCode defines ICMP Code ID (RFC 792).
	// if configured, this field matches against the ICMP/ICMPv6 header otherwise its ignored.
	// +optional
	// +kubebuilder:validation:Maximum:=255
	// +kubebuilder:validation:Minimum:=0
	ICMPCode uint8 `json:"icmpCode,omitempty"`
}

// IngressNodeFirewallProtoRule define ingress node firewall rule for TCP, UDP and SCTP protocols
type IngressNodeFirewallProtoRule struct {
	// ports defines either a single port or a range of ports to apply a protocol rule too.
	// To filter a single port, set a single port as an integer value. For example ports: 80.
	// To filter a range of ports, use a "start-end" range, string format. For example ports: "80-100".
	// +optional
	Ports intstr.IntOrString `json:"ports,omitempty"`
}

// IngressNodeProtocolConfig is a discriminated union of protocol's specific configuration.
// +union
type IngressNodeProtocolConfig struct {
	// protocol can be ICMP, ICMPv6, TCP, SCTP or UDP.
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum="ICMP";"ICMPv6";"TCP";"UDP";"SCTP"
	Protocol IngressNodeFirewallRuleProtocolType `json:"protocol"`

	// tcp defines an ingress node firewall rule for TCP protocol.
	// +optional
	TCP *IngressNodeFirewallProtoRule `json:"tcp,omitempty"`

	// udp defines an ingress node firewall rule for UDP protocol.
	// +optional
	UDP *IngressNodeFirewallProtoRule `json:"udp,omitempty"`

	// sctp defines an ingress node firewall rule for SCTP protocol.
	// +optional
	SCTP *IngressNodeFirewallProtoRule `json:"sctp,omitempty"`

	// icmp defines an ingress node firewall rule for ICMP protocol.
	// +optional
	ICMP *IngressNodeFirewallICMPRule `json:"icmp,omitempty"`

	// icmpv6 defines an ingress node firewall rule for ICMPv6 protocol.
	// +optional
	ICMPv6 *IngressNodeFirewallICMPRule `json:"icmpv6,omitempty"`
}

// IngressNodeFirewallProtocolRule defines an ingress node firewall rule per protocol.
type IngressNodeFirewallProtocolRule struct {
	// order defines the order of execution of ingress firewall rules.
	// The minimum order value is 1 and the values must be unique.
	// + index 0 is used internally as catch all for unclassified packets matching the same sourceCIDR.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum:=1
	Order uint32 `json:"order"`

	// protocolConfig is a discriminated union of a protocol's specific configuration for TCP, UDP, SCTP, ICMP and ICMPv6.
	// If not specified, packet matching will be based on the protocol value and protocol configuration, such as dstPort/type/code, will be ignored
	// +optional
	ProtocolConfig IngressNodeProtocolConfig `json:"protocolConfig"`

	// action can be Allow or Deny, default action is Deny.
	// +optional
	Action IngressNodeFirewallActionType `json:"action,omitempty"`
}

// ProtocolType defines the protocol types that are supported
type IngressNodeFirewallRuleProtocolType string

const (
	// ProtocolTypeICMP refers to the ICMP protocol.
	ProtocolTypeICMP IngressNodeFirewallRuleProtocolType = "ICMP"

	// ProtocolTypeICMP6 refers to the ICMPv6 protocol.
	ProtocolTypeICMP6 IngressNodeFirewallRuleProtocolType = "ICMPv6"

	// ProtocolTypeTCP refers to the TCP protocol, for either IPv4 or IPv6.
	ProtocolTypeTCP IngressNodeFirewallRuleProtocolType = "TCP"

	// ProtocolTypeUDP refers to the UDP protocol, for either IPv4 or IPv6.
	ProtocolTypeUDP IngressNodeFirewallRuleProtocolType = "UDP"

	// ProtocolTypeSCTP refers to the SCTP protocol, for either IPv4 or IPv6.
	ProtocolTypeSCTP IngressNodeFirewallRuleProtocolType = "SCTP"
)

// IngressNodeFirewallActionType indicates whether an IngressNodeFirewallRule allows or denies traffic.
// +kubebuilder:validation:Enum="Allow";"Deny"
type IngressNodeFirewallActionType string

const (
	IngressNodeFirewallAllow IngressNodeFirewallActionType = "Allow"
	IngressNodeFirewallDeny  IngressNodeFirewallActionType = "Deny"
)

// IngressNodeFirewallRules define ingress node firewall rule.
type IngressNodeFirewallRules struct {
	// sourceCIDRs defines the origin of packets that FirewallProtocolRules will be applied to.
	// +kubebuilder:validation:MinItems:=1
	SourceCIDRs []string `json:"sourceCIDRs"`
	// rules is a list of per protocol ingress node firewall rules.
	// +listType:=map
	// +listMapKey:=order
	FirewallProtocolRules []IngressNodeFirewallProtocolRule `json:"rules,omitempty"`
}

// IngressNodeFirewallSpec defines the desired state of IngressNodeFirewall.
type IngressNodeFirewallSpec struct {
	// nodeSelector Selects node(s) where ingress firewall rules will be applied to.
	// +optional
	NodeSelector metav1.LabelSelector `json:"nodeSelector,omitempty"`

	// ingress is a list of ingress firewall policy rules.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	Ingress []IngressNodeFirewallRules `json:"ingress,omitempty"`

	// interfaces is a list of interfaces where the ingress firewall policy will be applied on.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	Interfaces []string `json:"interfaces,omitempty"`
}

// IngressNodeFirewallStatus defines the observed state of IngressNodeFirewall.
type IngressNodeFirewallStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// IngressNodeFirewall is the Schema for the ingressnodefirewalls API.
type IngressNodeFirewall struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngressNodeFirewallSpec   `json:"spec,omitempty"`
	Status IngressNodeFirewallStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IngressNodeFirewallList contains a list of IngressNodeFirewall.
type IngressNodeFirewallList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngressNodeFirewall `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IngressNodeFirewall{}, &IngressNodeFirewallList{})
}
