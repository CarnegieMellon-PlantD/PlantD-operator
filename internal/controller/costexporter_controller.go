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
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/cost"
)

type Tag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type ExperimentTags struct {
	Name string `json:"name,omitempty"`
	Tags []Tag  `json:"tags,omitempty"`
}

// CostExporterReconciler reconciles a CostExporter object
type CostExporterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=costexporters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=costexporters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=costexporters/finalizers,verbs=update
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *CostExporterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the CostExporter resource
	costExporter := &windtunnelv1alpha1.CostExporter{}
	if err := r.Get(ctx, req.NamespacedName, costExporter); err != nil {
		if errors.IsNotFound(err) {
			// Custom resource not found, perform cleanup tasks here.
			// Created objects are automatically garbage collected.
			log.Info("CostExporter resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	experiments := &windtunnelv1alpha1.ExperimentList{}

	err := r.List(ctx, experiments)
	if err != nil {
		log.Error(err, "Unable to list experiments with cloud vendor "+costExporter.Spec.CloudServiceProvider)
		return ctrl.Result{}, err
	}

	earliestTime := time.Now()
	// Filter experiments with cloudVendor "aws"
	experimentsList := &windtunnelv1alpha1.ExperimentList{}
	for _, experiment := range experiments.Items {
		if experiment.Status.CloudVendor == costExporter.Spec.CloudServiceProvider {

			var experimentTime time.Time
			if !experiment.CreationTimestamp.IsZero() {
				experimentTime = experiment.CreationTimestamp.Time
			}
			if !experiment.Spec.ScheduledTime.IsZero() && experiment.Spec.ScheduledTime.Time.Before(experimentTime) {
				experimentTime = experiment.Spec.ScheduledTime.Time
			}

			// Check if this experiment's time is earlier than the current earliest time
			if experimentTime.Before(earliestTime) {
				earliestTime = experimentTime
			}
			experimentsList.Items = append(experimentsList.Items, experiment)
		}
	}

	resultList := []ExperimentTags{}
	for _, experiment := range experimentsList.Items {
		// Extract the Name from the experiment's metadata
		experimentName := experiment.Namespace + "/" + experiment.ObjectMeta.Name
		// Extract the Tags from the experiment's status field
		tags := experiment.Status.Tags

		// Create a list of Tag pairs for the current experiment's tags
		var tagsList []Tag
		for key, value := range tags {
			tagsList = append(tagsList, Tag{Key: key, Value: value})
		}
		// Create an ExperimentTags object for the current experiment
		experimentTags := ExperimentTags{
			Name: experimentName,
			Tags: tagsList,
		}
		// Append the ExperimentTags object to the resultList
		resultList = append(resultList, experimentTags)
	}

	resultListJSON, err := json.Marshal(resultList)
	if err != nil {
		log.Error(err, "Unable to marshal resultList to JSON")
		return ctrl.Result{}, err
	}

	// Convert the JSON byte slice to a string
	resultListStr := string(resultListJSON)
	costExporter.Status.Tags = resultListStr

	// update to include tags for the new pipelines; only if new pod is being created
	if err := r.Status().Update(ctx, costExporter); err != nil {
		log.Error(err, "Cannot update the status of Cost Service after creation.")
	}

	// TODO(user): your logic here
	currTime := time.Now()
	if costExporter.Status.JobCompletionTime != nil && (currTime.Sub(costExporter.Status.JobCompletionTime.Time) < (8 * time.Hour)) {
		return ctrl.Result{RequeueAfter: (8 * time.Hour) - currTime.Sub(costExporter.Status.JobCompletionTime.Time)}, nil
	}

	if costExporter.Status.PodName == "" {

		pod, _ := cost.CreateJobByCostServie(ctx, costExporter.Name+"-"+strconv.FormatInt(time.Now().Unix(), 10), costExporter, earliestTime)

		// Extract the Pod's name from the created Pod object
		costExporter.Status.PodName = pod.Name
		if err := ctrl.SetControllerReference(costExporter, pod, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, pod); err != nil {
			log.Error(err, "Cannot create cost job.")
		}

		if err := r.Status().Update(ctx, costExporter); err != nil {
			log.Error(err, "Cannot update the status of Cost Service after creation.")
		}
	} else {
		log.Info("checking if pod exists")
		// Pod name exists, fetch the Pod and check its status
		pod := &corev1.Pod{}
		if err := r.Get(ctx, types.NamespacedName{Namespace: costExporter.Namespace, Name: costExporter.Status.PodName}, pod); err != nil {
			log.Error(err, "Failed to get the Pod")
			return ctrl.Result{}, err
		}

		// Check the Pod's status
		switch pod.Status.Phase {
		case corev1.PodSucceeded:
			log.Info("Pod has succeeded")
			costExporter.Status.JobCompletionTime = &metav1.Time{Time: time.Now()}
			costExporter.Status.JobStatus = SUCCESS

			if err := r.Delete(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{
				Namespace: costExporter.Namespace,
				Name:      costExporter.Status.PodName,
			}}); err != nil {
				log.Error(err, "Failed to delete the cost service - "+costExporter.Status.PodName)
			}

			// CostExporter.Status.PodName = ""
			if err := r.Status().Update(ctx, costExporter); err != nil {
				log.Error(err, "Cannot update the status of Cost Service after pod success.")
			}

			// The Pod has succeeded
			// You can handle this case here
		case corev1.PodFailed:
			costExporter.Status.JobStatus = FAILED
			log.Info("Pod has failed")
			// The Pod has failed
			// You can handle this case here
		default:
			costExporter.Status.JobStatus = RUNNING
			// The Pod is still running or in an unknown state
			// You can handle this case here
		}
	}
	log.Info("Pod is still running")
	if err := r.Status().Update(ctx, costExporter); err != nil { // pod status update
		log.Error(err, "Cannot update the status of Cost Service after creation.")
	}
	// Requeue to re-run the reconiler. If it reaches here, it means the pod is still running
	return ctrl.Result{Requeue: true, RequeueAfter: time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CostExporterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.CostExporter{}).
		Complete(r)
}
