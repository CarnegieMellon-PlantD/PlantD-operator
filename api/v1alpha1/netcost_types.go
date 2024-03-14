/*
Copyright 2023.

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

// NetCostSpec defines the desired state of NetCost
type NetCostSpec struct {
	// NetCostPerMB defines the cost per MB of data transfer.
	NetCostPerMB string `json:"netCostPerMB,omitempty"`
	// RawDataStoreCostPerMB defines the cost per MB of raw data storage.
	RawDataStoreCostPerMB string `json:"rawDataStoreCostPerMB,omitempty"`
	// ProcessedDataStoreCostPerMB defines the cost per MB of processed data storage.
	ProcessedDataStoreCostPerMB string `json:"processedDataStoreCostPerMB,omitempty"`
	// RawDataRetentionPolicyMonths defines the months raw data is retained.
	RawDataRetentionPolicyMonths int `json:"rawDataRetentionPolicyMonths,omitempty"`
	// ProcessedDataRetentionPolicyMonths defines the months processed data is retained.
	ProcessedDataRetentionPolicyMonths int `json:"processedDataRetentionPolicyMonths,omitempty"`
}

// NetCostStatus defines the observed state of NetCost
type NetCostStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NetCost is the Schema for the netcosts API
type NetCost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetCostSpec   `json:"spec,omitempty"`
	Status NetCostStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NetCostList contains a list of NetCost
type NetCostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetCost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetCost{}, &NetCostList{})
}
