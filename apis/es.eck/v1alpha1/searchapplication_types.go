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
)

// SearchApplicationSpec defines the desired state of SearchApplication
type SearchApplicationSpec struct {
	// +optional
	TargetConfig CommonElasticsearchConfig `json:"targetInstance,omitempty"`
	// +optional
	Dependencies Dependencies `json:"dependencies,omitempty"`

	Body string `json:"body"`
}

// SearchApplicationStatus defines the observed state of SearchApplication
type SearchApplicationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SearchApplication is the Schema for the searchapplications API
type SearchApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SearchApplicationSpec   `json:"spec,omitempty"`
	Status SearchApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SearchApplicationList contains a list of SearchApplication
type SearchApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SearchApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SearchApplication{}, &SearchApplicationList{})
}
