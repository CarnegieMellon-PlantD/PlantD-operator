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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/simulation"
)

// SimulationReconciler reconciles a Simulation object
type SimulationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	SimulationPending  string = "Pending"  // Simulation object has been created
	SimulationRunning  string = "Running"  // Simulation is running
	SimulationFinished string = "Finished" // Simulation has finished running;
	SimulationFailed   string = "Failed"
)

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=simulations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=simulations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=simulations/finalizers,verbs=update
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=loadpatterns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=trafficmodels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=digitaltwins,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Modify the Reconcile function to compare the state specified by
// the Simulation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *SimulationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the DigitalTwin resource
	sim := &windtunnelv1alpha1.Simulation{}
	if err := r.Get(ctx, req.NamespacedName, sim); err != nil {
		if errors.IsNotFound(err) {
			// Custom resource not found, perform cleanup tasks here.
			// Created objects are automatically garbage collected.
			log.Info("Simulation resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	digitalTwin := &windtunnelv1alpha1.DigitalTwin{}
	digitalTwinName := types.NamespacedName{Namespace: sim.Spec.DigitalTwinRef.Namespace, Name: sim.Spec.DigitalTwinRef.Name}
	if err := r.Get(ctx, digitalTwinName, digitalTwin); err != nil {
		sim.Status.DigitalTwinState = "Error: No Digital Twin resource found"
		log.Error(err, "Cannot get Digital Twin: "+digitalTwinName.String())
		if err := r.Status().Update(ctx, sim); err != nil {
			log.Error(err, "Cannot update the status of Simulation: "+req.NamespacedName.String())
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	trafficModel := &windtunnelv1alpha1.TrafficModel{}
	trafficModelName := types.NamespacedName{Namespace: sim.Spec.TrafficModelRef.Namespace, Name: sim.Spec.TrafficModelRef.Name}
	if err := r.Get(ctx, trafficModelName, trafficModel); err != nil {
		sim.Status.TrafficModelState = "Error: No Traffic Model resource found"
		log.Error(err, "Cannot get Traffic Model: "+trafficModelName.String())
		if err := r.Status().Update(ctx, sim); err != nil {
			log.Error(err, "Cannot update the status of Simulation: "+req.NamespacedName.String())
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// List Experiments and LoadPatterns required for digital twin
	var experimentNames []string
	var loadPatternNames []string
	experimentList := &windtunnelv1alpha1.ExperimentList{}
	loadPatternList := &windtunnelv1alpha1.LoadPatternList{}
	for _, experimentRef := range digitalTwin.Spec.Experiments {
		experiment := &windtunnelv1alpha1.Experiment{}
		if err := r.Get(ctx, types.NamespacedName{Namespace: experimentRef.Namespace, Name: experimentRef.Name}, experiment); err != nil {
			return ctrl.Result{}, err
		}
		experimentNames = append(experimentNames, fmt.Sprintf("%s.%s", experiment.Namespace, experiment.Name))
		experimentList.Items = append(experimentList.Items, *experiment)

		for _, loadPatternConfig := range experiment.Spec.LoadPatterns {
			loadPattern := &windtunnelv1alpha1.LoadPattern{}
			if err := r.Get(ctx, types.NamespacedName{
				Namespace: loadPatternConfig.LoadPatternRef.Namespace,
				Name:      loadPatternConfig.LoadPatternRef.Name,
			}, loadPattern); err != nil {
				return ctrl.Result{}, err
			}
			loadPatternNames = append(loadPatternNames, fmt.Sprintf("%s.%s", loadPattern.Namespace, loadPattern.Name))
			loadPatternList.Items = append(loadPatternList.Items, *loadPattern)
		}
	}

	experimentListJSON, err := json.Marshal(experimentList)
	if err != nil {
		log.Error(err, "Unable to marshal experimentList to JSON")
		return ctrl.Result{}, err
	}

	loadPatternListJSON, err := json.Marshal(loadPatternList)
	if err != nil {
		log.Error(err, "Unable to marshal loadPatternList to JSON")
		return ctrl.Result{}, err
	}

	if sim.Status.PodName == "" {
		pod, _ := simulation.CreateJobBySimulation(sim.Name+"-"+strconv.FormatInt(time.Now().Unix(), 10), sim,
			digitalTwin, trafficModel, experimentNames, loadPatternNames, string(experimentListJSON), string(loadPatternListJSON))
		if err := ctrl.SetControllerReference(sim, pod, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, pod); err != nil {
			log.Error(err, "Cannot create simulation job.")
		}

		sim.Status.PodName = pod.Name
		sim.Status.JobStatus = SimulationPending
		if err := r.Status().Update(ctx, sim); err != nil {
			log.Error(err, "Cannot update the status of Simulation after creation.")
		}
	} else {
		log.Info("checking if pod exists")
		// Pod name exists, fetch the Pod and check its status
		pod := &corev1.Pod{}
		if err := r.Get(ctx, types.NamespacedName{Namespace: sim.Namespace, Name: sim.Status.PodName}, pod); err != nil {
			log.Error(err, "Failed to get the Pod")
			return ctrl.Result{}, err
		}

		// Check the Pod's status
		switch pod.Status.Phase {
		case corev1.PodSucceeded:
			log.Info("Pod has succeeded")

			if err := r.Delete(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{
				Namespace: sim.Namespace,
				Name:      sim.Status.PodName,
			}}); err != nil {
				log.Error(err, "Failed to delete the cost service - "+sim.Status.PodName)
				return ctrl.Result{}, err
			}

			sim.Status.JobStatus = SimulationFinished
			if err := r.Status().Update(ctx, sim); err != nil {
				log.Error(err, "Cannot update the status of Cost Service after pod succeeded.")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		case corev1.PodFailed:
			log.Info("Pod has failed")

			sim.Status.JobStatus = SimulationFailed
			if err := r.Status().Update(ctx, sim); err != nil {
				log.Error(err, "Cannot update the status of Cost Service after pod failed.")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		default:
			sim.Status.JobStatus = SimulationRunning
			if err := r.Status().Update(ctx, sim); err != nil {
				log.Error(err, "Cannot update the status of Cost Service when pod is running.")
				return ctrl.Result{}, err
			}
		}
	}
	// Requeue to re-run the reconiler. If it reaches here, it means the pod is still running
	return ctrl.Result{Requeue: true, RequeueAfter: time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SimulationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Simulation{}).
		Complete(r)
}
