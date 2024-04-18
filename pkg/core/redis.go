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
	redisName                 = config.GetString("core.redis.name")
	redisLabels               = config.GetStringMapString("core.redis.labels")
	redisDefaultReplicas      = config.GetInt32("core.redis.defaultReplicas")
	redisDefaultImage         = config.GetString("core.redis.defaultImage")
	redisDefaultCPURequest    = config.GetString("core.redis.defaultCPURequest")
	redisDefaultMemoryRequest = config.GetString("core.redis.defaultMemoryRequest")
	redisDefaultCPULimit      = config.GetString("core.redis.defaultCPULimit")
	redisDefaultMemoryLimit   = config.GetString("core.redis.defaultMemoryLimit")
	redisContainerPortName    = config.GetString("core.redis.containerPortName")
	redisContainerPort        = config.GetInt32("core.redis.containerPort")
	redisServicePortName      = config.GetString("core.redis.servicePortName")
	redisServicePort          = config.GetInt32("core.redis.servicePort")
)

// GetRedisDeployment creates the Deployment for Redis.
func GetRedisDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	replicas := plantDCore.Spec.RedisConfig.Replicas
	if replicas == 0 {
		replicas = redisDefaultReplicas
	}

	image := plantDCore.Spec.RedisConfig.Image
	if image == "" {
		image = redisDefaultImage
	}

	resources := getResources(
		&plantDCore.Spec.RedisConfig.Resources,
		redisDefaultCPURequest,
		redisDefaultMemoryRequest,
		redisDefaultCPULimit,
		redisDefaultMemoryLimit,
	)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      redisName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: redisLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: redisLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:      redisName,
							Image:     image,
							Resources: resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          redisContainerPortName,
									ContainerPort: redisContainerPort,
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

// GetRedisService creates the Service for Redis.
func GetRedisService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      redisName,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       redisServicePortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       redisServicePort,
					TargetPort: intstr.FromString(redisContainerPortName),
				},
			},
			Selector: redisLabels,
		},
	}

	return service
}
