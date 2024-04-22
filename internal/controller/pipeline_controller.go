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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/monitor"
)

const (
	pipelineFinalizerName = "pipeline.windtunnel.plantd.org/finalizer"
)

var (
	metricsServiceLabelKeyPipeline = config.GetString("monitor.service.labelKeys.pipeline")
)

// PipelineReconciler reconciles a Pipeline object
type PipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines/finalizers,verbs=update
//
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete

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

	// Check if the Pipeline is being deleted
	if pipeline.DeletionTimestamp.IsZero() {
		// Add finalizer to the Pipeline
		if !controllerutil.ContainsFinalizer(pipeline, pipelineFinalizerName) {
			controllerutil.AddFinalizer(pipeline, pipelineFinalizerName)
			if err := r.Update(ctx, pipeline); err != nil {
				logger.Error(err, "Cannot add finalizer")
				return ctrl.Result{}, err
			}
			logger.Info("Added finalizer")
		}
	} else {
		if controllerutil.ContainsFinalizer(pipeline, pipelineFinalizerName) {
			if err := r.removePipelineLabel(ctx, pipeline); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(pipeline, pipelineFinalizerName)
			if err := r.Update(ctx, pipeline); err != nil {
				logger.Error(err, "Cannot remove finalizer")
				return ctrl.Result{}, err
			}
			logger.Info("Removed finalizer")
		}
		return ctrl.Result{}, nil
	}

	// Initialize the Pipeline
	if pipeline.Status.Availability == "" {
		if err := r.initializeMonitor(ctx, pipeline); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Pipeline is already initialized, no need to re-queue
	return ctrl.Result{}, nil
}

// initializeMonitor creates monitoring resources for the Pipeline.
func (r *PipelineReconciler) initializeMonitor(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline) error {
	logger := log.FromContext(ctx)

	if !containMetricsEndpoint(pipeline) {
		logger.Info("No metrics endpoint is specified")
		return nil
	}

	if pipeline.Spec.InCluster {
		// For in-cluster Pipeline, the user creates the metrics Service
		service := &corev1.Service{}
		serviceName := types.NamespacedName{
			Namespace: pipeline.Namespace,
			Name:      pipeline.Spec.MetricsEndpoint.ServiceRef.Name,
		}
		if err := r.Get(ctx, serviceName, service); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot get metrics Service \"%s\"", serviceName))
			return err
		}

		// Add label so that the ServiceMonitor can select it
		if service.Labels == nil {
			service.Labels = make(map[string]string, 1)
		}
		service.Labels[metricsServiceLabelKeyPipeline] = pipeline.Name
		if err := r.Update(ctx, service); err != nil {
			logger.Error(err, fmt.Sprintf("Cannot add Pipeline label to metrics Service \"%s\"", serviceName))
			return err
		}
		logger.Info(fmt.Sprintf("Added Pipeline label to metrics Service \"%s\"", serviceName))
	} else {
		// For out-cluster Pipeline, we need to create the metrics Service of type ExternalName
		// in the same Namespace as the Pipeline.
		service, err := monitor.CreateExternalNameService(pipeline)
		if err != nil {
			logger.Error(err, "Cannot create manifest for ExternalName Service")
			return err
		}
		if err := ctrl.SetControllerReference(pipeline, service, r.Scheme); err != nil {
			logger.Error(err, "Cannot set controller reference for ExternalName Service")
			return err
		}
		if err := r.Create(ctx, service); err != nil {
			logger.Error(err, "Cannot create ExternalName Service")
			return err
		} else if err == nil {
			logger.Info("Created ExternalName Service")
		}
	}

	// Create the ServiceMonitor
	serviceMonitor, err := monitor.CreateServiceMonitor(pipeline)
	if err != nil {
		logger.Error(err, "Cannot create manifest for ServiceMonitor")
		return err
	}
	if err := ctrl.SetControllerReference(pipeline, serviceMonitor, r.Scheme); err != nil {
		logger.Error(err, "Cannot set controller reference for ServiceMonitor")
		return err
	}
	if err := r.Create(ctx, serviceMonitor); err != nil {
		logger.Error(err, "Cannot create ServiceMonitor")
		return err
	} else if err == nil {
		logger.Info("Created ServiceMonitor")
	}

	// Update the status
	pipeline.Status.Availability = windtunnelv1alpha1.PipelineReady
	if err := r.Status().Update(ctx, pipeline); err != nil {
		logger.Error(err, "Cannot update the status")
		return err
	}

	return nil
}

// removePipelineLabel removes the Pipeline label from the metrics Service.
// For in-cluster Pipeline only.
func (r *PipelineReconciler) removePipelineLabel(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline) error {
	logger := log.FromContext(ctx)

	// No need to clean up Pipeline label for out-cluster Pipeline,
	// as the Service will be deleted along with the Pipeline.
	if !pipeline.Spec.InCluster {
		return nil
	}

	service := &corev1.Service{}
	serviceName := types.NamespacedName{
		Namespace: pipeline.Namespace,
		Name:      pipeline.Spec.MetricsEndpoint.ServiceRef.Name,
	}
	if err := r.Get(ctx, serviceName, service); err != nil {
		logger.Error(err, fmt.Sprintf("Lost metrics Service \"%s\"", serviceName))
		return nil
	}

	if service.Labels == nil {
		return nil
	}
	if val, ok := service.Labels[metricsServiceLabelKeyPipeline]; !ok || val != pipeline.Name {
		return nil
	}

	delete(service.Labels, metricsServiceLabelKeyPipeline)
	if err := r.Update(ctx, service); err != nil {
		logger.Error(err, fmt.Sprintf("Cannot remove Pipeline label from metrics Service \"%s\"", serviceName))
		return err
	}
	logger.Info(fmt.Sprintf("Removed Pipeline label from metrics Service \"%s\"", serviceName))

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Pipeline{}).
		Complete(r)
}
