package datagen

import (
	"encoding/json"
	"strconv"

	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

var (
	image        = config.GetString("dataGenerator.image")
	backoffLimit = config.GetInt32("dataGenerator.backoffLimit")
	cpuRequest   = config.GetString("dataGenerator.requests.cpu")
	cpuLimit     = config.GetString("dataGenerator.limits.cpu")
	storageSize  = config.GetString("pvc.requests.storage")
)

// CreateJobByDataSet creates a Job based on the DataSet configuration.
func CreateJobByDataSet(jobName string, pvcName string, dataSet *windtunnelv1alpha1.DataSet, schemaMap map[string]*windtunnelv1alpha1.Schema) (*kbatch.Job, error) {
	// Calculate static step size and parallel jobs
	staticStepSize := dataSet.Spec.NumberOfFiles / dataSet.Spec.ParallelJobs
	parallelJobs := dataSet.Spec.ParallelJobs

	completionMode := kbatch.IndexedCompletion

	// Get data generator name and volume name
	dgName := dataSet.Name
	volumeName := utils.GetDataSetVolumeName(dgName)

	// Marshal dataset and schema map to JSON
	datasetBytes, err := json.Marshal(dataSet)
	if err != nil {
		return nil, err
	}

	schemaMapBytes, err := json.Marshal(schemaMap)
	if err != nil {
		return nil, err
	}

	// Create the Kubernetes Job object
	job := &kbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        jobName,
			Namespace:   dataSet.Namespace,
		},
		Spec: kbatch.JobSpec{
			CompletionMode: &completionMode,
			Completions:    &parallelJobs,
			Parallelism:    &parallelJobs,
			BackoffLimit:   &backoffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  jobName,
							Image: image,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse(cpuRequest),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse(cpuLimit),
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "JOB_STEP_SIZE",
									Value: strconv.FormatInt(int64(staticStepSize), 10),
								},
								{
									Name:  "MAX_REPEAT",
									Value: strconv.FormatInt(int64(dataSet.Spec.NumberOfFiles), 10),
								},
								{
									Name:  "DG_NAMESPACE",
									Value: dataSet.Namespace,
								},
								{
									Name:  "DG_NAME",
									Value: dgName,
								},
								{
									Name:  "DATASET",
									Value: string(datasetBytes),
								},
								{
									Name:  "SCHEMA_MAP",
									Value: string(schemaMapBytes),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: path,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: volumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvcName,
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

// CreatePVC creates a PersistentVolumeClaim for the data generator job.
func CreatePVC(name types.NamespacedName) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(storageSize),
				},
			},
		},
	}

	return pvc
}
