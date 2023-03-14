/*
Copyright 2023 Juned Memon.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MicroServiceSpec defines the desired state of MicroService
type MicroServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Owner    string `json:"owner,omitempty"`
	Template string `json:"template"`
}

// MicroServiceStatus defines the observed state of MicroService
type MicroServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Github Repository of the Service.
	Repository string `json:"repository,omitempty"`
	// Registry  of the Service.
	Registry string `json:"registry,omitempty"`
	// Pipeline of the Service.
	Pipeline string `json:"pipeline,omitempty"`
	// Conditions holds the conditions for the HelmRelease.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Failures is the reconciliation failure count against the latest desired
	// state. It is reset after a successful reconciliation.
	// +optional
	Failures int64 `json:"failures,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MicroService is the Schema for the microservices API
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:printcolumn:name=age,type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name=owner,type=string,JSONPath=.spec.owner
// +kubebuilder:printcolumn:name=repository,type=string,JSONPath=.status.repository
// +kubebuilder:printcolumn:name=registry,type=string,JSONPath=.status.registry
// +kubebuilder:printcolumn:name=pipeline,type=string,JSONPath=.status.pipeline
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
type MicroService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MicroServiceSpec   `json:"spec,omitempty"`
	Status MicroServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MicroServiceList contains a list of MicroService
type MicroServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroService{}, &MicroServiceList{})
}
