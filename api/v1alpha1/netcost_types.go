package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetCostSpec defines the desired state of NetCost.
type NetCostSpec struct {
	// The cost per MB of data transfer.
	// The value should be a float number in string format.
	NetCostPerMB string `json:"netCostPerMB"`
	// The cost per MB per month of raw data storage.
	// The value should be a float number in string format.
	RawDataStoreCostPerMBMonth string `json:"rawDataStoreCostPerMBMonth"`
	// The cost per MB per month of processed data storage.
	// The value should be a float number in string format.
	ProcessedDataStoreCostPerMBMonth string `json:"processedDataStoreCostPerMBMonth"`
	// The number of months the raw data is retained.
	RawDataRetentionPolicyMonths int `json:"rawDataRetentionPolicyMonths"`
	// The number of months the processed data is retained.
	ProcessedDataRetentionPolicyMonths int `json:"processedDataRetentionPolicyMonths"`
}

// NetCostStatus defines the observed state of NetCost.
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
