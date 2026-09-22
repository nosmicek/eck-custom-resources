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

// QueryRulesetSpec defines the desired state of QueryRuleset
type QueryRulesetSpec struct {
	// +optional
	TargetConfig CommonElasticsearchConfig `json:"targetInstance,omitempty"`

	Body string `json:"body"`
}

// QueryRulesetStatus defines the observed state of QueryRuleset
type QueryRulesetStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// QueryRuleset is the Schema for the queryrulesets API
type QueryRuleset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueryRulesetSpec   `json:"spec,omitempty"`
	Status QueryRulesetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QueryRulesetList contains a list of QueryRuleset
type QueryRulesetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueryRuleset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QueryRuleset{}, &QueryRulesetList{})
}
