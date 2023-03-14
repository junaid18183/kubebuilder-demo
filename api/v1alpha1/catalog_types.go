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

// CrossNamespaceObjectReference contains enough information to let you locate
// the typed referenced object at cluster level.
type CrossNamespaceObjectReference struct {
	// APIVersion of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referent.
	// +kubebuilder:validation:Enum=HelmRepository;GitRepository;Bucket
	// +required
	Kind string `json:"kind,omitempty"`

	// Name of the referent.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +required
	Name string `json:"name"`

	// Namespace of the referent.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Optional
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

type ComponentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// Role which can access the Catalog.
	Name string `json:"name"`
	// The name and namespace of the gitrepositories.source.toolkit.fluxcd.io the Component is available at.
	SourceRef CrossNamespaceObjectReference `json:"sourceRef,omitempty"`
	// image_path of the Catalog.
	ImagePath string `json:"image_path,omitempty"`
	// variable_file_path of the Catalog.
	VariableFilePath string `json:"variable_file_path,omitempty"`
	// tool_type of the Component.
	ToolType string `json:"tool_type,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CatalogSpec defines the desired state of Catalog
type CatalogSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// description which can access the Catalog.
	Description string `json:"description,omitempty"`
	// Role which can access the Catalog.
	Role string `json:"role"`
	// Type of the Catalog.
	Type string `json:"type"`
	// The name and namespace of the gitrepositories.source.toolkit.fluxcd.io the Component is available at.
	// +required
	SourceRef CrossNamespaceObjectReference `json:"sourceRef"`
	// values_folder_path of the Catalog.
	ValuesFolderPath string `json:"values_folder_path,omitempty"`
	// values_folder_path of the Catalog.
	SecretsFolderPath string `json:"secrets_folder_path,omitempty"`
	// image_path of the Catalog.
	ImagePath string `json:"image_path,omitempty"`
	// readme_file_path of the Catalog.
	ReadmeFilePath string `json:"readme_file_path,omitempty"`
	// ref of the Catalog.
	Sops bool `json:"sops,omitempty"`
	// components of the Catalog
	Components []ComponentSpec `json:"components",omitempty`
}

// CatalogStatus defines the observed state of Catalog
type CatalogStatus struct {
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

// Catalog is the Schema for the catalogs API
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:printcolumn:name=type,type=string,JSONPath=.spec.type
// +kubebuilder:printcolumn:name=source_type,type=string,JSONPath=.spec.sourceRef.kind
// +kubebuilder:printcolumn:name=source_name,type=string,JSONPath=.spec.sourceRef.name
// +kubebuilder:printcolumn:name=role,type=string,JSONPath=.spec.role
// +kubebuilder:printcolumn:name=age,type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.c[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
type Catalog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CatalogSpec   `json:"spec,omitempty"`
	Status CatalogStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CatalogList contains a list of Catalog
type CatalogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Catalog `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Catalog{}, &CatalogList{})
}
