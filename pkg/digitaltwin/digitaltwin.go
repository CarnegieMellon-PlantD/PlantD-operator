package digitaltwin

import (
	"context"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var image string

func init() {
	image = config.GetString("digitalTwin.image")
}

// CreateJobByDigitalTwin creates a Kubernetes Job based on the Digital Twinconfiguration.
func CreateJobByDigitalTwin(ctx context.Context, jobName string, digitalTwin *windtunnelv1alpha1.DigitalTwin,
	experimentListJson string, loadPatternListJson string) (*corev1.Pod, error) {

	// Create the Kubernetes Job object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        jobName,
			Namespace:   digitalTwin.Namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,

			Containers: []corev1.Container{
				{
					Name:            jobName,
					Image:           image,
					ImagePullPolicy: corev1.PullAlways,
					Env: []corev1.EnvVar{
						{
							Name:  "TWIN_NAME",
							Value: digitalTwin.Name,
						},
						{
							Name:  "PIPELINE_NAMESPACE",
							Value: "plantd-operator-system",
						},
						{
							Name:  "MODEL_TYPE",
							Value: string(digitalTwin.Spec.ModelType),
						},
						{
							Name:  "REDIS_PASSWORD",
							Value: "",
						},
						{
							Name:  "REDIS_HOST",
							Value: config.GetString("database.redis.host"),
						},
						{
							Name:  "PROMETHEUS_HOST",
							Value: config.GetString("database.prometheus.url"),
						},

						{
							Name:  "OPENCOST_ENDPOINT",
							Value: "http://opencost.opencost.svc.cluster.local:9003",
						},
						{
							Name:  "PROMETHEUS_PASSWORD",
							Value: "",
						},
						{
							Name:  "LOAD_PATTERN_NAMES",
							Value: string(digitalTwin.Spec.LoadPatternNames),
						},
						{
							Name:  "EXPERIMENT_NAMES",
							Value: string(digitalTwin.Spec.ExperimentNames),
						},
						{
							Name:  "LOAD_PATTERN_METADATA",
							Value: loadPatternListJson,
						},
						{
							Name:  "EXPERIMENT_METADATA",
							Value: experimentListJson,
						},
					},
				},
			},
		},
	}
	return pod, nil
}
