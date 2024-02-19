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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SimulationSpec defines the desired state of Simulation
type SimulationSpec struct {
	// TrafficModelRef defines a reference of the TrafficModel object.
	TrafficModelRef corev1.ObjectReference `json:"trafficModelRef"`
	// DigitalTwinRef defines a reference of the DigitalTwin object.
	DigitalTwinRef corev1.ObjectReference `json:"digitalTwinRef"`
	// ScheduledTime defines the scheduled time for the Simulation.
	ScheduledTime metav1.Time `json:"scheduledTime,omitempty"`
}

// SimulationStatus defines the observed state of Simulation
type SimulationStatus struct {
	DigitalTwinState  string `json:"digitalTwinState,omitempty"`
	TrafficModelState string `json:"trafficModalState,omitempty"`

	PodName string `json:"podName,omitempty"`
	// JobStatus defines the status of the cost calculation job.
	JobStatus string `json:"jobStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="DigitalTwinState",type="string",JSONPath=".status.digitaltwinState"
//+kubebuilder:printcolumn:name="TrafficModelState",type="string",JSONPath=".status.trafficmodalState"
//+kubebuilder:printcolumn:name="PodName",type="string",JSONPath=".status.podName"
//+kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"

// Simulation is the Schema for the simulations API
type Simulation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SimulationSpec   `json:"spec,omitempty"`
	Status SimulationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SimulationList contains a list of Simulation
type SimulationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Simulation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Simulation{}, &SimulationList{})
}
