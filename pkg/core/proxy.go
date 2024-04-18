package core

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
)

var (
	proxyName                 = config.GetString("core.proxy.name")
	proxyLabels               = config.GetStringMapString("core.proxy.labels")
	proxyDefaultReplicas      = config.GetInt32("core.proxy.defaultReplicas")
	proxyDefaultImage         = config.GetString("core.proxy.defaultImage")
	proxyDefaultCPURequest    = config.GetString("core.proxy.defaultCPURequest")
	proxyDefaultMemoryRequest = config.GetString("core.proxy.defaultMemoryRequest")
	proxyDefaultCPULimit      = config.GetString("core.proxy.defaultCPULimit")
	proxyDefaultMemoryLimit   = config.GetString("core.proxy.defaultMemoryLimit")
	proxyContainerPortName    = config.GetString("core.proxy.containerPortName")
	proxyContainerPort        = config.GetInt32("core.proxy.containerPort")
	proxyServicePortName      = config.GetString("core.proxy.servicePortName")
	proxyServicePort          = config.GetInt32("core.proxy.servicePort")
)

// GetProxyDeployment creates the Deployment for PlantD-Proxy.
func GetProxyDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	replicas := plantDCore.Spec.ProxyConfig.Replicas
	if replicas == 0 {
		replicas = proxyDefaultReplicas
	}

	image := plantDCore.Spec.ProxyConfig.Image
	if image == "" {
		image = proxyDefaultImage
	}

	resources := getResources(
		&plantDCore.Spec.ProxyConfig.Resources,
		proxyDefaultCPURequest,
		proxyDefaultMemoryRequest,
		proxyDefaultCPULimit,
		proxyDefaultMemoryLimit,
	)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      proxyName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: proxyLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: proxyLabels,
				},
				Spec: corev1.PodSpec{
					// Share the ServiceAccount of controller-manager,
					// so that PlantD-Proxy has the same permissions.
					ServiceAccountName: serviceAccountControllerManager,
					Containers: []corev1.Container{
						{
							Name:      proxyName,
							Image:     image,
							Resources: resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          proxyContainerPortName,
									ContainerPort: proxyContainerPort,
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

// GetProxyService creates the Service for PlantD-Proxy.
func GetProxyService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      proxyName,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       proxyServicePortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       proxyServicePort,
					TargetPort: intstr.FromString(proxyContainerPortName),
				},
			},
			Selector: proxyLabels,
		},
	}

	return service
}
