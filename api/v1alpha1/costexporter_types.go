package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CostExporterSpec defines the desired state of CostExporter.
type CostExporterSpec struct {
	// Container image to use for cost exporter.
	Image string `json:"image,omitempty"`
	// Cloud service provider to calculate costs for. Available value is `aws`.
	// +kubebuilder:validation:Enum=aws
	CloudServiceProvider string `json:"cloudServiceProvider"`
	// Configuration for the cloud service provider.
	// For AWS, the configuration should be a JSON string with the following fields:
	// - `AWS_ACCESS_KEY`
	// - `AWS_SECRET_KEY`
	// - `S3_BUCKET_NAME`
	Config *corev1.SecretKeySelector `json:"config"`
}

// CostExporterStatus defines the observed state of CostExporter.
type CostExporterStatus struct {
	// Time when the last successful completion of the Job.
	LastCompletionTime *metav1.Time `json:"lastCompletionTime,omitempty"`
	// Time when the last failed completion of the Job.
	LastFailureTime *metav1.Time `json:"lastFailureTime,omitempty"`
	// Whether the Job is running. For internal use only.
	IsRunning bool `json:"isRunning,omitempty"`
}

// The name of the Pod in the CostExporter will be
// "<costexporter-name>-costexporter-<random 5 chars>".
// So, we have 44 characters for the name to meet the 63-character limit.

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="LastCompletionTime",type="string",JSONPath=".status.lastCompletionTime"

// CostExporter is the Schema for the costexporters API
// +kubebuilder:validation:XValidation:rule="size(self.metadata.name) <= 44",message="must contain at most 44 characters"
type CostExporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CostExporterSpec   `json:"spec,omitempty"`
	Status CostExporterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CostExporterList contains a list of CostExporter
type CostExporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CostExporter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CostExporter{}, &CostExporterList{})
}
