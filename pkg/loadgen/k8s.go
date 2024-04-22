package loadgen

import (
	"encoding/json"
	"fmt"

	kbatch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	k6v1alpha1 "github.com/grafana/k6-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	filenameScript          = config.GetString("loadGenerator.filename.script")
	filenameEndpoint        = config.GetString("loadGenerator.filename.endpoint")
	filenamePlainText       = config.GetString("loadGenerator.filename.plainText")
	filenameDataSet         = config.GetString("loadGenerator.filename.dataSet")
	filenameLoadPattern     = config.GetString("loadGenerator.filename.loadPattern")
	copierImage             = config.GetString("loadGenerator.copier.image")
	defaultRunnerImage      = config.GetString("loadGenerator.testRun.defaultRunnerImage")
	defaultStarterImage     = config.GetString("loadGenerator.testRun.defaultStarterImage")
	defaultInitializerImage = config.GetString("loadGenerator.testRun.defaultInitializerImage")
	testRunRWArgs           = config.GetString("loadGenerator.testRun.remoteWriteArgs")
	testRunRWEnvVarName     = config.GetString("loadGenerator.testRun.remoteWriteEnvVar.name")
	testRunRWEnvVarValue    = config.GetString("loadGenerator.testRun.remoteWriteEnvVar.value")
	defaultStorageSize      = config.GetString("dataGenerator.defaultStorageSize")
)

// CreateConfigMapWithPlainText creates a ConfigMap for EndpointSpec with plain text data.
func CreateConfigMapWithPlainText(experiment *windtunnelv1alpha1.Experiment, endpointIdx int, pipelineEndpoint *windtunnelv1alpha1.PipelineEndpoint, text string, loadPattern *windtunnelv1alpha1.LoadPattern, protocol windtunnelv1alpha1.EndpointProtocol) (*corev1.ConfigMap, error) {
	jsonEndpoint, err := json.Marshal(pipelineEndpoint)
	if err != nil {
		return nil, err
	}

	jsonLoadPattern, err := json.Marshal(loadPattern)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
		},
		Data: map[string]string{
			filenameScript:      config.GetString(fmt.Sprintf("loadGenerator.script.%s.plainText", protocol)),
			filenameEndpoint:    string(jsonEndpoint),
			filenamePlainText:   text,
			filenameLoadPattern: string(jsonLoadPattern),
		},
	}, nil
}

// CreateConfigMapWithDataSet creates a ConfigMap for EndpointSpec with DataSet.
func CreateConfigMapWithDataSet(experiment *windtunnelv1alpha1.Experiment, endpointIdx int, pipelineEndpoint *windtunnelv1alpha1.PipelineEndpoint, dataSet *windtunnelv1alpha1.DataSet, loadPattern *windtunnelv1alpha1.LoadPattern, protocol windtunnelv1alpha1.EndpointProtocol) (*corev1.ConfigMap, error) {
	jsonEndpoint, err := json.Marshal(pipelineEndpoint)
	if err != nil {
		return nil, err
	}

	jsonDataSet, err := json.Marshal(dataSet)
	if err != nil {
		return nil, err
	}

	jsonLoadPattern, err := json.Marshal(loadPattern)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
		},
		Data: map[string]string{
			filenameScript:      config.GetString(fmt.Sprintf("loadGenerator.script.%s.dataSet", protocol)),
			filenameEndpoint:    string(jsonEndpoint),
			filenameDataSet:     string(jsonDataSet),
			filenameLoadPattern: string(jsonLoadPattern),
		},
	}, nil
}

// CreatePVC creates PVC for the EndpointSpec. The PVC will be bound by the TestRun. For EndpointSpec with DataSet only.
func CreatePVC(experiment *windtunnelv1alpha1.Experiment, endpointIdx int, endpointSpec *windtunnelv1alpha1.EndpointSpec, dataSet *windtunnelv1alpha1.DataSet) *corev1.PersistentVolumeClaim {
	var storageSize resource.Quantity
	if endpointSpec.StorageSize != nil && !endpointSpec.StorageSize.IsZero() {
		storageSize = *endpointSpec.StorageSize
	} else if dataSet.Spec.StorageSize != nil && !dataSet.Spec.StorageSize.IsZero() {
		storageSize = *dataSet.Spec.StorageSize
	} else {
		storageSize = resource.MustParse(defaultStorageSize)
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: storageSize,
				},
			},
		},
	}

	return pvc
}

// CreateCopierJob creates a Job to copy the configuration and data for the EndpointSpec.
// For EndpointSpec that uses a DataSet only.
func CreateCopierJob(experiment *windtunnelv1alpha1.Experiment, endpointIdx int, endpointSpec *windtunnelv1alpha1.EndpointSpec, configMap *corev1.ConfigMap, dataSet *windtunnelv1alpha1.DataSet) *kbatch.Job {
	return &kbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunCopierJobName(experiment.Name, endpointIdx),
		},
		Spec: kbatch.JobSpec{
			BackoffLimit: ptr.To(int32(0)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "copier",
							Image: copierImage,
							// Note that when a ConfigMap is mounted as a volume, the files are symlinked to "/data",
							// using the "-R" flag will copy files recursively but also preserve symlinks by default,
							// so we need the "-L" flag to follow symlinks.
							Command: []string{"/bin/sh", "-c", "cp -RL /configmap/* /testrun && cp -RL /dataset/* /testrun"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "configmap",
									MountPath: "/configmap",
								},
								{
									Name:      "dataset",
									MountPath: "/dataset",
								},
								{
									Name:      "testrun",
									MountPath: "/testrun",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "configmap",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMap.Name,
									},
								},
							},
						},
						{
							Name: "dataset",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: utils.GetDataGeneratorName(dataSet.Name, dataSet.Generation),
								},
							},
						},
						{
							Name: "testrun",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: utils.GetTestRunName(experiment.Name, endpointIdx),
								},
							},
						},
					},
				},
			},
		},
	}
}

// CreateTestRun creates a TestRun for the EndpointSpec.
func CreateTestRun(experiment *windtunnelv1alpha1.Experiment, endpointIdx int, endpointSpec *windtunnelv1alpha1.EndpointSpec) *k6v1alpha1.TestRun {
	runnerImage := experiment.Spec.K6RunnerImage
	if runnerImage == "" {
		runnerImage = defaultRunnerImage
	}

	starterImage := experiment.Spec.K6StarterImage
	if starterImage == "" {
		starterImage = defaultStarterImage
	}

	initializerImage := experiment.Spec.K6InitializerImage
	if initializerImage == "" {
		initializerImage = defaultInitializerImage
	}

	return &k6v1alpha1.TestRun{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
		},
		Spec: k6v1alpha1.TestRunSpec{
			Parallelism: 1,
			Arguments: fmt.Sprintf("%s --tag experiment=%s/%s --tag endpoint=%s",
				testRunRWArgs, experiment.Namespace, experiment.Name, endpointSpec.EndpointName,
			),
			Runner: k6v1alpha1.Pod{
				Image: experiment.Spec.K6RunnerImage,
				Env: []corev1.EnvVar{
					{
						Name:  testRunRWEnvVarName,
						Value: testRunRWEnvVarValue,
					},
				},
			},
			Starter: k6v1alpha1.Pod{
				Image: experiment.Spec.K6StarterImage,
			},
			Initializer: &k6v1alpha1.Pod{
				Image: experiment.Spec.K6InitializerImage,
			},
		},
	}
}
