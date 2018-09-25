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
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster contains information about a cluster in a cluster registry.
// +k8s:openapi-gen=x-kubernetes-print-columns:custom-columns=NAME:.metadata.name,CIDR:.spec.kubernetesApiEndpoints.serverEndpoints[].clientCIDR,SERVER:.spec.kubernetesApiEndpoints.serverEndpoints[].serverAddress,CREATION TIME:.metadata.creationTimestamp
// +resource:path=clusters
type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec is the specification of the cluster. This may or may not be
	// reconciled by an active controller.
	// +optional
	Spec ClusterSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status is the status of the cluster.
	// +optional
	Status ClusterStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ClusterSpec contains the specification of a cluster.
type ClusterSpec struct {
	// KubernetesAPIEndpoints represents the endpoints of the API server for this
	// cluster.
	// +optional
	KubernetesAPIEndpoints KubernetesAPIEndpoints `json:"kubernetesApiEndpoints,omitempty" protobuf:"bytes,1,opt,name=kubernetesApiEndpoints"`

	// AuthInfo contains public information that can be used to authenticate
	// to and authorize with this cluster. It is not meant to store private
	// information (e.g., tokens or client certificates) and cluster registry
	// implementations are not expected to provide hardened storage for
	// secrets.
	// +optional
	AuthInfo AuthInfo `json:"authInfo,omitempty" protobuf:"bytes,2,opt,name=authInfo"`
}

// ClusterStatus contains the status of a cluster.
type ClusterStatus struct {
	// Conditions contains the different condition statuses for this cluster.
	Conditions []ClusterCondition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`

	// TODO https://github.com/kubernetes/cluster-registry/issues/28
}

// KubernetesAPIEndpoints represents the endpoints for one and only one
// Kubernetes API server.
type KubernetesAPIEndpoints struct {
	// ServerEndpoints specifies the address(es) of the Kubernetes API server’s
	// network identity or identities.
	// +optional
	ServerEndpoints []ServerAddressByClientCIDR `json:"serverEndpoints,omitempty" protobuf:"bytes,1,rep,name=serverEndpoints"`

	// CABundle contains the certificate authority information.
	// +optional
	CABundle []byte `json:"caBundle,omitempty" protobuf:"bytes,2,opt,name=caBundle"`
}

// ServerAddressByClientCIDR helps clients determine the server address that
// they should use, depending on the ClientCIDR that they match.
type ServerAddressByClientCIDR struct {
	// The CIDR with which clients can match their IP to figure out if they should
	// use the corresponding server address.
	// +optional
	ClientCIDR string `json:"clientCIDR,omitempty" protobuf:"bytes,1,opt,name=clientCIDR"`
	// Address of this server, suitable for a client that matches the above CIDR.
	// This can be a hostname, hostname:port, IP or IP:port.
	// +optional
	ServerAddress string `json:"serverAddress,omitempty" protobuf:"bytes,2,opt,name=serverAddress"`
}

// AuthInfo holds information that describes how a client can get
// credentials to access the cluster. For example, OAuth2 client registration
// endpoints and supported flows, or Kerberos server locations.
type AuthInfo struct {
	// User references an object that contains implementation-specific details
	// about how a user should authenticate against this cluster.
	// +optional
	User *ObjectReference `json:"user,omitempty" protobuf:"bytes,1,opt,name=user"`

	// Controller references an object that contains implementation-specific
	// details about how a controller should authenticate. A simple use case for
	// this would be to reference a secret in another namespace that stores a
	// bearer token that can be used to authenticate against this cluster's API
	// server.
	Controller *ObjectReference `json:"controller,omitempty" protobuf:"bytes,2,opt,name=controller"`
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// Kind contains the kind of the referent, e.g., Secret or ConfigMap
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// Name contains the name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`

	// Namespace contains the namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

// ClusterConditionType marks the kind of cluster condition being reported.
type ClusterConditionType string

const (
	// ClusterOK means that the cluster is "OK".
	//
	// Since the cluster registry does not have a standard status controller, the
	// meaning of this condition is defined by the environment in which the
	// cluster is running. It is expected to mean that the cluster is reachable by
	// a controller that is reporting on its status, and that the cluster is ready
	// to have workloads scheduled.
	ClusterOK ClusterConditionType = "OK"
)

// ClusterCondition contains condition information for a cluster.
type ClusterCondition struct {
	// Type is the type of the cluster condition.
	Type ClusterConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ClusterConditionType"`

	// Status is the status of the condition. One of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`

	// LastHeartbeatTime is the last time this condition was updated.
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty" protobuf:"bytes,3,opt,name=lastHeartbeatTime"`

	// Last time the condition was checked.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`

	// LastTransitionTime is the last time the condition changed from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`

	// Reason is a (brief) reason for the condition's last status change.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`

	// Message is a human-readable message indicating details about the last status change.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCredentials combines a Cluster record with a set of credentials
// (stored in a kubeconfig) that can be used to access the cluster. The status
// of ClusterCredentials reflects whether those credentials can be used - from
// within the cluster hosting the ClusterCredentials resource - to access the
// cluster's healthz endpoint.
//
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=clustercredentials
// +kubebuilder:subresource:status
type ClusterCredentials struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterCredentialsSpec   `json:"spec,omitempty"`
	Status ClusterCredentialsStatus `json:"status,omitempty"`
}

// ClusterCredentialsSpec defines the desired state of a ClusterCredentials.
type ClusterCredentialsSpec struct {
	// ClusterRef is a reference to a Cluster to use.
	//
	// An empty namespace field indicates a reference to a Cluster in the
	// enclosing namespace.
	ClusterRef v1.ObjectReference `json:"clusterRef,omitempty"`

	// SecretRef is a reference to the secret in the namespace that contains a
	// kubeconfig with credentials access the referenced cluster. The secret
	// must have a key named "kubeconfig" that contains a file formatted as a
	// proper subset of the kubeconfig file format. The data must contain a
	// kubeconfig file with only the following fields:
	//
	// - a 'users' section
	// - a single, valid entry in the 'users' section containing the auth
	//   information to use the referenced Cluster.
	//
	// TODO: explain which auth methods are supported
	// TODO: explain kubeconfig version
	//
	// This can be left empty if the cluster allows insecure access.
	// +optional
	SecretRef *apiv1.LocalObjectReference `json:"secretRef,omitempty"`
}

// ClusterCredentialsStatus contains information about the current status of a
// ClusterCredentials.
type ClusterCredentialsStatus struct {
	// Region is designation for a geographic area that contains a Cluster.
	Region string `json:"region,omitempty"`

	// AvailabilityZone is a distinct physical site within a single Region where
	// a Cluster may be deployed.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// Conditions is an array of current conditions.
	// +optional
	Conditions []ClusterCredentialsCondition `json:"conditions,omitempty"`
}

// ClusterCredentialsCondition contains condition information about a ClusterCredentials.
type ClusterCredentialsCondition struct {
	// Type is the type of the cluster condition.
	Type ClusterConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ClusterCredentialsConditionType"`

	// Status is the status of the condition. One of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`

	// LastHeartbeatTime is the last time this condition was updated.
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty" protobuf:"bytes,3,opt,name=lastHeartbeatTime"`

	// Last time the condition was checked.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`

	// LastTransitionTime is the last time the condition changed from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`

	// Reason is a (brief) reason for the condition's last status change.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`

	// Message is a human-readable message indicating details about the last status change.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// ClusterCredentialsConditionType marks the kind of ClusterCredentials condition
// being reported.
type ClusterCredentialsConditionType string

const (
	// ClusterCredentialsHealthy represents that a ClusterCredentials is healthy
	// as determined by the Cluster's healthz endpoint.
	ClusterCredentialsHealthy ClusterConditionType = "Healthy"
)
