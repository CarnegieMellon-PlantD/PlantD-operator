package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Formula defines the formula in column.
type Formula struct {
	// Name of the formula. Used together with the `args` field.
	// See https://plantd.org/docs/reference/formulas for available values.
	Name string `json:"name"`
	// Arguments to be passed to the formula. Used together with the `name` field.
	// See https://plantd.org/docs/reference/formulas for available values.
	Args []string `json:"args,omitempty"`
}

// Column defines the column in Schema.
type Column struct {
	// Name of the column.
	Name string `json:"name"`
	// Data type of the random data to be generated in the column. Used together with the `params` field.
	// It should be a valid function name in gofakeit, which can be parsed by gofakeit.GetFuncLookup().
	// `formula` field has precedence over this field.
	// See https://plantd.org/docs/reference/types-and-params for available values.
	Type string `json:"type,omitempty"`
	// Map of parameters for generating the data in the column. Used together with the `type` field.
	// For any parameters not provided but required by the data type, the default value will be used, if available.
	// Will ignore any parameters not used by the data type.
	// See https://plantd.org/docs/reference/types-and-params for available values.
	Params map[string]string `json:"params,omitempty"`
	// Formula to be applied for populating the data in the column.
	// This field has precedence over the `type` fields.
	Formula Formula `json:"formula,omitempty"`
}

// SchemaSpec defines the desired state of Schema.
type SchemaSpec struct {
	// List of columns in the Schema.
	Columns []Column `json:"columns"`
}

// SchemaStatus defines the observed state of Schema.
type SchemaStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Schema is the Schema for the schemas API
type Schema struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SchemaSpec   `json:"spec,omitempty"`
	Status SchemaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SchemaList contains a list of Schema
type SchemaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Schema `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Schema{}, &SchemaList{})
}
