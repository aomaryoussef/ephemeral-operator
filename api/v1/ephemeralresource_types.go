// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EphemeralResourceSpec defines the desired state of EphemeralResource
type EphemeralResourceSpec struct {
	// TTL in seconds. Resources exceeding this TTL will be deleted
	// +kubebuilder:validation:Minimum=0
	TTLSeconds int32 `json:"ttlSeconds"`

	// Resources defines the list of resources (like Deployments, Secrets) to monitor
	Resources []ResourceRef `json:"resources"`
}

// ResourceRef defines a reference to a Kubernetes resource
type ResourceRef struct {
	// The kind of the resource (e.g., Deployment, Secret)
	Kind string `json:"kind"`

	// The name of the resource
	Name string `json:"name"`
}

// EphemeralResourceStatus defines the observed state of EphemeralResource
type EphemeralResourceStatus struct {
	// Keep track of the time when the resource was created or last modified
	LastModified metav1.Time `json:"lastModified"`
}

// +kubebuilder:object:root=true

// EphemeralResource is the Schema for the ephemeralresources API
type EphemeralResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              EphemeralResourceSpec   `json:"spec,omitempty"`
	Status            EphemeralResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EphemeralResourceList contains a list of EphemeralResource
type EphemeralResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EphemeralResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EphemeralResource{}, &EphemeralResourceList{})
}
