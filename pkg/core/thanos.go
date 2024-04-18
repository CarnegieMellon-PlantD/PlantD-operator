package core

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

var (
	thanosDefaultImage   = config.GetString("core.thanos.defaultImage")
	thanosDefaultVersion = config.GetString("core.thanos.defaultVersion")

	thanosSidecarDefaultCPURequest    = config.GetString("core.thanos.sidecar.defaultCPURequest")
	thanosSidecarDefaultMemoryRequest = config.GetString("core.thanos.sidecar.defaultMemoryRequest")
	thanosSidecarDefaultCPULimit      = config.GetString("core.thanos.sidecar.defaultCPULimit")
	thanosSidecarDefaultMemoryLimit   = config.GetString("core.thanos.sidecar.defaultMemoryLimit")
	thanosSidecarServicePortName      = config.GetString("core.thanos.sidecar.servicePortName")
	thanosSidecarServicePort          = config.GetInt32("core.thanos.sidecar.servicePort")

	thanosStoreName                  = config.GetString("core.thanos.store.name")
	thanosStoreLabels                = config.GetStringMapString("core.thanos.store.labels")
	thanosStoreDefaultReplicas       = config.GetInt32("core.thanos.store.defaultReplicas")
	thanosStoreDefaultCPURequest     = config.GetString("core.thanos.store.defaultCPURequest")
	thanosStoreDefaultMemoryRequest  = config.GetString("core.thanos.store.defaultMemoryRequest")
	thanosStoreDefaultCPULimit       = config.GetString("core.thanos.store.defaultCPULimit")
	thanosStoreDefaultMemoryLimit    = config.GetString("core.thanos.store.defaultMemoryLimit")
	thanosStoreDefaultStorageSize    = config.GetString("core.thanos.store.defaultStorageSize")
	thanosStoreContainerGrpcPortName = config.GetString("core.thanos.store.containerGrpcPortName")
	thanosStoreContainerGrpcPort     = config.GetInt32("core.thanos.store.containerGrpcPort")
	thanosStoreContainerHttpPortName = config.GetString("core.thanos.store.containerHttpPortName")
	thanosStoreContainerHttpPort     = config.GetInt32("core.thanos.store.containerHttpPort")
	thanosStoreServiceGrpcPortName   = config.GetString("core.thanos.store.serviceGrpcPortName")
	thanosStoreServiceGrpcPort       = config.GetInt32("core.thanos.store.serviceGrpcPort")
	thanosStoreServiceHttpPortName   = config.GetString("core.thanos.store.serviceHttpPortName")
	thanosStoreServiceHttpPort       = config.GetInt32("core.thanos.store.serviceHttpPort")
	thanosStorePath                  = config.GetString("core.thanos.store.path")

	thanosCompactorName                  = config.GetString("core.thanos.compactor.name")
	thanosCompactorLabels                = config.GetStringMapString("core.thanos.compactor.labels")
	thanosCompactorDefaultReplicas       = config.GetInt32("core.thanos.compactor.defaultReplicas")
	thanosCompactorDefaultCPURequest     = config.GetString("core.thanos.compactor.defaultCPURequest")
	thanosCompactorDefaultMemoryRequest  = config.GetString("core.thanos.compactor.defaultMemoryRequest")
	thanosCompactorDefaultCPULimit       = config.GetString("core.thanos.compactor.defaultCPULimit")
	thanosCompactorDefaultMemoryLimit    = config.GetString("core.thanos.compactor.defaultMemoryLimit")
	thanosCompactorDefaultStorageSize    = config.GetString("core.thanos.compactor.defaultStorageSize")
	thanosCompactorContainerGrpcPortName = config.GetString("core.thanos.compactor.containerGrpcPortName")
	thanosCompactorContainerGrpcPort     = config.GetInt32("core.thanos.compactor.containerGrpcPort")
	thanosCompactorContainerHttpPortName = config.GetString("core.thanos.compactor.containerHttpPortName")
	thanosCompactorContainerHttpPort     = config.GetInt32("core.thanos.compactor.containerHttpPort")
	thanosCompactorServiceGrpcPortName   = config.GetString("core.thanos.compactor.serviceGrpcPortName")
	thanosCompactorServiceGrpcPort       = config.GetInt32("core.thanos.compactor.serviceGrpcPort")
	thanosCompactorServiceHttpPortName   = config.GetString("core.thanos.compactor.serviceHttpPortName")
	thanosCompactorServiceHttpPort       = config.GetInt32("core.thanos.compactor.serviceHttpPort")
	thanosCompactorPath                  = config.GetString("core.thanos.compactor.path")

	thanosQuerierName                  = config.GetString("core.thanos.querier.name")
	thanosQuerierLabels                = config.GetStringMapString("core.thanos.querier.labels")
	thanosQuerierDefaultReplicas       = config.GetInt32("core.thanos.querier.defaultReplicas")
	thanosQuerierDefaultCPURequest     = config.GetString("core.thanos.querier.defaultCPURequest")
	thanosQuerierDefaultMemoryRequest  = config.GetString("core.thanos.querier.defaultMemoryRequest")
	thanosQuerierDefaultCPULimit       = config.GetString("core.thanos.querier.defaultCPULimit")
	thanosQuerierDefaultMemoryLimit    = config.GetString("core.thanos.querier.defaultMemoryLimit")
	thanosQuerierContainerGrpcPortName = config.GetString("core.thanos.querier.containerGrpcPortName")
	thanosQuerierContainerGrpcPort     = config.GetInt32("core.thanos.querier.containerGrpcPort")
	thanosQuerierContainerHttpPortName = config.GetString("core.thanos.querier.containerHttpPortName")
	thanosQuerierContainerHttpPort     = config.GetInt32("core.thanos.querier.containerHttpPort")
	thanosQuerierServiceGrpcPortName   = config.GetString("core.thanos.querier.serviceGrpcPortName")
	thanosQuerierServiceGrpcPort       = config.GetInt32("core.thanos.querier.serviceGrpcPort")
	thanosQuerierServiceHttpPortName   = config.GetString("core.thanos.querier.serviceHttpPortName")
	thanosQuerierServiceHttpPort       = config.GetInt32("core.thanos.querier.serviceHttpPort")
)

// GetThanosStoreStatefulSet creates the StatefulSet for Thanos-Store.
func GetThanosStoreStatefulSet(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.StatefulSet {
	image := plantDCore.Spec.ThanosConfig.Image
	if image == "" {
		image = thanosDefaultImage
	}

	replicas := plantDCore.Spec.ThanosConfig.StoreConfig.Replicas
	if replicas == 0 {
		replicas = thanosStoreDefaultReplicas
	}

	resources := getResources(
		&plantDCore.Spec.ThanosConfig.StoreConfig.Resources,
		thanosStoreDefaultCPURequest,
		thanosStoreDefaultMemoryRequest,
		thanosStoreDefaultCPULimit,
		thanosStoreDefaultMemoryLimit,
	)

	storageSize := plantDCore.Spec.ThanosConfig.StoreConfig.StorageSize
	if storageSize.IsZero() {
		storageSize = resource.MustParse(thanosStoreDefaultStorageSize)
	}

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      thanosStoreName,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: thanosStoreName,
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: thanosStoreLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: thanosStoreLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  thanosStoreName,
							Image: image,
							Args: []string{
								"store",
								fmt.Sprintf("--data-dir=%s", thanosStorePath),
								fmt.Sprintf("--grpc-address=0.0.0.0:%d", thanosStoreContainerGrpcPort),
								fmt.Sprintf("--http-address=0.0.0.0:%d", thanosStoreContainerHttpPort),
								"--objstore.config=$(OBJSTORE_CONFIG)",
							},
							Env: []corev1.EnvVar{
								{
									Name: "OBJSTORE_CONFIG",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: plantDCore.Spec.ThanosConfig.ObjectStoreConfig,
									},
								},
							},
							Resources: resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          thanosStoreContainerGrpcPortName,
									ContainerPort: thanosStoreContainerGrpcPort,
								},
								{
									Name:          thanosStoreContainerHttpPortName,
									ContainerPort: thanosStoreContainerHttpPort,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: thanosStorePath,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "data-volume",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: storageSize,
							},
						},
					},
				},
			},
		},
	}

	return statefulSet
}

// GetThanosStoreService creates the Service for Thanos-Store.
func GetThanosStoreService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      thanosStoreName,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       thanosStoreServiceGrpcPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       thanosStoreServiceGrpcPort,
					TargetPort: intstr.FromString(thanosStoreContainerGrpcPortName),
				},
				{
					Name:       thanosStoreServiceHttpPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       thanosStoreServiceHttpPort,
					TargetPort: intstr.FromString(thanosStoreContainerHttpPortName),
				},
			},
			Selector: thanosStoreLabels,
		},
	}

	return service
}

// GetThanosCompactorStatefulSet creates the StatefulSet for Thanos-Compactor.
func GetThanosCompactorStatefulSet(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.StatefulSet {
	image := plantDCore.Spec.ThanosConfig.Image
	if image == "" {
		image = thanosDefaultImage
	}

	replicas := plantDCore.Spec.ThanosConfig.CompactorConfig.Replicas
	if replicas == 0 {
		replicas = thanosCompactorDefaultReplicas
	}

	resources := getResources(
		&plantDCore.Spec.ThanosConfig.CompactorConfig.Resources,
		thanosCompactorDefaultCPURequest,
		thanosCompactorDefaultMemoryRequest,
		thanosCompactorDefaultCPULimit,
		thanosCompactorDefaultMemoryLimit,
	)

	storageSize := plantDCore.Spec.ThanosConfig.CompactorConfig.StorageSize
	if storageSize.IsZero() {
		storageSize = resource.MustParse(thanosCompactorDefaultStorageSize)
	}

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      thanosCompactorName,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: thanosCompactorName,
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: thanosCompactorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: thanosCompactorLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  thanosCompactorName,
							Image: image,
							Args: []string{
								"compact",
								"--wait",
								fmt.Sprintf("--data-dir=%s", thanosCompactorPath),
								fmt.Sprintf("--grpc-address=0.0.0.0:%d", thanosCompactorContainerGrpcPort),
								fmt.Sprintf("--http-address=0.0.0.0:%d", thanosCompactorContainerHttpPort),
								"--objstore.config=$(OBJSTORE_CONFIG)",
							},
							Env: []corev1.EnvVar{
								{
									Name: "OBJSTORE_CONFIG",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: plantDCore.Spec.ThanosConfig.ObjectStoreConfig,
									},
								},
							},
							Resources: resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          thanosCompactorContainerGrpcPortName,
									ContainerPort: thanosCompactorContainerGrpcPort,
								},
								{
									Name:          thanosCompactorContainerHttpPortName,
									ContainerPort: thanosCompactorContainerHttpPort,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: thanosCompactorPath,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "data-volume",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: storageSize,
							},
						},
					},
				},
			},
		},
	}

	return statefulSet
}

// GetThanosCompactorService creates the Service for Thanos-Compactor.
func GetThanosCompactorService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      thanosCompactorName,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       thanosCompactorServiceGrpcPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       thanosCompactorServiceGrpcPort,
					TargetPort: intstr.FromString(thanosCompactorContainerGrpcPortName),
				},
				{
					Name:       thanosCompactorServiceHttpPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       thanosCompactorServiceHttpPort,
					TargetPort: intstr.FromString(thanosCompactorContainerHttpPortName),
				},
			},
			Selector: thanosCompactorLabels,
		},
	}

	return service
}

// GetThanosQuerierStatefulSet creates the StatefulSet for Thanos-Querier.
func GetThanosQuerierStatefulSet(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.StatefulSet {
	image := plantDCore.Spec.ThanosConfig.Image
	if image == "" {
		image = thanosDefaultImage
	}

	replicas := plantDCore.Spec.ThanosConfig.QuerierConfig.Replicas
	if replicas == 0 {
		replicas = thanosQuerierDefaultReplicas
	}

	resources := getResources(
		&plantDCore.Spec.ThanosConfig.QuerierConfig.Resources,
		thanosQuerierDefaultCPURequest,
		thanosQuerierDefaultMemoryRequest,
		thanosQuerierDefaultCPULimit,
		thanosQuerierDefaultMemoryLimit,
	)

	containerArgs := []string{
		"query",
		fmt.Sprintf("--grpc-address=0.0.0.0:%d", thanosQuerierContainerGrpcPort),
		fmt.Sprintf("--http-address=0.0.0.0:%d", thanosQuerierContainerHttpPort),
		// Always query Thanos-Sidecar for data from Prometheus
		fmt.Sprintf("--endpoint=dnssrv+%s",
			utils.GetServiceSRVRecord(thanosSidecarServicePortName, "tcp", prometheusName, plantDCore.Namespace),
		),
	}
	// When upload is enabled, add Thanos-Store to the query endpoints
	if plantDCore.Spec.ThanosConfig.ObjectStoreConfig != nil {
		containerArgs = append(containerArgs, fmt.Sprintf("--endpoint=dnssrv+%s",
			utils.GetServiceSRVRecord(thanosStoreServiceGrpcPortName, "tcp", thanosStoreName, plantDCore.Namespace),
		))
	}

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      thanosQuerierName,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: thanosQuerierName,
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: thanosQuerierLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: thanosQuerierLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:      thanosQuerierName,
							Image:     image,
							Args:      containerArgs,
							Resources: resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          thanosQuerierContainerGrpcPortName,
									ContainerPort: thanosQuerierContainerGrpcPort,
								},
								{
									Name:          thanosQuerierContainerHttpPortName,
									ContainerPort: thanosQuerierContainerHttpPort,
								},
							},
						},
					},
				},
			},
		},
	}

	return statefulSet
}

// GetThanosQuerierService creates the Service for Thanos-Querier.
func GetThanosQuerierService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: plantDCore.Namespace,
			Name:      thanosQuerierName,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       thanosQuerierServiceGrpcPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       thanosQuerierServiceGrpcPort,
					TargetPort: intstr.FromString(thanosQuerierContainerGrpcPortName),
				},
				{
					Name:       thanosQuerierServiceHttpPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       thanosQuerierServiceHttpPort,
					TargetPort: intstr.FromString(thanosQuerierContainerHttpPortName),
				},
			},
			Selector: thanosQuerierLabels,
		},
	}
	return service
}
