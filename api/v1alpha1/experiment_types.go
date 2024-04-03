package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExperimentJobStatus defines the status of the load generator job.
type ExperimentJobStatus string

const (
	ExperimentScheduled       ExperimentJobStatus = "Scheduled"
	ExperimentWaitingDataSet  ExperimentJobStatus = "Waiting for DataSet"
	ExperimentWaitingPipeline ExperimentJobStatus = "Waiting for Pipeline"
	ExperimentInitializing    ExperimentJobStatus = "Initializing"
	ExperimentRunning         ExperimentJobStatus = "Running"
	ExperimentCompleted       ExperimentJobStatus = "Completed"
	ExperimentFailed          ExperimentJobStatus = "Failed"
)

// EndpointProtocol defines the protocol used by a PipelineEndpoint.
type EndpointProtocol string

const (
	EndpointProtocolHTTP EndpointProtocol = "http"
)

// EndpointDataOption defines the data option used by an EndpointSpec.
type EndpointDataOption string

const (
	EndpointDataOptionPlainText EndpointDataOption = "plainText"
	EndpointDataOptionDataSet   EndpointDataOption = "dataSet"
)

// DataSpec defines the data to be sent to an endpoint.
type DataSpec struct {
	// PlainText data to be sent. `dataSetRef` field has precedence over this field.
	PlainText string `json:"plainText,omitempty"`
	// Reference to the DataSet to be sent. This field has precedence over the `plainText` field.
	DataSetRef corev1.ObjectReference `json:"dataSetRef,omitempty"`
}

// EndpointSpec defines the test upon an endpoint.
type EndpointSpec struct {
	// Name of endpoint.
	// It should be an existing endpoint defined in the Pipeline used by the Experiment.
	EndpointName string `json:"endpointName"`
	// Data to be sent to the endpoint.
	DataSpec DataSpec `json:"dataSpec"`
	// LoadPattern to follow for the endpoint.
	LoadPatternRef corev1.ObjectReference `json:"loadPatternRef"`
	// Size of the PVC for the load generator job.
	// Only effective when `dataSpec.dataSetRef` is set.
	// Default to the PVC size of the DataSet.
	StorageSize resource.Quantity `json:"storageSize,omitempty"`
}

// ExperimentSpec defines the desired state of Experiment.
type ExperimentSpec struct {
	// Reference to the Pipeline to use for the Experiment.
	PipelineRef corev1.ObjectReference `json:"pipelineRef"`
	// List of tests upon endpoints.
	// +kubebuilder:validation:MinItems=1
	EndpointSpecs []EndpointSpec `json:"endpointSpecs"`
	// Scheduled time to run the Experiment.
	ScheduledTime metav1.Time `json:"scheduledTime,omitempty"`
}

// ExperimentStatus defines the observed state of Experiment.
type ExperimentStatus struct {
	// Calculated duration of each endpoint.
	Durations map[string]metav1.Duration `json:"durations,omitempty"`
	// Status of the load generator job.
	JobStatus ExperimentJobStatus `json:"jobStatus,omitempty"`
	// Time when the Experiment started.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// Time when the Experiment completed.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
	// Error message.
	Error string `json:"error,omitempty"`
	// Pipeline used by the Experiment.
	// For internal use only.
	Pipeline *Pipeline `json:"pipeline,omitempty"`
	// Map from endpoint name to the PipelineEndpoint, which is referenced by the EndpointSpec.
	// For internal use only.
	EndpointMap map[string]*PipelineEndpoint `json:"endpointMap,omitempty"`
	// Map from endpoint name to the protocol used by the PipelineEndpoint, which is referenced by the EndpointSpec.
	// For internal use only.
	ProtocolMap map[string]EndpointProtocol `json:"protocolMap,omitempty"`
	// Map from endpoint name to the data option used by the EndpointSpec.
	// For internal use only.
	DataOptionMap map[string]EndpointDataOption `json:"dataOptionMap,omitempty"`
	// Map from endpoint name to the DataSet used by the EndpointSpec.
	// For internal use only.
	DataSetMap map[string]*DataSet `json:"dataSetMap,omitempty"`
	// Map from endpoint name to the LoadPattern used by the EndpointSpec.
	// For internal use only.
	LoadPatternMap map[string]*LoadPattern `json:"loadPatternMap,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Duration",type="string",JSONPath=".status.durations"
//+kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"
//+kubebuilder:printcolumn:name="StartTime",type="string",JSONPath=".status.startTime"
//+kubebuilder:printcolumn:name="CompletionTime",type="string",JSONPath=".status.completionTime"
//+kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.error"

// Experiment is the Schema for the experiments API
type Experiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Experiment is immutable"
	Spec   ExperimentSpec   `json:"spec,omitempty"`
	Status ExperimentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExperimentList contains a list of Experiments.
type ExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Experiment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
}
