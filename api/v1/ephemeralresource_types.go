// api/v1/ephemeralresource_types.go
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EphemeralResourceSpec defines the desired state of EphemeralResource
type EphemeralResourceSpec struct {
	TTLSeconds int64 `json:"ttlSeconds,omitempty"`
	Resources  []Resource `json:"resources"`
}

// Resource defines a Kubernetes resource to be managed (e.g., Deployment, Secret)
type Resource struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// EphemeralResourceStatus defines the observed state of EphemeralResource
type EphemeralResourceStatus struct {
	// Add any observed state here.
	// Example:
	// LastChecked metav1.Time `json:"lastChecked,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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
