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
			"Deployment",
			curProxyDeployment,
			core.GetProxyDeployment(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"Service",
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
			"Deployment",
			curStudioDeployment,
			core.GetStudioDeployment(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"Service",
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
			"Prometheus",
			curPrometheusObject,
			core.GetPrometheusObject(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"Service",
			&corev1.Service{},
			core.GetPrometheusService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setPrometheusComponentStatus(&plantDCore.Status.PrometheusStatus, true, curPrometheusObject)
	}

	hasThanosObjStoreCfg := plantDCore.Spec.ThanosConfig.ObjectStoreConfig != nil

	// Thanos-Store
	{
		curThanosStoreStatefulSet := &appsv1.StatefulSet{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			"StatefulSet",
			curThanosStoreStatefulSet,
			core.GetThanosStoreStatefulSet(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			"Service",
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
			"StatefulSet",
			curThanosCompactorStatefulSet,
			core.GetThanosCompactorStatefulSet(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			hasThanosObjStoreCfg,
			"Service",
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
			true,
			"StatefulSet",
			curThanosQuerierStatefulSet,
			core.GetThanosQuerierStatefulSet(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"Service",
			&corev1.Service{},
			core.GetThanosQuerierService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setStatefulSetComponentStatus(&plantDCore.Status.ThanosQuerierStatus, true, curThanosQuerierStatefulSet)
	}

	// Redis
	{
		curRedisDeployment := &appsv1.Deployment{}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"Deployment",
			curRedisDeployment,
			core.GetRedisDeployment(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"Service",
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
			"Deployment",
			curOpenCostDeployment,
			core.GetOpenCostDeployment(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"Service",
			&corev1.Service{},
			core.GetOpenCostService(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"ServiceMonitor",
			&monitoringv1.ServiceMonitor{},
			core.GetOpenCostServiceMonitor(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(
			ctx,
			plantDCore,
			true,
			"ServiceMonitor",
			&monitoringv1.ServiceMonitor{},
			core.GetCAdvisorServiceMonitor(plantDCore),
		); err != nil {
			return ctrl.Result{}, err
		}
		r.setDeploymentComponentStatus(&plantDCore.Status.OpenCostStatus, true, curOpenCostDeployment)
	}

	if err := r.Status().Update(ctx, plantDCore); err != nil {
		logger.Error(err, "Cannot update the status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: plantDCorePollingInterval}, nil
}

// reconcileObject ensures the current state of the object matches the desired state.
func (r *PlantDCoreReconciler) reconcileObject(ctx context.Context, plantDCore *windtunnelv1alpha1.PlantDCore, shouldExist bool, kind string, curObj client.Object, desiredObj client.Object) error {
	logger := log.FromContext(ctx)

	// Get current object
	curObjName := types.NamespacedName{
		Namespace: desiredObj.GetNamespace(),
		Name:      desiredObj.GetName(),
	}
	if err := r.Get(ctx, curObjName, curObj); err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Error(err, fmt.Sprintf("Cannot get %s/%s", kind, curObjName))
			return err
		}

		// Object does not exist, create it if necessary
		if shouldExist {
			// Setting last-applied annotation before setting controller reference,
			// because a later comparison happens between the last-applied annotation and the desired object,
			// and the desired object does not have the controller reference.
			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(desiredObj); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set last-applied annotation for %s", getObjectName(kind, desiredObj)))
				return err
			}
			if err := ctrl.SetControllerReference(plantDCore, desiredObj, r.Scheme); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot set controller reference for %s", getObjectName(kind, desiredObj)))
				return err
			}
			if err := r.Create(ctx, desiredObj); err != nil {
				logger.Error(err, fmt.Sprintf("Cannot create %s", getObjectName(kind, desiredObj)))
				return err
			}
			logger.Info(fmt.Sprintf("Created %s", getObjectName(kind, desiredObj)))
		}

		return nil
	}

	// Object already exists, delete it if necessary
	if !shouldExist {
		if err := r.Delete(ctx, curObj); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot delete %s", getObjectName(kind, curObj)))
			return err
		}
		logger.Info(fmt.Sprintf("Deleted %s", getObjectName(kind, curObj)))
		return nil
	}

	// Object already exists, compare and update if necessary
	compareOpts := []patch.CalculateOption{
		patch.IgnoreStatusFields(),
		patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus(),
	}
	patchResult, err := patch.DefaultPatchMaker.Calculate(curObj, desiredObj, compareOpts...)
	if err != nil {
		logger.Error(err, fmt.Sprintf("Cannot calculate patch for %s", getObjectName(kind, desiredObj)))
		return err
	}
	if patchResult.IsEmpty() {
		return nil
	}

	if err = patch.DefaultAnnotator.SetLastAppliedAnnotation(desiredObj); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot set last-applied annotation for %s", getObjectName(kind, desiredObj)))
		return err
	}
	if err = ctrl.SetControllerReference(plantDCore, desiredObj, r.Scheme); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot set controller reference for %s", getObjectName(kind, desiredObj)))
		return err
	}
	// Avoid "metadata.resourceVersion: Invalid value: 0x0: must be specified for an update" error in some cases,
	// see https://github.com/argoproj/argo-cd/issues/3657.
	desiredObj.SetResourceVersion(curObj.GetResourceVersion())
	if err = r.Update(ctx, desiredObj); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot update %s", getObjectName(kind, desiredObj)))
		return err
	}
	logger.Info(fmt.Sprintf("Updated %s", getObjectName(kind, desiredObj)))

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
	if deployment.Status.ReadyReplicas > 0 && deployment.Status.ReadyReplicas == deployment.Status.Replicas {
		componentStatus.Text = windtunnelv1alpha1.ComponentReady
	} else {
		componentStatus.Text = windtunnelv1alpha1.ComponentNotReady
	}
	componentStatus.NumReady = deployment.Status.ReadyReplicas
	componentStatus.NumDesired = deployment.Status.Replicas
}

// setStatefulSetComponentStatus sets the ComponentStatus based on the StatefulSet.
func (r *PlantDCoreReconciler) setStatefulSetComponentStatus(componentStatus *windtunnelv1alpha1.ComponentStatus, shouldExist bool, statefulSet *appsv1.StatefulSet) {
	if !shouldExist {
		componentStatus.Text = windtunnelv1alpha1.ComponentSkipped
		componentStatus.NumReady = 0
		componentStatus.NumDesired = 0
		return
	}
	if statefulSet.Status.ReadyReplicas > 0 && statefulSet.Status.ReadyReplicas == statefulSet.Status.Replicas {
		componentStatus.Text = windtunnelv1alpha1.ComponentReady
	} else {
		componentStatus.Text = windtunnelv1alpha1.ComponentNotReady
	}
	componentStatus.NumReady = statefulSet.Status.ReadyReplicas
	componentStatus.NumDesired = statefulSet.Status.Replicas
}

// setPrometheusComponentStatus sets the ComponentStatus based on the Prometheus object.
func (r *PlantDCoreReconciler) setPrometheusComponentStatus(componentStatus *windtunnelv1alpha1.ComponentStatus, shouldExist bool, prometheus *monitoringv1.Prometheus) {
	if !shouldExist {
		componentStatus.Text = windtunnelv1alpha1.ComponentSkipped
		componentStatus.NumReady = 0
		componentStatus.NumDesired = 0
		return
	}
	if prometheus.Status.AvailableReplicas > 0 && prometheus.Status.AvailableReplicas == prometheus.Status.Replicas {
		componentStatus.Text = windtunnelv1alpha1.ComponentReady
	} else {
		componentStatus.Text = windtunnelv1alpha1.ComponentNotReady
	}
	componentStatus.NumReady = prometheus.Status.AvailableReplicas
	componentStatus.NumDesired = prometheus.Status.Replicas
}

// getObjectName returns the string representation of the object, containing its kind, namespace, and name.
func getObjectName(kind string, obj client.Object) string {
	return fmt.Sprintf("%s/%s/%s", kind, obj.GetNamespace(), obj.GetName())
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlantDCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.PlantDCore{}).
		Complete(r)
}
