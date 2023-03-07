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

type ComponentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// Role which can access the Catalog.
	Name string `json:"name"`
	// Github Repository of the Catalog.
	Repository string `json:"repository,omitempty"`
	// Gitlab ProjectID  of the Catalog.
	ProjectId int16 `json:"project_id,omitempty"`
	// readme_file_path of the Catalog.
	ReadmeFilePath string `json:"readme_file_path,omitempty"`
	// image_path of the Catalog.
	ImagePath string `json:"image_path,omitempty"`
	// variable_file_path of the Catalog.
	VariableFilePath string `json:"variable_file_path,omitempty"`
	// tool_type of the Component.
	ToolType string `json:"tool_type,omitempty"`
	// ref of the Catalog.
	Ref string `json:"ref,omitempty"`
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
	// Github Repository of the Catalog.
	Repository string `json:"repository,omitempty"`
	// Gitlab ProjectID  of the Catalog.
	ProjectId int16 `json:"project_id,omitempty"`
	// values_folder_path of the Catalog.
	ValuesFolderPath string `json:"values_folder_path,omitempty"`
	// values_folder_path of the Catalog.
	SecretsFolderPath string `json:"secrets_folder_path,omitempty"`
	// image_path of the Catalog.
	ImagePath string `json:"image_path,omitempty"`
	// readme_file_path of the Catalog.
	ReadmeFilePath string `json:"readme_file_path,omitempty"`
	// ref of the Catalog.
	Ref string `json:"ref,omitempty"`
	// ref of the Catalog.
	Sops bool `json:"sops,omitempty"`
	// components of the Catalog
	Components []ComponentSpec `json:"components",omitempty`
}

// CatalogStatus defines the observed state of Catalog
type CatalogStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Catalog is the Schema for the catalogs API
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:printcolumn:name=type,type=string,JSONPath=.spec.type
// +kubebuilder:printcolumn:name=repository,type=string,JSONPath=.spec.repository
// +kubebuilder:printcolumn:name=ref,type=string,JSONPath=.spec.ref
// +kubebuilder:printcolumn:name=role,type=string,JSONPath=.spec.role
// +kubebuilder:printcolumn:name=age,type=date,JSONPath=.metadata.creationTimestamp
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
