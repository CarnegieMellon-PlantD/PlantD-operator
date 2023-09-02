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

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/monitor"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

var (
	experimentLabelKey string
	metricsLabelKey    string
)

func init() {
	experimentLabelKey = config.GetString("monitor.jobLabel")
	metricsLabelKey = config.GetString("monitor.metricsService.labels.key")
}

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
//+kubebuilder:rbac:groups="",resources=endpoints,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Get the pipeline object in the queue.
	pipeline := &windtunnelv1alpha1.Pipeline{}
	if err := r.Get(ctx, req.NamespacedName, pipeline); err != nil {
		if apierrors.IsNotFound(err) {
			// Custom resource not found, perform cleanup tasks here.
			// Created objects are automatically garbage collected.
			log.Info("Pipeline is deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Cannot get Pipeline: "+req.NamespacedName.String())
		return ctrl.Result{}, err
	}

	// Initialize the state of the state of the pipeline.
	if pipeline.Status.PipelineState == "" {
		if err := r.UpdateStatus(ctx, pipeline, windtunnelv1alpha1.PipelineInitializing, ""); err != nil {
			log.Error(err, "Cannot update the status of the Pipeline: "+req.NamespacedName.String())
			return ctrl.Result{}, err
		}

		if err := r.Initialize(ctx, pipeline); err != nil {
			log.Error(err, "Cannot Initialize the Pipeline: "+req.NamespacedName.String())
			return ctrl.Result{}, err
		}
	}

	// Lock the pipeline to the specific experiment object.
	if pipeline.Spec.ExperimentRef.Name != "" && pipeline.Status.PipelineState == windtunnelv1alpha1.PipelineAvailable {
		if err := r.HealthCheck(ctx, pipeline); err != nil {
			log.Error(err, "Health Check Failled: "+req.NamespacedName.String())
			return ctrl.Result{}, r.UpdateStatus(ctx, pipeline, "", windtunnelv1alpha1.PipelineFailed)
		}

		if err := r.InitializeExp(ctx, pipeline); err != nil {
			log.Error(err, "Failled to initialize the experiment: "+req.NamespacedName.String())
			return ctrl.Result{}, err
		}
	}

	// Unlock the pipeline.
	if pipeline.Spec.ExperimentRef.Name == "" && pipeline.Status.PipelineState == windtunnelv1alpha1.PipelineEngaged {
		pipeline.Status.PipelineState = windtunnelv1alpha1.PipelineAvailable
		return ctrl.Result{}, r.UpdateStatus(ctx, pipeline, windtunnelv1alpha1.PipelineAvailable, "")
	}

	return ctrl.Result{}, nil
}

func (r *PipelineReconciler) UpdateStatus(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline, pipelineState string, statusCheck string) error {
	log := log.FromContext(ctx)
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		newPipeline := &windtunnelv1alpha1.Pipeline{}
		pipelineName := types.NamespacedName{Namespace: pipeline.Namespace, Name: pipeline.Name}
		if err := r.Get(ctx, pipelineName, newPipeline); err != nil {
			log.Error(err, "Cannot get Pipeline while updating the Pipeline status: "+pipelineName.String())
			return err
		}
		if pipelineState != "" {
			newPipeline.Status.PipelineState = pipelineState
			pipeline.Status.PipelineState = pipelineState
		}
		if statusCheck != "" {
			newPipeline.Status.StatusCheck = statusCheck
			pipeline.Status.StatusCheck = statusCheck
		}
		return r.Status().Update(ctx, newPipeline)
	})
}

func (r *PipelineReconciler) Initialize(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline) error {
	log := log.FromContext(ctx)
	// TODO: Initialize the metrics exporters(YACE and etc.).
	// Get the service object for the metrics endpoint.
	if pipeline.Spec.InCluster {
		service := &corev1.Service{}
		serviceName := types.NamespacedName{Namespace: pipeline.Namespace, Name: pipeline.Spec.MetricsEndpoint.ServiceRef.Name}
		if err := r.Get(ctx, serviceName, service); err != nil {
			log.Error(err, "Cannot find metrics service")
			return err
		}
		if service.Labels == nil {
			service.Labels = make(map[string]string, 1)
		}
		// Append pipeline label for selecting.
		service.Labels[metricsLabelKey] = utils.GetNamespacedName(pipeline)
		if err := r.Update(ctx, service); err != nil {
			return err
		}
	}

	// Create the service monitor manifest.
	serviceMonitor, err := monitor.CreateServiceMonitor(pipeline)
	if err != nil {
		log.Error(err, "Cannot create service monitor")
		return err
	}
	if err := ctrl.SetControllerReference(pipeline, serviceMonitor, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, serviceMonitor); err != nil {
		log.Error(err, "Cannot create ServiceMonitor")
		return err
	}

	// Create externalName service and endpoint for the out-cluster pipeline-under-test.
	if !pipeline.Spec.InCluster {
		if err := r.CreateService(ctx, pipeline); err != nil {
			log.Error(err, "Cannot create Service")
			return err
		}
	}

	//TODO: append the metrics endpoint label to the metrics endpoint service.

	// Set the state of the pipeline to be available.
	if err := r.UpdateStatus(ctx, pipeline, windtunnelv1alpha1.PipelineAvailable, ""); err != nil {
		log.Error(err, "Cannot update the status of the Pipeline")
		return err
	}
	return nil
}

func GetMetricsServiceName(pipelineName string) string {
	return pipelineName + "-plantd-metrics"
}

func (r *PipelineReconciler) InitializeExp(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline) error {
	log := log.FromContext(ctx)
	// Get the metrics service
	serviceName := types.NamespacedName{Namespace: pipeline.Namespace, Name: GetMetricsServiceName(pipeline.Name)}
	if pipeline.Spec.InCluster {
		serviceName = types.NamespacedName{Namespace: pipeline.Spec.MetricsEndpoint.ServiceRef.Namespace, Name: pipeline.Spec.MetricsEndpoint.ServiceRef.Name}
	}
	service := &corev1.Service{}
	if err := r.Get(ctx, serviceName, service); err != nil {
		log.Error(err, "Cannot get the Service: "+serviceName.String())
		return err
	}
	if service.Labels == nil {
		service.Labels = make(map[string]string, 1)
	}
	service.Labels[experimentLabelKey] = pipeline.Spec.ExperimentRef.Name
	// TODO(DEBUG): Add the metrics label to the label set of the metrics service
	if err := r.Update(ctx, service); err != nil {
		log.Error(err, "Cannot update the Service: "+serviceName.String())
		return err
	}
	return r.UpdateStatus(ctx, pipeline, windtunnelv1alpha1.PipelineEngaged, "")
}

func (r *PipelineReconciler) HealthCheck(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline) error {
	if pipeline.Spec.HealthCheckEndpoints != nil {
		for _, endpoint := range pipeline.Spec.HealthCheckEndpoints {
			ok, err := utils.HealthCheck(endpoint)
			if err != nil || !ok {
				pipeline.Status.StatusCheck = windtunnelv1alpha1.PipelineFailed
				return err
			}
		}

	}
	return r.UpdateStatus(ctx, pipeline, "", windtunnelv1alpha1.PipelineOK)
}

// Create externalName service and endpoint for pipeline endpoints and the metrics endpoint.
func (r *PipelineReconciler) CreateService(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline) error {
	log := log.FromContext(ctx)
	// Create externalName service and endpoint for pipeline endpoints.
	for _, endpoint := range pipeline.Spec.PipelineEndpoints {
		name := pipeline.Name + "-" + endpoint.Name
		serivce, endpoints, err := monitor.CreateExternalNameService(name, pipeline.Namespace, &endpoint)
		if err != nil {
			log.Error(err, "Cannot create service manifests for outside-cluster Pipeline")
			return err
		}
		if err := ctrl.SetControllerReference(pipeline, serivce, r.Scheme); err != nil {
			return err
		}
		if err := ctrl.SetControllerReference(pipeline, endpoints, r.Scheme); err != nil {
			return err
		}
		if err := r.Create(ctx, serivce); err != nil {
			log.Error(err, "Cannot create ExternalName service for outside-cluster Pipeline")
			return err
		}
		if err := r.Create(ctx, endpoints); err != nil {
			log.Error(err, "Cannot create Endpoints for outside-cluster Pipeline")
			return err
		}
	}

	// Create externalName service and endpoint for the metrics endpoint.
	name := GetMetricsServiceName(pipeline.Name)
	serivce, endpoints, err := monitor.CreateExternalNameService(name, pipeline.Namespace, &pipeline.Spec.MetricsEndpoint)
	if err != nil {
		log.Error(err, "Cannot create service manifests for outside-cluster Pipeline")
		return err
	}
	if err := ctrl.SetControllerReference(pipeline, serivce, r.Scheme); err != nil {
		return err
	}
	if err := ctrl.SetControllerReference(pipeline, endpoints, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, serivce); err != nil {
		log.Error(err, "Cannot create ExternalName service for outside-cluster Pipeline")
		return err
	}
	if err := r.Create(ctx, endpoints); err != nil {
		log.Error(err, "Cannot create Endpoints for outside-cluster Pipeline")
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Pipeline{}).
		Complete(r)
}
