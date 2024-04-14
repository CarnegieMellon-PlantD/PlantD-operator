package controller

import (
	"context"
	"fmt"
	"time"

	k6v1alpha1 "github.com/grafana/k6-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/loadgen"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

const (
	experimentPollingInterval = 5 * time.Second
)

var (
	filenameScript                   = config.GetViper().GetString("loadGenerator.filename.script")
	metricsServiceLabelKeyExperiment = config.GetViper().GetString("monitor.service.labelKeys.experiment")
)

// ExperimentReconciler reconciles a Experiment object
type ExperimentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments/finalizers,verbs=update
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=datasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=loadpatterns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k6.io,resources=testruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *ExperimentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the requested Experiment
	experiment := &windtunnelv1alpha1.Experiment{}
	if err := r.Get(ctx, req.NamespacedName, experiment); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unable to fetch Pipeline")
		return ctrl.Result{}, err
	}

	stop, result, err := r.getRelatedResources(ctx, experiment)
	if stop {
		if err := r.Status().Update(ctx, experiment); err != nil {
			logger.Error(err, "Cannot update the status")
			return ctrl.Result{}, err
		}
		return result, err
	}

	if experiment.Status.JobStatus == "" || experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentScheduled {
		stop, result, err := r.reconcileCreatedOrScheduled(ctx, experiment)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentWaitingDataSet {
		stop, result, err := r.reconcileWaitingDataSet(ctx, experiment)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentWaitingPipeline {
		stop, result, err := r.reconcileWaitingPipeline(ctx, experiment)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentInitializing {
		stop, result, err := r.reconcileInitializing(ctx, experiment)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentRunning {
		stop, result, err := r.reconcileRunning(ctx, experiment)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentCompleted {
		stop, result, err := r.reconcileCompleted(ctx, experiment)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	return ctrl.Result{}, nil
}

// getRelatedResources gets the related resources used by the Experiment and update the status fields.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) getRelatedResources(ctx context.Context, experiment *windtunnelv1alpha1.Experiment) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Pipeline used by the Experiment
	pipeline := &windtunnelv1alpha1.Pipeline{}
	pipelineName := types.NamespacedName{
		Namespace: experiment.Spec.PipelineRef.Namespace,
		Name:      experiment.Spec.PipelineRef.Name,
	}
	if err := r.Get(ctx, pipelineName, pipeline); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot get Pipeline: %s", pipelineName))
		experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
		experiment.Status.Error = fmt.Sprintf("Cannot find Pipeline: %s", err)
		return true, ctrl.Result{}, nil
	}
	experiment.Status.Pipeline = pipeline

	experiment.Status.EndpointMap = make(map[string]*windtunnelv1alpha1.PipelineEndpoint, len(experiment.Spec.EndpointSpecs))
	experiment.Status.ProtocolMap = make(map[string]windtunnelv1alpha1.EndpointProtocol, len(experiment.Spec.EndpointSpecs))
	experiment.Status.DataOptionMap = make(map[string]windtunnelv1alpha1.EndpointDataOption, len(experiment.Spec.EndpointSpecs))
	experiment.Status.DataSetMap = make(map[string]*windtunnelv1alpha1.DataSet, len(experiment.Spec.EndpointSpecs))
	experiment.Status.LoadPatternMap = make(map[string]*windtunnelv1alpha1.LoadPattern, len(experiment.Spec.EndpointSpecs))
	experiment.Status.Durations = make(map[string]metav1.Duration, len(experiment.Spec.EndpointSpecs))
	for _, endpointSpec := range experiment.Spec.EndpointSpecs {
		// Find the PipelineEndpoint referenced by the EndpointSpec
		pipelineEndpoint := getPipelineEndpoint(pipeline, endpointSpec.EndpointName)
		if pipelineEndpoint == nil {
			logger.Error(nil, fmt.Sprintf("Cannot find endpoint: %s", endpointSpec.EndpointName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Cannot find endpoint: %s", endpointSpec.EndpointName)
			return true, ctrl.Result{}, nil
		}
		experiment.Status.EndpointMap[endpointSpec.EndpointName] = pipelineEndpoint

		// Determine the protocol used by the PipelineEndpoint
		protocol := getPipelineEndpointProtocol(pipelineEndpoint)
		if protocol == "" {
			logger.Error(nil, fmt.Sprintf("Unspecified protocol in endpoint: %s", endpointSpec.EndpointName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Unspecified protocol in endpoint: %s", endpointSpec.EndpointName)
			return true, ctrl.Result{}, nil
		}
		experiment.Status.ProtocolMap[endpointSpec.EndpointName] = protocol

		// Determine the data option used by the EndpointSpec
		experiment.Status.DataOptionMap[endpointSpec.EndpointName] = getEndpointSpecDataOption(&endpointSpec)

		// Fetch the DataSet used by the EndpointSpec
		if getEndpointSpecDataOption(&endpointSpec) == windtunnelv1alpha1.EndpointDataOptionDataSet {
			dataSet := &windtunnelv1alpha1.DataSet{}
			dataSetName := types.NamespacedName{
				Namespace: endpointSpec.DataSpec.DataSetRef.Namespace,
				Name:      endpointSpec.DataSpec.DataSetRef.Name,
			}
			if err := r.Get(ctx, dataSetName, dataSet); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot get DataSet: %s", dataSetName))
				experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
				experiment.Status.Error = fmt.Sprintf("Cannot find DataSet: %s", err)
				return true, ctrl.Result{}, nil
			}
			experiment.Status.DataSetMap[endpointSpec.EndpointName] = dataSet
		}

		// Fetch the LoadPattern used by the EndpointSpec
		loadPattern := &windtunnelv1alpha1.LoadPattern{}
		loadPatternName := types.NamespacedName{
			Namespace: endpointSpec.LoadPatternRef.Namespace,
			Name:      endpointSpec.LoadPatternRef.Name,
		}
		if err := r.Get(ctx, loadPatternName, loadPattern); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get LoadPattern: %s", loadPatternName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Cannot find LoadPattern: %s", err)
			return true, ctrl.Result{}, nil
		}
		experiment.Status.LoadPatternMap[endpointSpec.EndpointName] = loadPattern

		// Calculate the duration of LoadPattern
		duration, err := getLoadPatternDuration(loadPattern)
		if err != nil {
			logger.Error(err, fmt.Sprintf("Cannot calculate the duration of LoadPattern: %s", loadPatternName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Cannot calculate the duration for endpoint %s: %s",
				endpointSpec.EndpointName, err,
			)
			return true, ctrl.Result{}, nil
		}
		experiment.Status.Durations[endpointSpec.EndpointName] = duration
	}

	// Proceed
	return false, ctrl.Result{}, nil
}

// reconcileCreatedOrScheduled reconciles the Experiment when it is created or scheduled.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileCreatedOrScheduled(ctx context.Context, experiment *windtunnelv1alpha1.Experiment) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Check if the scheduled time has been reached
	curTime := time.Now()
	if curTime.Before(experiment.Spec.ScheduledTime.Time) {
		// Time has not been reached, calculate the waiting time
		waitTime := experiment.Spec.ScheduledTime.Time.Sub(curTime)
		logger.Info(fmt.Sprintf("Scheduled time has not been reached yet, waiting for %s", waitTime))
		experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentScheduled
		return true, ctrl.Result{RequeueAfter: waitTime}, nil
	}

	// Time has been reached, proceed to the next step
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentWaitingDataSet
	return false, ctrl.Result{}, nil
}

// reconcileWaitingDataSet reconciles the Experiment when it is waiting for the DataSet.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileWaitingDataSet(ctx context.Context, experiment *windtunnelv1alpha1.Experiment) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	for _, dataSet := range experiment.Status.DataSetMap {
		// Check if the DataSet is in "Ready" status
		if dataSet.Status.JobStatus != windtunnelv1alpha1.DataSetJobSuccess {
			logger.Info(fmt.Sprintf("DataSet %s is in status %s, waiting for it",
				utils.GetNamespacedName(dataSet), dataSet.Status.JobStatus,
			))
			return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
		}
	}

	// Proceed to the next step
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentWaitingPipeline
	return false, ctrl.Result{}, nil
}

// reconcileWaitingPipeline reconciles the Experiment when it is waiting for the Pipeline.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileWaitingPipeline(ctx context.Context, experiment *windtunnelv1alpha1.Experiment) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Check if the Pipeline is in "Ready" status
	if experiment.Status.Pipeline.Status.Availability != windtunnelv1alpha1.PipelineReady {
		logger.Info(fmt.Sprintf("Pipeline %s is in status %s, waiting for it",
			utils.GetNamespacedName(experiment.Status.Pipeline), experiment.Status.Pipeline.Status.Availability,
		))
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	}

	// Try to lock the Pipeline by setting its status to "In-Use"
	experiment.Status.Pipeline.Status.Availability = windtunnelv1alpha1.PipelineInUse
	if err := r.Status().Update(ctx, experiment.Status.Pipeline); apierrors.IsConflict(err) {
		// Re-queue the request if conflict occurs
		logger.Info("Conflict occurs when trying to lock the Pipeline, retrying")
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	} else if err != nil {
		logger.Error(err, "Cannot update the status of Pipeline")
		return true, ctrl.Result{}, err
	}

	// Set the Experiment label for the metrics Service
	metricsService := &corev1.Service{}
	var metricsServiceName types.NamespacedName
	if experiment.Status.Pipeline.Spec.InCluster {
		metricsServiceName = types.NamespacedName{
			Namespace: experiment.Status.Pipeline.Namespace,
			Name:      experiment.Status.Pipeline.Spec.MetricsEndpoint.ServiceRef.Name,
		}
	} else {
		metricsServiceName = types.NamespacedName{
			Namespace: experiment.Status.Pipeline.Namespace,
			Name:      utils.GetPipelineMetricsServiceName(experiment.Status.Pipeline.Name),
		}
	}
	if err := r.Get(ctx, metricsServiceName, metricsService); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot get metrics Service: %s", metricsServiceName))
		experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
		experiment.Status.Error = fmt.Sprintf("Cannot find metrics Service: %s", err)
		return true, ctrl.Result{}, nil
	}
	if metricsService.Labels == nil {
		metricsService.Labels = make(map[string]string, 1)
	}
	metricsService.Labels[metricsServiceLabelKeyExperiment] = experiment.Name
	if err := r.Update(ctx, metricsService); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot add Experiment label to metrics Service: %s", metricsServiceName))
		return true, ctrl.Result{}, err
	}
	logger.Info(fmt.Sprintf("Added Experiment label to metrics Service: %s", metricsServiceName))

	// Pipeline is in "Ready" status and locked, proceed to the next step
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentInitializing
	return false, ctrl.Result{}, nil
}

// reconcileInitializing reconciles the Experiment when it is initializing.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileInitializing(ctx context.Context, experiment *windtunnelv1alpha1.Experiment) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Perform health check
	for _, healthCheckURL := range experiment.Status.Pipeline.Spec.HealthCheckURLs {
		if err := utils.CheckHealth(healthCheckURL); err != nil {
			logger.Error(err, "Pipeline health check failed")
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Pipeline health check failed: %s", err)
			return true, ctrl.Result{}, nil
		}
	}

	// Create ConfigMap, PVC and copier Pod for each endpoint
	doneCounter := 0
	for _, endpointSpec := range experiment.Spec.EndpointSpecs {
		switch experiment.Status.DataOptionMap[endpointSpec.EndpointName] {
		case windtunnelv1alpha1.EndpointDataOptionPlainText:
			// ConfigMap
			configMap, err := loadgen.CreateConfigMapWithPlainText(
				experiment, &endpointSpec,
				experiment.Status.EndpointMap[endpointSpec.EndpointName],
				experiment.Status.LoadPatternMap[endpointSpec.EndpointName],
				experiment.Status.ProtocolMap[endpointSpec.EndpointName],
			)
			if err != nil {
				logger.Error(err, fmt.Sprintf("Cannot prepare ConfigMap to create for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			}
			if err := ctrl.SetControllerReference(experiment, configMap, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for ConfigMap for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			}
			if err := r.Create(ctx, configMap); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create ConfigMap for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created ConfigMap for endpoint %s", endpointSpec.EndpointName))
			}

			doneCounter++

		case windtunnelv1alpha1.EndpointDataOptionDataSet:
			// ConfigMap
			configMap, err := loadgen.CreateConfigMapWithDataSet(
				experiment, &endpointSpec,
				experiment.Status.EndpointMap[endpointSpec.EndpointName],
				experiment.Status.DataSetMap[endpointSpec.EndpointName],
				experiment.Status.LoadPatternMap[endpointSpec.EndpointName],
				experiment.Status.ProtocolMap[endpointSpec.EndpointName],
			)
			if err != nil {
				logger.Error(err, fmt.Sprintf("Cannot prepare ConfigMap to create for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			}
			if err := ctrl.SetControllerReference(experiment, configMap, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for ConfigMap for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			}
			if err := r.Create(ctx, configMap); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create ConfigMap for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created ConfigMap for endpoint %s", endpointSpec.EndpointName))
			}

			// PVC
			pvc := loadgen.CreatePVC(experiment, &endpointSpec, experiment.Status.DataSetMap[endpointSpec.EndpointName])
			if err := ctrl.SetControllerReference(experiment, pvc, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for PVC for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			}
			if err := r.Create(ctx, pvc); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create PVC for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created PVC for endpoint %s", endpointSpec.EndpointName))
			}

			// Copier Pod
			copierPod := &corev1.Pod{}
			copierPodName := types.NamespacedName{
				Namespace: experiment.Namespace,
				Name:      utils.GetTestRunCopierPodName(experiment.Name, endpointSpec.EndpointName),
			}
			if err := r.Get(ctx, copierPodName, copierPod); client.IgnoreNotFound(err) != nil {
				logger.Error(err, fmt.Sprintf("Lost copier Pod: %s", copierPodName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				if copierPod.Status.Phase == corev1.PodFailed {
					experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
					experiment.Status.Error = fmt.Sprintf("Copier Pod failed: %s", copierPodName)
					return true, ctrl.Result{}, nil
				} else if copierPod.Status.Phase == corev1.PodSucceeded {
					doneCounter++
					continue
				}
			}

			copierPod = loadgen.CreateCopierPod(experiment, &endpointSpec, configMap, experiment.Status.DataSetMap[endpointSpec.EndpointName])
			if err := ctrl.SetControllerReference(experiment, copierPod, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for copier Pod for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			}
			if err := r.Create(ctx, copierPod); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create copier Pod for endpoint %s", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created copier Pod for endpoint %s", endpointSpec.EndpointName))
			}
		}
	}
	if doneCounter < len(experiment.Spec.EndpointSpecs) {
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	}

	// Create TestRun for each endpoint
	for _, endpointSpec := range experiment.Spec.EndpointSpecs {
		testRun := loadgen.CreateTestRun(experiment, &endpointSpec)
		switch experiment.Status.DataOptionMap[endpointSpec.EndpointName] {
		case windtunnelv1alpha1.EndpointDataOptionPlainText:
			testRun.Spec.Script = k6v1alpha1.K6Script{
				ConfigMap: k6v1alpha1.K6Configmap{
					Name: utils.GetTestRunConfigMapName(experiment.Name, endpointSpec.EndpointName),
					File: filenameScript,
				},
			}
		case windtunnelv1alpha1.EndpointDataOptionDataSet:
			testRun.Spec.Script = k6v1alpha1.K6Script{
				VolumeClaim: k6v1alpha1.K6VolumeClaim{
					Name: utils.GetTestRunPVCName(experiment.Name, endpointSpec.EndpointName),
					File: filenameScript,
				},
			}
		}
		if err := ctrl.SetControllerReference(experiment, testRun, r.Scheme); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot set controller reference for TestRun for endpoint %s", endpointSpec.EndpointName))
			return true, ctrl.Result{}, err
		}
		if err := r.Create(ctx, testRun); client.IgnoreAlreadyExists(err) != nil {
			logger.Error(err, fmt.Sprintf("Cannot create TestRun for endpoint %s", endpointSpec.EndpointName))
			return true, ctrl.Result{}, err
		} else {
			logger.Info(fmt.Sprintf("Created TestRun for endpoint %s", endpointSpec.EndpointName))
		}
	}

	// Proceed to the next step
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentRunning
	experiment.Status.StartTime = ptr.To(metav1.Now())
	return false, ctrl.Result{}, nil
}

// reconcileRunning reconciles the Experiment when it is running.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileRunning(ctx context.Context, experiment *windtunnelv1alpha1.Experiment) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Check if all TestRuns are finished
	doneCounter := 0
	for _, endpointSpec := range experiment.Spec.EndpointSpecs {
		testRun := &k6v1alpha1.TestRun{}
		testRunName := types.NamespacedName{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunName(experiment.Name, endpointSpec.EndpointName),
		}
		if err := r.Get(ctx, testRunName, testRun); err != nil {
			logger.Error(err, fmt.Sprintf("Lost TestRun: %s", testRunName))
			return true, ctrl.Result{}, err
		} else if err == nil {
			if testRun.Status.Stage == "error" {
				experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
				experiment.Status.Error = fmt.Sprintf("TestRun failed: %s", testRunName)
				return true, ctrl.Result{}, nil
			} else if testRun.Status.Stage == "finished" {
				doneCounter++
			}
		}
	}
	if doneCounter < len(experiment.Spec.EndpointSpecs) {
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	}

	// Proceed to the next step
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentCompleted
	experiment.Status.CompletionTime = ptr.To(metav1.Now())
	return false, ctrl.Result{}, nil
}

// reconcileCompleted reconciles the Experiment when it is completed.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileCompleted(ctx context.Context, experiment *windtunnelv1alpha1.Experiment) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO: cleanup the ConfigMap, PVC, copier Pod and TestRun for each endpoint

	// Try to unlock the Pipeline by setting its status to "Ready"
	experiment.Status.Pipeline.Status.Availability = windtunnelv1alpha1.PipelineReady
	if err := r.Status().Update(ctx, experiment.Status.Pipeline); apierrors.IsConflict(err) {
		// Re-queue the request if conflict occurs
		logger.Info("Conflict occurs when trying to unlock the Pipeline, retrying")
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	} else if err != nil {
		logger.Error(err, "Cannot update the status of Pipeline")
		return true, ctrl.Result{}, err
	}

	// Stop the current reconciliation loop anyway
	return true, ctrl.Result{}, nil
}

// getPipelineEndpoint finds the PipelineEndpoint with the given name in the Pipeline.
// It returns nil if no PipelineEndpoint is found.
func getPipelineEndpoint(pipeline *windtunnelv1alpha1.Pipeline, endpointName string) *windtunnelv1alpha1.PipelineEndpoint {
	for _, pipelineEndpoint := range pipeline.Spec.PipelineEndpoints {
		if pipelineEndpoint.Name == endpointName {
			return &pipelineEndpoint
		}
	}
	return nil
}

// getPipelineEndpointProtocol returns the protocol used by the PipelineEndpoint.
// It returns an empty string if no protocol is specified.
func getPipelineEndpointProtocol(pipelineEndpoint *windtunnelv1alpha1.PipelineEndpoint) windtunnelv1alpha1.EndpointProtocol {
	if pipelineEndpoint.HTTP.URL != "" && pipelineEndpoint.HTTP.Method != "" {
		return windtunnelv1alpha1.EndpointProtocolHTTP
	}
	return ""
}

// getEndpointSpecDataOption returns the data option used by the EndpointSpec.
func getEndpointSpecDataOption(endpointSpec *windtunnelv1alpha1.EndpointSpec) windtunnelv1alpha1.EndpointDataOption {
	if endpointSpec.DataSpec.DataSetRef.Namespace != "" && endpointSpec.DataSpec.DataSetRef.Name != "" {
		return windtunnelv1alpha1.EndpointDataOptionDataSet
	}
	// Fallback to plain text
	return windtunnelv1alpha1.EndpointDataOptionPlainText
}

// getLoadPatternDuration calculates the duration of LoadPattern.
func getLoadPatternDuration(loadPattern *windtunnelv1alpha1.LoadPattern) (metav1.Duration, error) {
	duration := time.Duration(0)
	for _, stage := range loadPattern.Spec.Stages {
		stageDuration, err := time.ParseDuration(stage.Duration)
		if err != nil {
			return metav1.Duration{}, err
		}
		duration += stageDuration
	}
	return metav1.Duration{Duration: duration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExperimentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Experiment{}).
		Complete(r)
}
