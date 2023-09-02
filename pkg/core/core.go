package core

import (
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

var kubeProxyImage string
var frontendImage string
var prometheusObjectName string
var prometheusMetricLabelSelector string

// var prometheuDefaultScrapInterval model.Duration

func init() {
	kubeProxyImage = config.GetString("plantdCore.kubeProxyImage")
	frontendImage = config.GetString("plantdCore.frontendImage")
	prometheusObjectName = config.GetString("plantdCore.prometheusObjectName")

	prometheusMetricLabelSelector = config.GetString("plantdCore.prometheusMetricLabelSelector")
	// var prometheusScrapeIntervalString = config.GetString("plantdCore.prometheuDefaultScrapInterval")
	// prometheuDefaultScrapInterval, _ = model.ParseDuration(prometheusScrapeIntervalString)

}

// SetupKubeProxyDeployment creates kube proxy deployment with Cluster IP
func SetupProxyDeployment(plantD *windtunnelv1alpha1.PlantDCore) (*appsv1.Deployment, *corev1.Service) {
	// Create the Kubernetes Job object
	// Define labels and annotations as needed
	labels := map[string]string{
		"app": "kube-proxy",
	}
	numReplicas := int32(1)
	// Define the pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "kube-proxy",
					Image:           kubeProxyImage,
					ImagePullPolicy: corev1.PullAlways,
				},
			},
		},
	}

	// Define the Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kube-proxy-deployment",
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

	// Define the Service with a LoadBalancer
	// Define the Service with a ClusterIP
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kube-proxy-service",
			Namespace: plantD.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       5000, // Specify the port as needed
					TargetPort: intstr.FromInt(5000),
				},
			},
			Type: corev1.ServiceTypeClusterIP, // Use ClusterIP type
		},
	}

	return deployment, service
}

// SetupFrontendDeployment creates a PlantD Frontend deployment
func SetupFrontendDeployment(plantD *windtunnelv1alpha1.PlantDCore, proxyFQDN string) (*appsv1.Deployment, *corev1.Service) {
	// Create the Kubernetes Job object
	// Define labels and annotations as needed
	labels := map[string]string{
		"app": "frontend",
	}
	numReplicas := int32(1)
	// Define the pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "frontend",
					Image:           frontendImage,
					ImagePullPolicy: corev1.PullAlways,
					// Add environment variables as needed
					Env: []corev1.EnvVar{
						{
							Name:  "KUBE_PROXY_URL",
							Value: string(proxyFQDN),
						},
						// Add more environment variables here...
					},
				},
			},
		},
	}

	// Define the Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "frontend-deployment",
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

	// Define the Service with a LoadBalancer
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "frontend-service",
			Namespace: plantD.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80, // Specify the port as needed
					TargetPort: intstr.FromInt(8080),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer, // Use LoadBalancer type
		},
	}

	return deployment, service
}

// SetupFrontendDeployment creates a PlantD Frontend deployment
func SetupCreatePrometheusSMObject(plantD *windtunnelv1alpha1.PlantDCore) (*monitoringv1.Prometheus, *corev1.Service) {
	// Define the Prometheus resource
	prometheus := &monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusObjectName,
			Namespace: plantD.Namespace,
		},
		Spec: monitoringv1.PrometheusSpec{
			CommonPrometheusFields: monitoringv1.CommonPrometheusFields{
				ServiceAccountName: "prometheus-operator",
				ServiceMonitorSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"component": prometheusMetricLabelSelector,
					},
				},
				ServiceMonitorNamespaceSelector: &metav1.LabelSelector{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceMemory: resource.MustParse("1Gi"),
					},
				},
				EnableRemoteWriteReceiver: true,
				ScrapeInterval:            "15s",
			},
			EnableAdminAPI: false,
		},
	}

	// Define the Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus",
			Namespace: plantD.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       "web",
					NodePort:   30900,
					Port:       9090,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("web"),
				},
			},
			Selector: map[string]string{
				"prometheus": "prometheus",
			},
		},
	}

	return prometheus, service
}

// SetupFrontendDeployment creates a PlantD Frontend deployment
func SetupRoleBindingsForPrometheus(plantD *windtunnelv1alpha1.PlantDCore) (*rbac.ClusterRole, *rbac.ClusterRoleBinding) {
	clusterRole := &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "plantd-prometheus-role",
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
			Name: "plantd-prometheus-role-binding",
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "plantd-prometheus-role-binding",
		},
		Subjects: []rbac.Subject{
			{
				Kind:     "ServiceAccount",
				Name:     "system:serviceaccounts",
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
	}
	return clusterRole, clusterRoleBinding
}
