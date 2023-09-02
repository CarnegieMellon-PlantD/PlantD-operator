package datagen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	"go.uber.org/zap"
	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/brianvoe/gofakeit/v6"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

// DataGeneratorJob is an interface for generating data.
type DataGeneratorJob interface {
	GenerateData() error
}

// JobConfig holds the configuration for a data generator job.
type JobConfig struct {
	RepeatStart int
	RepeatEnd   int
}

// BuildBasedDataGeneratorJob is a data generator job based on the Build strategy.
type BuildBasedDataGeneratorJob struct {
	RepeatStart int
	RepeatEnd   int
	Namespace   string
	DGName      string
	Dataset     *windtunnelv1alpha1.DataSet
	SchemaMap   map[string]*windtunnelv1alpha1.Schema
}

// NewBuildBasedDataGeneratorJob creates a new BuildBasedDataGeneratorJob instance.
func NewBuildBasedDataGeneratorJob(start int, end int, dgNamespace string, dgName string, dataset *windtunnelv1alpha1.DataSet, schemaMap map[string]*windtunnelv1alpha1.Schema) DataGeneratorJob {
	return &BuildBasedDataGeneratorJob{
		RepeatStart: start,
		RepeatEnd:   end,
		Namespace:   dgNamespace,
		DGName:      dgName,
		Dataset:     dataset,
		SchemaMap:   schemaMap,
	}
}

// MakeOutputDir creates the output directory for a schema in the dataset.
func MakeOutputDir(dataGeneratorConfig *windtunnelv1alpha1.DataSet, seqNum int) error {
	schPath := filepath.Join(path, dataGeneratorConfig.Spec.Schemas[seqNum].Name)
	err := os.RemoveAll(schPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(schPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// ApplyOperations applies the operations defined in the output builder to generate the final output.
func ApplyOperations(outputBuilder *OutputBuilder, seqNum int) error {
	var err error
	for _, op := range outputBuilder.Operations {
		err = op(outputBuilder, seqNum)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateData generates the data using the Build strategy.
func (dg *BuildBasedDataGeneratorJob) GenerateData() error {
	var err error

	// Create schema builders and populate the schema builder cache
	for _, schemaSelector := range dg.Dataset.Spec.Schemas {
		schemaName := schemaSelector.Name
		schemaObj := dg.SchemaMap[schemaName]
		schBldr, err := NewSchemaBuilder(schemaObj)
		if err != nil {
			return err
		}
		PutSchemaBuilder(schemaName, schBldr)
	}

	// Create the output builder
	outputBuilder, err := NewOutputBuilder(dg.Dataset)
	if err != nil {
		return err
	}

	// Set the seed for random number generation
	seed := gofakeit.New(0).Rand

	// Create output directories for each schema
	scheNum := len(outputBuilder.SchBuilders)
	for i := 0; i < scheNum; i++ {
		err := MakeOutputDir(dg.Dataset, i)
		if err != nil {
			return err
		}
	}

	startTime := time.Now()

	// Generate data for each repeat
	for i := dg.RepeatStart; i < dg.RepeatEnd; i++ {
		// Build data for each schema
		for _, schBldr := range outputBuilder.SchBuilders {
			err := schBldr.Build(seed)
			if err != nil {
				return err
			}
		}
		// Apply operations to generate the final output
		err = ApplyOperations(outputBuilder, i)
		if err != nil {
			return err
		}
	}
	endTime := time.Now()

	dgDuration := endTime.Sub(startTime)
	zap.S().Info(dgDuration)
	return nil
}

// CreateJobByDataSet creates a Kubernetes Job based on the DataSet configuration.
func CreateJobByDataSet(jobName string, pvcName string, dataGenerator *windtunnelv1alpha1.DataSet, schemaMap map[string]*windtunnelv1alpha1.Schema) (*kbatch.Job, error) {
	// Calculate static step size and parallel jobs
	staticStepSize := dataGenerator.Spec.NumberOfFiles / dataGenerator.Spec.ParallelJobs
	parallelJobs := dataGenerator.Spec.ParallelJobs

	// Get backoff limit and completion mode from configuration
	backoffLimit := config.GetInt32("dataGenerator.backoffLimit")
	completionMode := kbatch.IndexedCompletion

	// Get data generator name and volume name
	dgName := dataGenerator.Name
	volumeName := utils.GetVolumeName(dgName)

	// Marshal dataset and schema map to JSON
	datasetBytes, err := json.Marshal(*dataGenerator)
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
			Namespace:   dataGenerator.Namespace,
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
							Image: config.GetString("dataGenerator.image"),
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse(config.GetString("dataGenerator.requests.cpu")),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse(config.GetString("dataGenerator.limits.cpu")),
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "JOB_STEP_SIZE",
									Value: strconv.FormatInt(int64(staticStepSize), 10),
								},
								{
									Name:  "MAX_REPEAT",
									Value: strconv.FormatInt(int64(dataGenerator.Spec.NumberOfFiles), 10),
								},
								{
									Name:  "DG_NAMESPACE",
									Value: dataGenerator.Namespace,
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
					corev1.ResourceStorage: resource.MustParse(config.GetString("pvc.requests.storage")),
				},
			},
		},
	}

	return pvc
}
