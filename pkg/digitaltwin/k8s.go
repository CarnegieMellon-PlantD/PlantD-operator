package digitaltwin

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	kbatch "k8s.io/api/batch/v1"
	"k8s.io/utils/ptr"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultImage   = config.GetString("digitalTwin.defaultImage")
	redisHost      = utils.GetServiceARecord(config.GetString("core.redis.name"), config.GetString("core.namespace"))
	redisPort      = config.GetInt32("core.redis.servicePort")
	redisPwd       = ""
	prometheusHost = fmt.Sprintf("http://%s:%d",
		utils.GetServiceARecord(config.GetString("core.thanos.querier.name"), config.GetString("core.namespace")),
		config.GetInt32("core.prometheus.servicePort"),
	)
	prometheusPwd      = ""
	prometheusEndpoint = fmt.Sprintf("%s/api/v1/query", prometheusHost)
	openCostEndpoint   = fmt.Sprintf("http://%s:%d/allocation",
		utils.GetServiceARecord(config.GetString("core.openCost.name"), config.GetString("core.namespace")),
		config.GetInt32("core.openCost.servicePort"),
	)
)

// CreateSimulationJob creates a Job for the Simulation.
func CreateSimulationJob(simulation *windtunnelv1alpha1.Simulation, digitalTwin *windtunnelv1alpha1.DigitalTwin,
	trafficModel *windtunnelv1alpha1.TrafficModel, netCost *windtunnelv1alpha1.NetCost, scenario *windtunnelv1alpha1.Scenario,
	pipeline *windtunnelv1alpha1.Pipeline,
	dataSets *windtunnelv1alpha1.DataSetList, loadPatterns *windtunnelv1alpha1.LoadPatternList, experiments *windtunnelv1alpha1.ExperimentList,
) (*kbatch.Job, error) {
	image := simulation.Spec.Image
	if image == "" {
		image = defaultImage
	}

	env := []corev1.EnvVar{
		{Name: "REDIS_HOST", Value: redisHost},
		{Name: "REDIS_PORT", Value: strconv.FormatInt(int64(redisPort), 10)},
		{Name: "REDIS_PASSWORD", Value: redisPwd},
		{Name: "PROMETHEUS_HOST", Value: prometheusHost},
		{Name: "PROMETHEUS_PASSWORD", Value: prometheusPwd},
		{Name: "PROMETHEUS_ENDPOINT", Value: prometheusEndpoint},

		{Name: "OPENCOST_ENDPOINT", Value: openCostEndpoint},

		{Name: "SIM_NAME", Value: fmt.Sprintf("%s.%s", simulation.Namespace, simulation.Name)},

		{Name: "TRAFFIC_MODEL_NAME", Value: trafficModel.Name},
		{Name: "TRAFFIC_MODEL", Value: trafficModel.Spec.Config},
	}

	if digitalTwin != nil {
		env = append(env, []corev1.EnvVar{
			{Name: "TWIN_NAME", Value: fmt.Sprintf("%s.%s", digitalTwin.Namespace, digitalTwin.Name)},
			{Name: "MODEL_TYPE", Value: digitalTwin.Spec.ModelType},
			{Name: "DIGITAL_TWIN_TYPE", Value: digitalTwin.Spec.DigitalTwinType},
		}...)
	}

	if pipeline != nil {
		pipelineLabelKeys := make([]string, 0, len(pipeline.Spec.Tags))
		pipelineLabelValues := make([]string, 0, len(pipeline.Spec.Tags))
		for k, v := range pipeline.Spec.Tags {
			pipelineLabelKeys = append(pipelineLabelKeys, k)
			pipelineLabelValues = append(pipelineLabelValues, v)
		}
		env = append(env, []corev1.EnvVar{
			{Name: "PIPELINE_LABEL_KEYS", Value: strings.Join(pipelineLabelKeys, ",")},
			{Name: "PIPELINE_LABEL_VALUES", Value: strings.Join(pipelineLabelValues, ",")},
		}...)
	}

	if experiments != nil {
		experimentNames := make([]string, 0, len(experiments.Items))
		for _, experiment := range experiments.Items {
			experimentNames = append(experimentNames, fmt.Sprintf("%s.%s", experiment.Namespace, experiment.Name))
		}

		jsonExperiments, err := json.Marshal(experiments)
		if err != nil {
			return nil, err
		}

		env = append(env, []corev1.EnvVar{
			{Name: "EXPERIMENT_NAMES", Value: strings.Join(experimentNames, ",")},
			{Name: "EXPERIMENT_JSON", Value: string(jsonExperiments)},
		}...)
	} else {
		env = append(env, []corev1.EnvVar{
			{Name: "EXPERIMENT_NAMES", Value: ""},
		}...)
	}

	if dataSets != nil {
		jsonDataSets, err := json.Marshal(dataSets)
		if err != nil {
			return nil, err
		}

		env = append(env, []corev1.EnvVar{
			{Name: "DATASET_JSON", Value: string(jsonDataSets)},
		}...)
	}

	if loadPatterns != nil {
		jsonLoadPatterns, err := json.Marshal(loadPatterns)
		if err != nil {
			return nil, err
		}

		env = append(env, []corev1.EnvVar{
			{Name: "LOAD_PATTERN_JSON", Value: string(jsonLoadPatterns)},
		}...)
	}

	if netCost != nil {
		jsonNetCost, err := json.Marshal(netCost)
		if err != nil {
			return nil, err
		}
		env = append(env, []corev1.EnvVar{
			{Name: "NETCOSTS", Value: string(jsonNetCost)},
		}...)
	} else {
		env = append(env, []corev1.EnvVar{
			{Name: "NETCOSTS", Value: ""},
		}...)
	}

	if scenario != nil {
		jsonScenario, err := json.Marshal(scenario)
		if err != nil {
			return nil, err
		}
		env = append(env, []corev1.EnvVar{
			{Name: "SCENARIO_NAME", Value: fmt.Sprintf("%s.%s", scenario.Namespace, scenario.Name)},
			{Name: "SCENARIO", Value: string(jsonScenario)},
		}...)
	}

	job := &kbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: simulation.Namespace,
			Name:      utils.GetSimulationJobName(simulation.Name),
		},
		Spec: kbatch.JobSpec{
			BackoffLimit: ptr.To(int32(0)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "simulation",
							Image:   image,
							Command: []string{"python3"},
							Args:    []string{"/digitaltwin/main.py", "sim_all"},
							Env:     env,
						},
					},
				},
			},
		},
	}

	return job, nil
}

// CreateEndDetectorJob creates a Job for the end detector.
func CreateEndDetectorJob(experiment *windtunnelv1alpha1.Experiment, debouncePeriod, queryWindow, podDetachAdjustment int64) (*kbatch.Job, error) {
	image := experiment.Spec.EndDetectionImage
	if image == "" {
		image = defaultImage
	}

	experimentList := &windtunnelv1alpha1.ExperimentList{}
	experimentList.Items = append(experimentList.Items, *experiment)
	jsonExperiments, err := json.Marshal(experimentList)
	if err != nil {
		return nil, err
	}

	env := []corev1.EnvVar{
		{Name: "PROMETHEUS_HOST", Value: prometheusHost},
		{Name: "PROMETHEUS_PASSWORD", Value: prometheusPwd},
		{Name: "PROMETHEUS_ENDPOINT", Value: prometheusEndpoint},

		{Name: "EXPERIMENT_NAMES", Value: fmt.Sprintf("%s.%s", experiment.Namespace, experiment.Name)},
		{Name: "EXPERIMENT_JSON", Value: string(jsonExperiments)},

		{Name: "DEBOUNCE_PERIOD", Value: strconv.FormatInt(debouncePeriod, 10)},
		{Name: "QUERY_WINDOW", Value: strconv.FormatInt(queryWindow, 10)},
		{Name: "POD_DETACH_ADJUSTMENT", Value: strconv.FormatInt(podDetachAdjustment, 10)},
	}

	job := &kbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: experiment.Namespace,
			Name:      utils.GetEndDetectorJobName(experiment.Name),
		},
		Spec: kbatch.JobSpec{
			BackoffLimit: ptr.To(int32(0)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "end-detector",
							Image:   image,
							Command: []string{"python3"},
							Args:    []string{"/digitaltwin/end_detect.py"},
							Env:     env,
						},
					},
				},
			},
		},
	}

	return job, nil
}
