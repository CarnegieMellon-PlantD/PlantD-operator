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

	"github.com/cisco-open/k8s-objectmatcher/patch"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/core"
)

const (
	finalizerName = "plantdcore.windtunnel.plantd.org/finalizer"
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

// reconcileObject ensures the actual state of the object matches the desired state. It creates the object if it does
// not exist, and otherwise updates it if necessary. Also, controller reference is set to delete the object when the
// owner object is deleted. Update behavior can be enabled or disabled. It returns a bool to show if any action is taken
// and the error if any.
func (r *PlantDCoreReconciler) reconcileObject(ctx context.Context, plantDCore *windtunnelv1alpha1.PlantDCore, curObj client.Object, desiredObj client.Object, allowUpdate bool) (bool, error) {
	// Get current object
	err := r.Get(ctx, types.NamespacedName{
		Namespace: desiredObj.GetNamespace(),
		Name:      desiredObj.GetName(),
	}, curObj)
	if err != nil {
		if !errors.IsNotFound(err) {
			return false, fmt.Errorf("failed to get object: %s", err)
		}

		// Object does not exist, create it
		// Setting last applied annotation before setting controller reference since it excludes the
		// "metadata.ownerReferences" from the annotation. Since a later comparison happens between the annotation
		// and the "desired" object, both of them should not contain the controller reference.
		if err = patch.DefaultAnnotator.SetLastAppliedAnnotation(desiredObj); err != nil {
			return false, fmt.Errorf("failed to set last applied annotation: %s", err)
		}
		if err = ctrl.SetControllerReference(plantDCore, desiredObj, r.Scheme); err != nil {
			return false, fmt.Errorf("failed to set controller reference: %s", err)
		}

		if err = r.Create(ctx, desiredObj); err != nil {
			return false, fmt.Errorf("failed to create object: %s", err)
		}

		return true, nil
	}

	if !allowUpdate {
		return false, nil
	}

	// Object exists, compare and update if necessary
	compareOpts := []patch.CalculateOption{
		patch.IgnoreStatusFields(),
	}
	patchResult, err := patch.DefaultPatchMaker.Calculate(curObj, desiredObj, compareOpts...)
	if err != nil {
		return false, fmt.Errorf("failed to compare objects: %s", err)
	}
	if !patchResult.IsEmpty() {
		// Setting last applied annotation before setting controller reference since it excludes the
		// "metadata.ownerReferences" from the annotation. Since a later comparison happens between the annotation
		// and the "desired" object, both of them should not contain the controller reference.
		if err = patch.DefaultAnnotator.SetLastAppliedAnnotation(desiredObj); err != nil {
			return false, fmt.Errorf("failed to set last applied annotation: %s", err)
		}
		if err = ctrl.SetControllerReference(plantDCore, desiredObj, r.Scheme); err != nil {
			return false, fmt.Errorf("failed to set controller reference: %s", err)
		}

		// Avoid "metadata.resourceVersion: Invalid value: 0x0: must be specified for an update" error in some cases,
		// see https://github.com/argoproj/argo-cd/issues/3657.
		desiredObj.SetResourceVersion(curObj.GetResourceVersion())

		if err = r.Update(ctx, desiredObj); err != nil {
			return false, fmt.Errorf("failed to update object: %s", err)
		}

		return true, nil
	}

	return false, nil
}

// reconcileClusterObject ensures the actual state of the cluster-level object matches the desired state. It creates the
// object if it does not exist, and otherwise updates it if necessary. Update behavior can be enabled or disabled. It
// returns a bool to show if any action is taken and the error if any.
func (r *PlantDCoreReconciler) reconcileClusterObject(ctx context.Context, curObj client.Object, desiredObj client.Object, allowUpdate bool) (bool, error) {
	// Get current object
	err := r.Get(ctx, types.NamespacedName{
		Name: desiredObj.GetName(),
	}, curObj)
	if err != nil {
		if !errors.IsNotFound(err) {
			return false, fmt.Errorf("failed to get object: %s", err)
		}

		// Object does not exist, create it
		if err = patch.DefaultAnnotator.SetLastAppliedAnnotation(desiredObj); err != nil {
			return false, fmt.Errorf("failed to set last applied annotation: %s", err)
		}

		if err = r.Create(ctx, desiredObj); err != nil {
			return false, fmt.Errorf("failed to create object: %s", err)
		}

		return true, nil
	}

	if !allowUpdate {
		return false, nil
	}

	// Object exists, compare and update if necessary
	compareOpts := []patch.CalculateOption{
		patch.IgnoreStatusFields(),
	}
	patchResult, err := patch.DefaultPatchMaker.Calculate(curObj, desiredObj, compareOpts...)
	if err != nil {
		return false, fmt.Errorf("failed to compare objects: %s", err)
	}
	if !patchResult.IsEmpty() {
		if err = patch.DefaultAnnotator.SetLastAppliedAnnotation(desiredObj); err != nil {
			return false, fmt.Errorf("failed to set last applied annotation: %s", err)
		}

		// Avoid "metadata.resourceVersion: Invalid value: 0x0: must be specified for an update" error in some cases,
		// see https://github.com/argoproj/argo-cd/issues/3657.
		desiredObj.SetResourceVersion(curObj.GetResourceVersion())

		if err = r.Update(ctx, desiredObj); err != nil {
			return false, fmt.Errorf("failed to update object: %s", err)
		}

		return true, nil
	}

	return false, nil
}

// finalizeResources cleans up the resources before deletion. It returns the error if any.
func (r *PlantDCoreReconciler) finalizeResources(ctx context.Context, plantDCore *windtunnelv1alpha1.PlantDCore) error {
	logger := log.FromContext(ctx)

	_, clusterRole, clusterRoleBinding := core.GetPrometheusRBACResources(plantDCore)

	if err := r.Delete(ctx, clusterRole); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to delete Prometheus ClusterRole")
		return err
	}
	logger.Info("deleted Prometheus ClusterRole")

	if err := r.Delete(ctx, clusterRoleBinding); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "failed to delete Prometheus ClusterRoleBinding")
		return err
	}
	logger.Info("deleted Prometheus ClusterRoleBinding")

	return nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PlantDCore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:ec
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *PlantDCoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the PlantD resource
	plantDCore := &windtunnelv1alpha1.PlantDCore{}
	if err := r.Get(ctx, req.NamespacedName, plantDCore); err != nil {
		// Ignore not-found errors as we can get them on delete requests
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Finalizer
	if plantDCore.ObjectMeta.DeletionTimestamp.IsZero() {
		// Object is not being deleted, add finalizer if not present
		if !controllerutil.ContainsFinalizer(plantDCore, finalizerName) {
			controllerutil.AddFinalizer(plantDCore, finalizerName)
			if err := r.Update(ctx, plantDCore); err != nil {
				logger.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		// Object is being deleted
		if controllerutil.ContainsFinalizer(plantDCore, finalizerName) {
			// Finalizer presents
			if err := r.finalizeResources(ctx, plantDCore); err != nil {
				// Retry on errors
				return ctrl.Result{}, err
			}

			// Remove finalizer from the list and update it.
			controllerutil.RemoveFinalizer(plantDCore, finalizerName)
			if err := r.Update(ctx, plantDCore); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	curKubeProxyDeployment := &appsv1.Deployment{}
	curKubeProxyService := &corev1.Service{}
	curStudioDeployment := &appsv1.Deployment{}
	curStudioService := &corev1.Service{}
	curPrometheusServiceAccount := &corev1.ServiceAccount{}
	curPrometheusClusterRole := &rbacv1.ClusterRole{}
	curPrometheusClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	curPrometheusObject := &monitoringv1.Prometheus{}
	curPrometheusService := &corev1.Service{}
	curRedisDeployment := &appsv1.Deployment{}
	curRedisService := &corev1.Service{}

	desiredKubeProxyDeployment := core.GetKubeProxyDeployment(plantDCore)
	desiredKubeProxyService := core.GetKubeProxyService(plantDCore)
	desiredStudioDeployment := core.GetStudioDeployment(plantDCore)
	desiredStudioService := core.GetStudioService(plantDCore)
	desiredPrometheusServiceAccount, desiredPrometheusClusterRole, desiredPrometheusClusterRoleBinding := core.GetPrometheusRBACResources(plantDCore)
	desiredPrometheusObject := core.GetPrometheusObject(plantDCore)
	desiredPrometheusService := core.GetPrometheusService(plantDCore)
	desiredRedisDeployment := core.GetRedisDeployment(plantDCore)
	desiredRedisService := core.GetRedisService(plantDCore)

	hasModification := false

	if modified, err := r.reconcileObject(ctx, plantDCore, curKubeProxyDeployment, desiredKubeProxyDeployment, true); err != nil {
		logger.Error(err, "failed to reconcile Kube Proxy Deployment")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created or updated Kube Proxy Deployment")
		hasModification = true
	}
	if modified, err := r.reconcileObject(ctx, plantDCore, curKubeProxyService, desiredKubeProxyService, false); err != nil {
		logger.Error(err, "failed to reconcile Kube Proxy Service")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created Kube Proxy Service")
		hasModification = true
	}
	if modified, err := r.reconcileObject(ctx, plantDCore, curStudioDeployment, desiredStudioDeployment, true); err != nil {
		logger.Error(err, "failed to reconcile Studio Deployment")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created or updated Studio Deployment")
		hasModification = true
	}
	if modified, err := r.reconcileObject(ctx, plantDCore, curStudioService, desiredStudioService, false); err != nil {
		logger.Error(err, "failed to reconcile Studio Service")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created Studio Service")
		hasModification = true
	}
	if modified, err := r.reconcileObject(ctx, plantDCore, curPrometheusServiceAccount, desiredPrometheusServiceAccount, false); err != nil {
		logger.Error(err, "failed to reconcile Prometheus ServiceAccount")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created Prometheus ServiceAccount")
		hasModification = true
	}
	if modified, err := r.reconcileClusterObject(ctx, curPrometheusClusterRole, desiredPrometheusClusterRole, false); err != nil {
		logger.Error(err, "failed to reconcile Prometheus ClusterRole")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created Prometheus ClusterRole")
		hasModification = true
	}
	if modified, err := r.reconcileClusterObject(ctx, curPrometheusClusterRoleBinding, desiredPrometheusClusterRoleBinding, false); err != nil {
		logger.Error(err, "failed to reconcile Prometheus ClusterRoleBinding")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created Prometheus ClusterRoleBinding")
		hasModification = true
	}
	if modified, err := r.reconcileObject(ctx, plantDCore, curPrometheusObject, desiredPrometheusObject, true); err != nil {
		logger.Error(err, "failed to reconcile Prometheus object")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created or updated Prometheus object")
		hasModification = true
	}
	if modified, err := r.reconcileObject(ctx, plantDCore, curPrometheusService, desiredPrometheusService, false); err != nil {
		logger.Error(err, "failed to reconcile Prometheus Service")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created Prometheus Service")
		hasModification = true
	}
	if modified, err := r.reconcileObject(ctx, plantDCore, curRedisDeployment, desiredRedisDeployment, true); err != nil {
		logger.Error(err, "failed to reconcile Redis Deployment")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created or updated Redis Deployment")
		hasModification = true
	}
	if modified, err := r.reconcileObject(ctx, plantDCore, curRedisService, desiredRedisService, false); err != nil {
		logger.Error(err, "failed to reconcile Redis Service")
		return ctrl.Result{}, err
	} else if modified {
		logger.Info("created Redis Service")
		hasModification = true
	}

	_ = hasModification

	return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlantDCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.PlantDCore{}).
		Complete(r)
}
