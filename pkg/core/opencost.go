package core

import (
	"fmt"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

var (
	openCostName            = config.GetString("core.openCost.name")
	openCostLabels          = config.GetStringMapString("core.openCost.labels")
	openCostDefaultReplicas = config.GetInt32("core.openCost.defaultReplicas")

	openCostDefaultImage         = config.GetString("core.openCost.defaultImage")
	openCostDefaultCPURequest    = config.GetString("core.openCost.defaultCPURequest")
	openCostDefaultMemoryRequest = config.GetString("core.openCost.defaultMemoryRequest")
	openCostDefaultCPULimit      = config.GetString("core.openCost.defaultCPULimit")
	openCostDefaultMemoryLimit   = config.GetString("core.openCost.defaultMemoryLimit")
	openCostContainerPortName    = config.GetString("core.openCost.containerPortName")
	openCostContainerPort        = config.GetInt32("core.openCost.containerPort")
	openCostServicePortName      = config.GetString("core.openCost.servicePortName")
	openCostServicePort          = config.GetInt32("core.openCost.servicePort")

	openCostUIDefaultImage         = config.GetString("core.openCost.ui.defaultImage")
	openCostUIDefaultCPURequest    = config.GetString("core.openCost.ui.defaultCPURequest")
	openCostUIDefaultMemoryRequest = config.GetString("core.openCost.ui.defaultMemoryRequest")
	openCostUIDefaultCPULimit      = config.GetString("core.openCost.ui.defaultCPULimit")
	openCostUIDefaultMemoryLimit   = config.GetString("core.openCost.ui.defaultMemoryLimit")
	openCostUIContainerPortName    = config.GetString("core.openCost.ui.containerPortName")
	openCostUIContainerPort        = config.GetInt32("core.openCost.ui.containerPort")
	openCostUIServicePortName      = config.GetString("core.openCost.ui.servicePortName")
	openCostUIServicePort          = config.GetInt32("core.openCost.ui.servicePort")
)

// GetOpenCostDeployment creates the Deployment for OpenCost.
func GetOpenCostDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	replicas := plantDCore.Spec.OpenCostConfig.Replicas
	if replicas == 0 {
		replicas = openCostDefaultReplicas
	}

	image := plantDCore.Spec.OpenCostConfig.Image
	if image == "" {
		image = openCostDefaultImage
	}

	resources := getResources(
		&plantDCore.Spec.OpenCostConfig.Resources,
		openCostDefaultCPURequest,
		openCostDefaultMemoryRequest,
		openCostDefaultCPULimit,
		openCostDefaultMemoryLimit,
	)

	uiImage := plantDCore.Spec.OpenCostConfig.UIImage
	if uiImage == "" {
		uiImage = openCostUIDefaultImage
	}

	uiResources := getResources(
		&plantDCore.Spec.OpenCostConfig.UIResources,
		openCostUIDefaultCPURequest,
		openCostUIDefaultMemoryRequest,
		openCostUIDefaultCPULimit,
		openCostUIDefaultMemoryLimit,
	)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      openCostName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: openCostLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: openCostLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountOpenCost,
					Containers: []corev1.Container{
						{
							Name:      openCostName,
							Image:     image,
							Resources: resources,
							Env: []corev1.EnvVar{
								{
									Name: "PROMETHEUS_SERVER_ENDPOINT",
									Value: fmt.Sprintf("http://%s:%d",
										utils.GetServiceARecord(prometheusName, plantDCore.Namespace),
										prometheusServicePort,
									),
								},
								{
									Name:  "CLOUD_PROVIDER_API_KEY",
									Value: "",
								},
								// Default cluster ID to use if cluster_id is not set in Prometheus metrics
								{
									Name:  "CLUSTER_ID",
									Value: "default-cluster-id",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          openCostContainerPortName,
									ContainerPort: openCostContainerPort,
								},
							},
						},
						{
							Name:      openCostName + "-ui",
							Image:     uiImage,
							Resources: uiResources,
							Ports: []corev1.ContainerPort{
								{
									Name:          openCostUIContainerPortName,
									ContainerPort: openCostUIContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

// GetOpenCostService creates the Service for OpenCost.
func GetOpenCostService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      openCostName,
			// Add labels so that OpenCost ServiceMonitor can select it
			Labels: openCostLabels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       openCostServicePortName,
					Port:       openCostServicePort,
					TargetPort: intstr.FromString(openCostContainerPortName),
				},
				{
					Name:       openCostUIServicePortName,
					Port:       openCostUIServicePort,
					TargetPort: intstr.FromString(openCostUIContainerPortName),
				},
			},
			Selector: openCostLabels,
		},
	}

	return service
}

// GetOpenCostServiceMonitor creates the ServiceMonitor for OpenCost.
func GetOpenCostServiceMonitor(plantDCore *windtunnelv1alpha1.PlantDCore) *monitoringv1.ServiceMonitor {
	serviceMonitor := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      openCostName,
			// Set the labels so that Prometheus can select it
			Labels: prometheusServiceMonitorLabels,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			JobLabel: "experiment",
			// Use labels to select OpenCost Service
			Selector: metav1.LabelSelector{
				MatchLabels: openCostLabels,
			},
			Endpoints: []monitoringv1.Endpoint{
				{
					Port:        openCostServicePortName,
					HonorLabels: true,
				},
			},
		},
	}
	return serviceMonitor
}

// GetCAdvisorServiceMonitor creates the ServiceMonitor to monitor cAdvisor.
func GetCAdvisorServiceMonitor(plantDCore *windtunnelv1alpha1.PlantDCore) *monitoringv1.ServiceMonitor {
	serviceMonitor := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      openCostName + "-cadvisor",
			// Set the labels so that Prometheus can select it
			Labels: prometheusServiceMonitorLabels,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			JobLabel: "kubelet",
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "kubelet",
				},
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{
					"kube-system",
				},
			},
			Endpoints: []monitoringv1.Endpoint{
				{
					Port:            "https-metrics",
					Scheme:          "https",
					Interval:        "30s",
					TLSConfig:       &monitoringv1.TLSConfig{SafeTLSConfig: monitoringv1.SafeTLSConfig{InsecureSkipVerify: true}},
					BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
					HonorLabels:     true,
				},
				{
					Port:            "https-metrics",
					Path:            "/metrics/cadvisor",
					Scheme:          "https",
					Interval:        "30s",
					TLSConfig:       &monitoringv1.TLSConfig{SafeTLSConfig: monitoringv1.SafeTLSConfig{InsecureSkipVerify: true}},
					BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
					HonorLabels:     true,
				},
			},
		},
	}
	return serviceMonitor
}
