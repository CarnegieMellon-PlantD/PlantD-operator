package cost

import (
	"encoding/json"
	"strconv"
	"time"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

var (
	defaultImage = config.GetString("costService.defaultImage")
	redisHost    = utils.GetServiceARecord(config.GetString("core.redis.name"), config.GetString("core.namespace"))
	redisPort    = config.GetInt32("core.redis.servicePort")
)

// Tag defines a key-value pair for tagging cloud resources
type Tag struct {
	// Tag key
	Key string `json:"key,omitempty"`
	// Tag value
	Value string `json:"value,omitempty"`
}

// ExperimentTags defines the struct for storing the Experiment name and its CSP tags
type ExperimentTags struct {
	// Experiment name
	Name string `json:"name,omitempty"`
	// Experiment tags
	Tags []*Tag `json:"tags,omitempty"`
}

// CreateCostExporterJob creates a Job for the cost exporter.
func CreateCostExporterJob(costExporter *windtunnelv1alpha1.CostExporter, experimentTags []*ExperimentTags, earliestTime time.Time) (*kbatch.Job, error) {
	image := costExporter.Spec.Image
	if image == "" {
		image = defaultImage
	}

	jsonExperimentTags, err := json.Marshal(experimentTags)
	if err != nil {
		return nil, err
	}

	job := &kbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: costExporter.Namespace,
			Name:      costExporter.Name,
		},
		Spec: kbatch.JobSpec{
			BackoffLimit: ptr.To(int32(0)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "cost-exporter",
							Image: image,
							Env: []corev1.EnvVar{
								{
									Name:  "REDIS_HOST",
									Value: redisHost,
								},
								{
									Name:  "REDIS_PORT",
									Value: strconv.FormatInt(int64(redisPort), 10),
								},
								{
									Name:  "CLOUD_SERVICE_PROVIDER",
									Value: costExporter.Spec.CloudServiceProvider,
								},
								{
									Name:  "EXPERIMENT_TAGS",
									Value: string(jsonExperimentTags),
								},
								{
									Name:  "EARLIEST_EXPERIMENT",
									Value: earliestTime.Format("2006-01-02 15:04:05"),
								},
								{
									Name: "CSP_CONFIG",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: costExporter.Spec.Config,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return job, nil
}
