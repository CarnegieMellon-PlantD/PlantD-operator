package controller

import (
	"context"
	"fmt"
	"time"

	k6v1alpha1 "github.com/grafana/k6-operator/api/v1alpha1"
	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/digitaltwin"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/loadgen"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

const (
	experimentFinalizerName         = "experiment.windtunnel.plantd.org/finalizer"
	experimentPollingInterval       = 5 * time.Second
	experimentEndDetectorDebounce   = 30 // Seconds
	experimentEndDetectorWindow     = 90 // Seconds
	experimentEndDetectorAdjustment = 60 // Seconds
)

var (
	filenameScript                   = config.GetString("loadGenerator.filename.script")
	metricsServiceLabelKeyExperiment = config.GetString("monitor.service.labelKeys.experiment")
)

// ExperimentReconciler reconciles a Experiment object
type ExperimentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// ExperimentReconcilerContext contains the context for the ExperimentReconciler
type ExperimentReconcilerContext struct {
	Pipeline             *windtunnelv1alpha1.Pipeline
	Endpoints            map[string]*windtunnelv1alpha1.PipelineEndpoint
	EndpointProtocols    map[string]windtunnelv1alpha1.EndpointProtocol
	EndpointDataOptions  map[string]windtunnelv1alpha1.EndpointDataOption
	EndpointDataSets     map[string]*windtunnelv1alpha1.DataSet
	EndpointLoadPatterns map[string]*windtunnelv1alpha1.LoadPattern
}

// NewExperimentReconcilerContext creates a new ExperimentReconcilerContext
func NewExperimentReconcilerContext() *ExperimentReconcilerContext {
	return &ExperimentReconcilerContext{
		Pipeline:             nil,
		Endpoints:            make(map[string]*windtunnelv1alpha1.PipelineEndpoint),
		EndpointProtocols:    make(map[string]windtunnelv1alpha1.EndpointProtocol),
		EndpointDataOptions:  make(map[string]windtunnelv1alpha1.EndpointDataOption),
		EndpointDataSets:     make(map[string]*windtunnelv1alpha1.DataSet),
		EndpointLoadPatterns: make(map[string]*windtunnelv1alpha1.LoadPattern),
	}
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
		logger.Error(err, "Unable to fetch Experiment")
		return ctrl.Result{}, err
	}

	// Check if the Experiment is being deleted
	if experiment.DeletionTimestamp.IsZero() {
		// Add finalizer to the Experiment
		if !controllerutil.ContainsFinalizer(experiment, experimentFinalizerName) {
			controllerutil.AddFinalizer(experiment, experimentFinalizerName)
			if err := r.Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot add finalizer")
				return ctrl.Result{}, err
			}
			logger.Info("Added finalizer")
		}
	} else {
		if controllerutil.ContainsFinalizer(experiment, experimentFinalizerName) {
			// Try to release the Pipeline
			pipeline := &windtunnelv1alpha1.Pipeline{}
			pipelineName := types.NamespacedName{
				Namespace: experiment.Namespace,
				Name:      experiment.Spec.PipelineRef.Name,
			}
			if err := r.Get(ctx, pipelineName, pipeline); err != nil {
				logger.Error(err, fmt.Sprintf("Lost Pipeline \"%s\"", pipelineName))
			} else {
				if err := r.releasePipeline(ctx, pipeline); err != nil {
					return ctrl.Result{}, err
				}
			}

			// Remove finalizer from the Experiment
			controllerutil.RemoveFinalizer(experiment, experimentFinalizerName)
			if err := r.Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot remove finalizer")
				return ctrl.Result{}, err
			}
			logger.Info("Removed finalizer")
		}
		return ctrl.Result{}, nil
	}

	// No need to reconcile if the Experiment is completed or failed
	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentCompleted || experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentFailed {
		return ctrl.Result{}, nil
	}

	// Initiate the reconciler context
	rc := NewExperimentReconcilerContext()
	stop, result, err := r.getRelatedResources(ctx, experiment, rc)
	if stop {
		if err := r.Status().Update(ctx, experiment); err != nil {
			logger.Error(err, "Cannot update the status")
			return ctrl.Result{}, err
		}
		return result, err
	}

	if experiment.Status.JobStatus == "" {
		stop, result, err := r.reconcileCreated(ctx, experiment, rc)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentScheduled {
		stop, result, err := r.reconciledScheduled(ctx, experiment, rc)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentWaitingDataSet {
		stop, result, err := r.reconcileWaitingDataSet(ctx, experiment, rc)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentWaitingPipeline {
		stop, result, err := r.reconcileWaitingPipeline(ctx, experiment, rc)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentInitializing {
		stop, result, err := r.reconcileInitializing(ctx, experiment, rc)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentRunning {
		stop, result, err := r.reconcileRunning(ctx, experiment, rc)
		if stop {
			if err := r.Status().Update(ctx, experiment); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
			return result, err
		}
	}

	if experiment.Status.JobStatus == windtunnelv1alpha1.ExperimentDraining {
		stop, result, err := r.reconcileDraining(ctx, experiment, rc)
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

// getRelatedResources gets the related resources used by the Experiment and update the reconciler context.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) getRelatedResources(ctx context.Context, experiment *windtunnelv1alpha1.Experiment, rc *ExperimentReconcilerContext) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Pipeline used by the Experiment
	pipeline := &windtunnelv1alpha1.Pipeline{}
	pipelineName := types.NamespacedName{
		Namespace: experiment.Namespace,
		Name:      experiment.Spec.PipelineRef.Name,
	}
	if err := r.Get(ctx, pipelineName, pipeline); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot get Pipeline \"%s\"", pipelineName))
		experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
		experiment.Status.Error = fmt.Sprintf("Cannot find Pipeline \"%s\": %s", pipelineName, err)
		return true, ctrl.Result{}, nil
	}
	rc.Pipeline = pipeline

	for _, endpointSpec := range experiment.Spec.EndpointSpecs {
		// Find the PipelineEndpoint referenced by the EndpointSpec
		pipelineEndpoint := getPipelineEndpoint(pipeline, endpointSpec.EndpointName)
		if pipelineEndpoint == nil {
			logger.Error(nil, fmt.Sprintf("Cannot find endpoint \"%s\"", endpointSpec.EndpointName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Cannot find endpoint \"%s\"", endpointSpec.EndpointName)
			return true, ctrl.Result{}, nil
		}
		rc.Endpoints[endpointSpec.EndpointName] = pipelineEndpoint

		// Determine the protocol used by the PipelineEndpoint
		protocol := getPipelineEndpointProtocol(pipelineEndpoint)
		if protocol == "" {
			logger.Error(nil, fmt.Sprintf("Unspecified protocol in endpoint \"%s\"", endpointSpec.EndpointName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Unspecified protocol in endpoint \"%s\"", endpointSpec.EndpointName)
			return true, ctrl.Result{}, nil
		}
		rc.EndpointProtocols[endpointSpec.EndpointName] = protocol

		// Determine the data option used by the EndpointSpec
		dataOption := getEndpointSpecDataOption(&endpointSpec)
		if dataOption == "" {
			logger.Error(nil, fmt.Sprintf("Unspecified data option in endpoint \"%s\"", endpointSpec.EndpointName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Unspecified data option in endpoint \"%s\"", endpointSpec.EndpointName)
			return true, ctrl.Result{}, nil
		}
		rc.EndpointDataOptions[endpointSpec.EndpointName] = dataOption

		// Fetch the DataSet used by the EndpointSpec
		if dataOption == windtunnelv1alpha1.EndpointDataOptionDataSet {
			dataSet := &windtunnelv1alpha1.DataSet{}
			dataSetName := types.NamespacedName{
				Namespace: experiment.Namespace,
				Name:      endpointSpec.DataSpec.DataSetRef.Name,
			}
			if err := r.Get(ctx, dataSetName, dataSet); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot get DataSet \"%s\" for endpoint \"%s\"",
					dataSetName, endpointSpec.EndpointName,
				))
				experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
				experiment.Status.Error = fmt.Sprintf("Cannot find DataSet \"%s\" for endpoint \"%s\": %s",
					dataSetName, endpointSpec.EndpointName, err,
				)
				return true, ctrl.Result{}, nil
			}
			rc.EndpointDataSets[endpointSpec.EndpointName] = dataSet
		}

		// Fetch the LoadPattern used by the EndpointSpec
		loadPattern := &windtunnelv1alpha1.LoadPattern{}
		loadPatternName := types.NamespacedName{
			Namespace: endpointSpec.LoadPatternRef.Namespace,
			Name:      endpointSpec.LoadPatternRef.Name,
		}
		if err := r.Get(ctx, loadPatternName, loadPattern); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get LoadPattern \"%s\" for endpoint \"%s\"",
				loadPatternName, endpointSpec.EndpointName,
			))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Cannot find LoadPattern \"%s\" for endpoint \"%s\": %s",
				loadPatternName, endpointSpec.EndpointName, err,
			)
			return true, ctrl.Result{}, nil
		}
		rc.EndpointLoadPatterns[endpointSpec.EndpointName] = loadPattern
	}

	// Proceed
	return false, ctrl.Result{}, nil
}

// reconcileCreated reconciles the Experiment when it is created.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileCreated(ctx context.Context, experiment *windtunnelv1alpha1.Experiment, rc *ExperimentReconcilerContext) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Copy values from the Pipeline into the status
	experiment.Status.EnableCostCalculation = rc.Pipeline.Spec.EnableCostCalculation
	experiment.Status.CloudProvider = rc.Pipeline.Spec.CloudProvider
	experiment.Status.Tags = rc.Pipeline.Spec.Tags

	// Calculate the durations in the status
	experiment.Status.Durations = make(map[string]*metav1.Duration, len(experiment.Spec.EndpointSpecs))
	for endpointName, endpointLoadPattern := range rc.EndpointLoadPatterns {
		duration, err := getLoadPatternDuration(endpointLoadPattern)
		if err != nil {
			logger.Error(err, fmt.Sprintf("Cannot calculate the duration for endpoint \"%s\"", endpointName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Cannot calculate the duration for endpoint \"%s\": %s", endpointName, err)
			return true, ctrl.Result{}, nil
		}
		experiment.Status.Durations[endpointName] = duration
	}

	// Proceed to the next state
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentScheduled
	return false, ctrl.Result{}, nil
}

// reconciledScheduled reconciles the Experiment when it is scheduled.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconciledScheduled(ctx context.Context, experiment *windtunnelv1alpha1.Experiment, rc *ExperimentReconcilerContext) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Check if the scheduled time has been reached
	curTime := time.Now()
	if experiment.Spec.ScheduledTime != nil && curTime.Before(experiment.Spec.ScheduledTime.Time) {
		// Time has not been reached, calculate the waiting time
		waitTime := experiment.Spec.ScheduledTime.Time.Sub(curTime)
		logger.Info(fmt.Sprintf("Scheduled time has not been reached yet, waiting for \"%s\"", waitTime))
		return true, ctrl.Result{RequeueAfter: waitTime}, nil
	}

	// Time has been reached, proceed to the next state
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentWaitingDataSet
	return false, ctrl.Result{}, nil
}

// reconcileWaitingDataSet reconciles the Experiment when it is waiting for the DataSet.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileWaitingDataSet(ctx context.Context, experiment *windtunnelv1alpha1.Experiment, rc *ExperimentReconcilerContext) (bool, ctrl.Result, error) {
	for _, dataSet := range rc.EndpointDataSets {
		// Check if the DataSet is in "Success" status
		if dataSet.Status.JobStatus != windtunnelv1alpha1.DataSetJobSuccess {
			return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
		}
	}

	// Proceed to the next state
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentWaitingPipeline
	return false, ctrl.Result{}, nil
}

// reconcileWaitingPipeline reconciles the Experiment when it is waiting for the Pipeline.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileWaitingPipeline(ctx context.Context, experiment *windtunnelv1alpha1.Experiment, rc *ExperimentReconcilerContext) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Check if the Pipeline is in "Ready" status
	if rc.Pipeline.Status.Availability != windtunnelv1alpha1.PipelineReady {
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	}

	// Lock the Pipeline by setting its status to "In-Use"
	rc.Pipeline.Status.Availability = windtunnelv1alpha1.PipelineInUse
	if err := r.Status().Update(ctx, rc.Pipeline); err != nil {
		logger.Error(err, "Cannot update the status of the Pipeline")
		return true, ctrl.Result{}, err
	}
	logger.Info("Set the Pipeline status to \"In-Use\"")

	// Set the Experiment label for the metrics Service
	if containMetricsEndpoint(rc.Pipeline) {
		metricsService := &corev1.Service{}
		var metricsServiceName types.NamespacedName
		if rc.Pipeline.Spec.InCluster {
			metricsServiceName = types.NamespacedName{
				Namespace: rc.Pipeline.Namespace,
				Name:      rc.Pipeline.Spec.MetricsEndpoint.ServiceRef.Name,
			}
		} else {
			metricsServiceName = types.NamespacedName{
				Namespace: rc.Pipeline.Namespace,
				Name:      utils.GetMetricsServiceName(rc.Pipeline.Name),
			}
		}
		if err := r.Get(ctx, metricsServiceName, metricsService); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get metrics Service \"%s\"", metricsServiceName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Cannot find metrics Service \"%s\": %s", metricsServiceName, err)
			return true, ctrl.Result{}, nil
		}
		if metricsService.Labels == nil {
			metricsService.Labels = make(map[string]string, 1)
		}
		metricsService.Labels[metricsServiceLabelKeyExperiment] = experiment.Name
		if err := r.Update(ctx, metricsService); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot add Experiment label to metrics Service \"%s\"", metricsServiceName))
			return true, ctrl.Result{}, err
		}
		logger.Info(fmt.Sprintf("Added Experiment label to metrics Service \"%s\"", metricsServiceName))
	}

	// Proceed to the next state
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentInitializing
	return false, ctrl.Result{}, nil
}

// reconcileInitializing reconciles the Experiment when it is initializing.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileInitializing(ctx context.Context, experiment *windtunnelv1alpha1.Experiment, rc *ExperimentReconcilerContext) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Perform health check
	for _, healthCheckURL := range rc.Pipeline.Spec.HealthCheckURLs {
		if err := utils.CheckHealth(healthCheckURL); err != nil {
			logger.Error(err, "Pipeline health check failed")
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Pipeline health check failed: %s", err)
			return true, ctrl.Result{}, nil
		}
	}

	// Create ConfigMap, PVC and copier Pod for each endpoint
	doneCounter := 0
	for endpointIdx, endpointSpec := range experiment.Spec.EndpointSpecs {
		switch rc.EndpointDataOptions[endpointSpec.EndpointName] {
		case windtunnelv1alpha1.EndpointDataOptionPlainText:
			// ConfigMap
			configMap, err := loadgen.CreateConfigMapWithPlainText(
				experiment,
				endpointIdx,
				rc.Endpoints[endpointSpec.EndpointName],
				endpointSpec.DataSpec.PlainText,
				rc.EndpointLoadPatterns[endpointSpec.EndpointName],
				rc.EndpointProtocols[endpointSpec.EndpointName],
			)
			if err != nil {
				logger.Error(err, fmt.Sprintf("Cannot create manifest for ConfigMap for endpoint \"%s\"",
					endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			}
			if err := ctrl.SetControllerReference(experiment, configMap, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for ConfigMap for endpoint \"%s\"",
					endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			}
			if err := r.Create(ctx, configMap); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create ConfigMap for endpoint \"%s\"", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created ConfigMap for endpoint \"%s\"", endpointSpec.EndpointName))
			}

			doneCounter++

		case windtunnelv1alpha1.EndpointDataOptionDataSet:
			// ConfigMap
			configMap, err := loadgen.CreateConfigMapWithDataSet(
				experiment,
				endpointIdx,
				rc.Endpoints[endpointSpec.EndpointName],
				rc.EndpointDataSets[endpointSpec.EndpointName],
				rc.EndpointLoadPatterns[endpointSpec.EndpointName],
				rc.EndpointProtocols[endpointSpec.EndpointName],
			)
			if err != nil {
				logger.Error(err, fmt.Sprintf("Cannot create manifest for ConfigMap for endpoint \"%s\"",
					endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			}
			if err := ctrl.SetControllerReference(experiment, configMap, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for ConfigMap for endpoint \"%s\"",
					endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			}
			if err := r.Create(ctx, configMap); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create ConfigMap for endpoint \"%s\"", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created ConfigMap for endpoint \"%s\"", endpointSpec.EndpointName))
			}

			// PVC
			pvc := loadgen.CreatePVC(experiment, endpointIdx, &endpointSpec, rc.EndpointDataSets[endpointSpec.EndpointName])
			if err := ctrl.SetControllerReference(experiment, pvc, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for PVC for endpoint \"%s\"",
					endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			}
			if err := r.Create(ctx, pvc); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create PVC for endpoint \"%s\"", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created PVC for endpoint \"%s\"", endpointSpec.EndpointName))
			}

			// Copier Job
			copierJob := &kbatch.Job{}
			copierJobName := types.NamespacedName{
				Namespace: experiment.Namespace,
				Name:      utils.GetTestRunCopierJobName(experiment.Name, endpointIdx),
			}
			if err := r.Get(ctx, copierJobName, copierJob); client.IgnoreNotFound(err) != nil {
				logger.Error(err, fmt.Sprintf("Lost copier Job \"%s\" for endpoint \"%s\"",
					copierJobName, endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			} else if err == nil {
				jobFinished, jobConditionType := isJobFinished(copierJob)
				if jobFinished {
					logger.Info(fmt.Sprintf("Copier Job \"%s\" for endpoint \"%s\" finished",
						copierJobName, endpointSpec.EndpointName,
					))
					switch jobConditionType {
					case kbatch.JobComplete:
						doneCounter++
						continue
					case kbatch.JobFailed:
						experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
						experiment.Status.Error = fmt.Sprintf("Copier Job \"%s\" for endpoint \"%s\" failed",
							copierJobName, endpointSpec.EndpointName,
						)
						return true, ctrl.Result{}, nil
					}
				}
			}
			copierJob = loadgen.CreateCopierJob(experiment, endpointIdx, &endpointSpec, configMap, rc.EndpointDataSets[endpointSpec.EndpointName])
			if err := ctrl.SetControllerReference(experiment, copierJob, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for copier Job for endpoint \"%s\"",
					endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			}
			if err := r.Create(ctx, copierJob); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create copier Job for endpoint \"%s\"", endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created copier Job for endpoint \"%s\"", endpointSpec.EndpointName))
			}
		}
	}
	if doneCounter < len(experiment.Spec.EndpointSpecs) {
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	}

	// Create TestRun for each endpoint
	for endpointIdx, endpointSpec := range experiment.Spec.EndpointSpecs {
		testRun := loadgen.CreateTestRun(experiment, endpointIdx, &endpointSpec)
		switch rc.EndpointDataOptions[endpointSpec.EndpointName] {
		case windtunnelv1alpha1.EndpointDataOptionPlainText:
			testRun.Spec.Script = k6v1alpha1.K6Script{
				ConfigMap: k6v1alpha1.K6Configmap{
					Name: utils.GetTestRunName(experiment.Name, endpointIdx),
					File: filenameScript,
				},
			}

		case windtunnelv1alpha1.EndpointDataOptionDataSet:
			testRun.Spec.Script = k6v1alpha1.K6Script{
				VolumeClaim: k6v1alpha1.K6VolumeClaim{
					Name: utils.GetTestRunName(experiment.Name, endpointIdx),
					File: filenameScript,
				},
			}
		}
		if err := ctrl.SetControllerReference(experiment, testRun, r.Scheme); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot set controller reference for TestRun for endpoint \"%s\"",
				endpointSpec.EndpointName,
			))
			return true, ctrl.Result{}, err
		}
		if err := r.Create(ctx, testRun); client.IgnoreAlreadyExists(err) != nil {
			logger.Error(err, fmt.Sprintf("Cannot create TestRun for endpoint \"%s\"", endpointSpec.EndpointName))
			return true, ctrl.Result{}, err
		} else {
			logger.Info(fmt.Sprintf("Created TestRun for endpoint \"%s\"", endpointSpec.EndpointName))
		}
	}

	// Set the start time first because end detector Job relies on it
	experiment.Status.StartTime = ptr.To(metav1.Now())

	// Create end detector Job if needed
	if experiment.Spec.UseEndDetection {
		endDetectorJob, err := digitaltwin.CreateEndDetectorJob(experiment, experimentEndDetectorDebounce, experimentEndDetectorWindow, experimentEndDetectorAdjustment)
		if err != nil {
			logger.Error(err, "Cannot create manifest for end detector Job")
			return true, ctrl.Result{}, err
		}
		if err := ctrl.SetControllerReference(experiment, endDetectorJob, r.Scheme); err != nil {
			logger.Error(err, "Cannot set controller reference for end detector Job")
			return true, ctrl.Result{}, err
		}
		if err := r.Create(ctx, endDetectorJob); client.IgnoreAlreadyExists(err) != nil {
			logger.Error(err, "Cannot create end detector Job")
			return true, ctrl.Result{}, err
		} else {
			logger.Info("Created end detector Job")
		}
	}

	// Proceed to the next state
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentRunning
	return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
}

// reconcileRunning reconciles the Experiment when it is running.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileRunning(ctx context.Context, experiment *windtunnelv1alpha1.Experiment, rc *ExperimentReconcilerContext) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Check if all TestRuns are finished
	doneCounter := 0
	for endpointIdx, endpointSpec := range experiment.Spec.EndpointSpecs {
		testRun := &k6v1alpha1.TestRun{}
		testRunName := types.NamespacedName{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
		}
		if err := r.Get(ctx, testRunName, testRun); err != nil {
			logger.Error(err, fmt.Sprintf("Lost TestRun \"%s\" for endpoint \"%s\"", testRunName, endpointSpec.EndpointName))
			return true, ctrl.Result{}, err
		} else if err == nil {
			if testRun.Status.Stage == "error" {
				experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
				experiment.Status.Error = fmt.Sprintf("TestRun \"%s\" for endpoint \"%s\" failed",
					testRunName, endpointSpec.EndpointName,
				)
				return true, ctrl.Result{}, nil
			} else if testRun.Status.Stage == "finished" {
				doneCounter++
			}
		}
	}
	if doneCounter < len(experiment.Spec.EndpointSpecs) {
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	}

	// Remove the resources created for the Experiment
	for endpointIdx, endpointSpec := range experiment.Spec.EndpointSpecs {
		// TestRun
		testRun := &k6v1alpha1.TestRun{}
		testRunName := types.NamespacedName{
			Namespace: experiment.Namespace,
			Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
		}
		if err := r.Get(ctx, testRunName, testRun); client.IgnoreNotFound(err) != nil {
			logger.Error(err, fmt.Sprintf("Lost TestRun \"%s\" for endpoint \"%s\"", testRunName, endpointSpec.EndpointName))
			return true, ctrl.Result{}, err
		} else if err == nil {
			if err := r.Delete(ctx, testRun); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot delete TestRun \"%s\" for endpoint \"%s\"", testRunName, endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			}
			logger.Info(fmt.Sprintf("Deleted TestRun \"%s\" for endpoint \"%s\"", testRunName, endpointSpec.EndpointName))
		}

		switch rc.EndpointDataOptions[endpointSpec.EndpointName] {
		case windtunnelv1alpha1.EndpointDataOptionPlainText:
			// ConfigMap
			configMap := &corev1.ConfigMap{}
			configMapName := types.NamespacedName{
				Namespace: experiment.Namespace,
				Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
			}
			if err := r.Get(ctx, configMapName, configMap); client.IgnoreNotFound(err) != nil {
				logger.Error(err, fmt.Sprintf("Lost ConfigMap \"%s\" for endpoint \"%s\"",
					configMapName, endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			} else if err == nil {
				if err := r.Delete(ctx, configMap); err != nil {
					logger.Error(err, fmt.Sprintf("Cannot delete ConfigMap \"%s\" for endpoint \"%s\"",
						configMapName, endpointSpec.EndpointName,
					))
					return true, ctrl.Result{}, err
				}
				logger.Info(fmt.Sprintf("Deleted ConfigMap \"%s\" for endpoint \"%s\"",
					configMapName, endpointSpec.EndpointName,
				))
			}

		case windtunnelv1alpha1.EndpointDataOptionDataSet:
			// Copier Job
			copierJob := &kbatch.Job{}
			copierJobName := types.NamespacedName{
				Namespace: experiment.Namespace,
				Name:      utils.GetTestRunCopierJobName(experiment.Name, endpointIdx),
			}
			if err := r.Get(ctx, copierJobName, copierJob); client.IgnoreNotFound(err) != nil {
				logger.Error(err, fmt.Sprintf("Lost copier Job \"%s\" for endpoint \"%s\"",
					copierJobName, endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			} else if err == nil {
				// By default, the Pod of the Job will be reserved after the Job is deleted,
				// and Kubernetes will raise a warning.
				// Set the propagation policy to "Background" to avoid the warning and delete the Pod.
				if err := r.Delete(ctx, copierJob, &client.DeleteOptions{
					PropagationPolicy: ptr.To(metav1.DeletePropagationBackground),
				}); err != nil {
					logger.Error(err, fmt.Sprintf("Cannot delete copier Job \"%s\" for endpoint \"%s\"",
						copierJobName, endpointSpec.EndpointName,
					))
					return true, ctrl.Result{}, err
				}
				logger.Info(fmt.Sprintf("Deleted copier Job \"%s\" for endpoint \"%s\"",
					copierJobName, endpointSpec.EndpointName,
				))
			}

			// ConfigMap
			configMap := &corev1.ConfigMap{}
			configMapName := types.NamespacedName{
				Namespace: experiment.Namespace,
				Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
			}
			if err := r.Get(ctx, configMapName, configMap); client.IgnoreNotFound(err) != nil {
				logger.Error(err, fmt.Sprintf("Lost ConfigMap \"%s\" for endpoint \"%s\"",
					configMapName, endpointSpec.EndpointName,
				))
				return true, ctrl.Result{}, err
			} else if err == nil {
				if err := r.Delete(ctx, configMap); err != nil {
					logger.Error(err, fmt.Sprintf("Cannot delete ConfigMap \"%s\" for endpoint \"%s\"",
						configMapName, endpointSpec.EndpointName,
					))
					return true, ctrl.Result{}, err
				}
				logger.Info(fmt.Sprintf("Deleted ConfigMap \"%s\" for endpoint \"%s\"",
					configMapName, endpointSpec.EndpointName,
				))
			}

			// PVC
			pvc := &corev1.PersistentVolumeClaim{}
			pvcName := types.NamespacedName{
				Namespace: experiment.Namespace,
				Name:      utils.GetTestRunName(experiment.Name, endpointIdx),
			}
			if err := r.Get(ctx, pvcName, pvc); client.IgnoreNotFound(err) != nil {
				logger.Error(err, fmt.Sprintf("Lost PVC \"%s\" for endpoint \"%s\"", pvcName, endpointSpec.EndpointName))
				return true, ctrl.Result{}, err
			} else if err == nil {
				if err := r.Delete(ctx, pvc); err != nil {
					logger.Error(err, fmt.Sprintf("Cannot delete PVC \"%s\" for endpoint \"%s\"", pvcName, endpointSpec.EndpointName))
					return true, ctrl.Result{}, err
				}
				logger.Info(fmt.Sprintf("Deleted PVC \"%s\" for endpoint \"%s\"", pvcName, endpointSpec.EndpointName))
			}
		}
	}

	// Proceed to the next state
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentDraining
	experiment.Status.DrainingStartTime = ptr.To(metav1.Now())
	return false, ctrl.Result{}, nil
}

// reconcileDraining reconciles the Experiment when it is draining.
// It returns a flag of whether the current reconciliation loop should stop,
// the reconciliation result, and an error, if any.
func (r *ExperimentReconciler) reconcileDraining(ctx context.Context, experiment *windtunnelv1alpha1.Experiment, rc *ExperimentReconcilerContext) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)

	curTime := time.Now()

	// Check if the end detector Job is finished
	if experiment.Spec.UseEndDetection {
		endDetectorJob := &kbatch.Job{}
		endDetectorJobName := types.NamespacedName{
			Namespace: experiment.Namespace,
			Name:      utils.GetEndDetectorJobName(experiment.Name),
		}
		if err := r.Get(ctx, endDetectorJobName, endDetectorJob); err != nil {
			logger.Error(err, fmt.Sprintf("Lost end detector Job \"%s\"", endDetectorJobName))
			experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
			experiment.Status.Error = fmt.Sprintf("Lost end detector Job \"%s\": %s", endDetectorJobName, err)
		}
		jobFinished, jobConditionType := isJobFinished(endDetectorJob)
		if jobFinished {
			logger.Info(fmt.Sprintf("End detector Job \"%s\" finished", endDetectorJobName))
			switch jobConditionType {
			case kbatch.JobComplete:
				experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentCompleted
				calculatedCompletionTime := curTime.Add(-experimentEndDetectorAdjustment * time.Second)
				experiment.Status.CompletionTime = &metav1.Time{Time: calculatedCompletionTime}
			case kbatch.JobFailed:
				logger.Error(nil, fmt.Sprintf("End detector Job \"%s\" failed", endDetectorJobName))
				experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentFailed
				experiment.Status.Error = fmt.Sprintf("End detector Job \"%s\" failed", endDetectorJobName)
			}
		}
		return true, ctrl.Result{RequeueAfter: experimentPollingInterval}, nil
	}

	// Check if the draining time has been fulfilled
	if !experiment.Spec.UseEndDetection && experiment.Spec.DrainingTime != nil && curTime.Before(experiment.Status.DrainingStartTime.Time.Add(experiment.Spec.DrainingTime.Duration)) {
		// Time has not been fulfilled, calculate the waiting time
		waitTime := experiment.Status.DrainingStartTime.Time.Add(experiment.Spec.DrainingTime.Duration).Sub(curTime)
		logger.Info(fmt.Sprintf("Pipeline is draining, waiting for \"%s\"", waitTime))
		return true, ctrl.Result{RequeueAfter: waitTime}, nil
	}

	// Release the Pipeline
	if err := r.releasePipeline(ctx, rc.Pipeline); err != nil {
		return true, ctrl.Result{}, err
	}

	// Stop the reconciliation loop
	experiment.Status.JobStatus = windtunnelv1alpha1.ExperimentCompleted
	experiment.Status.CompletionTime = ptr.To(metav1.Now())
	return true, ctrl.Result{}, nil
}

// releasePipeline unlocks the Pipeline by setting its status to "Ready".
// It also removes the Experiment label from the metrics Service.
func (r *ExperimentReconciler) releasePipeline(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline) error {
	logger := log.FromContext(ctx)

	// Unlock the Pipeline by setting its status to "Ready"
	pipeline.Status.Availability = windtunnelv1alpha1.PipelineReady
	if err := r.Status().Update(ctx, pipeline); err != nil {
		logger.Error(err, "Cannot update the status of the Pipeline")
		return err
	}
	logger.Info("Set the Pipeline status to \"Ready\"")

	// Remove the Experiment label for the metrics Service
	if containMetricsEndpoint(pipeline) {
		metricsService := &corev1.Service{}
		var metricsServiceName types.NamespacedName
		if pipeline.Spec.InCluster {
			metricsServiceName = types.NamespacedName{
				Namespace: pipeline.Namespace,
				Name:      pipeline.Spec.MetricsEndpoint.ServiceRef.Name,
			}
		} else {
			metricsServiceName = types.NamespacedName{
				Namespace: pipeline.Namespace,
				Name:      utils.GetMetricsServiceName(pipeline.Name),
			}
		}
		if err := r.Get(ctx, metricsServiceName, metricsService); err != nil {
			logger.Error(err, fmt.Sprintf("Lost metrics Service \"%s\"", metricsService))
		} else if metricsService.Labels != nil {
			if _, ok := metricsService.Labels[metricsServiceLabelKeyExperiment]; ok {
				delete(metricsService.Labels, metricsServiceLabelKeyExperiment)
				if err := r.Update(ctx, metricsService); err != nil {
					logger.Error(err, fmt.Sprintf("Cannot remove Experiment label from metrics Service \"%s\"", metricsServiceName))
					return err
				}
				logger.Info(fmt.Sprintf("Removed Experiment label from metrics Service \"%s\"", metricsServiceName))
			}
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExperimentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Experiment{}).
		Complete(r)
}
