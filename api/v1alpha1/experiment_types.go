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
	ExperimentDraining        ExperimentJobStatus = "Draining"
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
	// PlainText data to be sent.
	// `dataSetRef` field has precedence over this field.
	PlainText string `json:"plainText,omitempty"`
	// Reference to the DataSet to be sent.
	// The DataSet must be in the same namespace as the Experiment.
	// This field has precedence over the `plainText` field.
	DataSetRef *corev1.LocalObjectReference `json:"dataSetRef,omitempty"`
}

// EndpointSpec defines the test upon an endpoint.
type EndpointSpec struct {
	// Name of endpoint.
	// It should be an existing endpoint defined in the Pipeline used by the Experiment.
	EndpointName string `json:"endpointName"`
	// Data to be sent to the endpoint.
	DataSpec *DataSpec `json:"dataSpec"`
	// LoadPattern to follow for the endpoint.
	LoadPatternRef *corev1.ObjectReference `json:"loadPatternRef"`
	// Size of the PVC for the load generator job.
	// Only effective when `dataSpec.dataSetRef` is set.
	// Default to the PVC size of the DataSet.
	StorageSize *resource.Quantity `json:"storageSize,omitempty"`
}

// ExperimentSpec defines the desired state of Experiment.
type ExperimentSpec struct {
	// Reference to the Pipeline to use for the Experiment.
	PipelineRef *corev1.LocalObjectReference `json:"pipelineRef"`
	// List of tests upon endpoints.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=65535
	EndpointSpecs []EndpointSpec `json:"endpointSpecs"`
	// Scheduled time to run the Experiment.
	ScheduledTime *metav1.Time `json:"scheduledTime,omitempty"`
	// Time to wait after the load generator job is completed before finishing the Experiment.
	// It allows the pipeline-under-test to finish its processing.
	// Default to no draining time.
	DrainingTime *metav1.Duration `json:"drainingTime,omitempty"`
}

// ExperimentStatus defines the observed state of Experiment.
type ExperimentStatus struct {
	// Calculated duration of each endpoint.
	Durations map[string]*metav1.Duration `json:"durations,omitempty"`
	// Status of the load generator job.
	JobStatus ExperimentJobStatus `json:"jobStatus,omitempty"`
	// Time when the Experiment started.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// Time when the Experiment completed.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
	// Error message.
	Error string `json:"error,omitempty"`
	// Time when the pipeline-under-test started draining. For internal use only.
	DrainingStartTime *metav1.Time `json:"drainingStartTime,omitempty"`
	// Whether to enable cost calculation.
	// Copied from the Pipeline used by the Experiment. For internal use only.
	EnableCostCalculation bool `json:"enableCostCalculation,omitempty"`
	// Cloud provider. Available values are `aws`, `azure`, and `gcp`.
	// Copied from the Pipeline used by the Experiment. For internal use only.
	CloudProvider string `json:"cloudProvider,omitempty"`
	// Map of tags to select cloud resources. Equivalent to the tags in the cloud service provider.
	// Copied from the Pipeline used by the Experiment. For internal use only.
	Tags map[string]string `json:"tags,omitempty"`
}

// The longest name of the Pod in the Experiment will be
// "<experiment-name>-loadgen-<up to 4 digits of endpoint index>-initializer-<random 5 characters>".
// So, we have 32 characters for the name to meet the 63-character limit.

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"
//+kubebuilder:printcolumn:name="Durations",type="string",JSONPath=".status.durations"
//+kubebuilder:printcolumn:name="Draining",type="string",JSONPath=".spec.drainingTime"
//+kubebuilder:printcolumn:name="ScheduledTime",type="string",JSONPath=".spec.scheduledTime"
//+kubebuilder:printcolumn:name="StartTime",type="string",JSONPath=".status.startTime"
//+kubebuilder:printcolumn:name="CompletionTime",type="string",JSONPath=".status.completionTime"

// Experiment is the Schema for the experiments API
// +kubebuilder:validation:XValidation:rule="size(self.metadata.name) <= 32",message="must contain at most 32 characters"
type Experiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

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
