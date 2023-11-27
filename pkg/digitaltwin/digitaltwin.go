package digitaltwin

import (
	"context"
	"strconv"
	"time"

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
							Name:  "MODEL_TYPE",
							Value: string(digitalTwin.Spec.ModelType),
						},
						{
							Name:  "REDIS_PASSWORD",
							Value: "",
						},
						{
							Name:  "REDIS_HOST",
							Value: "redis.plantd-operator-system.svc.cluster.local",
						},
						{
							Name:  "PROMETHEUS_HOST",
							Value: "prometheus.plantd-operator-system.svc.cluster.local",
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
							Value: string(loadPatternListJson),
						},
						{
							Name:  "EXPERIMENT_METADATA",
							Value: string(experimentListJson),
						},
					},
				},
			},
		},
	}
	return pod, nil
}