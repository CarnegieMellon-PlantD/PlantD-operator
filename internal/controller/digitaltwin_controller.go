package controller

import (
	"context"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
)

const (
	digitalTwinPollingInterval     = 10 * time.Second
	digitalTwinLoadPatternDuration = 600 // Seconds
)

// DigitalTwinReconciler reconciles a DigitalTwin object
type DigitalTwinReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=digitaltwins,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=digitaltwins/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=digitaltwins/finalizers,verbs=update

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

	// Initialize the DigitalTwin
	if digitalTwin.Status.JobStatus == "" {
		if err := r.Status().Update(ctx, digitalTwin); err != nil {
			logger.Error(err, "Cannot update the status")
			return ctrl.Result{}, err
		}
	}

	/**
	for _, task := range digitalTwin.Spec.Tasks {
		dataSetName := fmt.Sprintf("%s-dataset-pure-%s", digitalTwin.Name, task.Name)
		loadPatternName := fmt.Sprintf("%s-loadpattern-%s", digitalTwin.Name, task.Name)
		experimentName := fmt.Sprintf("%s-experiment-pure-%s", digitalTwin.Name, task.Name)

		// Create a DataSet
		dataSet := &windtunnelv1alpha1.DataSet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: digitalTwin.Namespace,
				Name:      dataSetName,
			},
			Spec: windtunnelv1alpha1.DataSetSpec{
				CompressPerSchema:    digitalTwin.Spec.DataSetConfig.CompressPerSchema,
				CompressedFileFormat: digitalTwin.Spec.DataSetConfig.CompressedFileFormat,
				FileFormat:           digitalTwin.Spec.DataSetConfig.FileFormat,
				Parallelism:          1,
				NumberOfFiles:        int32(DATASET_SIZE),
				Schemas: []windtunnelv1alpha1.SchemaSelector{
					{
						Name: task.Name,
						NumRecords: windtunnelv1alpha1.NaturalIntRange{
							Min: 1,
							Max: 1,
						},
						NumFilesPerCompressedFile: windtunnelv1alpha1.NaturalIntRange{
							Min: 1,
							Max: 1,
						},
					},
				},
			},
		}
		if err := r.Create(ctx, dataSet); client.IgnoreAlreadyExists(err) != nil {
			logger.Error(err, "Failed to create DataSet")
			return ctrl.Result{}, err
		}
		logger.Info("Created DataSet", "name", dataSetName, "namespace", digitalTwin.Namespace)

		// Create LoadPattern
		loadPattern := &windtunnelv1alpha1.LoadPattern{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: digitalTwin.Namespace,
				Name:      loadPatternName,
			},
			Spec: windtunnelv1alpha1.LoadPatternSpec{
				PreAllocatedVUs: 30,
				MaxVUs:          100,
				TimeUnit:        "1s",
				StartRate:       0,
				Stages: []windtunnelv1alpha1.Stage{
					{
						Duration: fmt.Sprintf("%ds", EXPERIMENT_DURATION),
						Target:   int64(maxRate),
					},
				},
			},
		}
		if err := r.Create(ctx, loadPattern); client.IgnoreAlreadyExists(err) != nil {
			logger.Error(err, "Failed to create LoadPattern")
			return ctrl.Result{}, err
		}
		logger.Info("Created LoadPattern", "name", loadPatternName, "namespace", digitalTwin.Namespace)

		// Create Experiment
		experiment := &windtunnelv1alpha1.Experiment{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: digitalTwin.Namespace,
				Name:      experimentName,
			},
			Spec: windtunnelv1alpha1.ExperimentSpec{
				PipelineRef: &digitalTwin.Spec.PipelineRef,
				EndpointSpecs: []windtunnelv1alpha1.EndpointSpec{
					{
						EndpointName: "upload",
						DataSpec: &windtunnelv1alpha1.DataSpec{
							DataSetRef: &corev1.ObjectReference{
								Namespace: digitalTwin.Namespace,
								Name:      dataSetName,
							},
						},
						LoadPatternRef: &corev1.ObjectReference{
							Namespace: digitalTwin.Namespace,
							Name:      loadPatternName,
						},
					},
				},
			},
		}
		if err := r.Create(ctx, experiment); client.IgnoreAlreadyExists(err) != nil {
			logger.Error(err, "Failed to create Experiment")
			return ctrl.Result{}, err
		}
		logger.Info("Created Experiment", "name", experimentName, "namespace", digitalTwin.Namespace)
	}

	// Update the DigitalTwin status
	digitalTwin.Status.IsPopulated = true
	if err := r.Status().Update(ctx, digitalTwin); err != nil {
		logger.Error(err, "Failed to update DigitalTwin status")
		return ctrl.Result{}, err
	}
	*/

	return ctrl.Result{}, nil
}

// reconcileCreated reconciles the DigitalTwin when it is created.
func (r *DigitalTwinReconciler) reconcileCreated(ctx context.Context, digitalTwin *windtunnelv1alpha1.DigitalTwin) (ctrl.Result, error) {
	switch digitalTwin.Spec.DigitalTwinType {
	case "regular":
		digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinCompleted

	case "schemaware":
		// Get the DataSet

		digitalTwin.Status.JobStatus = windtunnelv1alpha1.DigitalTwinRunning
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DigitalTwinReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.DigitalTwin{}).
		Complete(r)
}
