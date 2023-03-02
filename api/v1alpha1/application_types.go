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

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Owner of the Application.
	// +kubebuilder:validation:Required
	Owner string `json:"owner"`
	// Github Repository of the Application.
	Repository string `json:"repository,omitempty"`
	// Infrastrcture  of the Application.
	Infrastrcture string `json:"infrastrcture,omitempty"`
	// logs of the Application.
	Logs string `json:"logs,omitempty"`
	// traces of the Application.
	Traces string `json:"traces,omitempty"`
	// Dashboard of the Application.
	Dashboard string `json:"dashboard,omitempty"`
}

// ApplicationStatus defines the observed state of Application
type ApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Application is the Schema for the applications API
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:printcolumn:name=owner,type=string,JSONPath=.spec.owner
// +kubebuilder:printcolumn:name=repository,type=string,JSONPath=.spec.repository
// +kubebuilder:printcolumn:name=infrastrcture,type=string,JSONPath=.spec.infrastrcture
// +kubebuilder:printcolumn:name=logs,type=string,JSONPath=.spec.infrastrcture
// +kubebuilder:printcolumn:name=traces,type=string,JSONPath=.spec.traces
// +kubebuilder:printcolumn:name=age,type=date,JSONPath=.metadata.creationTimestamp
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
