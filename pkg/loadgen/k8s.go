package loadgen

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	k6v1alpha1 "github.com/grafana/k6-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	filenameScript       = config.GetViper().GetString("loadGenerator.filename.script")
	filenameEndpoint     = config.GetViper().GetString("loadGenerator.filename.endpoint")
	filenamePlainText    = config.GetViper().GetString("loadGenerator.filename.plainText")
	filenameDataSet      = config.GetViper().GetString("loadGenerator.filename.dataSet")
	filenameLoadPattern  = config.GetViper().GetString("loadGenerator.filename.loadPattern")
	defaultStorageSize   = config.GetViper().GetString("dataGenerator.defaultStorageSize")
	testRunRWArgs        = config.GetViper().GetString("loadGenerator.testRun.remoteWriteArgs")
	testRunRWEnvVarName  = config.GetViper().GetString("loadGenerator.testRun.remoteWriteEnvVar.name")
	testRunRWEnvVarValue = config.GetViper().GetString("loadGenerator.testRun.remoteWriteEnvVar.value")
)

// CreateConfigMapWithPlainText creates a ConfigMap for EndpointSpec with plain text data.
func CreateConfigMapWithPlainText(experiment *windtunnelv1alpha1.Experiment, endpointSpec *windtunnelv1alpha1.EndpointSpec, pipelineEndpoint *windtunnelv1alpha1.PipelineEndpoint, loadPattern *windtunnelv1alpha1.LoadPattern, protocol windtunnelv1alpha1.EndpointProtocol) (*corev1.ConfigMap, error) {
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
			Name:      utils.GetTestRunConfigMapName(experiment.Name, endpointSpec.EndpointName),
		},
		Data: map[string]string{
			filenameScript:      config.GetViper().GetString(fmt.Sprintf("loadGenerator.script.%s.plainText", protocol)),
			filenameEndpoint:    string(jsonEndpoint),
			filenamePlainText:   endpointSpec.DataSpec.PlainText,
			filenameLoadPattern: string(jsonLoadPattern),
		},
	}, nil
}

// CreateConfigMapWithDataSet creates a ConfigMap for EndpointSpec with DataSet.
func CreateConfigMapWithDataSet(experiment *windtunnelv1alpha1.Experiment, endpointSpec *windtunnelv1alpha1.EndpointSpec, pipelineEndpoint *windtunnelv1alpha1.PipelineEndpoint, dataSet *windtunnelv1alpha1.DataSet, loadPattern *windtunnelv1alpha1.LoadPattern, protocol windtunnelv1alpha1.EndpointProtocol) (*corev1.ConfigMap, error) {
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
			Name:      utils.GetTestRunConfigMapName(experiment.Name, endpointSpec.EndpointName),
		},
		Data: map[string]string{
			filenameScript:      config.GetViper().GetString(fmt.Sprintf("loadGenerator.script.%s.dataSet", protocol)),
			filenameEndpoint:    string(jsonEndpoint),
			filenameDataSet:     string(jsonDataSet),
			filenameLoadPattern: string(jsonLoadPattern),
		},
	}, nil
}

// CreatePVC creates PVC for the EndpointSpec. The PVC will be bound by the TestRun. For EndpointSpec with DataSet only.
func CreatePVC(experiment *windtunnelv1alpha1.Experiment, endpointSpec *windtunnelv1alpha1.EndpointSpec, dataSet *windtunnelv1alpha1.DataSet) *corev1.PersistentVolumeClaim {
	storageSize := endpointSpec.StorageSize
	if storageSize.IsZero() {
		storageSize = dataSet.Spec.StorageSize
	}
	if storageSize.IsZero() {
		storageSize = resource.MustParse(defaultStorageSize)
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunPVCName(experiment.Name, endpointSpec.EndpointName),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: storageSize,
				},
			},
		},
	}

	return pvc
}

// CreateCopierPod creates a Pod to copy the configuration and data for the EndpointSpec.
// For EndpointSpec with DataSet only.
func CreateCopierPod(experiment *windtunnelv1alpha1.Experiment, endpointSpec *windtunnelv1alpha1.EndpointSpec, configMap *corev1.ConfigMap, dataSet *windtunnelv1alpha1.DataSet) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunCopierPodName(experiment.Name, endpointSpec.EndpointName),
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:    "copier",
					Image:   "busybox:1.36.1",
					Command: []string{"/bin/sh", "-c", "cp /configmap/* /testrun && cp -r /dataset/* /testrun"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "configmap-volume",
							MountPath: "/configmap",
						},
						{
							Name:      "dataset-volume",
							MountPath: "/dataset",
						},
						{
							Name:      "testrun-volume",
							MountPath: "/testrun",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "configmap-volume",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: configMap.Name,
							},
						},
					},
				},
				{
					Name: "dataset-volume",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: utils.GetDataSetPVCName(dataSet.Name, dataSet.Generation),
						},
					},
				},
				{
					Name: "testrun-volume",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: utils.GetTestRunPVCName(experiment.Name, endpointSpec.EndpointName),
						},
					},
				},
			},
		},
	}
}

// CreateTestRun creates a TestRun for the EndpointSpec.
func CreateTestRun(experiment *windtunnelv1alpha1.Experiment, endpointSpec *windtunnelv1alpha1.EndpointSpec) *k6v1alpha1.TestRun {
	return &k6v1alpha1.TestRun{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunName(experiment.Name, endpointSpec.EndpointName),
		},
		Spec: k6v1alpha1.TestRunSpec{
			Parallelism: 1,
			Arguments: fmt.Sprintf("%s --tag experiment=%s/%s --tag endpoint=%s",
				testRunRWArgs, experiment.Namespace, experiment.Name, endpointSpec.EndpointName,
			),
			Runner: k6v1alpha1.Pod{
				Env: []corev1.EnvVar{
					{
						Name:  testRunRWEnvVarName,
						Value: testRunRWEnvVarValue,
					},
				},
			},
		},
	}
}
