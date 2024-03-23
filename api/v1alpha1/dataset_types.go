package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DataSetJobStatus defines the status of the data generating job.
type DataSetJobStatus string

const (
	DataSetJobRunning DataSetJobStatus = "Running"
	DataSetJobFailed  DataSetJobStatus = "Failed"
	DataSetJobSuccess DataSetJobStatus = "Success"
)

// DataSetErrorType defines the type of error occurred.
type DataSetErrorType string

const (
	DataSetControllerError DataSetErrorType = "controller"
	DataSetJobError        DataSetErrorType = "job"
)

// SchemaSelector defines the reference to a Schema and its usage in the DataSet.
type SchemaSelector struct {
	// Name of the Schema. Note that the Schema must be present in the same namespace as the DataSet.
	Name string `json:"name"`
	// Range of number of rows to be generated in each output file.
	// Should be a map containing `min` and `max` keys. For each output file, a random number is picked
	// from the specified range.
	NumRecords map[string]int `json:"numRecords,omitempty"`
	// Range of number of files to be generated in the compressed file.
	// Take effect only if `compressedFileFormat` is set in the DataSet.
	// Should be a map containing `min` and `max` keys. A random number is picked from the specified range.
	NumberOfFilesPerCompressedFile map[string]int `json:"numFilesPerCompressedFile,omitempty"`
}

// DataSetSpec defines the desired state of DataSet.
type DataSetSpec struct {
	// Format of the output file containing generated data. Available values are `csv` and `binary`.
	FileFormat string `json:"fileFormat,omitempty"`
	// Format of the compressed file containing output files. Available value is `zip`.
	// Leave empty to disable compression.
	CompressedFileFormat string `json:"compressedFileFormat,omitempty"`
	// Flag for compression behavior. Takes effect only if `compressedFileFormat` is set.
	// When set to `false` (default), files from all Schemas will be compressed into a single
	// compressed file in each repetition.
	// When set to `true`, files from each Schema will be compressed into a separate compressed
	// file in each repetition.
	CompressPerSchema bool `json:"compressPerSchema,omitempty"`
	// Number of repetitions of the data generation process.
	// If `compressedFileFormat` is unset, this is the number of files for each Schema.
	// If `compressedFileFormat` is set and `compressPerSchema` is `false`, this is the number of
	// compressed files for each Schema.
	// If `compressedFileFormat` is set and `compressPerSchema` is `true`, this is the total
	// number of compressed files.
	NumberOfFiles int32 `json:"numFiles"`
	// List of Schemas in the DataSet.
	Schemas []SchemaSelector `json:"schemas"`
	// Number of parallel jobs when generating the dataset.
	ParallelJobs int32 `json:"parallelJobs,omitempty"`
}

// DataSetStatus defines the observed state of DataSet.
type DataSetStatus struct {
	// Status of the data generating job.
	JobStatus DataSetJobStatus `json:"jobStatus,omitempty"`
	// Status of the PVC of the data generating job.
	PVCStatus v1.PersistentVolumeClaimPhase `json:"pvcStatus,omitempty"`
	// Time when the data generating job started.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// Time when the data generating job completed.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
	// Number of errors occurred.
	ErrorCount int32 `json:"errorCount,omitempty"`
	// List of errors occurred, which is a map from error type to list of error messages.
	Errors map[DataSetErrorType][]string `json:"errors,omitempty"`
	// Last generation of the DataSet object. For internal use only.
	LastGeneration int64 `json:"lastGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="StartTime",type="string",JSONPath=".status.startTime"
// +kubebuilder:printcolumn:name="CompletionTime",type="string",JSONPath=".status.completionTime"
// +kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"
// +kubebuilder:printcolumn:name="VolumeStatus",type="string",JSONPath=".status.pvcStatus"
// +kubebuilder:printcolumn:name="ErrorCount",type="integer",JSONPath=".status.errorCount"

// DataSet is the Schema for the datasets API
type DataSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataSetSpec   `json:"spec,omitempty"`
	Status DataSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataSetList contains a list of DataSet
type DataSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataSet{}, &DataSetList{})
}
