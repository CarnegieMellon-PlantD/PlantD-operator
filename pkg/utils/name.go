package utils

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetDataSetJobName returns the name of the Job for the DataSet.
func GetDataSetJobName(dataSetName string, generation int64) string {
	return fmt.Sprintf("%s-%d-job", dataSetName, generation)
}

// GetDataSetPVCName returns the name of the PVC for the DataSet.
func GetDataSetPVCName(dataSetName string, generation int64) string {
	return fmt.Sprintf("%s-%d-pvc", dataSetName, generation)
}

// GetDataSetVolumeName returns the name of the volume in the DataSet Job.
func GetDataSetVolumeName(dataSetName string) string {
	return fmt.Sprintf("%s-volume", dataSetName)
}

func GetNamespacedName(obj client.Object) string {
	return fmt.Sprintf("%s-%s", obj.GetNamespace(), obj.GetName())
}

func GetTestRunName(expName string, endpointName string) string {
	return fmt.Sprintf("%s-%s", expName, endpointName)
}

func GetMetricsServiceName(pipelineName string) string {
	return pipelineName + "-plantd-metrics"
}

func GetPipelineEndpointServiceName(pipelineName string, endpointName string) string {
	return fmt.Sprintf("%s-%s", pipelineName, endpointName)
}
