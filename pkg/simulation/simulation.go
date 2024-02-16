package simulation

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
func CreateJobBySimulation(ctx context.Context, jobName string, simulation *windtunnelv1alpha1.Simulation,
	digitalTwin *windtunnelv1alpha1.DigitalTwin, trafficModel *windtunnelv1alpha1.TrafficModel,
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
							Name:  "TRAFFIC_MODEL",
							Value: trafficModel.Spec.Config,
						},
						{
							Name:  "TRAFFIC_MODEL_NAME",
							Value: trafficModel.Name,
						},
						{
							Name:  "SIM_NAME",
							Value: simulation.Namespace + "." + simulation.Name,
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
							Value: config.GetString("costService.opencost.url"),
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
