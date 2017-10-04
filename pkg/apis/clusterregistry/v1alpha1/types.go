/*
Copyright 2017 The Kubernetes Authors.

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
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Information about a cluster in a cluster registry.
type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// State is the state of the cluster. This is not a specification, and is not
	// meant to be used by actively-reconciling controllers; it is also not
	// a status, as it contains fields that do not necessarily describe the
	// status of the cluster, and is not necessarily updated by an active
	// controller.
	// +optional
	State ClusterState `json:"state,omitempty" protobuf:"bytes,1,opt,name=state"`
}

// ClusterState contains the state of a cluster.
type ClusterState struct {

	// KubeAPIServer represents the API server for this cluster.
	// +optional
	KubeAPIServer KubeAPIServer `json:"kubeApiServer,omitempty" protobuf:"bytes,1,opt,name=kubeApiServer"`

	// AuthInfo contains public information that can be used to authenticate
	// to and authorize with this cluster.
	// +optional
	AuthInfo AuthInfo `json:"authInfo,omitempty" protobuf:"bytes,2,opt,name=authInfo"`

	// CloudProvider contains information about the cloud provider this cluster
	// is running on.
	// +optional
	CloudProvider CloudProvider `json:"cloudProvider,omitempty" protobuf:"bytes,3,opt,name=cloudProvider"`
}

type URL string

// KubeAPIServer represents one and only one Kubernetes API server.
type KubeAPIServer struct {
	// Servers specifies the addresses of the Kubernetes(s) API server’s network
	// identity or identities. They can be any valid HTTP URL including the
	// IP:Port combination or the host name.
	Servers []URL `json:"servers" protobuf:"bytes,1,rep,name=servers"`

	// CABundle is the certificate authority information.
	// +optional
	CABundle []byte `json:"caBundle,omitempty" protobuf:"bytes,2,rep,name=caBundle"`
}

// AuthInfo holds public information that describes how a client can get
// credentials to access the cluster. For example, OAuth2 client registration
// endpoints and supported flows, Kerberos servers locations, etc.
//
// It should not hold any private or sensitive information.
type AuthInfo struct {

	// AuthProviders is a list of configurations for auth providers.
	AuthProviders []AuthProviderConfig `json:"authProviders" protobuf:"bytes,1,rep,name=authProviders"`
}

// AuthProviderConfig contains the information necessary for a client to
// authenticate with a Kubernetes API server. It is modeled after
// k8s.io/client-go/tools/clientcmd/api/v1.AuthProviderConfig.
type AuthProviderConfig struct {
	// Name is the name of this configuration.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Config is a map of values that contains the information necessary for a
	// client to authenticate to a Kubernetes API server.
	// +optional
	Config map[string]string `json:"config,omitempty" protobuf:"bytes,2,rep,name=config"`
}

// CloudProvider contains information about the cloud provider this cluster is
// running on.
type CloudProvider struct {
	// Name is the name of the cloud provider for this cluster.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A list of Kubernetes clusters in the cluster registry.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of Cluster objects.
	Items []Cluster `json:"items" protobuf:"bytes,2,rep,name=items"`
}
