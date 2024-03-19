package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetCostSpec defines the desired state of NetCost
type NetCostSpec struct {
	// NetCostPerMB defines the cost per MB of data transfer.
	NetCostPerMB resource.Quantity `json:"netCostPerMB,omitempty"`
	// RawDataStoreCostPerMBMonth defines the cost per MB per month of raw data storage.
	RawDataStoreCostPerMBMonth resource.Quantity `json:"rawDataStoreCostPerMBMonth,omitempty"`
	// ProcessedDataStoreCostPerMBMonth defines the cost per MB per month of processed data storage.
	ProcessedDataStoreCostPerMBMonth resource.Quantity `json:"processedDataStoreCostPerMBMonth,omitempty"`
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
