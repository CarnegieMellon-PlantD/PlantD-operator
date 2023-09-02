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

	plantdv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	plantdcore "github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/core"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PlantDCoreReconciler reconciles a PlantDCore object
type PlantDCoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=plantdcores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=plantdcores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=plantdcores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PlantDCore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *PlantDCoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the PlantD resource
	plantdCore := &plantdv1alpha1.PlantDCore{}
	if err := r.Get(ctx, req.NamespacedName, plantdCore); err != nil {
		if errors.IsNotFound(err) {
			// Custom resource not found, perform cleanup tasks here.
			// Created objects are automatically garbage collected.
			log.Info("CostService resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Deploy Proxy as a Deployment and a Service with a ClusterIP
	proxyDeployment, proxyService := plantdcore.SetupProxyDeployment(plantdCore)

	r.CreateObject(ctx, r.Client, proxyDeployment)
	r.CreateObject(ctx, r.Client, proxyService)

	// Deploy Studio as a Deployment and a Service with a LoadBalancer
	studioDeployment, studioService := plantdcore.SetupFrontendDeployment(plantdCore, fmt.Sprintf("%s.%s.svc.%s", proxyService.Name, proxyService.Namespace, "cluster.local"))
	r.CreateObject(ctx, r.Client, studioDeployment)
	r.CreateObject(ctx, r.Client, studioService)

	// Check if the Prometheus CRD exists before deploying the prometheus service monitor object
	prometheus := &monitoringv1.Prometheus{}
	if err := r.Get(ctx, types.NamespacedName{Name: "prometheus", Namespace: plantdCore.Namespace}, prometheus); err != nil {
		if errors.IsNotFound(err) {
			// Handle successful creation
			return ctrl.Result{}, nil
		}
	}

	prometheus, prometheusService := plantdcore.SetupCreatePrometheusSMObject(plantdCore)
	r.CreateObject(ctx, r.Client, prometheus)
	r.CreateObject(ctx, r.Client, prometheusService)

	// Deploy ClusterRole and ClusterRoleBinding for Prometheus
	clusterRole, clusterRoleBinding := plantdcore.SetupRoleBindingsForPrometheus(plantdCore)
	r.CreateObject(ctx, r.Client, clusterRole)
	r.CreateObject(ctx, r.Client, clusterRoleBinding)

	// Deploy RedisStack

	return ctrl.Result{}, nil
}

// CreateObject creates a new object of the provided kind.
func (r *PlantDCoreReconciler) CreateObject(ctx context.Context, c client.Client, newObj client.Object) (objectExists, creationFailed error) {

	if err := ctrl.SetControllerReference(&plantdv1alpha1.PlantDCore{}, newObj, r.Scheme); err != nil {
		return nil, err
	}

	key := types.NamespacedName{Name: newObj.GetName(), Namespace: newObj.GetNamespace()}

	err := c.Get(ctx, key, newObj)
	if err == nil {
		return nil, err
	}

	if err := c.Create(ctx, newObj); err != nil {
		return nil, err
	}
	return nil, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlantDCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&plantdv1alpha1.PlantDCore{}).
		Complete(r)
}
