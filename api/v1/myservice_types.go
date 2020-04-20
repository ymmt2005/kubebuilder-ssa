/*
Copyright 2020 Hirotaka Yamamoto.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServicePort is ...
type ServicePort struct {
	// Port number.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port uint16 `json:"port"`

	// Protocol name.
	// +kubebuilder:validation:Enum=TCP;UDP
	Protocol string `json:"protocol"`

	// Target port number.  If not given, the target is the same as port.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	TargetPort *uint16 `json:"targetPort,omitempty"`
}

// MyServiceSpec defines the desired state of MyService
type MyServiceSpec struct {
	String  string  `json:"string,omitempty"`
	Pointer *string `json:"pointer,omitempty"`

	// The list of ports that are exposed by this service.
	// +listType=map
	// +listMapKey=port
	// +listMapKey=protocol
	Ports []ServicePort `json:"ports,omitempty"`
}

// MyServiceStatus defines the observed state of MyService
type MyServiceStatus struct {
	Count     int32        `json:"count"`
	Timestamp *metav1.Time `json:"timestamp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MyService is the Schema for the myservices API
type MyService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyServiceSpec   `json:"spec,omitempty"`
	Status MyServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MyServiceList contains a list of MyService
type MyServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyService{}, &MyServiceList{})
}
