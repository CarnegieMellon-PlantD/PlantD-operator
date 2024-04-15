package utils

import (
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNamespacedName returns the namespace and name of the resource in string representation.
func GetNamespacedName(obj v1.Object) string {
	return fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
}

// GetDataSetJobName returns the name of the Job for the DataSet.
func GetDataSetJobName(dataSetName string, generation int64) string {
	return fmt.Sprintf("dataset-%s-%d", dataSetName, generation)
}

// GetDataSetPVCName returns the name of the PVC for the DataSet.
func GetDataSetPVCName(dataSetName string, generation int64) string {
	return fmt.Sprintf("dataset-%s-%d", dataSetName, generation)
}

// GetPipelineMetricsServiceName returns the name of the metrics Service for the Pipeline.
// For out-cluster Pipeline only, for whom we need to create Service of type ExternalName.
func GetPipelineMetricsServiceName(pipelineName string) string {
	return fmt.Sprintf("%s-metrics", pipelineName)
}

// GetTestRunConfigMapName returns the name of the ConfigMap for the TestRun.
func GetTestRunConfigMapName(experimentName string, endpointName string) string {
	return fmt.Sprintf("experiment-%s-%s", experimentName, endpointName)
}

// GetTestRunPVCName returns the name of the PVC for the TestRun.
// For data option "dataSet" only, which requires a PVC.
func GetTestRunPVCName(experimentName string, endpointName string) string {
	return fmt.Sprintf("experiment-%s-%s", experimentName, endpointName)
}

// GetTestRunCopierJobName returns the name of the copier Job for the TestRun.
// The copier Job is used to copy the configuration and data for the TestRun.
func GetTestRunCopierJobName(experimentName string, endpointName string) string {
	return fmt.Sprintf("experiment-%s-%s-copier", experimentName, endpointName)
}

// GetTestRunName returns the name of the TestRun.
func GetTestRunName(experimentName string, endpointName string) string {
	return fmt.Sprintf("experiment-%s-%s", experimentName, endpointName)
}
