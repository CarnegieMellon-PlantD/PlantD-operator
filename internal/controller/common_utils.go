package controller

import (
	"time"

	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
)

// isJobFinished checks if the Job is finished and returns the condition type.
func isJobFinished(job *kbatch.Job) (bool, kbatch.JobConditionType) {
	for _, c := range job.Status.Conditions {
		if (c.Type == kbatch.JobComplete || c.Type == kbatch.JobFailed) && c.Status == corev1.ConditionTrue {
			return true, c.Type
		}
	}
	return false, ""
}

// containMetricsEndpoint returns whether the Pipeline contains a valid specification of the metrics endpoint.
func containMetricsEndpoint(pipeline *windtunnelv1alpha1.Pipeline) bool {
	if pipeline.Spec.MetricsEndpoint == nil {
		return false
	}

	if pipeline.Spec.InCluster {
		return pipeline.Spec.MetricsEndpoint.ServiceRef != nil && pipeline.Spec.MetricsEndpoint.ServiceRef.Name != ""
	} else {
		return pipeline.Spec.MetricsEndpoint.HTTP != nil && pipeline.Spec.MetricsEndpoint.HTTP.URL != ""
	}
}

// getPipelineEndpoint finds the PipelineEndpoint with the given name in the Pipeline.
// It returns nil if no PipelineEndpoint is found.
func getPipelineEndpoint(pipeline *windtunnelv1alpha1.Pipeline, endpointName string) *windtunnelv1alpha1.PipelineEndpoint {
	for _, pipelineEndpoint := range pipeline.Spec.PipelineEndpoints {
		if pipelineEndpoint.Name == endpointName {
			return &pipelineEndpoint
		}
	}
	return nil
}

// getPipelineEndpointProtocol returns the protocol used by the PipelineEndpoint.
// It returns an empty string if no protocol is specified.
func getPipelineEndpointProtocol(pipelineEndpoint *windtunnelv1alpha1.PipelineEndpoint) windtunnelv1alpha1.EndpointProtocol {
	if pipelineEndpoint.HTTP != nil && pipelineEndpoint.HTTP.URL != "" && pipelineEndpoint.HTTP.Method != "" {
		return windtunnelv1alpha1.EndpointProtocolHTTP
	}
	return ""
}

// getEndpointSpecDataOption returns the data option used by the EndpointSpec.
// It returns an empty string if no data option is specified.
func getEndpointSpecDataOption(endpointSpec *windtunnelv1alpha1.EndpointSpec) windtunnelv1alpha1.EndpointDataOption {
	if endpointSpec.DataSpec == nil {
		return ""
	}

	if endpointSpec.DataSpec.DataSetRef != nil && endpointSpec.DataSpec.DataSetRef.Namespace != "" && endpointSpec.DataSpec.DataSetRef.Name != "" {
		return windtunnelv1alpha1.EndpointDataOptionDataSet
	}

	// Fallback to plain text
	return windtunnelv1alpha1.EndpointDataOptionPlainText
}

// getLoadPatternDuration calculates the duration of LoadPattern.
// It returns an error if the duration cannot be parsed.
func getLoadPatternDuration(loadPattern *windtunnelv1alpha1.LoadPattern) (*metav1.Duration, error) {
	duration := time.Duration(0)
	for _, stage := range loadPattern.Spec.Stages {
		stageDuration, err := time.ParseDuration(stage.Duration)
		if err != nil {
			return &metav1.Duration{}, err
		}
		duration += stageDuration
	}
	return &metav1.Duration{Duration: duration}, nil
}
