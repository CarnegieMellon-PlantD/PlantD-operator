package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FormulaSpec defines the specification of the formula.
type FormulaSpec struct {
	// Name defines the name of the formula. Should match the name with one of the provided formulas.
	Name string `json:"name"`
	// Args defines the arugments for calling the formula.
	Args []string `json:"args,omitempty"`
}

// Column defines the metadata of the column data.
type Column struct {
	// Name defines the name of the column.
	Name string `json:"name"`
	// Type defines the data type of the column. Should match the type with one of the provided types.
	Type string `json:"type,omitempty"`
	// Params defines the parameters for constructing the data give certain data type.
	Params map[string]string `json:"params,omitempty"`
	// Formula defines the formula applies to the column data.
	Formula FormulaSpec `json:"formula,omitempty"`
}

// SchemaSpec defines the desired state of Schema
type SchemaSpec struct {
	// Columns defines a list of column specifications.
	Columns []Column `json:"columns"`
}

// SchemaStatus defines the observed state of Schema
type SchemaStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Schema is the Schema for the schemas API
type Schema struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specifications of the Schema.
	Spec SchemaSpec `json:"spec,omitempty"`
	// Status defines the status of the Schema.
	Status SchemaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SchemaList contains a list of Schema
type SchemaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items defines a list of Schemas.
	Items []Schema `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Schema{}, &SchemaList{})
}
