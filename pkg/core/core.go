package core

import (
	"log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var plantDCoreServiceAccountName string
var kubeProxyImage string
var kubeProxyContainer string
var kubeProxyLabel map[string]string
var kubeProxySelectorKey string
var kubeProxySelectorValue string
var kubeProxyReplicas int
var kubeProxyDeploymentName string
var kubeProxyServiceName string
var kubeProxyPort int32
var kubeProxyTargetPort int32

var studioImage string
var studioContainer string
var studioLabel map[string]string
var studioSelectorKey string
var studioSelectorValue string
var studioReplicas int
var studioDeploymentName string
var studioServiceName string
var studioPortName string
var studioPort int32
var studioTargetPort int32

var prometheusServiceAccountName string
var prometheusObjectName string
var prometheusResourceMemory string
var prometheusScrapeInterval string
var prometheusClusterRoleName string
var prometheusSelectorKey string
var prometheusSelectorValue string
var prometheusClusterRoleBindingName string
var prometheusMetricSpecSelector map[string]string
var prometheusServicePort int32
var prometheusServiceNodePort int32

// var prometheuDefaultScrapInterval model.Duration

func init() {
	plantDCoreServiceAccountName = config.GetString("plantdCore.serviceAccountName")

	kubeProxyImage = config.GetString("plantdCore.kubeProxy.image")
	kubeProxyLabel = config.GetStringMapString("plantdCore.kubeProxy.label")
	kubeProxySelectorKey = config.GetString("plantdCore.kubeProxy.selector.key")
	kubeProxySelectorValue = config.GetString("plantdCore.kubeProxy.selector.value")
	kubeProxyContainer = config.GetString("plantdCore.kubeProxy.containerName")
	kubeProxyReplicas = config.GetInt("plantdCore.kubeProxy.replicas")
	kubeProxyDeploymentName = config.GetString("plantdCore.kubeProxy.deploymentName")
	kubeProxyServiceName = config.GetString("plantdCore.kubeProxy.serviceName")
	kubeProxyPort = config.GetInt32("plantdCore.kubeProxy.port")
	kubeProxyTargetPort = config.GetInt32("plantdCore.kubeProxy.targetPort")

	studioImage = config.GetString("plantdCore.studio.image")
	studioLabel = config.GetStringMapString("plantdCore.studio.label")
	studioSelectorKey = config.GetString("plantdCore.studio.selector.key")
	studioSelectorValue = config.GetString("plantdCore.studio.selector.value")
	studioContainer = config.GetString("plantdCore.studio.containerName")
	studioReplicas = config.GetInt("plantdCore.studio.replicas")
	studioDeploymentName = config.GetString("plantdCore.studio.deploymentName")
	studioServiceName = config.GetString("plantdCore.studio.serviceName")
	studioPortName = config.GetString("plantdCore.studio.portName")
	studioPort = config.GetInt32("plantdCore.studio.port")
	studioTargetPort = config.GetInt32("plantdCore.studio.targetPort")

	prometheusServiceAccountName = config.GetString("plantdCore.prometheus.serviceAccount")
	prometheusObjectName = config.GetString("plantdCore.prometheus.name")
	prometheusResourceMemory = config.GetString("plantdCore.prometheus.resourceMemory")
	prometheusScrapeInterval = config.GetString("plantdCore.prometheus.scrapeInterval")
	prometheusMetricSpecSelector = config.GetStringMapString("plantdCore.prometheus.specSelector")
	prometheusSelectorKey = config.GetString("plantdCore.prometheus.selector.key")
	prometheusSelectorValue = config.GetString("plantdCore.prometheus.selector.value")
	prometheusServicePort = config.GetInt32("plantdCore.prometheus.service.port")
	prometheusServiceNodePort = config.GetInt32("plantdCore.prometheus.service.nodePort")
	prometheusClusterRoleName = config.GetString("plantdCore.prometheus.clusteRoleName")
	prometheusClusterRoleBindingName = config.GetString("plantdCore.prometheus.clusterRoleBindingName")
}

// SetupKubeProxyDeployment creates kube proxy deployment with Cluster IP
func SetupProxyDeployment(plantD *windtunnelv1alpha1.PlantDCore) (*appsv1.Deployment, *corev1.Service) {

	numReplicas := int32(kubeProxyReplicas)
	// Define the pod template

	labels := map[string]string{
		kubeProxySelectorKey: kubeProxySelectorValue,
	}
	log.Println(plantDCoreServiceAccountName)
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: plantDCoreServiceAccountName,
			Containers: []corev1.Container{
				{
					Name:            kubeProxyContainer,
					Image:           kubeProxyImage,
					ImagePullPolicy: corev1.PullAlways,
				},
			},
		},
	}

	// Define the Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeProxyDeploymentName,
			Namespace: plantD.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numReplicas, // Set the number of replicas as needed
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: podTemplate,
		},
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeProxyServiceName,
			Namespace: plantD.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: kubeProxyLabel,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       kubeProxyPort,
					TargetPort: intstr.FromInt(int(kubeProxyTargetPort)),
				},
			},
		},
	}
	return deployment, service
}

// SetupFrontendDeployment creates a PlantD Frontend deployment
func SetupFrontendDeployment(plantD *windtunnelv1alpha1.PlantDCore, proxyFQDN string) (*appsv1.Deployment, *corev1.Service) {

	numReplicas := int32(studioReplicas)
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: studioLabel,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            studioContainer,
					Image:           studioImage,
					ImagePullPolicy: corev1.PullAlways,
					Env: []corev1.EnvVar{
						{
							Name:  "KUBE_PROXY_URL",
							Value: string(proxyFQDN),
						},
					},
				},
			},
		},
	}

	// Define the Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      studioDeploymentName,
			Namespace: plantD.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					studioSelectorKey: studioSelectorValue,
				},
			},
			Template: podTemplate,
		},
	}

	// Define the Service with a LoadBalancer
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      studioServiceName,
			Namespace: plantD.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: studioLabel,
			Ports: []corev1.ServicePort{
				{
					Name:       studioPortName,
					Port:       studioPort, // Specify the port as needed
					TargetPort: intstr.FromInt(int(studioTargetPort)),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer, // Use LoadBalancer type
		},
	}
	return deployment, service
}

// SetupFrontendDeployment creates a PlantD Frontend deployment
func SetupPrometheusObject(plantD *windtunnelv1alpha1.PlantDCore) (*monitoringv1.Prometheus, *corev1.Service) {
	// Define the Prometheus resource
	var scrapeInterval = monitoringv1.Duration(prometheusScrapeInterval)
	var resourceMemory = resource.MustParse(prometheusResourceMemory)

	if plantD.Spec.PrometheusConfiguration.ScrapeInterval != "" {
		scrapeInterval = monitoringv1.Duration(plantD.Spec.PrometheusConfiguration.ScrapeInterval)
	}

	if plantD.Spec.PrometheusConfiguration.ResourceMemory.Limits != nil {
		if !plantD.Spec.PrometheusConfiguration.ResourceMemory.Limits.Memory().IsZero() {
			resourceMemory = resource.MustParse(plantD.Spec.PrometheusConfiguration.ResourceMemory.Limits.Memory().String())
		}
	}

	prometheus := &monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusObjectName,
			Namespace: plantD.Namespace,
		},
		Spec: monitoringv1.PrometheusSpec{
			CommonPrometheusFields: monitoringv1.CommonPrometheusFields{
				ServiceAccountName: prometheusServiceAccountName,
				ServiceMonitorSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						prometheusSelectorKey: prometheusSelectorValue,
					},
				},
				ServiceMonitorNamespaceSelector: &metav1.LabelSelector{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceMemory: resourceMemory,
					},
				},
				EnableRemoteWriteReceiver: true,
				ScrapeInterval:            scrapeInterval,
			},
			EnableAdminAPI: false,
		},
	}

	// Define the Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusObjectName,
			Namespace: plantD.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       "web",
					NodePort:   prometheusServiceNodePort,
					Port:       prometheusServicePort,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("web"),
				},
			},
			Selector: prometheusMetricSpecSelector,
		},
	}

	return prometheus, service
}

// SetupFrontendDeployment creates a PlantD Frontend deployment
func SetupRoleBindingsForPrometheus(plantD *windtunnelv1alpha1.PlantDCore) (*corev1.ServiceAccount, *rbac.ClusterRole, *rbac.ClusterRoleBinding) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusServiceAccountName,
			Namespace: plantD.Namespace,
		},
	}

	clusterRole := &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusClusterRoleName,
		},
		Rules: []rbac.PolicyRule{
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

	clusterRoleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusClusterRoleBindingName,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     prometheusClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind,
				Name:      prometheusServiceAccountName,
				Namespace: plantD.Namespace,
			},
		},
	}
	return sa, clusterRole, clusterRoleBinding
}
