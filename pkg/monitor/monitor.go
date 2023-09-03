package monitor

import (
	"fmt"
	"net"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	pipelineLabels       map[string]string
	serviceMonitorLabels map[string]string
	metricsLabelKey      string
	metricsPortName      string
	experimentLabelKey   string
)

func init() {
	pipelineLabels = config.GetStringMapString("monitor.pipelineEndpoint.labels")
	serviceMonitorLabels = config.GetStringMapString("monitor.serviceMonitor.labels")
	metricsLabelKey = config.GetString("monitor.metricsService.labels.key")
	metricsPortName = config.GetString("monitor.metrics.port.name")
	experimentLabelKey = config.GetString("monitor.jobLabel")
}

func CreateServiceMonitor(pipeline *windtunnelv1alpha1.Pipeline) (*monitoringv1.ServiceMonitor, error) {
	var endpoints []monitoringv1.Endpoint
	// For in-cluster pipeline-under-test, user need to create the k8s service and provide port name.
	if pipeline.Spec.InCluster {
		endpoints = append(endpoints, pipeline.Spec.MetricsEndpoint.Endpoint)
	} else {
		metricsPath, err := utils.GetHTTPPath(pipeline.Spec.MetricsEndpoint.HTTP.URL)
		if err != nil {
			return nil, err
		}
		endpoint := monitoringv1.Endpoint{
			Port: metricsPortName,
			Path: metricsPath,
		}
		endpoints = append(endpoints, endpoint)
	}

	return &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pipeline.Name,
			Namespace: pipeline.Namespace,
			Labels:    serviceMonitorLabels,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			JobLabel: experimentLabelKey,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					metricsLabelKey: utils.GetNamespacedName(pipeline),
				},
			},
			Endpoints: endpoints,
		},
	}, nil
}

func CreateExternalNameService(name string, namespace string, endpoint *windtunnelv1alpha1.Endpoint) (*corev1.Service, *corev1.Endpoints, error) {
	hostname, err := utils.GetHostname(endpoint.HTTP.URL)
	if err != nil {
		return nil, nil, err
	}
	port, err := utils.GetPort(endpoint.HTTP.URL)
	if err != nil {
		return nil, nil, err
	}

	ipAddresses, err := net.LookupIP(hostname)
	if err != nil {
		return nil, nil, err
	}
	addresses := make([]corev1.EndpointAddress, len(ipAddresses))
	for i, ip := range ipAddresses {
		addresses[i].IP = ip.String()
	}

	portName := endpoint.Name
	label := pipelineLabels
	// Metrics Endpoint should not have an endpoint name
	if portName == "" {
		portName = metricsPortName
		label = map[string]string{metricsLabelKey: fmt.Sprintf("%s-%s", namespace, name)}
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    label,
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: hostname,
			Ports: []corev1.ServicePort{
				{
					Name:     portName,
					Protocol: corev1.ProtocolTCP,
					Port:     port,
				},
			},
		},
	}
	endpoints := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    label,
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: addresses,
				Ports: []corev1.EndpointPort{
					{
						Name:     portName,
						Port:     port,
						Protocol: corev1.ProtocolTCP,
					},
				},
			},
		},
	}
	return service, endpoints, nil
}
