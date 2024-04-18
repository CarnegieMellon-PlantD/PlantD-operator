package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/cisco-open/k8s-objectmatcher/patch"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/core"
)

const (
	plantDCorePollingInterval = 10 * time.Second
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
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheuses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:ec
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *PlantDCoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the requested PlantDCore
	plantDCore := &windtunnelv1alpha1.PlantDCore{}
	if err := r.Get(ctx, req.NamespacedName, plantDCore); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unable to fetch PlantDCore")
		return ctrl.Result{}, err
	}

	// PlantD-Proxy
	{
		curProxyDeployment := &appsv1.Deployment{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			curProxyDeployment,
			core.GetProxyDeployment(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			&corev1.Service{},
			core.GetProxyService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setDeploymentComponentStatus(&plantDCore.Status.ProxyStatus, true, curProxyDeployment)
	}

	// PlantD-Studio
	{
		curStudioDeployment := &appsv1.Deployment{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			curStudioDeployment,
			core.GetStudioDeployment(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			&corev1.Service{},
			core.GetStudioService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setDeploymentComponentStatus(&plantDCore.Status.StudioStatus, true, curStudioDeployment)
	}

	// Prometheus
	{
		curPrometheusObject := &monitoringv1.Prometheus{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			curPrometheusObject,
			core.GetPrometheusObject(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			&corev1.Service{},
			core.GetPrometheusService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setPrometheusComponentStatus(&plantDCore.Status.StudioStatus, true, curPrometheusObject)
	}

	hasThanosObjStoreCfg := plantDCore.Spec.ThanosConfig.ObjectStoreConfig != nil

	// Thanos-Store
	{
		curThanosStoreStatefulSet := &appsv1.StatefulSet{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			curThanosStoreStatefulSet,
			core.GetThanosStoreStatefulSet(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			&corev1.Service{},
			core.GetThanosStoreService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setStatefulSetComponentStatus(&plantDCore.Status.ThanosStoreStatus, hasThanosObjStoreCfg, curThanosStoreStatefulSet)
	}

	// Thanos-Compactor
	{
		curThanosCompactorStatefulSet := &appsv1.StatefulSet{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			curThanosCompactorStatefulSet,
			core.GetThanosCompactorStatefulSet(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			&corev1.Service{},
			core.GetThanosCompactorService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setStatefulSetComponentStatus(&plantDCore.Status.ThanosCompactorStatus, hasThanosObjStoreCfg, curThanosCompactorStatefulSet)
	}

	// Thanos-Querier
	{
		curThanosQuerierStatefulSet := &appsv1.StatefulSet{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			curThanosQuerierStatefulSet,
			core.GetThanosQuerierStatefulSet(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			&corev1.Service{},
			core.GetThanosQuerierService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setStatefulSetComponentStatus(&plantDCore.Status.ThanosQuerierStatus, hasThanosObjStoreCfg, curThanosQuerierStatefulSet)
	}

	// Redis
	{
		curRedisDeployment := &appsv1.Deployment{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			curRedisDeployment,
			core.GetRedisDeployment(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			&corev1.Service{},
			core.GetRedisService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setDeploymentComponentStatus(&plantDCore.Status.RedisStatus, true, curRedisDeployment)
	}

	// OpenCost
	{
		curOpenCostDeployment := &appsv1.Deployment{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			curOpenCostDeployment,
			core.GetOpenCostDeployment(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			&corev1.Service{},
			core.GetOpenCostService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			&monitoringv1.ServiceMonitor{},
			core.GetOpenCostServiceMonitor(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			&monitoringv1.ServiceMonitor{},
			core.GetCAdvisorServiceMonitor(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setDeploymentComponentStatus(&plantDCore.Status.OpenCostStatus, true, curOpenCostDeployment)
	}

	return ctrl.Result{RequeueAfter: plantDCorePollingInterval}, nil
}

// reconcileObject ensures the current state of the object matches the desired state.
func (r *PlantDCoreReconciler) reconcileObject(ctx context.Context, plantDCore *windtunnelv1alpha1.PlantDCore, shouldExist bool, curObj client.Object, desiredObj client.Object) error {
	logger := log.FromContext(ctx)

	// Get current object
	curObjName := types.NamespacedName{
		Namespace: desiredObj.GetNamespace(),
		Name:      desiredObj.GetName(),
	}
	if err := r.Get(ctx, curObjName, curObj); client.IgnoreNotFound(err) != nil {
		logger.Error(err, fmt.Sprintf("Cannot get %s/%s", curObj.GetObjectKind().GroupVersionKind().Kind, curObjName))
		return err
	} else if apierrors.IsNotFound(err) {
		// Object does not exist, create it if necessary
		if shouldExist {
			// Setting last-applied annotation before setting controller reference,
			// because a later comparison happens between the last-applied annotation and the desired object,
			// and the desired object does not have the controller reference.
			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(desiredObj); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set last-applied annotation for %s", getObjectName(desiredObj)))
				return err
			}
			if err := ctrl.SetControllerReference(plantDCore, desiredObj, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for %s", getObjectName(desiredObj)))
				return err
			}
			if err := r.Create(ctx, desiredObj); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot create %s", getObjectName(desiredObj)))
				return err
			}
			logger.Info(fmt.Sprintf("Created %s", getObjectName(desiredObj)))
			return nil
		}
	} else {
		// Object already exists, delete it if necessary
		if !shouldExist {
			if err := r.Delete(ctx, curObj); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot delete %s", getObjectName(curObj)))
				return err
			}
			logger.Info(fmt.Sprintf("Deleted %s", getObjectName(curObj)))
			return nil
		}
	}

	// Object already exists, compare and update if necessary
	compareOpts := []patch.CalculateOption{
		patch.IgnoreStatusFields(),
	}
	patchResult, err := patch.DefaultPatchMaker.Calculate(curObj, desiredObj, compareOpts...)
	if err != nil {
		logger.Error(err, fmt.Sprintf("Cannot calculate patch for %s", getObjectName(desiredObj)))
		return err
	}
	if patchResult.IsEmpty() {
		return nil
	}

	if err = patch.DefaultAnnotator.SetLastAppliedAnnotation(desiredObj); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot set last-applied annotation for %s", getObjectName(desiredObj)))
		return err
	}
	if err = ctrl.SetControllerReference(plantDCore, desiredObj, r.Scheme); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot set controller reference for %s", getObjectName(desiredObj)))
		return err
	}
	// Avoid error "metadata.resourceVersion: Invalid value: 0x0: must be specified for an update".
	// See https://github.com/argoproj/argo-cd/issues/3657.
	desiredObj.SetResourceVersion(curObj.GetResourceVersion())
	if err = r.Update(ctx, desiredObj); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot update %s", getObjectName(desiredObj)))
		return err
	}
	logger.Info(fmt.Sprintf("Updated %s", getObjectName(desiredObj)))

	return nil
}

// setDeploymentComponentStatus sets the ComponentStatus based on the deployment.
func (r *PlantDCoreReconciler) setDeploymentComponentStatus(componentStatus *windtunnelv1alpha1.ComponentStatus, shouldExist bool, deployment *appsv1.Deployment) {
	if !shouldExist {
		componentStatus.Text = windtunnelv1alpha1.ComponentSkipped
		componentStatus.NumReady = 0
		componentStatus.NumDesired = 0
		return
	}
	if deployment.Status.ReadyReplicas > 0 && deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
		componentStatus.Text = windtunnelv1alpha1.ComponentReady
	} else {
		componentStatus.Text = windtunnelv1alpha1.ComponentNotReady
	}
	componentStatus.NumReady = deployment.Status.ReadyReplicas
	componentStatus.NumDesired = *deployment.Spec.Replicas
}

// setStatefulSetComponentStatus sets the ComponentStatus based on the StatefulSet.
func (r *PlantDCoreReconciler) setStatefulSetComponentStatus(componentStatus *windtunnelv1alpha1.ComponentStatus, shouldExist bool, statefulSet *appsv1.StatefulSet) {
	if !shouldExist {
		componentStatus.Text = windtunnelv1alpha1.ComponentSkipped
		componentStatus.NumReady = 0
		componentStatus.NumDesired = 0
		return
	}
	if statefulSet.Status.ReadyReplicas > 0 && statefulSet.Status.ReadyReplicas == *statefulSet.Spec.Replicas {
		componentStatus.Text = windtunnelv1alpha1.ComponentReady
	} else {
		componentStatus.Text = windtunnelv1alpha1.ComponentNotReady
	}
	componentStatus.NumReady = statefulSet.Status.ReadyReplicas
	componentStatus.NumDesired = *statefulSet.Spec.Replicas
}

// setPrometheusComponentStatus sets the ComponentStatus based on the Prometheus object.
func (r *PlantDCoreReconciler) setPrometheusComponentStatus(componentStatus *windtunnelv1alpha1.ComponentStatus, shouldExist bool, prometheus *monitoringv1.Prometheus) {
	if !shouldExist {
		componentStatus.Text = windtunnelv1alpha1.ComponentSkipped
		componentStatus.NumReady = 0
		componentStatus.NumDesired = 0
		return
	}
	if prometheus.Status.AvailableReplicas > 0 && prometheus.Status.AvailableReplicas == *prometheus.Spec.Replicas {
		componentStatus.Text = windtunnelv1alpha1.ComponentReady
	} else {
		componentStatus.Text = windtunnelv1alpha1.ComponentNotReady
	}
	componentStatus.NumReady = prometheus.Status.AvailableReplicas
	componentStatus.NumDesired = *prometheus.Spec.Replicas
}

// getObjectName returns the string representation of the object, containing its kind, namespace, and name.
func getObjectName(obj client.Object) string {
	return fmt.Sprintf("%s/%s/%s", obj.GetObjectKind().GroupVersionKind().Kind, obj.GetNamespace(), obj.GetName())
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlantDCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.PlantDCore{}).
		Complete(r)
}
