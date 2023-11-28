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
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/digitaltwin"
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
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DigitalTwin object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *DigitalTwinReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the DigitalTwin resource
	digitalTwin := &windtunnelv1alpha1.DigitalTwin{}
	if err := r.Get(ctx, req.NamespacedName, digitalTwin); err != nil {
		if errors.IsNotFound(err) {
			// Custom resource not found, perform cleanup tasks here.
			// Created objects are automatically garbage collected.
			log.Info("Digital Twun resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	experiments := &windtunnelv1alpha1.ExperimentList{}

	err := r.List(ctx, experiments)
	if err != nil {
		log.Error(err, "Unable to list experiments when running digital twin")
		return ctrl.Result{}, err
	}

	// List experiments required for digital twin
	experimentsList := &windtunnelv1alpha1.ExperimentList{}
	for _, experiment := range experiments.Items {
		experimentNameCheck := experiment.Namespace + "." + experiment.Name
		if strings.Contains(digitalTwin.Spec.ExperimentNames, experimentNameCheck) {
			experimentsList.Items = append(experimentsList.Items, experiment)
		}
	}

	loadPatterns := &windtunnelv1alpha1.LoadPatternList{}

	err1 := r.List(ctx, loadPatterns)
	if err1 != nil {
		log.Error(err, "Unable to list load patterns when running digital twin")
		return ctrl.Result{}, err
	}

	// List experiments required for digital twin
	loadPatternList := &windtunnelv1alpha1.LoadPatternList{}
	for _, loadPattern := range loadPatterns.Items {
		loadPatternNameCheck := loadPattern.Namespace + "." + loadPattern.Name
		if strings.Contains(digitalTwin.Spec.LoadPatternNames, loadPatternNameCheck) {
			loadPatternList.Items = append(loadPatternList.Items, loadPattern)
		}
	}
	experimentListJSON, err := json.Marshal(experimentsList)
	if err != nil {
		log.Error(err, "Unable to marshal experimentList to JSON")
		return ctrl.Result{}, err
	}

	loadPatternListJSON, err := json.Marshal(loadPatternList)
	if err != nil {
		log.Error(err, "Unable to marshal loadPatternList to JSON")
		return ctrl.Result{}, err
	}

	pod, _ := digitaltwin.CreateJobByDigitalTwin(ctx, digitalTwin.Name+"-"+strconv.FormatInt(time.Now().Unix(), 10), digitalTwin, string(experimentListJSON), string(loadPatternListJSON))

	if err := ctrl.SetControllerReference(digitalTwin, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, pod); err != nil {
		log.Error(err, "Cannot create digital twin job.")
	}

	if err := r.Status().Update(ctx, digitalTwin); err != nil {
		log.Error(err, "Cannot update the status of Digital Twin after creation.")
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DigitalTwinReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.DigitalTwin{}).
		Complete(r)
}
