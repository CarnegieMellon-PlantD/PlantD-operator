/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"math"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
)

// ScenarioReconciler reconciles a Scenario object
type ScenarioReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=scenarios,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=scenarios/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=scenarios/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scenario object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *ScenarioReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the requested Scenario
	scenario := &windtunnelv1alpha1.Scenario{}
	if err := r.Get(ctx, req.NamespacedName, scenario); err != nil {
		// We'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Skip if the Scenario is already populated
	if scenario.Status.IsPopulated {
		return ctrl.Result{}, nil
	}

	// Figure out max rate the Pipeline will need to handle
	var maxRate float64
	for _, task := range scenario.Spec.Tasks {
		taskMaxRate := float64(task.SendingDevices["max"]*task.PushFrequencyPerMonth["max"]) / (30 * 24 * 3600)
		maxRate = max(maxRate, taskMaxRate)
	}
	// Go a bit higher than what the user requested
	maxRate *= 1.5

	EXPERIMENT_DURATION := 600
	DATASET_SIZE := int(math.Ceil(float64(EXPERIMENT_DURATION) * maxRate / 2))

	for _, task := range scenario.Spec.Tasks {
		dataSetName := fmt.Sprintf("%s-dataset-pure-%s", scenario.Name, task.Name)
		loadPatternName := fmt.Sprintf("%s-loadpattern-%s", scenario.Name, task.Name)
		experimentName := fmt.Sprintf("%s-experiment-pure-%s", scenario.Name, task.Name)

		// Create a DataSet
		dataSet := &windtunnelv1alpha1.DataSet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: scenario.Namespace,
				Name:      dataSetName,
			},
			Spec: windtunnelv1alpha1.DataSetSpec{
				CompressPerSchema:    scenario.Spec.DataSetConfig.CompressPerSchema,
				CompressedFileFormat: scenario.Spec.DataSetConfig.CompressedFileFormat,
				FileFormat:           scenario.Spec.DataSetConfig.FileFormat,
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
		logger.Info("Created DataSet", "name", dataSetName, "namespace", scenario.Namespace)

		// Create LoadPattern
		loadPattern := &windtunnelv1alpha1.LoadPattern{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: scenario.Namespace,
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
		logger.Info("Created LoadPattern", "name", loadPatternName, "namespace", scenario.Namespace)

		// Create Experiment
		experiment := &windtunnelv1alpha1.Experiment{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: scenario.Namespace,
				Name:      experimentName,
			},
			Spec: windtunnelv1alpha1.ExperimentSpec{
				PipelineRef: scenario.Spec.PipelineRef,
				EndpointSpecs: []windtunnelv1alpha1.EndpointSpec{
					{
						EndpointName: "upload",
						DataSpec: windtunnelv1alpha1.DataSpec{
							DataSetRef: corev1.ObjectReference{
								Namespace: scenario.Namespace,
								Name:      dataSetName,
							},
						},
						LoadPatternRef: corev1.ObjectReference{
							Namespace: scenario.Namespace,
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
		logger.Info("Created Experiment", "name", experimentName, "namespace", scenario.Namespace)
	}

	// Update the Scenario status
	scenario.Status.IsPopulated = true
	if err := r.Status().Update(ctx, scenario); err != nil {
		logger.Error(err, "Failed to update Scenario status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScenarioReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Scenario{}).
		Complete(r)
}
