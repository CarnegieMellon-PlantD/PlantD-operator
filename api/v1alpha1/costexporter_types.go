package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CostExporterSpec defines the desired state of CostExporter
type CostExporterSpec struct {
	// S3Bucket defines the AWS S3 bucket name where stores the cost logs.
	S3Bucket string `json:"s3Bucket,omitempty"`
	// CloudServiceProvider defines the target cloud service provide for calculating cost.
	CloudServiceProvider string `json:"cloudServiceProvider,omitempty"`
	// SecretRef defines the reference to the Kubernetes Secret where stores the credentials of cloud service provider
	SecretRef corev1.ObjectReference `json:"secretRef,omitempty"`
}

// CostExporterStatus defines the observed state of CostExporter
type CostExporterStatus struct {
	// JobCompletionTime defines the completion time of the cost calculation job.
	JobCompletionTime *metav1.Time `json:"jobCompletionTime,omitempty"`
	// PodName defines the name of the cost calculation pod.
	PodName string `json:"podName,omitempty"`
	// JobStatus defines the status of the cost calculation job.
	JobStatus string `json:"jobStatus,omitempty"`
	// Tags defines the json string of using tags.
	Tags string `json:"tags,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="JobCompletionTime",type="string",JSONPath=".status.jobCompletionTime"
// +kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"
// +kubebuilder:printcolumn:name="PodName",type="string",JSONPath=".status.podName"

// CostExporter is the Schema for the costexporters API
type CostExporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specifications of the CostExporter.
	Spec CostExporterSpec `json:"spec,omitempty"`
	// Status defines the status of the CostExporter.
	Status CostExporterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CostExporterList contains a list of CostExporter
type CostExporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items defines a list of CostExporters.
	Items []CostExporter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CostExporter{}, &CostExporterList{})
}
