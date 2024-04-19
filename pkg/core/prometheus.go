package core

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
)

var (
	prometheusServiceMonitorLabels  = config.GetStringMapString("monitor.serviceMonitor.labels")
	prometheusName                  = config.GetString("core.prometheus.name")
	prometheusDefaultScrapeInterval = config.GetString("core.prometheus.defaultScrapeInterval")
	prometheusDefaultReplicas       = config.GetInt32("core.prometheus.defaultReplicas")
	prometheusDefaultCPURequest     = config.GetString("core.prometheus.defaultCPURequest")
	prometheusDefaultMemoryRequest  = config.GetString("core.prometheus.defaultMemoryRequest")
	prometheusDefaultCPULimit       = config.GetString("core.prometheus.defaultCPULimit")
	prometheusDefaultMemoryLimit    = config.GetString("core.prometheus.defaultMemoryLimit")
	prometheusServicePortName       = config.GetString("core.prometheus.servicePortName")
	prometheusServicePort           = config.GetInt32("core.prometheus.servicePort")
)

// GetPrometheusObject creates the Prometheus object for Prometheus.
func GetPrometheusObject(plantDCore *windtunnelv1alpha1.PlantDCore) *monitoringv1.Prometheus {
	replicas := plantDCore.Spec.PrometheusConfig.Replicas
	if replicas == 0 {
		replicas = prometheusDefaultReplicas
	}

	scrapeInterval := plantDCore.Spec.PrometheusConfig.ScrapeInterval
	if scrapeInterval == "" {
		scrapeInterval = monitoringv1.Duration(prometheusDefaultScrapeInterval)
	}

	resources := getResources(
		&plantDCore.Spec.PrometheusConfig.Resources,
		prometheusDefaultCPURequest,
		prometheusDefaultMemoryRequest,
		prometheusDefaultCPULimit,
		prometheusDefaultMemoryLimit,
	)

	thanosSidecarImage := plantDCore.Spec.ThanosConfig.Image
	if thanosSidecarImage == "" {
		thanosSidecarImage = thanosDefaultImage
	}

	thanosSidecarVersion := plantDCore.Spec.ThanosConfig.Version
	if thanosSidecarVersion == "" {
		thanosSidecarVersion = thanosDefaultVersion
	}

	thanosSidecarResources := getResources(
		&plantDCore.Spec.ThanosConfig.SidecarConfig.Resources,
		thanosSidecarDefaultCPURequest,
		thanosSidecarDefaultMemoryRequest,
		thanosSidecarDefaultCPULimit,
		thanosSidecarDefaultMemoryLimit,
	)

	prometheus := &monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      prometheusName,
		},
		Spec: monitoringv1.PrometheusSpec{
			CommonPrometheusFields: monitoringv1.CommonPrometheusFields{
				ServiceAccountName: serviceAccountPrometheus,
				// Use the labels to select ServiceMonitors
				ServiceMonitorSelector: &metav1.LabelSelector{
					MatchLabels: prometheusServiceMonitorLabels,
				},
				// The ServiceMonitors can be in any namespaces
				ServiceMonitorNamespaceSelector: &metav1.LabelSelector{},
				ScrapeInterval:                  scrapeInterval,
				Replicas:                        &replicas,
				Resources:                       resources,
				EnableRemoteWriteReceiver:       true,
				// Prometheus runs as user 65534 (nobody) by default, and Thanos runs as user 1001 by default.
				// Set the `SecurityContext.RunAsUser` to change the user Prometheus runs as,
				// so that Thanos can read the data Prometheus writes.
				// See https://github.com/prometheus-operator/prometheus-operator/issues/4664.
				SecurityContext: &corev1.PodSecurityContext{
					RunAsUser: ptr.To(int64(1001)),
				},
			},
			// Always deploy Thanos-Sidecar, as Thanos-Querier uses it to query data from Prometheus
			Thanos: &monitoringv1.ThanosSpec{
				BaseImage: &thanosSidecarImage,
				Version:   &thanosSidecarVersion,
				Resources: thanosSidecarResources,
			},
			EnableAdminAPI: false,
		},
	}

	// When upload is enabled, add object store config to Thanos-Sidecar
	if plantDCore.Spec.ThanosConfig.ObjectStoreConfig != nil {
		prometheus.Spec.Thanos.ObjectStorageConfig = plantDCore.Spec.ThanosConfig.ObjectStoreConfig
	}

	return prometheus
}

// GetPrometheusService creates the Service for Prometheus.
func GetPrometheusService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      prometheusName,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:     prometheusServicePortName,
					Protocol: corev1.ProtocolTCP,
					Port:     prometheusServicePort,
					// Prometheus container exposes its web interface on a port called "web"
					TargetPort: intstr.FromString("web"),
				},
				{
					Name:     thanosSidecarServicePortName,
					Protocol: corev1.ProtocolTCP,
					Port:     thanosSidecarServicePort,
					// Thanos-Sidecar container exposes its gRPC interface on a port called "grpc"
					TargetPort: intstr.FromString("grpc"),
				},
			},
			// Prometheus operator adds the "prometheus" label to Pods,
			// where the label value is the name of the Prometheus object that manages the Pods.
			Selector: map[string]string{
				"prometheus": prometheusName,
			},
		},
	}
	return service
}
