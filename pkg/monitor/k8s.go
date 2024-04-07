package monitor

import (
	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	serviceLabelKeyPipeline   = config.GetViper().GetString("monitor.service.labelKeys.pipeline")
	serviceLabelKeyExperiment = config.GetViper().GetString("monitor.service.labelKeys.experiment")
	servicePortName           = config.GetViper().GetString("monitor.service.portName")
	serviceMonitorLabels      = config.GetViper().GetStringMapString("monitor.serviceMonitor.labels")
	defaultEndpointPort       = config.GetViper().GetString("monitor.serviceMonitor.endpoint.defaultPort")
	defaultEndpointPath       = config.GetViper().GetString("monitor.serviceMonitor.endpoint.defaultPath")
)

// CreateExternalNameService creates a metrics Service of type ExternalName. For out-cluster Pipeline only.
func CreateExternalNameService(pipeline *windtunnelv1alpha1.Pipeline) (*corev1.Service, error) {
	hostname, err := pipeline.Spec.MetricsEndpoint.HTTP.GetHostname()
	if err != nil {
		return nil, err
	}
	port, err := pipeline.Spec.MetricsEndpoint.HTTP.GetPort()
	if err != nil {
		return nil, err
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: pipeline.Namespace,
			Name:      utils.GetPipelineMetricsServiceName(pipeline.Name),
			// Set the Pipeline label so that ServiceMonitor can select it
			Labels: map[string]string{
				serviceLabelKeyPipeline: pipeline.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: hostname,
			Ports: []corev1.ServicePort{
				{
					Name:     servicePortName,
					Protocol: corev1.ProtocolTCP,
					Port:     int32(port),
				},
			},
		},
	}
	return service, nil
}

// CreateServiceMonitor creates a ServiceMonitor for a Pipeline's metrics Service.
func CreateServiceMonitor(pipeline *windtunnelv1alpha1.Pipeline) (*monitoringv1.ServiceMonitor, error) {
	var metricsServiceNamespace string
	if pipeline.Spec.InCluster {
		// If the Pipeline is in-cluster, the user creates the metrics Service, and it can be in any namespace
		metricsServiceNamespace = pipeline.Spec.MetricsEndpoint.ServiceRef.Namespace
	} else {
		// If the Pipeline is out-cluster, we have created the metrics Service of type ExternalName,
		// and it is in the same namespace as the Pipeline
		metricsServiceNamespace = pipeline.Namespace
	}

	endpointPort := pipeline.Spec.MetricsEndpoint.Port
	if endpointPort == "" {
		endpointPort = defaultEndpointPort
	}
	endpointPath := pipeline.Spec.MetricsEndpoint.Path
	if endpointPath == "" {
		endpointPath = defaultEndpointPath
	}

	serviceMonitor := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: pipeline.Namespace,
			Name:      pipeline.Name,
			// Set the labels so that Prometheus can select it
			Labels: serviceMonitorLabels,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			// Set `job` label of all Prometheus metrics to Experiment label value
			JobLabel: serviceLabelKeyExperiment,
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{metricsServiceNamespace},
			},
			// Use the Pipeline label to select Service
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					serviceLabelKeyPipeline: pipeline.Name,
				},
			},
			Endpoints: []monitoringv1.Endpoint{
				{
					Port: endpointPort,
					Path: endpointPath,
				},
			},
		},
	}
	return serviceMonitor, nil
}
