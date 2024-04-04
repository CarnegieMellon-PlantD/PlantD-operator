package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/monitor"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

var (
	metricsServiceLabelKeyPipeline = config.GetViper().GetString("monitor.service.labelKeys.pipeline")
)

// PipelineReconciler reconciles a Pipeline object
type PipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the requested Pipeline
	pipeline := &windtunnelv1alpha1.Pipeline{}
	if err := r.Get(ctx, req.NamespacedName, pipeline); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unable to fetch Pipeline")
		return ctrl.Result{}, err
	}

	// Initialize the Pipeline
	if pipeline.Status.Availability == "" {
		return r.reconcileCreated(ctx, pipeline)
	}

	// Pipeline is already initialized, no action needed
	return ctrl.Result{}, nil
}

// reconcileCreated reconciles the Pipeline when it is created.
func (r *PipelineReconciler) reconcileCreated(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if pipeline.Spec.InCluster {
		// For in-cluster Pipeline, the user creates the metrics Service.
		// Add label so that the ServiceMonitor can select it.
		service := &corev1.Service{}
		serviceName := types.NamespacedName{
			Namespace: pipeline.Spec.MetricsEndpoint.ServiceRef.Namespace,
			Name:      pipeline.Spec.MetricsEndpoint.ServiceRef.Name,
		}
		if err := r.Get(ctx, serviceName, service); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get metrics Service: %s", serviceName))
			return ctrl.Result{}, err
		}

		if service.Labels == nil {
			service.Labels = make(map[string]string, 1)
		}
		service.Labels[metricsServiceLabelKeyPipeline] = pipeline.Name
		if err := r.Update(ctx, service); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot add Pipeline label to metrics Service: %s", serviceName))
			return ctrl.Result{}, err
		}
		logger.Info(fmt.Sprintf("Added Pipeline label to metrics Service: %s", serviceName))
	} else {
		// For out-cluster Pipeline, we need to create the metrics Service of type ExternalName
		// in the same Namespace as the Pipeline.
		service, err := monitor.CreateExternalNameService(pipeline)
		if err != nil {
			logger.Error(err, "Cannot prepare ExternalName Service object to create")
			return ctrl.Result{}, err
		}
		if err := ctrl.SetControllerReference(pipeline, service, r.Scheme); err != nil {
			logger.Error(err, "Cannot set controller reference for ExternalName Service")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, service); err != nil {
			logger.Error(err, "Cannot create ExternalName Service")
			return ctrl.Result{}, err
		} else if err == nil {
			logger.Info("Created ExternalName Service")
		}
	}

	// Create the ServiceMonitor
	serviceMonitor, err := monitor.CreateServiceMonitor(pipeline)
	if err != nil {
		logger.Error(err, "Cannot prepare ServiceMonitor object to create")
		return ctrl.Result{}, err
	}
	if err := ctrl.SetControllerReference(pipeline, serviceMonitor, r.Scheme); err != nil {
		logger.Error(err, "Cannot set controller reference for ServiceMonitor")
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, serviceMonitor); err != nil {
		logger.Error(err, "Cannot create ServiceMonitor")
		return ctrl.Result{}, err
	} else if err == nil {
		logger.Info(fmt.Sprintf("Created ServiceMonitor: %s", utils.GetNamespacedName(serviceMonitor)))
	}

	// Update the status
	pipeline.Status.Availability = windtunnelv1alpha1.PipelineReady
	if err := r.Status().Update(ctx, pipeline); err != nil {
		logger.Error(err, "Cannot update the status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Pipeline{}).
		Complete(r)
}
