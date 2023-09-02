package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SchemaSelector defines a list of Schemas and the required numbers and format.
type SchemaSelector struct {
	// Name defines the name of the Schame. Should match the name of existing Schema in the same namespace as the DataSet.
	Name string `json:"name"`
	// NumRecords defines the number of records to be generated in each output file. A random number is picked from the specified range.
	NumRecords map[string]int `json:"numRecords,omitempty"`
	// NumberOfFilesPerCompressedFile defines the number of intermediate files to be compressed into a single compressed file.
	// A random number is picked from the specified range.
	NumberOfFilesPerCompressedFile map[string]int `json:"numFilesPerCompressedFile,omitempty"`
}

// DataSetSpec defines the desired state of DataSet
type DataSetSpec struct {
	// FileFormat defines the file format of the each file containing the generated data.
	// This may or may not be the output file format based on whether you want to compress these files.
	FileFormat string `json:"fileFormat"`
	// CompressedFileFormat defines the file format for the compressed files.
	// Each file inside the compressed file is of "fileFormat" format specified above.
	// This is the output format if specified for the files.
	CompressedFileFormat string `json:"compressedFileFormat,omitempty"`
	// CompressPerSchema defines the flag of compression.
	// If you wish files from all the different schemas to compressed into one compressed file leave this field as false.
	// If you wish to have a different compressed file for every schema, mark this field as true.
	CompressPerSchema bool `json:"compressPerSchema,omitempty"`
	// NumberOfFiles defines the total number of output files irrespective of compression.
	// Unless "compressPerSchema" is false, this field is applicable per schema.
	NumberOfFiles int32 `json:"numFiles"`
	// Schemas defines a list of Schemas.
	Schemas []SchemaSelector `json:"schemas"`
	// ParallelJobs defines the number of parallel jobs when generating the dataset.
	// TODO: Infer the optimal number of parallel jobs automatically.
	ParallelJobs int32 `json:"parallelJobs,omitempty"`
}

// DataSetStatus defines the observed state of DataSet
type DataSetStatus struct {
	// JobStatus defines the status of the data generating job.
	JobStatus string `json:"jobStatus,omitempty"`
	// PVCStatus defines the status of the PVC mount to the data generating pod.
	PVCStatus string `json:"pvcStatus,omitempty"`
	// StartTime defines the start time of the data generating job.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// CompletionTime defines the duration of the data generating job.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
	// LastGeneration defines the last generation of the DataSet object.
	LastGeneration int64 `json:"lastGeneration,omitempty"`
	// ErrorCount defines the number of errors raised by the controller or data generating job.
	ErrorCount int `json:"errorCount,omitempty"`
	// Errors defines the map of error messages.
	Errors map[string][]string `json:"errors,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="StartTime",type="string",JSONPath=".status.startTime"
// +kubebuilder:printcolumn:name="CompletionTime",type="string",JSONPath=".status.completionTime"
// +kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"
// +kubebuilder:printcolumn:name="VolumeStatus",type="string",JSONPath=".status.pvcStatus"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="ErrorCount",type="integer",JSONPath=".status.errorCount"
// +kubebuilder:printcolumn:name="Errors",type="string",JSONPath=".status.errorsString"

// DataSet is the Schema for the datasets API
type DataSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specifications of the DataSet.
	Spec DataSetSpec `json:"spec,omitempty"`
	// Status defines the status of the DataSet.
	Status DataSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataSetList contains a list of DataSet
type DataSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items defines a list of DataSets.
	Items []DataSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataSet{}, &DataSetList{})
}
