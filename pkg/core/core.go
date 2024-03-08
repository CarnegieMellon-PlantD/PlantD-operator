package core

import (
	"fmt"
	"k8s.io/utils/pointer"

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
	plantDCoreServiceAccountName string

	kubeProxyLabels         map[string]string
	kubeProxyContainerName  string
	kubeProxyImage          string
	kubeProxyDeploymentName string
	kubeProxyReplicas       int32
	kubeProxyServiceName    string
	kubeProxyPort           int32
	kubeProxyTargetPort     int32

	studioLabels         map[string]string
	studioContainerName  string
	studioImage          string
	studioDeploymentName string
	studioReplicas       int32
	studioServiceName    string
	studioPort           int32
	studioTargetPort     int32

	prometheusLabels                 map[string]string
	prometheusServiceMonitorSelector map[string]string
	prometheusServiceAccountName     string
	prometheusClusterRoleName        string
	prometheusClusterRoleBindingName string
	prometheusObjectName             string
	prometheusScrapeInterval         string
	prometheusReplicas               int32
	prometheusMemoryLimit            string
	prometheusServiceName            string
	prometheusPort                   int32
	prometheusTargetPort             int32
	prometheusNodePort               int32

	redisLabels         map[string]string
	redisContainerName  string
	redisImage          string
	redisDeploymentName string
	redisReplicas       int32
	redisServiceName    string
	redisPort           int32
	redisTargetPort     int32
)

func init() {
	plantDCoreServiceAccountName = config.GetString("plantDCore.serviceAccountName")

	kubeProxyLabels = config.GetStringMapString("plantDCore.kubeProxy.labels")
	kubeProxyContainerName = config.GetString("plantDCore.kubeProxy.containerName")
	kubeProxyImage = config.GetString("plantDCore.kubeProxy.image")
	kubeProxyDeploymentName = config.GetString("plantDCore.kubeProxy.deploymentName")
	kubeProxyReplicas = config.GetInt32("plantDCore.kubeProxy.replicas")
	kubeProxyServiceName = config.GetString("plantDCore.kubeProxy.serviceName")
	kubeProxyPort = config.GetInt32("plantDCore.kubeProxy.port")
	kubeProxyTargetPort = config.GetInt32("plantDCore.kubeProxy.targetPort")

	studioLabels = config.GetStringMapString("plantDCore.studio.labels")
	studioContainerName = config.GetString("plantDCore.studio.containerName")
	studioImage = config.GetString("plantDCore.studio.image")
	studioDeploymentName = config.GetString("plantDCore.studio.deploymentName")
	studioReplicas = config.GetInt32("plantDCore.studio.replicas")
	studioServiceName = config.GetString("plantDCore.studio.serviceName")
	studioPort = config.GetInt32("plantDCore.studio.port")
	studioTargetPort = config.GetInt32("plantDCore.studio.targetPort")

	prometheusLabels = config.GetStringMapString("plantDCore.prometheus.labels")
	prometheusServiceMonitorSelector = config.GetStringMapString("plantDCore.prometheus.serviceMonitorSelector")
	prometheusServiceAccountName = config.GetString("plantDCore.prometheus.serviceAccountName")
	prometheusClusterRoleName = config.GetString("plantDCore.prometheus.clusterRoleName")
	prometheusClusterRoleBindingName = config.GetString("plantDCore.prometheus.clusterRoleBindingName")
	prometheusObjectName = config.GetString("plantDCore.prometheus.name")
	prometheusScrapeInterval = config.GetString("plantDCore.prometheus.scrapeInterval")
	prometheusReplicas = config.GetInt32("plantDCore.prometheus.replicas")
	prometheusMemoryLimit = config.GetString("plantDCore.prometheus.memoryLimit")
	prometheusServiceName = config.GetString("plantDCore.prometheus.serviceName")
	prometheusPort = config.GetInt32("plantDCore.prometheus.port")
	prometheusTargetPort = config.GetInt32("plantDCore.prometheus.targetPort")
	prometheusNodePort = config.GetInt32("plantDCore.prometheus.nodePort")

	redisLabels = config.GetStringMapString("plantDCore.redis.labels")
	redisContainerName = config.GetString("plantDCore.redis.containerName")
	redisImage = config.GetString("plantDCore.redis.image")
	redisDeploymentName = config.GetString("plantDCore.redis.deploymentName")
	redisReplicas = config.GetInt32("plantDCore.redis.replicas")
	redisServiceName = config.GetString("plantDCore.redis.serviceName")
	redisPort = config.GetInt32("plantDCore.redis.port")
	redisTargetPort = config.GetInt32("plantDCore.redis.targetPort")
}

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

	thanosBaseImage := "quay.io/thanos/thanos"
	thanosVersion := "v0.34.0"

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
					RunAsUser:    pointer.Int64Ptr(1000),
					RunAsNonRoot: pointer.BoolPtr(true),
					FSGroup:      pointer.Int64Ptr(2000),
					RunAsGroup:   pointer.Int64Ptr(2000),
				},
			},
			EnableAdminAPI: false,
		},
	}
	if plantDCore.Spec.ThanosEnabled {
		prometheus.Spec.Thanos = &monitoringv1.ThanosSpec{
			BaseImage: &thanosBaseImage,
			Version:   &thanosVersion,
			ObjectStorageConfig: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "thanos-objstore-config",
				},
				Key: "thanos.yaml",
			},
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
			Name:      "thanos-querier",
			Namespace: plantDCore.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       9090,
					TargetPort: intstr.FromInt(9090),
				},
			},
			Selector: map[string]string{
				"app": "thanos-querier",
			},
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
					Name:       "grpc",
					Port:       10901,
					TargetPort: intstr.FromInt(10901),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"prometheus": "prometheus",
			},
		},
	}
	return service
}
func GetThanosQuerierDeployment(plantDCore *windtunnelv1alpha1.PlantDCore) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "thanos-querier",
			Namespace: plantDCore.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "thanos-querier",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "thanos-querier",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "thanos-querier",
							Image: "quay.io/thanos/thanos:v0.34.0",
							Args: []string{
								"query",
								"--http-address=0.0.0.0:9090",
								"--grpc-address=0.0.0.0:10901",
								"--store=thanos-sidecar-grpc.plantd-operator-system.svc.cluster.local:10901",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 9090,
								},
								{
									Name:          "grpc",
									ContainerPort: 10901,
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
func int32Ptr(i int32) *int32 { return &i }
