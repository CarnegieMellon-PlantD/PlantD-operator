package core

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

var (
	studioName                 = config.GetString("core.studio.name")
	studioLabels               = config.GetStringMapString("core.studio.labels")
	studioDefaultReplicas      = config.GetInt32("core.studio.defaultReplicas")
	studioDefaultImage         = config.GetString("core.studio.defaultImage")
	studioDefaultCPURequest    = config.GetString("core.studio.defaultCPURequest")
	studioDefaultMemoryRequest = config.GetString("core.studio.defaultMemoryRequest")
	studioDefaultCPULimit      = config.GetString("core.studio.defaultCPULimit")
	studioDefaultMemoryLimit   = config.GetString("core.studio.defaultMemoryLimit")
	studioContainerPortName    = config.GetString("core.studio.containerPortName")
	studioContainerPort        = config.GetInt32("core.studio.containerPort")
	studioServicePortName      = config.GetString("core.studio.servicePortName")
	studioServicePort          = config.GetInt32("core.studio.servicePort")
)

// GetStudioDeployment creates the Deployment for PlantD-Studio.
func GetStudioDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	replicas := plantDCore.Spec.StudioConfig.Replicas
	if replicas == 0 {
		replicas = studioDefaultReplicas
	}

	image := plantDCore.Spec.StudioConfig.Image
	if image == "" {
		image = studioDefaultImage
	}

	resources := getResources(
		&plantDCore.Spec.StudioConfig.Resources,
		studioDefaultCPURequest,
		studioDefaultMemoryRequest,
		studioDefaultCPULimit,
		studioDefaultMemoryLimit,
	)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      studioName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: studioLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: studioLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  studioName,
							Image: image,
							Env: []corev1.EnvVar{
								{
									Name: "KUBE_PROXY_URL",
									Value: fmt.Sprintf("http://%s:%d",
										utils.GetServiceARecord(proxyName, plantDCore.Namespace),
										proxyServicePort,
									),
								},
							},
							Resources: resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          studioContainerPortName,
									ContainerPort: studioContainerPort,
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

// GetStudioService creates the Service for PlantD-Studio.
func GetStudioService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      studioName,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{
					Name:       studioServicePortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       studioServicePort,
					TargetPort: intstr.FromString(studioContainerPortName),
				},
			},
			Selector: studioLabels,
		},
	}

	return service
}
