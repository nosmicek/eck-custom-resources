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

// InferenceEndpointSpec defines the desired state of InferenceEndpoint
type InferenceEndpointSpec struct {
	// +optional
	TargetConfig CommonElasticsearchConfig `json:"targetInstance,omitempty"`

	// TaskType forms part of the endpoint's path in Elasticsearch, so it cannot be
	// changed without recreating the endpoint.
	// +kubebuilder:validation:Enum=sparse_embedding;text_embedding;rerank;completion;chat_completion;embedding
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="taskType is immutable"
	TaskType string `json:"taskType"`

	Body string `json:"body"`
}

// InferenceEndpointStatus defines the observed state of InferenceEndpoint
type InferenceEndpointStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// InferenceEndpoint is the Schema for the inferenceendpoints API
type InferenceEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InferenceEndpointSpec   `json:"spec,omitempty"`
	Status InferenceEndpointStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// InferenceEndpointList contains a list of InferenceEndpoint
type InferenceEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InferenceEndpoint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InferenceEndpoint{}, &InferenceEndpointList{})
}
