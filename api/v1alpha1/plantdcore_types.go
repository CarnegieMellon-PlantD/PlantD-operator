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
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type PlantDPrometheusSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ResourceMemory v1.ResourceRequirements `json:"resourceMemory,omitempty"`
	ScrapeInterval monitoringv1.Duration   `json:"scrapeInterval,omitempty"`
}

// PlantDCoreSpec defines the desired state of PlantDCore
type PlantDCoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	PrometheusConfiguration PlantDPrometheusSpec `json:"prometheusConfiguration,omitempty"`
}

// PlantDCoreStatus defines the observed state of PlantDCore
type PlantDCoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ProxyStatus shows the status of the PlantD Proxy
	ProxyStatus string `json:"kubeProxyStatus,omitempty"`
	// StudioStatus shows the status of the PlantD Studio
	StudioStatus string `json:"studioStatus,omitempty"`
	// PrometheusStatus shows the status of the Prometheus
	PrometheusStatus string `json:"prometheusStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="KubeProxyStatus",type="string",JSONPath=".status.kubeProxyStatus"
//+kubebuilder:printcolumn:name="StudioStatus",type="string",JSONPath=".status.studioStatus"
//+kubebuilder:printcolumn:name="PrometheusStatus",type="string",JSONPath=".status.prometheusStatus"

// PlantDCore is the Schema for the plantdcores API
type PlantDCore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlantDCoreSpec   `json:"spec,omitempty"`
	Status PlantDCoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PlantDCoreList contains a list of PlantDCore
type PlantDCoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlantDCore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PlantDCore{}, &PlantDCoreList{})
}
