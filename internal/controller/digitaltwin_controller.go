package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

const (
	digitalTwinPollingInterval    = 5 * time.Second
	digitalTwinExperimentDuration = 600 // Seconds
)

// DigitalTwinReconciler reconciles a DigitalTwin object
type DigitalTwinReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=digitaltwins,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=digitaltwins/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=digitaltwins/finalizers,verbs=update
//
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=datasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=loadpatterns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *DigitalTwinReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the requested DigitalTwin
	digitalTwin := &windtunnelv1alpha1.DigitalTwin{}
	if err := r.Get(ctx, req.NamespacedName, digitalTwin); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unable to fetch DigitalTwin")
		return ctrl.Result{}, err
	}

	if digitalTwin.Status.JobStatus == "" {
		result, err := r.reconcileCreated(ctx, digitalTwin)
		if err == nil {
			if err := r.Status().Update(ctx, digitalTwin); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
		}
		return result, err
	}

	if digitalTwin.Status.JobStatus == windtunnelv1alpha1.DigitalTwinRunning {
		result, err := r.reconcileRunning(ctx, digitalTwin)
		if err == nil {
			if err := r.Status().Update(ctx, digitalTwin); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
		}
		return result, err
	}

	return ctrl.Result{}, nil
}

// reconcileCreated reconciles the DigitalTwin when it is created.
func (r *DigitalTwinReconciler) reconcileCreated(ctx context.Context, digitalTwin *windtunnelv1alpha1.DigitalTwin) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	switch digitalTwin.Spec.DigitalTwinType {
	case "regular":
		digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinCompleted
		return ctrl.Result{}, nil

	case "schemaaware":
		// Get the DataSet
		dataSet := &windtunnelv1alpha1.DataSet{}
		dataSetName := types.NamespacedName{
			Namespace: digitalTwin.Namespace,
			Name:      digitalTwin.Spec.DataSet.Name,
		}
		if err := r.Get(ctx, dataSetName, dataSet); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get DataSet \"%s\"", dataSetName))
			digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
			digitalTwin.Status.Error = fmt.Sprintf("Cannot find DataSet \"%s\": %s", dataSetName, err)
			return ctrl.Result{}, nil
		}

		// Create DataSets
		for schemaIdx, schemaSelector := range dataSet.Spec.Schemas {
			biasDataSetName := utils.GetBiasDataSetName(digitalTwin.Name, schemaIdx)
			biasDataSet := dataSet.DeepCopy()

			// Avoid error "resourceVersion should not be set on objects to be created"
			biasDataSet.ResourceVersion = ""
			biasDataSet.Name = biasDataSetName
			biasDataSet.Spec.NumberOfFiles = 100
			for schemaSelectorIdx, _ := range biasDataSet.Spec.Schemas {
				if schemaSelectorIdx == schemaIdx {
					biasDataSet.Spec.Schemas[schemaSelectorIdx].NumRecords.Min = 100
					biasDataSet.Spec.Schemas[schemaSelectorIdx].NumRecords.Max = 100
				} else {
					biasDataSet.Spec.Schemas[schemaSelectorIdx].NumRecords.Min = 1
					biasDataSet.Spec.Schemas[schemaSelectorIdx].NumRecords.Max = 1
				}
				biasDataSet.Spec.Schemas[schemaSelectorIdx].NumFilesPerCompressedFile.Min = 1
				biasDataSet.Spec.Schemas[schemaSelectorIdx].NumFilesPerCompressedFile.Max = 1
			}

			if err := ctrl.SetControllerReference(digitalTwin, biasDataSet, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for bias DataSet for Schema \"%s\"", schemaSelector.Name))
				digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
				digitalTwin.Status.Error = fmt.Sprintf("Cannot set controller reference for bias DataSet for Schema \"%s\": %s", schemaSelector.Name, err)
				return ctrl.Result{}, nil
			}
			if err := r.Create(ctx, biasDataSet); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create bias DataSet for Schema \"%s\"", schemaSelector.Name))
				digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
				digitalTwin.Status.Error = fmt.Sprintf("Cannot create bias DataSet for Schema \"%s\": %s", schemaSelector.Name, err)
				return ctrl.Result{}, nil
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created bias DataSet for Schema \"%s\"", schemaSelector.Name))
			}
		}

		// Create LoadPattern
		loadPattern := &windtunnelv1alpha1.LoadPattern{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: digitalTwin.Namespace,
				Name:      utils.GetBiasLoadPatternName(digitalTwin.Name),
			},
			Spec: windtunnelv1alpha1.LoadPatternSpec{
				Stages: []windtunnelv1alpha1.Stage{
					{
						Duration: fmt.Sprintf("%ds", digitalTwinExperimentDuration),
						Target:   int64(digitalTwin.Spec.PipelineCapacity),
					},
				},
				PreAllocatedVUs: 30,
				StartRate:       0,
				TimeUnit:        "1s",
				MaxVUs:          100,
			},
		}
		if err := ctrl.SetControllerReference(digitalTwin, loadPattern, r.Scheme); err != nil {
			logger.Error(err, "Cannot set controller reference for LoadPattern")
			digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
			digitalTwin.Status.Error = fmt.Sprintf("Cannot set controller reference for LoadPattern: %s", err)
			return ctrl.Result{}, nil
		}
		if err := r.Create(ctx, loadPattern); client.IgnoreAlreadyExists(err) != nil {
			logger.Error(err, "Cannot create LoadPattern")
			digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
			digitalTwin.Status.Error = fmt.Sprintf("Cannot create LoadPattern: %s", err)
			return ctrl.Result{}, nil
		}

		// Create Experiments
		for schemaIdx, schemaSelector := range dataSet.Spec.Schemas {
			biasExperiment := &windtunnelv1alpha1.Experiment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: digitalTwin.Namespace,
					Name:      utils.GetBiasExperimentName(digitalTwin.Name, schemaIdx),
				},
				Spec: windtunnelv1alpha1.ExperimentSpec{
					PipelineRef: digitalTwin.Spec.Pipeline,
					EndpointSpecs: []windtunnelv1alpha1.EndpointSpec{
						{
							EndpointName: "upload",
							DataSpec: &windtunnelv1alpha1.DataSpec{
								DataSetRef: &corev1.LocalObjectReference{
									Name: utils.GetBiasDataSetName(digitalTwin.Name, schemaIdx),
								},
							},
							LoadPatternRef: &corev1.ObjectReference{
								Namespace: digitalTwin.Namespace,
								Name:      utils.GetBiasLoadPatternName(digitalTwin.Name),
							},
						},
					},
					UseEndDetection: true,
				},
			}
			if err := ctrl.SetControllerReference(digitalTwin, biasExperiment, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for bias Experiment for Schema \"%s\"", schemaSelector.Name))
				digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
				digitalTwin.Status.Error = fmt.Sprintf("Cannot set controller reference for bias Experiment for Schema \"%s\": %s", schemaSelector.Name, err)
				return ctrl.Result{}, nil
			}
			if err := r.Create(ctx, biasExperiment); client.IgnoreAlreadyExists(err) != nil {
				logger.Error(err, fmt.Sprintf("Cannot create bias Experiment for Schema \"%s\"", schemaSelector.Name))
				digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
				digitalTwin.Status.Error = fmt.Sprintf("Cannot create bias Experiment for Schema \"%s\": %s", schemaSelector.Name, err)
				return ctrl.Result{}, nil
			} else if err == nil {
				logger.Info(fmt.Sprintf("Created bias Experiment for Schema \"%s\"", schemaSelector.Name))
			}
		}

		digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinRunning
		return ctrl.Result{RequeueAfter: digitalTwinPollingInterval}, nil
	}

	return ctrl.Result{}, nil
}

// reconcileRunning reconciles the DigitalTwin when it is running.
func (r *DigitalTwinReconciler) reconcileRunning(ctx context.Context, digitalTwin *windtunnelv1alpha1.DigitalTwin) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	switch digitalTwin.Spec.DigitalTwinType {
	case "regular":
		digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinCompleted
		return ctrl.Result{}, nil

	case "schemaaware":
		// Get the DataSet
		dataSet := &windtunnelv1alpha1.DataSet{}
		dataSetName := types.NamespacedName{
			Namespace: digitalTwin.Namespace,
			Name:      digitalTwin.Spec.DataSet.Name,
		}
		if err := r.Get(ctx, dataSetName, dataSet); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get DataSet \"%s\"", dataSetName))
			digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
			digitalTwin.Status.Error = fmt.Sprintf("Cannot find DataSet \"%s\": %s", dataSetName, err)
			return ctrl.Result{}, nil
		}

		// Check if any Experiment is completed or failed
		for schemaIdx, _ := range dataSet.Spec.Schemas {
			biasExperimentName := types.NamespacedName{
				Namespace: digitalTwin.Namespace,
				Name:      utils.GetBiasExperimentName(digitalTwin.Name, schemaIdx),
			}
			biasExperiment := &windtunnelv1alpha1.Experiment{}
			if err := r.Get(ctx, biasExperimentName, biasExperiment); err != nil {
				logger.Error(err, fmt.Sprintf("Lost bias Experiment \"%s\"", biasExperimentName))
				digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
				digitalTwin.Status.Error = fmt.Sprintf("Lost bias Experiment \"%s\": %s", biasExperimentName, err)
				return ctrl.Result{}, nil
			}

			if biasExperiment.Status.JobStatus == windtunnelv1alpha1.ExperimentCompleted {
				continue
			} else if biasExperiment.Status.JobStatus == windtunnelv1alpha1.ExperimentFailed {
				logger.Info(fmt.Sprintf("Bias Experiment \"%s\" failed", biasExperimentName))
				digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinFailed
				digitalTwin.Status.Error = fmt.Sprintf("Experiment \"%s\" failed", biasExperimentName)
				return ctrl.Result{}, nil
			} else {
				return ctrl.Result{RequeueAfter: digitalTwinPollingInterval}, nil
			}
		}

		// All Experiments are completed
		logger.Info("All bias Experiments are completed")
		digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinCompleted
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DigitalTwinReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.DigitalTwin{}).
		Complete(r)
}
