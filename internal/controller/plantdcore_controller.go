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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	plantdcore "github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/core"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

const (
	ModuleCreated string = "Created"
	ModuleRunning string = "Running"
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
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch
//+kubebuilder:rbac:urls=/metrics,verbs=get

// reconcileProxy reconciles the proxy module. It tries to create all necessary resources and return the status
// of the proxy module and the error if any.
func (r *PlantDCoreReconciler) reconcileProxy(ctx context.Context, plantDCore *windtunnelv1alpha1.PlantDCore) (string, error) {
	logger := log.FromContext(ctx)

	hasCreation := false

	// Prepare resources
	newProxyDeploy, newProxySvc := plantdcore.GetProxyResources(plantDCore)
	if err := ctrl.SetControllerReference(plantDCore, newProxyDeploy, r.Scheme); err != nil {
		logger.Error(err, "failed to set controller reference for proxy Deployment")
		return "", err
	}
	if err := ctrl.SetControllerReference(plantDCore, newProxySvc, r.Scheme); err != nil {
		logger.Error(err, "failed to set controller reference for proxy Service")
		return "", err
	}

	// Find proxy Deployment and create if not exist
	curProxyDeploy := &appsv1.Deployment{}
	proxyDeployKey := client.ObjectKey{
		Name:      newProxyDeploy.Name,
		Namespace: newProxyDeploy.Namespace,
	}
	if err := r.Get(ctx, proxyDeployKey, curProxyDeploy); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of proxy Deployment")
		return "", err
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newProxyDeploy); err != nil {
			logger.Error(err, "failed to create proxy Deployment")
			return "", err
		}
		logger.Info("created proxy Deployment")
		hasCreation = true
	}

	// Find proxy Service and create if not exist
	curProxySvc := &corev1.Service{}
	proxySvcKey := client.ObjectKey{
		Name:      newProxySvc.Name,
		Namespace: newProxySvc.Namespace,
	}
	if err := r.Get(ctx, proxySvcKey, curProxySvc); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of proxy Service")
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newProxySvc); err != nil && !errors.IsAlreadyExists(err) {
			logger.Error(err, "failed to create proxy Service")
			return "", err
		}
		logger.Info("created proxy Service")
		hasCreation = true
	}

	if hasCreation {
		return ModuleCreated, nil
	}
	return fmt.Sprintf("%s (%d/%d)", ModuleRunning, curProxyDeploy.Status.AvailableReplicas, curProxyDeploy.Status.AvailableReplicas+curProxyDeploy.Status.UnavailableReplicas), nil
}

// reconcileStudio reconciles the studio module. It tries to create all necessary resources and return the status of the
// studio module and the error if any.
func (r *PlantDCoreReconciler) reconcileStudio(ctx context.Context, plantDCore *windtunnelv1alpha1.PlantDCore) (string, error) {
	logger := log.FromContext(ctx)

	hasCreation := false

	// Prepare resources
	_, newProxySvc := plantdcore.GetProxyResources(plantDCore)
	proxySvcFQDN := fmt.Sprintf("http://%s.%s.svc.%s:5000", newProxySvc.Name, newProxySvc.Namespace, "cluster.local")
	newStudioDeploy, newStudioSvc := plantdcore.GetStudioResources(plantDCore, proxySvcFQDN)
	if err := ctrl.SetControllerReference(plantDCore, newStudioDeploy, r.Scheme); err != nil {
		logger.Error(err, "failed to set controller reference for studio Deployment")
		return "", err
	}
	if err := ctrl.SetControllerReference(plantDCore, newStudioSvc, r.Scheme); err != nil {
		logger.Error(err, "failed to set controller reference for studio Service")
		return "", err
	}

	// Find studio Deployment and create if not exist
	curStudioDeploy := &appsv1.Deployment{}
	studioDeployKey := client.ObjectKey{
		Name:      newStudioDeploy.Name,
		Namespace: newStudioDeploy.Namespace,
	}
	if err := r.Get(ctx, studioDeployKey, curStudioDeploy); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of studio Deployment")
		return "", err
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newStudioDeploy); err != nil {
			logger.Error(err, "failed to create studio Deployment")
			return "", err
		}
		logger.Info("created studio Deployment")
		hasCreation = true
	}

	// Find studio Service and create if not exist
	curStudioSvc := &corev1.Service{}
	studioSvcKey := client.ObjectKey{
		Name:      newStudioSvc.Name,
		Namespace: newStudioSvc.Namespace,
	}
	if err := r.Get(ctx, studioSvcKey, curStudioSvc); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of studio Service")
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newStudioSvc); err != nil {
			logger.Error(err, "failed to create studio Service")
			return "", err
		}
		logger.Info("created studio Service")
		hasCreation = true
	}

	if hasCreation {
		return ModuleCreated, nil
	}
	return fmt.Sprintf("%s (%d/%d)", ModuleRunning, curStudioDeploy.Status.AvailableReplicas, curStudioDeploy.Status.AvailableReplicas+curStudioDeploy.Status.UnavailableReplicas), nil
}

// reconcilePrometheus reconciles the Prometheus module. It tried to create all necessary resources and return the
// status of the Prometheus module and the error if any.
func (r *PlantDCoreReconciler) reconcilePrometheus(ctx context.Context, plantDCore *windtunnelv1alpha1.PlantDCore) (string, error) {
	logger := log.FromContext(ctx)

	hasCreation := false

	// Prepare resources
	newPromSA, newPromCR, newPromCRB := plantdcore.GetPrometheusRoleBindings(plantDCore)
	newProm, newPromSvc := plantdcore.GetPrometheusResources(plantDCore)
	if err := ctrl.SetControllerReference(plantDCore, newPromSA, r.Scheme); err != nil {
		logger.Error(err, "failed to set controller reference for prometheus ServiceAccount")
		return "", err
	}
	if err := ctrl.SetControllerReference(plantDCore, newProm, r.Scheme); err != nil {
		logger.Error(err, "failed to set controller reference for Prometheus")
		return "", err
	}
	if err := ctrl.SetControllerReference(plantDCore, newPromSvc, r.Scheme); err != nil {
		logger.Error(err, "failed to set controller reference for prometheus Service")
		return "", err
	}

	// Find prometheus ServiceAccount and create if not exist
	curPromSA := &corev1.ServiceAccount{}
	promSAKey := client.ObjectKey{
		Name:      newPromSA.Name,
		Namespace: newPromSA.Namespace,
	}
	if err := r.Get(ctx, promSAKey, curPromSA); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of prometheus ServiceAccount")
		return "", err
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newPromSA); err != nil {
			logger.Error(err, "failed to create prometheus ServiceAccount")
			return "", err
		}
		logger.Info("created prometheus ServiceAccount")
		hasCreation = true
	}

	// Find prometheus ClusterRole and create if not exist
	curPromCR := &rbacv1.ClusterRole{}
	promCRKey := client.ObjectKey{
		Name: newPromCR.Name,
	}
	if err := r.Get(ctx, promCRKey, curPromCR); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of prometheus ClusterRole")
		return "", err
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newPromCR); err != nil {
			logger.Error(err, "failed to create prometheus ClusterRole")
			return "", err
		}
		logger.Info("created prometheus ClusterRole")
		hasCreation = true
	}

	// Find prometheus ClusterRoleBinding and create if not exist
	curPromCRB := &rbacv1.ClusterRoleBinding{}
	promCRBKey := client.ObjectKey{
		Name: newPromCRB.Name,
	}
	if err := r.Get(ctx, promCRBKey, curPromCRB); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of prometheus ClusterRoleBinding")
		return "", err
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newPromCRB); err != nil && !errors.IsAlreadyExists(err) {
			logger.Error(err, "failed to create prometheus ClusterRoleBinding")
			return "", err
		}
		logger.Info("created prometheus ClusterRoleBinding")
		hasCreation = true
	}

	// Find Prometheus and create if not exist
	curProm := &monitoringv1.Prometheus{}
	promKey := client.ObjectKey{
		Name:      newProm.Name,
		Namespace: newProm.Namespace,
	}
	if err := r.Get(ctx, promKey, curProm); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of Prometheus")
		return "", err
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newProm); err != nil {
			logger.Error(err, "failed to create Prometheus")
			return "", err
		}
		logger.Info("created Prometheus")
		hasCreation = true
	}

	// Find prometheus Service and create if not exist
	curPromSvc := &corev1.Service{}
	promSvcKey := client.ObjectKey{
		Name:      newPromSvc.Name,
		Namespace: newPromSvc.Namespace,
	}
	if err := r.Get(ctx, promSvcKey, curPromSvc); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to check existence of prometheus Service")
		return "", err
	} else if errors.IsNotFound(err) {
		if err := r.Create(ctx, newPromSvc); err != nil {
			logger.Error(err, "failed to create prometheus Service")
			return "", err
		}
		logger.Info("created prometheus Service")
		hasCreation = true
	}

	if hasCreation {
		return ModuleCreated, nil
	}
	return fmt.Sprintf("%s (%d/%d)", ModuleRunning, curProm.Status.AvailableReplicas, curProm.Status.AvailableReplicas+curProm.Status.UnavailableReplicas), nil
}

// finalizePrometheus cleans up the Prometheus resources that cannot be deleted automatically.
func (r *PlantDCoreReconciler) finalizePrometheus(ctx context.Context, plantDCore *windtunnelv1alpha1.PlantDCore) error {
	logger := log.FromContext(ctx)

	// Prepare resources
	_, clusterRole, clusterRoleBinding := plantdcore.GetPrometheusRoleBindings(plantDCore)

	if err := r.Delete(ctx, clusterRole); err != nil {
		logger.Error(err, "failed to delete prometheus ClusterRole")
		return err
	}
	logger.Info("deleted prometheus ClusterRole")

	if err := r.Delete(ctx, clusterRoleBinding); err != nil {
		logger.Error(err, "failed to delete prometheus ClusterRoleBinding")
		return err
	}
	logger.Info("deleted prometheus ClusterRoleBinding")

	return nil
}

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
	logger := log.FromContext(ctx)

	// Fetch the PlantD resource
	plantDCore := &windtunnelv1alpha1.PlantDCore{}
	if err := r.Get(ctx, req.NamespacedName, plantDCore); err != nil {
		// Ignore not-found errors as we can get them on delete requests
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Finalizers
	prometheusFinalizerName := "plantdcore.windtunnel.plantd.org/prometheusFinalizer"
	if plantDCore.ObjectMeta.DeletionTimestamp.IsZero() {
		// Object is not being deleted, add finalizer if not present
		if !controllerutil.ContainsFinalizer(plantDCore, prometheusFinalizerName) {
			controllerutil.AddFinalizer(plantDCore, prometheusFinalizerName)
			if err := r.Update(ctx, plantDCore); err != nil {
				logger.Error(err, "failed to add prometheus finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		// Object is being deleted
		if controllerutil.ContainsFinalizer(plantDCore, prometheusFinalizerName) {
			// Finalizer presents
			if err := r.finalizePrometheus(ctx, plantDCore); err != nil {
				// Retry on errors
				return ctrl.Result{}, err
			}

			// Remove finalizer from the list and update it.
			controllerutil.RemoveFinalizer(plantDCore, prometheusFinalizerName)
			if err := r.Update(ctx, plantDCore); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// Reconcile proxy
	statusKubeProxy, err := r.reconcileProxy(ctx, plantDCore)
	if err != nil {
		logger.Error(err, "failed to reconcile proxy")
		return ctrl.Result{}, err
	}
	plantDCore.Status.ProxyStatus = statusKubeProxy

	// Reconcile studio
	statusStudio, err := r.reconcileStudio(ctx, plantDCore)
	if err != nil {
		logger.Error(err, "failed to reconcile studio")
		return ctrl.Result{}, err
	}
	plantDCore.Status.StudioStatus = statusStudio

	// Reconcile prometheus
	statusPrometheus, err := r.reconcilePrometheus(ctx, plantDCore)
	if err != nil {
		logger.Error(err, "failed to reconcile prometheus")
		return ctrl.Result{}, err
	}
	plantDCore.Status.PrometheusStatus = statusPrometheus

	// Update status
	if err := r.Status().Update(ctx, plantDCore); err != nil {
		logger.Error(err, "failed to update PlantDCore status")
		return ctrl.Result{}, err
	}

	// TODO: Deploy RedisStack

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlantDCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.PlantDCore{}).
		Complete(r)
}
