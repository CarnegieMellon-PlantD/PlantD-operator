package loadgen

import (
	"encoding/json"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	k6v1alpha1 "github.com/grafana/k6-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTestRunManifest(name string, namespace string) *k6v1alpha1.K6 {
	return &k6v1alpha1.K6{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: k6v1alpha1.K6Spec{
			Parallelism: 1,
			Arguments:   config.GetString("k6.arguments"),
			Runner: k6v1alpha1.Pod{
				Env: []corev1.EnvVar{
					{
						Name:  config.GetString("k6.remoteWriteURL.name"),
						Value: config.GetString("k6.remoteWriteURL.value"),
					},
				},
			},
		},
	}
}

func CreateConfigMap(name string, protocol string, exp *windtunnelv1alpha1.Experiment, endpoint *windtunnelv1alpha1.Endpoint, loadPattern *windtunnelv1alpha1.LoadPattern) (*corev1.ConfigMap, error) {
	pipelineSpec, err := json.Marshal(endpoint)
	if err != nil {
		return nil, err
	}
	loadPatternSpec, err := json.Marshal(&loadPattern.Spec)
	if err != nil {
		return nil, err
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: exp.Namespace,
		},
		Data: map[string]string{
			config.GetString("k6.config.script"):      config.GetString("k6.script." + protocol),
			config.GetString("k6.config.pipeline"):    string(pipelineSpec),
			config.GetString("k6.config.loadPattern"): string(loadPatternSpec),
		},
	}, nil
}

func CreateConfigMapWithDataSet(name, protocol string, exp *windtunnelv1alpha1.Experiment, endpoint *windtunnelv1alpha1.Endpoint, loadPattern *windtunnelv1alpha1.LoadPattern, dataset *windtunnelv1alpha1.DataSet) (*corev1.ConfigMap, error) {
	pipelineSpec, err := json.Marshal(endpoint)
	if err != nil {
		return nil, err
	}
	loadPatternSpec, err := json.Marshal(&loadPattern.Spec)
	if err != nil {
		return nil, err
	}
	jsonDataset, err := json.Marshal(dataset)
	if err != nil {
		return nil, err
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: exp.Namespace,
		},
		Data: map[string]string{
			config.GetString("k6.config.script"):      config.GetString("k6.script." + protocol),
			config.GetString("k6.config.pipeline"):    string(pipelineSpec),
			config.GetString("k6.config.loadPattern"): string(loadPatternSpec),
			config.GetString("k6.config.dataset"):     string(jsonDataset),
		},
	}, nil
}

func CreateCopyPod(dataset *windtunnelv1alpha1.DataSet, configMap *corev1.ConfigMap) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap.Name,
			Namespace: configMap.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "cp",
					Image:   "busybox",
					Command: []string{"/bin/sh", "-c", "cp /config/* /pvc/"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config-volume",
							MountPath: "/config",
						},
						{
							Name:      "pvc-volume",
							MountPath: "/pvc",
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
			Volumes: []corev1.Volume{
				{
					Name: "config-volume",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: configMap.Name,
							},
						},
					},
				},
				{
					Name: "pvc-volume",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: utils.GetPVCName(dataset.Name, dataset.Generation),
						},
					},
				},
			},
		},
	}
}
