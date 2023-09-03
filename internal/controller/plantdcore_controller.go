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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
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
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheuses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=nodes/metrics,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=endpoints,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch
//+kubebuilder:rbac:urls=/metrics,verbs=get

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
	plantdCoreObj := &plantdv1alpha1.PlantDCore{}
	if err := r.Get(ctx, req.NamespacedName, plantdCoreObj); err != nil {
		if errors.IsNotFound(err) {
			// Custom resource not found, perform cleanup tasks here.
			// Created objects are automatically garbage collected.
			log.Info("PlantDCore object not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Deploy Proxy as a Deployment and a Service with a ClusterIP
	proxyDeployment, proxyService := plantdcore.SetupProxyDeployment(plantdCoreObj)

	if err := ctrl.SetControllerReference(plantdCoreObj, proxyDeployment, r.Scheme); err != nil {
		log.Error(err, "error creating proxy deployment")
		return ctrl.Result{}, err
	}

	if err := ctrl.SetControllerReference(plantdCoreObj, proxyService, r.Scheme); err != nil {
		log.Error(err, "error creating proxy deployment")
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, proxyDeployment); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating proxy deployment")
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, proxyService); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating proxy service")
		return ctrl.Result{}, err
	}

	// Deploy Studio as a Deployment and a Service with a LoadBalancer
	studioDeployment, studioService := plantdcore.SetupFrontendDeployment(plantdCoreObj, fmt.Sprintf("http://%s.%s.svc.%s:5000", proxyService.Name, proxyService.Namespace, "cluster.local"))
	if err := ctrl.SetControllerReference(plantdCoreObj, studioDeployment, r.Scheme); err != nil {
		log.Error(err, "error creating studio deployment controller reference")
		return ctrl.Result{}, err
	}

	if err := ctrl.SetControllerReference(plantdCoreObj, studioService, r.Scheme); err != nil {
		log.Error(err, "error creating studio controller reference")
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, studioDeployment); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating studio deployment")
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, studioService); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating studio service")
		return ctrl.Result{}, err
	}

	// Deploy ClusterRole and ClusterRoleBinding for Prometheus
	sa, clusterRole, clusterRoleBinding := plantdcore.SetupRoleBindingsForPrometheus(plantdCoreObj)

	if err := ctrl.SetControllerReference(plantdCoreObj, sa, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, sa); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating prometheus service account")
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, clusterRole); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating prometheus cluster role")
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, clusterRoleBinding); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating prometheus cluster role binding")
		return ctrl.Result{}, err
	}

	// Create Prometheus Object
	prometheus, prometheusService := plantdcore.SetupPrometheusObject(plantdCoreObj)
	if err := ctrl.SetControllerReference(plantdCoreObj, prometheus, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := ctrl.SetControllerReference(plantdCoreObj, prometheusService, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, prometheus); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating prometheus resource")
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, prometheusService); err != nil && !errors.IsAlreadyExists(err) {
		log.Error(err, "error creating prometheus service")
		return ctrl.Result{}, err
	}

	//TODO: Deploy RedisStack

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlantDCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&plantdv1alpha1.PlantDCore{}).
		Complete(r)
}
