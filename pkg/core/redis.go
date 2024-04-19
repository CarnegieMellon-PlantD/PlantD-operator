package core

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

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
	redisDefaultStorageSize   = config.GetString("core.redis.defaultStorageSize")
	redisContainerPortName    = config.GetString("core.redis.containerPortName")
	redisContainerPort        = config.GetInt32("core.redis.containerPort")
	redisServicePortName      = config.GetString("core.redis.servicePortName")
	redisServicePort          = config.GetInt32("core.redis.servicePort")
	redisPath                 = config.GetString("core.redis.path")
)

// GetRedisStatefulSet creates the StatefulSet for Redis.
func GetRedisStatefulSet(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.StatefulSet {
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

	var storageSize resource.Quantity
	if plantDCore.Spec.RedisConfig.StorageSize != nil && !plantDCore.Spec.RedisConfig.StorageSize.IsZero() {
		storageSize = *plantDCore.Spec.RedisConfig.StorageSize
	} else {
		storageSize = resource.MustParse(redisDefaultStorageSize)
	}

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      redisName,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: redisName,
			Replicas:    &replicas,
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
							Name:  redisName,
							Image: image,
							Env: []corev1.EnvVar{
								{
									Name:  "REDIS_ARGS",
									Value: "--save 3600 1 --save 1800 10 --save 600 100",
								},
							},
							Resources: resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          redisContainerPortName,
									ContainerPort: redisContainerPort,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: redisPath,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "data",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: storageSize,
							},
						},
						VolumeMode: ptr.To(corev1.PersistentVolumeFilesystem),
					},
				},
			},
		},
	}

	return statefulSet
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
