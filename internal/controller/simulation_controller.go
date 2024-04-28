package controller

import (
	"context"
	"fmt"
	"time"

	kbatch "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/digitaltwin"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

const (
	simulationPollingInterval = 5 * time.Second
)

// SimulationReconciler reconciles a Simulation object
type SimulationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=simulations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=simulations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=simulations/finalizers,verbs=update
//
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=datasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=loadpatterns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=digitaltwins,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=simulations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=trafficmodels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=netcosts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=scenarios,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *SimulationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the requested Simulation
	simulation := &windtunnelv1alpha1.Simulation{}
	if err := r.Get(ctx, req.NamespacedName, simulation); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unable to fetch Simulation")
		return ctrl.Result{}, err
	}

	if simulation.Status.JobStatus == "" {
		result, err := r.reconcileCreated(ctx, simulation)
		if err == nil {
			if err := r.Status().Update(ctx, simulation); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
		}
		return result, err
	}

	if simulation.Status.JobStatus == windtunnelv1alpha1.SimulationRunning {
		result, err := r.reconcileRunning(ctx, simulation)
		if err == nil {
			if err := r.Status().Update(ctx, simulation); err != nil {
				logger.Error(err, "Cannot update the status")
				return ctrl.Result{}, err
			}
		}
		return result, err
	}

	return ctrl.Result{}, nil
}

// reconcileCreated reconciles the DigitalTwin when it is created.
func (r *SimulationReconciler) reconcileCreated(ctx context.Context, simulation *windtunnelv1alpha1.Simulation) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get the DigitalTwin
	var digitalTwin *windtunnelv1alpha1.DigitalTwin
	if simulation.Spec.DigitalTwinRef != nil {
		digitalTwin = &windtunnelv1alpha1.DigitalTwin{}
		digitalTwinName := types.NamespacedName{
			Namespace: simulation.Spec.DigitalTwinRef.Namespace,
			Name:      simulation.Spec.DigitalTwinRef.Name,
		}
		if err := r.Get(ctx, digitalTwinName, digitalTwin); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get DigitalTwin \"%s\"", digitalTwinName))
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot find DigitalTwin \"%s\": %s", digitalTwinName, err)
			return ctrl.Result{}, nil
		}
	}

	// Get the TrafficModel
	trafficModel := &windtunnelv1alpha1.TrafficModel{}
	trafficModelName := types.NamespacedName{
		Namespace: simulation.Spec.TrafficModelRef.Namespace,
		Name:      simulation.Spec.TrafficModelRef.Name,
	}
	if err := r.Get(ctx, trafficModelName, trafficModel); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot get TrafficModel \"%s\"", trafficModelName))
		simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
		simulation.Status.Error = fmt.Sprintf("Cannot find TrafficModel \"%s\": %s", trafficModelName, err)
		return ctrl.Result{}, nil
	}

	var job *kbatch.Job
	if digitalTwin == nil {
		// Get the NetCost
		var netCost *windtunnelv1alpha1.NetCost
		if simulation.Spec.NetCostRef != nil {
			netCost = &windtunnelv1alpha1.NetCost{}
			netCostName := types.NamespacedName{
				Namespace: simulation.Spec.NetCostRef.Namespace,
				Name:      simulation.Spec.NetCostRef.Name,
			}
			if err := r.Get(ctx, netCostName, netCost); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot get NetCost \"%s\"", netCostName))
				simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
				simulation.Status.Error = fmt.Sprintf("Cannot find NetCost \"%s\": %s", netCostName, err)
				return ctrl.Result{}, nil
			}
		}

		// Get the Scenario
		scenario := &windtunnelv1alpha1.Scenario{}
		scenarioName := types.NamespacedName{
			Namespace: simulation.Spec.ScenarioRef.Namespace,
			Name:      simulation.Spec.ScenarioRef.Name,
		}
		if err := r.Get(ctx, scenarioName, scenario); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get Scenario \"%s\"", scenarioName))
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot find Scenario \"%s\": %s", scenarioName, err)
			return ctrl.Result{}, nil
		}

		// Create the Job
		var err error
		job, err = digitaltwin.CreateSimulationJob(simulation, nil,
			trafficModel, netCost, scenario,
			nil, nil, nil, nil,
		)
		if err != nil {
			logger.Error(err, "Cannot create manifest for Job")
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot create manifest for Job: %s", err)
			return ctrl.Result{}, nil
		}
	} else if digitalTwin.Spec.DigitalTwinType == "regular" {
		// Get the Experiments
		experimentList := &windtunnelv1alpha1.ExperimentList{}
		for _, experimentRef := range digitalTwin.Spec.Experiments {
			experiment := &windtunnelv1alpha1.Experiment{}
			experimentName := types.NamespacedName{
				Namespace: experimentRef.Namespace,
				Name:      experimentRef.Name,
			}
			if err := r.Get(ctx, experimentName, experiment); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot get Experiment \"%s\"", experimentName))
				simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
				simulation.Status.Error = fmt.Sprintf("Cannot find Experiment \"%s\": %s", experimentName, err)
				return ctrl.Result{}, nil
			}
			experimentList.Items = append(experimentList.Items, *experiment)
		}

		// Get the Pipeline
		pipeline := &windtunnelv1alpha1.Pipeline{}
		pipelineName := types.NamespacedName{
			Namespace: experimentList.Items[0].Namespace,
			Name:      experimentList.Items[0].Spec.PipelineRef.Name,
		}
		if err := r.Get(ctx, pipelineName, pipeline); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get Pipeline \"%s\"", pipelineName))
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot find Pipeline \"%s\": %s", pipelineName, err)
			return ctrl.Result{}, nil
		}

		// Check if the Pipeline is the same, and get the DataSets and LoadPatterns
		dataSetList := &windtunnelv1alpha1.DataSetList{}
		loadPatternList := &windtunnelv1alpha1.LoadPatternList{}
		for _, experiment := range experimentList.Items {
			if experiment.Namespace != experimentList.Items[0].Namespace || experiment.Spec.PipelineRef.Name != experimentList.Items[0].Spec.PipelineRef.Name {
				logger.Error(nil, "Experiments have different Pipelines")
				simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
				simulation.Status.Error = "Experiments have different Pipelines"
				return ctrl.Result{}, nil
			}

			for _, endpointSpec := range experiment.Spec.EndpointSpecs {
				if getEndpointSpecDataOption(&endpointSpec) == windtunnelv1alpha1.EndpointDataOptionDataSet {
					dataSet := &windtunnelv1alpha1.DataSet{}
					dataSetName := types.NamespacedName{
						Namespace: experiment.Namespace,
						Name:      endpointSpec.DataSpec.DataSetRef.Name,
					}
					if err := r.Get(ctx, dataSetName, dataSet); err != nil {
						logger.Error(err, fmt.Sprintf("Cannot get DataSet \"%s\"", dataSetName))
						simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
						simulation.Status.Error = fmt.Sprintf("Cannot find DataSet \"%s\": %s", dataSetName, err)
						return ctrl.Result{}, nil
					}
					dataSetList.Items = append(dataSetList.Items, *dataSet)
				}

				loadPattern := &windtunnelv1alpha1.LoadPattern{}
				loadPatternName := types.NamespacedName{
					Namespace: endpointSpec.LoadPatternRef.Namespace,
					Name:      endpointSpec.LoadPatternRef.Name,
				}
				if err := r.Get(ctx, loadPatternName, loadPattern); err != nil {
					logger.Error(err, fmt.Sprintf("Cannot get LoadPattern \"%s\"", loadPatternName))
					simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
					simulation.Status.Error = fmt.Sprintf("Cannot find LoadPattern \"%s\": %s", loadPatternName, err)
					return ctrl.Result{}, nil
				}
				loadPatternList.Items = append(loadPatternList.Items, *loadPattern)
			}
		}

		// Create the Job
		var err error
		job, err = digitaltwin.CreateSimulationJob(simulation, digitalTwin,
			trafficModel, nil, nil,
			pipeline, dataSetList, loadPatternList, experimentList,
		)
		if err != nil {
			logger.Error(err, "Cannot create manifest for Job")
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot create manifest for Job: %s", err)
			return ctrl.Result{}, nil
		}
	} else if digitalTwin.Spec.DigitalTwinType == "schemaaware" {
		// Get the NetCost
		var netCost *windtunnelv1alpha1.NetCost
		if simulation.Spec.NetCostRef != nil {
			netCost = &windtunnelv1alpha1.NetCost{}
			netCostName := types.NamespacedName{
				Namespace: simulation.Spec.NetCostRef.Namespace,
				Name:      simulation.Spec.NetCostRef.Name,
			}
			if err := r.Get(ctx, netCostName, netCost); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot get NetCost \"%s\"", netCostName))
				simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
				simulation.Status.Error = fmt.Sprintf("Cannot find NetCost \"%s\": %s", netCostName, err)
				return ctrl.Result{}, nil
			}
		}

		// Get the Scenario
		scenario := &windtunnelv1alpha1.Scenario{}
		scenarioName := types.NamespacedName{
			Namespace: simulation.Spec.ScenarioRef.Namespace,
			Name:      simulation.Spec.ScenarioRef.Name,
		}
		if err := r.Get(ctx, scenarioName, scenario); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get Scenario \"%s\"", scenarioName))
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot find Scenario \"%s\": %s", scenarioName, err)
			return ctrl.Result{}, nil
		}

		// Get the original DataSet
		originalDataSet := &windtunnelv1alpha1.DataSet{}
		originalDataSetName := types.NamespacedName{
			Namespace: digitalTwin.Namespace,
			Name:      digitalTwin.Spec.DataSet.Name,
		}
		if err := r.Get(ctx, originalDataSetName, originalDataSet); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get original DataSet \"%s\"", originalDataSetName))
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot find original DataSet \"%s\": %s", originalDataSetName, err)
			return ctrl.Result{}, nil
		}

		// Get the Pipeline
		pipeline := &windtunnelv1alpha1.Pipeline{}
		pipelineName := types.NamespacedName{
			Namespace: digitalTwin.Namespace,
			Name:      digitalTwin.Spec.Pipeline.Name,
		}
		if err := r.Get(ctx, pipelineName, pipeline); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get Pipeline \"%s\"", pipelineName))
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot find Pipeline \"%s\": %s", pipelineName, err)
			return ctrl.Result{}, nil
		}

		// Get bias DataSets, LoadPatterns, and Experiments
		dataSetList := &windtunnelv1alpha1.DataSetList{}
		loadPatternList := &windtunnelv1alpha1.LoadPatternList{}
		experimentList := &windtunnelv1alpha1.ExperimentList{}
		for schemaIdx, _ := range originalDataSet.Spec.Schemas {
			dataSet := &windtunnelv1alpha1.DataSet{}
			dataSetName := types.NamespacedName{
				Namespace: digitalTwin.Namespace,
				Name:      utils.GetBiasDataSetName(digitalTwin.Name, schemaIdx),
			}
			if err := r.Get(ctx, dataSetName, dataSet); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot get bias DataSet \"%s\"", dataSetName))
				simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
				simulation.Status.Error = fmt.Sprintf("Cannot find bias DataSet \"%s\": %s", dataSetName, err)
				return ctrl.Result{}, nil
			}
			dataSetList.Items = append(dataSetList.Items, *dataSet)

			experiment := &windtunnelv1alpha1.Experiment{}
			experimentName := types.NamespacedName{
				Namespace: digitalTwin.Namespace,
				Name:      utils.GetBiasExperimentName(digitalTwin.Name, schemaIdx),
			}
			if err := r.Get(ctx, experimentName, experiment); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot get bias Experiment \"%s\"", experimentName))
				simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
				simulation.Status.Error = fmt.Sprintf("Cannot find bias Experiment \"%s\": %s", experimentName, err)
				return ctrl.Result{}, nil
			}
			experimentList.Items = append(experimentList.Items, *experiment)
		}
		loadPattern := &windtunnelv1alpha1.LoadPattern{}
		loadPatternName := types.NamespacedName{
			Namespace: digitalTwin.Namespace,
			Name:      utils.GetBiasLoadPatternName(digitalTwin.Name),
		}
		if err := r.Get(ctx, loadPatternName, loadPattern); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get bias LoadPattern \"%s\"", loadPatternName))
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot find bias LoadPattern \"%s\": %s", loadPatternName, err)
			return ctrl.Result{}, nil
		}
		loadPatternList.Items = append(loadPatternList.Items, *loadPattern)

		// Create the Job
		var err error
		job, err = digitaltwin.CreateSimulationJob(simulation, digitalTwin,
			trafficModel, netCost, scenario,
			pipeline, dataSetList, loadPatternList, experimentList,
		)
		if err != nil {
			logger.Error(err, "Cannot create manifest for Job")
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Cannot create manifest for Job: %s", err)
			return ctrl.Result{}, nil
		}
	}

	if err := ctrl.SetControllerReference(simulation, job, r.Scheme); err != nil {
		logger.Error(err, "Cannot set controller reference for Job")
		simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
		simulation.Status.Error = fmt.Sprintf("Cannot set controller reference for Job: %s", err)
		return ctrl.Result{}, nil
	}
	if err := r.Create(ctx, job); client.IgnoreAlreadyExists(err) != nil {
		logger.Error(err, "Cannot create Job")
		simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
		simulation.Status.Error = fmt.Sprintf("Cannot create Job: %s", err)
		return ctrl.Result{}, nil
	} else if err == nil {
		logger.Info("Created Job")
	}

	simulation.Status.JobStatus = windtunnelv1alpha1.SimulationRunning
	return ctrl.Result{RequeueAfter: simulationPollingInterval}, nil
}

// reconcileRunning reconciles the DigitalTwin when it is running.
func (r *SimulationReconciler) reconcileRunning(ctx context.Context, simulation *windtunnelv1alpha1.Simulation) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get the Job
	job := &kbatch.Job{}
	jobName := types.NamespacedName{
		Namespace: simulation.Namespace,
		Name:      utils.GetSimulationJobName(simulation.Name),
	}
	if err := r.Get(ctx, jobName, job); err != nil {
		logger.Error(err, fmt.Sprintf("Lost Job \"%s\"", jobName))
		simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
		simulation.Status.Error = fmt.Sprintf("Lost Job \"%s\": %s", jobName, err)
		return ctrl.Result{}, nil
	}

	// Check if the Job is completed
	jobFinished, jobConditionType := isJobFinished(job)
	if jobFinished {
		logger.Info(fmt.Sprintf("Job \"%s\" finished", jobName))
		switch jobConditionType {
		case kbatch.JobComplete:
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationCompleted
		case kbatch.JobFailed:
			simulation.Status.JobStatus = windtunnelv1alpha1.SimulationFailed
			simulation.Status.Error = fmt.Sprintf("Job \"%s\" failed", jobName)
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: simulationPollingInterval}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SimulationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Simulation{}).
		Complete(r)
}
