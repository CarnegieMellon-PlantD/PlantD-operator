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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TrafficModelSpec defines the desired state of TrafficModel
type TrafficModelSpec struct {
	// Config defines the configuration of the TrafficModel.
	Config string `json:"config,omitempty"`
}

// TrafficModelStatus defines the observed state of TrafficModel
type TrafficModelStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TrafficModel is the Schema for the trafficmodels API
type TrafficModel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specifications of the TrafficModel.
	Spec TrafficModelSpec `json:"spec,omitempty"`
	// Status defines the status of the TrafficModel.
	Status TrafficModelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TrafficModelList contains a list of TrafficModel
type TrafficModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items defines a list of TrafficModels.
	Items []TrafficModel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrafficModel{}, &TrafficModelList{})
}
