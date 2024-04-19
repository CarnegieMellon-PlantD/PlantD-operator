package utils

import (
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNamespacedName returns the namespace and name of the resource in string representation.
func GetNamespacedName(obj v1.Object) string {
	return fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
}

// GetDataGeneratorName returns the name of the data generator resources for the DataSet.
// Note that to shorten the name, only the last 4 hex digits of the generation number are used.
// It is safe because we always delete the old resources before creating new ones.
func GetDataGeneratorName(dataSetName string, generation int64) string {
	return fmt.Sprintf("%s-datagen-%x", dataSetName, generation%0x10000)
}

// GetMetricsServiceName returns the name of the metrics Service and ServiceMonitor for the Pipeline.
func GetMetricsServiceName(pipelineName string) string {
	return fmt.Sprintf("%s-metrics", pipelineName)
}

// GetTestRunName returns the name of the TestRun for the Experiment.
// Note that to shorten the name, only the last 4 hex digits of the endpoint index are used.
// It is safe because we limit the number of EndpointSpecs in the Experiment to be no more than 65535.
func GetTestRunName(experimentName string, endpointIdx int) string {
	return fmt.Sprintf("%s-loadgen-%x", experimentName, (endpointIdx+1)%0x10000)
}

// GetTestRunCopierJobName returns the name of the copier Job for the TestRun.
// The copier Job is used to copy the configuration and data for the TestRun.
// Note that to shorten the name, only the last 4 hex digits of the endpoint index are used.
// It is safe because we limit the number of EndpointSpecs in the Experiment to be no more than 65535.
func GetTestRunCopierJobName(experimentName string, endpointIdx int) string {
	return fmt.Sprintf("%s-loadgen-%x-copier", experimentName, (endpointIdx+1)%0x10000)
}
