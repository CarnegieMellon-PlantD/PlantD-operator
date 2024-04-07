package core

import (
	"fmt"

	"k8s.io/utils/ptr"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	plantDCoreServiceAccountName     = config.GetViper().GetString("plantDCore.serviceAccountName")
	kubeProxyLabels                  = config.GetViper().GetStringMapString("plantDCore.kubeProxy.labels")
	kubeProxyContainerName           = config.GetViper().GetString("plantDCore.kubeProxy.containerName")
	kubeProxyImage                   = config.GetViper().GetString("plantDCore.kubeProxy.image")
	kubeProxyDeploymentName          = config.GetViper().GetString("plantDCore.kubeProxy.deploymentName")
	kubeProxyReplicas                = config.GetViper().GetInt32("plantDCore.kubeProxy.replicas")
	kubeProxyServiceName             = config.GetViper().GetString("plantDCore.kubeProxy.serviceName")
	kubeProxyPort                    = config.GetViper().GetInt32("plantDCore.kubeProxy.port")
	kubeProxyTargetPort              = config.GetViper().GetInt32("plantDCore.kubeProxy.targetPort")
	studioLabels                     = config.GetViper().GetStringMapString("plantDCore.studio.labels")
	studioContainerName              = config.GetViper().GetString("plantDCore.studio.containerName")
	studioImage                      = config.GetViper().GetString("plantDCore.studio.image")
	studioDeploymentName             = config.GetViper().GetString("plantDCore.studio.deploymentName")
	studioReplicas                   = config.GetViper().GetInt32("plantDCore.studio.replicas")
	studioServiceName                = config.GetViper().GetString("plantDCore.studio.serviceName")
	studioPort                       = config.GetViper().GetInt32("plantDCore.studio.port")
	studioTargetPort                 = config.GetViper().GetInt32("plantDCore.studio.targetPort")
	prometheusLabels                 = config.GetViper().GetStringMapString("plantDCore.prometheus.labels")
	prometheusServiceMonitorSelector = config.GetViper().GetStringMapString("monitor.serviceMonitor.labels")
	prometheusServiceAccountName     = config.GetViper().GetString("plantDCore.prometheus.serviceAccountName")
	prometheusClusterRoleName        = config.GetViper().GetString("plantDCore.prometheus.clusterRoleName")
	prometheusClusterRoleBindingName = config.GetViper().GetString("plantDCore.prometheus.clusterRoleBindingName")
	prometheusObjectName             = config.GetViper().GetString("plantDCore.prometheus.name")
	prometheusScrapeInterval         = config.GetViper().GetString("plantDCore.prometheus.scrapeInterval")
	prometheusReplicas               = config.GetViper().GetInt32("plantDCore.prometheus.replicas")
	prometheusMemoryLimit            = config.GetViper().GetString("plantDCore.prometheus.memoryLimit")
	prometheusServiceName            = config.GetViper().GetString("plantDCore.prometheus.serviceName")
	prometheusPort                   = config.GetViper().GetInt32("plantDCore.prometheus.port")
	prometheusTargetPort             = config.GetViper().GetInt32("plantDCore.prometheus.targetPort")
	prometheusNodePort               = config.GetViper().GetInt32("plantDCore.prometheus.nodePort")
	redisLabels                      = config.GetViper().GetStringMapString("plantDCore.redis.labels")
	redisContainerName               = config.GetViper().GetString("plantDCore.redis.containerName")
	redisImage                       = config.GetViper().GetString("plantDCore.redis.image")
	redisDeploymentName              = config.GetViper().GetString("plantDCore.redis.deploymentName")
	redisReplicas                    = config.GetViper().GetInt32("plantDCore.redis.replicas")
	redisServiceName                 = config.GetViper().GetString("plantDCore.redis.serviceName")
	redisPort                        = config.GetViper().GetInt32("plantDCore.redis.port")
	redisTargetPort                  = config.GetViper().GetInt32("plantDCore.redis.targetPort")
)

// GetKubeProxyServiceFQDN returns the in-cluster DNS name of PlantD Kube Proxy Service
func GetKubeProxyServiceFQDN(plantDCore *windtunnelv1alpha1.PlantDCore) string {
	return fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", kubeProxyServiceName, plantDCore.Namespace, kubeProxyPort)
}

// GetStudioServiceFQDN returns the in-cluster DNS name of PlantD Studio Service
func GetStudioServiceFQDN(plantDCore *windtunnelv1alpha1.PlantDCore) string {
	return fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", studioServiceName, plantDCore.Namespace, studioPort)
}

// GetPrometheusServiceFQDN returns the in-cluster DNS name of Prometheus Service
func GetPrometheusServiceFQDN(plantDCore *windtunnelv1alpha1.PlantDCore) string {
	return fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", prometheusServiceName, plantDCore.Namespace, prometheusPort)
}

// GetRedisServiceFQDN returns the in-cluster DNS name of Redis Service
func GetRedisServiceFQDN(plantDCore *windtunnelv1alpha1.PlantDCore) string {
	return fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", redisServiceName, plantDCore.Namespace, redisPort)
}

// GetKubeProxyDeployment returns the Deployment for PlantD Kube Proxy
func GetKubeProxyDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	// Define the pod template
	image := plantDCore.Spec.KubeProxyConfig.Image
	if image == "" {
		image = kubeProxyImage
	}

	resourceRequirements := plantDCore.Spec.KubeProxyConfig.Resources

	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: kubeProxyLabels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: plantDCoreServiceAccountName,
			Containers: []corev1.Container{
				{
					Name:            kubeProxyContainerName,
					Image:           image,
					ImagePullPolicy: corev1.PullAlways,
					Resources:       resourceRequirements,
				},
			},
		},
	}

	// Define the Deployment
	numReplicas := plantDCore.Spec.KubeProxyConfig.Replicas
	if numReplicas == 0 {
		numReplicas = kubeProxyReplicas
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeProxyDeploymentName,
			Namespace: plantDCore.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: kubeProxyLabels,
			},
			Template: podTemplate,
		},
	}

	return deployment
}

// GetKubeProxyService returns the Service for PlantD Kube Proxy
func GetKubeProxyService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	// Define the Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeProxyServiceName,
			Namespace: plantDCore.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: kubeProxyLabels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       kubeProxyPort,
					TargetPort: intstr.FromInt32(kubeProxyTargetPort),
				},
			},
		},
	}

	return service
}

// GetStudioDeployment returns the Deployment for PlantD Studio
func GetStudioDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	// Define the pod template
	image := plantDCore.Spec.StudioConfig.Image
	if image == "" {
		image = studioImage
	}

	resourceRequirements := plantDCore.Spec.StudioConfig.Resources

	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: studioLabels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            studioContainerName,
					Image:           image,
					ImagePullPolicy: corev1.PullAlways,
					Env: []corev1.EnvVar{
						{
							Name:  "KUBE_PROXY_URL",
							Value: GetKubeProxyServiceFQDN(plantDCore),
						},
					},
					Resources: resourceRequirements,
				},
			},
		},
	}

	// Define the Deployment
	numReplicas := plantDCore.Spec.StudioConfig.Replicas
	if numReplicas == 0 {
		numReplicas = studioReplicas
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      studioDeploymentName,
			Namespace: plantDCore.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: studioLabels,
			},
			Template: podTemplate,
		},
	}

	return deployment
}

// GetStudioService returns the Service for PlantD Studio
func GetStudioService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	// Define the Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      studioServiceName,
			Namespace: plantDCore.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: studioLabels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       studioPort,
					TargetPort: intstr.FromInt32(studioTargetPort),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer, // Use LoadBalancer type
		},
	}

	return service
}

// GetPrometheusRBACResources returns the ServiceAccount, ClusterRole, and ClusterRoleBinding for Prometheus
func GetPrometheusRBACResources(plantDCore *windtunnelv1alpha1.PlantDCore) (*corev1.ServiceAccount, *rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusServiceAccountName,
			Namespace: plantDCore.Namespace,
		},
	}

	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusClusterRoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{
					"nodes",
					"nodes/metrics",
					"services",
					"endpoints",
					"pods",
				},
				Verbs: []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				NonResourceURLs: []string{"/metrics"},
				Verbs:           []string{"get"},
			},
		},
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusClusterRoleBindingName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     prometheusClusterRoleName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      prometheusServiceAccountName,
				Namespace: plantDCore.Namespace,
			},
		},
	}

	return sa, clusterRole, clusterRoleBinding
}

// GetPrometheusObject returns the Prometheus object for Prometheus
func GetPrometheusObject(plantDCore *windtunnelv1alpha1.PlantDCore) *monitoringv1.Prometheus {
	// Define the Prometheus
	scrapeInterval := plantDCore.Spec.PrometheusConfig.ScrapeInterval
	if scrapeInterval == "" {
		scrapeInterval = monitoringv1.Duration(prometheusScrapeInterval)
	}

	numReplicas := plantDCore.Spec.PrometheusConfig.Replicas
	if numReplicas == 0 {
		numReplicas = prometheusReplicas
	}

	resourceRequirements := plantDCore.Spec.PrometheusConfig.Resources
	if resourceRequirements.Limits == nil {
		resourceRequirements.Limits = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(prometheusMemoryLimit),
		}
	} else if resourceRequirements.Limits.Memory().IsZero() {
		resourceRequirements.Limits[corev1.ResourceMemory] = resource.MustParse(prometheusMemoryLimit)
	}

	prometheus := &monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusObjectName,
			Namespace: plantDCore.Namespace,
		},
		Spec: monitoringv1.PrometheusSpec{
			CommonPrometheusFields: monitoringv1.CommonPrometheusFields{
				ServiceAccountName: prometheusServiceAccountName,
				ServiceMonitorSelector: &metav1.LabelSelector{
					MatchLabels: prometheusServiceMonitorSelector,
				},
				ServiceMonitorNamespaceSelector: &metav1.LabelSelector{},
				ScrapeInterval:                  scrapeInterval,
				EnableRemoteWriteReceiver:       true,
				Replicas:                        &numReplicas,
				Resources:                       resourceRequirements,
				SecurityContext: &corev1.PodSecurityContext{
					RunAsUser:    ptr.To(config.GetViper().GetInt64("plantDCore.prometheus.securityContext.runAsUser")),
					RunAsNonRoot: ptr.To(config.GetViper().GetBool("plantDCore.prometheus.securityContext.runAsNonRoot")),
					FSGroup:      ptr.To(config.GetViper().GetInt64("plantDCore.prometheus.securityContext.fsGroup")),
					RunAsGroup:   ptr.To(config.GetViper().GetInt64("plantDCore.prometheus.securityContext.runAsGroup")),
				},
			},
			Thanos: &monitoringv1.ThanosSpec{
				BaseImage: ptr.To(config.GetViper().GetString("plantDCore.prometheus.thanos.thanosBaseImage")),
				Version:   ptr.To(config.GetViper().GetString("plantDCore.prometheus.thanos.thanosVersion")),
			},
			EnableAdminAPI: false,
		},
	}
	if plantDCore.Spec.ThanosEnabled {
		prometheus.Spec.Thanos.ObjectStorageConfig = &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: config.GetViper().GetString("plantDCore.prometheus.thanos.thanosConfig.name"),
			},
			Key: config.GetViper().GetString("plantDCore.prometheus.thanos.thanosConfig.key"),
		}
	}
	return prometheus
}

// GetPrometheusService returns the Service for Prometheus
func GetPrometheusService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	// Define the Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusServiceName,
			Namespace: plantDCore.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       prometheusPort,
					TargetPort: intstr.FromInt32(prometheusTargetPort),
					NodePort:   prometheusNodePort,
				},
			},
			Selector: prometheusLabels,
		},
	}
	return service
}

// GetRedisDeployment returns the Deployment for Redis
func GetRedisDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	// Define the pod template
	image := plantDCore.Spec.RedisConfig.Image
	if image == "" {
		image = redisImage
	}

	resourceRequirements := plantDCore.Spec.RedisConfig.Resources

	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: redisLabels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: plantDCoreServiceAccountName,
			Containers: []corev1.Container{
				{
					Name:            redisContainerName,
					Image:           image,
					ImagePullPolicy: corev1.PullAlways,
					Resources:       resourceRequirements,
				},
			},
		},
	}

	// Define the Deployment
	numReplicas := redisReplicas
	if plantDCore.Spec.RedisConfig.Replicas != 0 {
		numReplicas = plantDCore.Spec.RedisConfig.Replicas
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redisDeploymentName,
			Namespace: plantDCore.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: redisLabels,
			},
			Template: podTemplate,
		},
	}

	return deployment
}

// GetRedisService returns the Service for Redis
func GetRedisService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	// Define the Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redisServiceName,
			Namespace: plantDCore.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: redisLabels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       redisPort,
					TargetPort: intstr.FromInt32(redisTargetPort),
				},
			},
		},
	}

	return service
}
func GetThanosQuerierService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.GetViper().GetString("plantDCore.prometheus.thanosQuerier.deploymentName"),
			Namespace: plantDCore.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       config.GetViper().GetInt32("plantDCore.prometheus.thanosQuerier.httpPort"),
					TargetPort: intstr.FromInt32(config.GetViper().GetInt32("plantDCore.prometheus.thanosQuerier.httpPort")),
				},
			},
			Selector: config.GetViper().GetStringMapString("plantDCore.prometheus.thanosQuerier.labels"),
		},
	}
	return service
}
func GetThanosSidecarService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "thanos-sidecar-grpc",
			Namespace: plantDCore.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       config.GetViper().GetString("plantDCore.prometheus.thanosSidecarService.portName"),
					Protocol:   corev1.ProtocolTCP,
					Port:       config.GetViper().GetInt32("plantDCore.prometheus.thanosSidecarService.port"),
					TargetPort: intstr.FromInt32(config.GetViper().GetInt32("plantDCore.prometheus.thanosSidecarService.port")),
				},
			},
			Selector: config.GetViper().GetStringMapString("plantDCore.prometheus.labels"),
		},
	}
	return service
}
func GetThanosQuerierDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	containerArgs := []string{
		"query",
		fmt.Sprintf("--http-address=0.0.0.0:%s", config.GetViper().GetString("plantDCore.prometheus.thanosQuerier.httpPort")),
		fmt.Sprintf("--grpc-address=0.0.0.0:%s", config.GetViper().GetString("plantDCore.prometheus.thanosQuerier.grpcPort")),
		fmt.Sprintf("--store=%s", config.GetViper().GetString("plantDCore.prometheus.thanosQuerier.url")),
	}
	if plantDCore.Spec.ThanosEnabled {
		containerArgs = append(containerArgs, fmt.Sprintf("--store=%s", config.GetViper().GetString("plantDCore.prometheus.thanosStore.url")))
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "thanos-querier",
			Namespace: plantDCore.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(config.GetViper().GetInt32("plantDCore.prometheus.thanosQuerier.replicas")),
			Selector: &metav1.LabelSelector{
				MatchLabels: config.GetViper().GetStringMapString("plantDCore.prometheus.thanosQuerier.labels"),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: config.GetViper().GetStringMapString("plantDCore.prometheus.thanosQuerier.labels"),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "thanos-querier",
							Image: config.GetViper().GetString("plantDCore.prometheus.thanosQuerier.image"),
							Args:  containerArgs,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: config.GetViper().GetInt32("plantDCore.prometheus.thanosQuerier.httpPort"),
								},
								{
									Name:          "grpc",
									ContainerPort: config.GetViper().GetInt32("plantDCore.prometheus.thanosQuerier.grpcPort"),
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

// GetThanosStoreStatefulSet Thanos Store Stateful Set
func GetThanosStoreStatefulSet(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.StatefulSet {
	volumeClaimSpec := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(config.GetViper().GetString("plantDCore.prometheus.thanosStore.volumeSize")),
			},
		},
	}
	volumeClaimTemplates := []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "thanos-data",
			},
			Spec: volumeClaimSpec,
		},
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.GetViper().GetString("plantDCore.prometheus.thanosStore.name"),
			Namespace: plantDCore.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "thanos-store",
			Replicas:    ptr.To(config.GetViper().GetInt32("plantDCore.prometheus.thanosStore.replicas")),
			Selector: &metav1.LabelSelector{
				MatchLabels: config.GetViper().GetStringMapString("plantDCore.prometheus.thanosStore.labels"),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: config.GetViper().GetStringMapString("plantDCore.prometheus.thanosStore.labels"),
				},
				Spec: corev1.PodSpec{
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser:  ptr.To(config.GetViper().GetInt64("plantDCore.prometheus.thanosStore.securityContext.runAsUser")),
						RunAsGroup: ptr.To(config.GetViper().GetInt64("plantDCore.prometheus.thanosStore.securityContext.runAsGroup")),
						FSGroup:    ptr.To(config.GetViper().GetInt64("plantDCore.prometheus.thanosStore.securityContext.fsGroup")),
					},
					Containers: []corev1.Container{
						{
							Name:  config.GetViper().GetString("plantDCore.prometheus.thanosStore.name"),
							Image: config.GetViper().GetString("plantDCore.prometheus.thanosStore.image"),
							Args: []string{
								"store",
								fmt.Sprintf("--data-dir=%s", config.GetViper().GetString("plantDCore.prometheus.thanosStore.dataDir")),
								"--objstore.config=$(OBJSTORE_CONFIG)",
							},
							Ports: []corev1.ContainerPort{
								{Name: "grpc", ContainerPort: config.GetViper().GetInt32("plantDCore.prometheus.thanosStore.grpcPort")},
								{Name: "http", ContainerPort: config.GetViper().GetInt32("plantDCore.prometheus.thanosStore.httpPort")},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "thanos-data", MountPath: config.GetViper().GetString("plantDCore.prometheus.thanosStore.dataDir")},
							},
							Env: []corev1.EnvVar{
								{
									Name: "OBJSTORE_CONFIG",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: config.GetViper().GetString("plantDCore.prometheus.thanos.thanosConfig.name")},
											Key:                  config.GetViper().GetString("plantDCore.prometheus.thanos.thanosConfig.key"),
										},
									},
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: volumeClaimTemplates,
		},
	}
}
func GetThanosStoreService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.GetViper().GetString("plantDCore.prometheus.thanosStore.serviceName"),
			Namespace: plantDCore.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Protocol:   corev1.ProtocolTCP,
					Port:       config.GetViper().GetInt32("plantDCore.prometheus.thanosStore.grpcPort"),
					TargetPort: intstr.FromInt32(config.GetViper().GetInt32("plantDCore.prometheus.thanosStore.grpcPort")),
				},
			},
			Selector: config.GetViper().GetStringMapString("plantDCore.prometheus.thanosStore.labels"),
		},
	}
	return service
}

// Opencost resources
func GetOpencostRBACResources(plantDCore *windtunnelv1alpha1.PlantDCore) (*corev1.ServiceAccount, *rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.GetViper().GetString("opencost.serviceAccount"),
			Namespace: plantDCore.Namespace,
		},
	}

	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.GetViper().GetString("opencost.clusterRole"),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{
					"configmaps",
					"deployments",
					"nodes",
					"pods",
					"services",
					"resourcequotas",
					"replicationcontrollers",
					"limitranges",
					"persistentvolumeclaims",
					"persistentvolumes",
					"namespaces",
					"endpoints",
				},
				Verbs: []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"extensions"},
				Resources: []string{"daemonsets", "deployments", "replicasets"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"statefulsets", "deployments", "daemonsets", "replicasets"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{"cronjobs", "jobs"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"autoscaling"},
				Resources: []string{"horizontalpodautoscalers"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"poddisruptionbudgets"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.GetViper().GetString("opencost.clusterRole"),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     config.GetViper().GetString("opencost.clusterRole"),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      config.GetViper().GetString("opencost.clusterRole"),
				Namespace: plantDCore.Namespace,
			},
		},
	}
	return sa, clusterRole, clusterRoleBinding
}
func GetOpencostDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	labels := config.GetViper().GetStringMapString("opencost.deployment.labels")

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "opencost",
			Namespace: plantDCore.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(config.GetViper().GetInt32("opencost.deployment.replicas")),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       ptr.To(intstr.FromInt32(config.GetViper().GetInt32("opencost.deployment.maxSurge"))),
					MaxUnavailable: ptr.To(intstr.FromInt32(config.GetViper().GetInt32("opencost.deployment.maxUnavailable"))),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					RestartPolicy:      corev1.RestartPolicyAlways,
					ServiceAccountName: "opencost",
					Containers: []corev1.Container{
						{
							Name:  "opencost",
							Image: config.GetViper().GetString("opencost.deployment.image"),
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resourceQuantity("10m"),
									corev1.ResourceMemory: resourceQuantity("55M"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resourceQuantity("999m"),
									corev1.ResourceMemory: resourceQuantity("1G"),
								},
							},
							Env: []corev1.EnvVar{
								{Name: "PROMETHEUS_SERVER_ENDPOINT", Value: config.GetViper().GetString("database.prometheus.url")},
								{Name: "CLOUD_PROVIDER_API_KEY", Value: ""},
								{Name: "CLUSTER_ID", Value: "cluster-one"},
								{Name: "LOG_LEVEL", Value: "debug"},
							},
							ImagePullPolicy: corev1.PullAlways,
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: ptr.To(false),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
								Privileged:             ptr.To(false),
								ReadOnlyRootFilesystem: ptr.To(true),
								RunAsUser:              ptr.To(config.GetViper().GetInt64("opencost.deployment.runAsUser")),
							},
						},
						{
							Name:  "opencost-ui",
							Image: config.GetViper().GetString("opencost.deployment.ui-image"),
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resourceQuantity("10m"),
									corev1.ResourceMemory: resourceQuantity("55M"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resourceQuantity("999m"),
									corev1.ResourceMemory: resourceQuantity("1G"),
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
				},
			},
		},
	}
	return deployment
}
func GetOpencostService(plantDCore *windtunnelv1alpha1.PlantDCore) *corev1.Service {
	labels := config.GetViper().GetStringMapString("opencost.service.labels")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "opencost",
			Namespace: plantDCore.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "opencost",
					Port:       config.GetViper().GetInt32("opencost.service.port"),
					TargetPort: intstr.FromInt32(config.GetViper().GetInt32("opencost.service.port")),
				},
				{
					Name:       "opencost-ui",
					Port:       config.GetViper().GetInt32("opencost.service.ui-port"),
					TargetPort: intstr.FromInt32(config.GetViper().GetInt32("opencost.service.ui-port")),
				},
			},
		},
	}
	return service
}
func GetOpencostServiceMonitor(plantDCore *windtunnelv1alpha1.PlantDCore) *monitoringv1.ServiceMonitor {
	serviceMonitor := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.GetViper().GetString("opencost.opencostServiceMonitor.name"),
			Namespace: plantDCore.Namespace,
			Labels:    config.GetViper().GetStringMapString("opencost.opencostServiceMonitor.labels"),
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: []monitoringv1.Endpoint{
				{
					Port:        "opencost",
					HonorLabels: true,
				},
			},
			JobLabel: "experiment",
			Selector: metav1.LabelSelector{
				MatchLabels: config.GetViper().GetStringMapString("opencost.opencostServiceMonitor.selector"),
			},
		},
	}
	return serviceMonitor
}
func GetCadvisorServiceMonitor(plantDCore *windtunnelv1alpha1.PlantDCore) *monitoringv1.ServiceMonitor {
	serviceMonitor := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.GetViper().GetString("opencost.cadvisorServiceMonitor.name"),
			Namespace: plantDCore.Namespace,
			Labels:    config.GetViper().GetStringMapString("cadvisor.cadvisorServiceMonitor.labels"),
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: []monitoringv1.Endpoint{
				{
					BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
					HonorLabels:     true,
					Interval:        "30s",
					Port:            "https-metrics",
					Scheme:          "https",
					TLSConfig:       &monitoringv1.TLSConfig{SafeTLSConfig: monitoringv1.SafeTLSConfig{InsecureSkipVerify: true}},
				},
				{
					BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
					Port:            "https-metrics",
					Scheme:          "https",
					TLSConfig:       &monitoringv1.TLSConfig{SafeTLSConfig: monitoringv1.SafeTLSConfig{InsecureSkipVerify: true}},
					HonorLabels:     true,
					Interval:        "30s",
					Path:            "/metrics/cadvisor",
				},
			},
			JobLabel: "kubelet",
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{config.GetViper().GetString("opencost.cadvisorServiceMonitor.namespaceSelector")},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: config.GetViper().GetStringMapString("opencost.cadvisorServiceMonitor.selector"),
			},
		},
	}
	return serviceMonitor
}

func resourceQuantity(quantityStr string) resource.Quantity {
	quantity, _ := resource.ParseQuantity(quantityStr)
	return quantity
}
