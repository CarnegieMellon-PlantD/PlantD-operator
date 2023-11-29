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

	k6v1alpha1 "github.com/grafana/k6-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/loadgen"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

const (
	ExperimentPending                 string = "Pending"                 // experiment object has been created but no pipeline has been selected/attached
	ExperimentInitializing            string = "Initializing"            // experiment is initializing after pipeline is selected
	ExperimentWaitingForPipelineReady string = "WaitingForPipelineReady" // experiment is waiting for pipeline to pass health checks
	ExperimentReady                   string = "Ready"                   // experiment is ready to be run on the selected pipeline
	ExperimentRunning                 string = "Running"                 // experiment is running
	ExperimentFinished                string = "Finished"                // experiment has finished running; i.e. the k6 load generator has finished sending data
)

// ExperimentReconciler reconciles a Experiment object
type ExperimentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=experiments/finalizers,verbs=update
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=datasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=loadpatterns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k6.io,resources=k6s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *ExperimentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Get the experiment object in the queue.
	exp := &windtunnelv1alpha1.Experiment{}
	if err := r.Get(ctx, req.NamespacedName, exp); err != nil {
		// The object is deleted.
		if apierrors.IsNotFound(err) {
			// Custom resource not found, perform cleanup tasks here.
			// Created objects are automatically garbage collected.
			log.Info("Experiment is deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Cannot get Experiment: "+req.NamespacedName.String())
		return ctrl.Result{}, err
	}

	// Get the pipeline object referred in the experiment object.
	pipeline := &windtunnelv1alpha1.Pipeline{}
	pipelineName := types.NamespacedName{Namespace: exp.Spec.PipelineRef.Namespace, Name: exp.Spec.PipelineRef.Name}
	if err := r.Get(ctx, pipelineName, pipeline); err != nil {
		exp.Status.ExperimentState = "Error: No Pipeline found"
		log.Error(err, "Cannot get Pipeline: "+pipelineName.String())
		if err := r.Status().Update(ctx, exp); err != nil {
			log.Error(err, "Cannot update the status of Experiment: "+req.NamespacedName.String())
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Set the status of experiment object to pending.
	if exp.Status.ExperimentState == "" {
		exp.Status.ExperimentState = ExperimentPending
		// Initialize the protocol for each endpoint.
		SetProtocol(exp, pipeline)
		if err := r.SetDuration(ctx, exp); err != nil {
			return ctrl.Result{}, nil
		}
		if err := r.Status().Update(ctx, exp); err != nil {
			log.Error(err, "Cannot update the status of Experiment: "+req.NamespacedName.String())
			return ctrl.Result{}, err
		}
	}

	// Check to see if the scheduled time has been reached.
	currTime := time.Now()
	if currTime.Before(exp.Spec.ScheduledTime.Time) {
		return ctrl.Result{RequeueAfter: exp.Spec.ScheduledTime.Time.Sub(currTime)}, nil
	}

	// Initialize the infrastructure of the experiment.
	if err := r.InitializeExp(ctx, exp, pipeline); err != nil {
		return ctrl.Result{}, err
	}

	// Get the latest pipeline object.
	if err := r.Get(ctx, pipelineName, pipeline); err != nil {
		exp.Status.ExperimentState = "Error: No Pipeline found"
		log.Error(err, "Cannot get Pipeline: "+pipelineName.String())
		if err := r.Status().Update(ctx, exp); err != nil {
			log.Error(err, "Cannot update the status of Experiment: "+req.NamespacedName.String())
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Pipeline health check.
	if exp.Status.ExperimentState == ExperimentWaitingForPipelineReady {
		if pipeline.Status.StatusCheck == windtunnelv1alpha1.PipelineFailed {
			exp.Status.ExperimentState = "Error: Pipeline is unhealthy"
			return ctrl.Result{}, r.Status().Update(ctx, exp)
		} else if pipeline.Status.StatusCheck == windtunnelv1alpha1.PipelineOK && pipeline.Status.PipelineState == windtunnelv1alpha1.PipelineEngaged {
			exp.Status.ExperimentState = ExperimentReady
			if err := r.Status().Update(ctx, exp); err != nil {
				log.Error(err, "Cannot update the status of Experiment: "+req.NamespacedName.String())
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{RequeueAfter: time.Second * 3}, nil
		}
	}

	// Create the load generator(testRun).
	if err := r.CreateTestRun(ctx, exp, pipeline); err != nil {
		log.Info(err.Error())
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	// Check if the experiment is finished.
	if err := r.CheckFinished(ctx, exp); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ExperimentReconciler) SetDuration(ctx context.Context, exp *windtunnelv1alpha1.Experiment) error {
	log := log.FromContext(ctx)

	n := len(exp.Spec.LoadPatterns)
	exp.Status.Duration = make(map[string]metav1.Duration, n)
	for _, loadConfig := range exp.Spec.LoadPatterns {
		load := &windtunnelv1alpha1.LoadPattern{}
		loadName := types.NamespacedName{
			Namespace: GetLoadPatternNamespace(exp, &loadConfig.LoadPatternRef),
			Name:      loadConfig.LoadPatternRef.Name,
		}
		if err := r.Get(ctx, loadName, load); err != nil {
			log.Error(err, "Cannot get LoadPattern: "+loadName.String())
			return err
		}
		sum := time.Duration(0)
		for _, stage := range load.Spec.Stages {
			d, err := time.ParseDuration(stage.Duration)
			if err != nil {
				log.Error(err, "Cannot parse duration: "+stage.Duration)
				exp.Status.ExperimentState = "Error: Cannot calculate duration"
				r.Status().Update(ctx, exp)
				return err
			}
			sum += d
		}
		exp.Status.Duration[loadConfig.EndpointName] = metav1.Duration{Duration: sum}
	}
	return nil
}

func SetProtocol(exp *windtunnelv1alpha1.Experiment, pipeline *windtunnelv1alpha1.Pipeline) {
	n := len(pipeline.Spec.PipelineEndpoints)
	exp.Status.Protocols = make(map[string]string, n)

	for _, endpoint := range pipeline.Spec.PipelineEndpoints {
		http := endpoint.HTTP.URL != ""
		websocket := endpoint.WebSocket.URL != ""
		grpc := endpoint.GRPC.URL != ""
		var dataOption string
		if endpoint.HTTP.Body.DataSetRef.Name != "" {
			dataOption = windtunnelv1alpha1.WithDataSet
		} else {
			dataOption = windtunnelv1alpha1.WithData
		}

		if http && !websocket && !grpc {
			exp.Status.Protocols[endpoint.Name] = fmt.Sprintf("%s.%s", windtunnelv1alpha1.ProtocolHTTP, dataOption)
		} else if !http && websocket && !grpc {
			exp.Status.Protocols[endpoint.Name] = fmt.Sprintf("%s.%s", windtunnelv1alpha1.ProtocolWebSocket, dataOption)
		} else if !http && !websocket && grpc {
			exp.Status.Protocols[endpoint.Name] = fmt.Sprintf("%s.%s", windtunnelv1alpha1.ProtocolGRPC, dataOption)
		} else {
			exp.Status.ExperimentState = "Error: Invalid protocol"
		}
	}
}

func (r *ExperimentReconciler) InitializeExp(ctx context.Context, exp *windtunnelv1alpha1.Experiment, pipeline *windtunnelv1alpha1.Pipeline) error {
	// (Idempotency) If the status of the experiment is not pending, do nothing.
	if exp.Status.ExperimentState != ExperimentPending {
		return nil
	}
	if pipeline.Status.PipelineState != windtunnelv1alpha1.PipelineAvailable {
		exp.Status.ExperimentState = "Error: Pipeline is not available"
		return r.Status().Update(ctx, exp)
	}
	log := log.FromContext(ctx)
	exp.Status.ExperimentState = ExperimentInitializing
	if err := r.Status().Update(ctx, exp); err != nil {
		log.Error(err, "Cannot update the status of Experiment: "+exp.Name)
		return err
	}

	experimentRef := corev1.ObjectReference{
		APIVersion: exp.APIVersion,
		Kind:       exp.Kind,
		Namespace:  exp.Namespace,
		Name:       exp.Name,
		UID:        exp.UID,
	}
	// Set the ExperimentRef field in pipeline spec.
	if err := r.UpdatePipeline(ctx, pipeline, experimentRef); err != nil {
		log.Error(err, "Cannot update Pipeline")
		return err
	}

	exp.Status.ExperimentState = ExperimentWaitingForPipelineReady
	exp.Status.Tags = pipeline.Spec.ExtraMetrics.System.Tags
	exp.Status.CloudVendor = pipeline.Spec.CloudVendor
	exp.Status.EnableCostCalculation = pipeline.Spec.EnableCostCalculation
	return r.Status().Update(ctx, exp)
}

func (r *ExperimentReconciler) CreateTestRun(ctx context.Context, exp *windtunnelv1alpha1.Experiment, pipeline *windtunnelv1alpha1.Pipeline) error {
	if exp.Status.ExperimentState != ExperimentReady {
		return nil
	}
	for _, endpoint := range pipeline.Spec.PipelineEndpoints {
		var err error
		if requireDataSet(&endpoint) {
			err = r.CreateTestRunWithDataSet(ctx, exp, &endpoint)
		} else {
			err = r.CreateTestRunWithOutDataSet(ctx, exp, &endpoint)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func requireDataSet(endpoint *windtunnelv1alpha1.Endpoint) bool {
	return endpoint.HTTP.Body.DataSetRef.Name != ""
}

func FindLoadConfig(endpointName string, exp *windtunnelv1alpha1.Experiment) *windtunnelv1alpha1.LoadPatternConfig {
	for _, config := range exp.Spec.LoadPatterns {
		if config.EndpointName == endpointName {
			return &config
		}
	}
	return nil
}

func GetLoadPatternNamespace(exp *windtunnelv1alpha1.Experiment, ref *corev1.ObjectReference) string {
	if ref.Namespace != "" {
		return ref.Namespace
	}
	return exp.Namespace
}

func (r *ExperimentReconciler) CreateTestRunWithDataSet(ctx context.Context, exp *windtunnelv1alpha1.Experiment, endpoint *windtunnelv1alpha1.Endpoint) error {
	log := log.FromContext(ctx)

	targetEndpointName := endpoint.Name
	loadpatternConfig := FindLoadConfig(targetEndpointName, exp)
	if loadpatternConfig == nil {
		exp.Status.ExperimentState = "Error: No LoadPatternConfig found"
		log.Info("Cannot find Load Configuration for endpoint: " + targetEndpointName)
		return r.Status().Update(ctx, exp)
	}
	load := &windtunnelv1alpha1.LoadPattern{}

	loadName := types.NamespacedName{
		Namespace: GetLoadPatternNamespace(exp, &loadpatternConfig.LoadPatternRef),
		Name:      loadpatternConfig.LoadPatternRef.Name,
	}
	if err := r.Get(ctx, loadName, load); err != nil {
		exp.Status.ExperimentState = "Error: No LoadPattern found"
		log.Error(err, "Cannot get LoadPattern: "+loadName.String())
		return r.Status().Update(ctx, exp)
	}

	dataset := &windtunnelv1alpha1.DataSet{}
	datasetName := types.NamespacedName{
		Namespace: endpoint.HTTP.Body.DataSetRef.Namespace,
		Name:      endpoint.HTTP.Body.DataSetRef.Name,
	}
	if err := r.Get(ctx, datasetName, dataset); err != nil {
		exp.Status.ExperimentState = "Error: No DataSet found"
		log.Error(err, "Cannot get DataSet: "+datasetName.String())
		return r.Status().Update(ctx, exp)
	}

	if err := r.CreateConfigMapWithDataSet(ctx, exp, endpoint, load, dataset); err != nil {
		exp.Status.ExperimentState = "Error: Failed to create configMap"
		return r.Status().Update(ctx, exp)
	}
	// TODO: Fix the config map name
	testRunName := utils.GetTestRunName(exp.Name, targetEndpointName)
	configMapName := types.NamespacedName{Namespace: exp.Namespace, Name: testRunName}
	configMap := &corev1.ConfigMap{}
	if err := r.Get(ctx, configMapName, configMap); err != nil {
		log.Error(err, "Cannot get ConfigMap: "+configMapName.String())
		return err
	}

	if err := r.CopyConfigMap(ctx, dataset, configMap, exp); err != nil {
		return err
	}
	if err := r.CreateK6WithDataSet(ctx, exp, dataset, targetEndpointName); err != nil {
		return err
	}

	if exp.Status.ExperimentState == ExperimentReady {
		exp.Status.ExperimentState = ExperimentRunning
		currTime := &metav1.Time{Time: time.Now()}
		exp.Status.StartTime = currTime
	}

	return r.Status().Update(ctx, exp)
}

func (r *ExperimentReconciler) CreateTestRunWithOutDataSet(ctx context.Context, exp *windtunnelv1alpha1.Experiment, endpoint *windtunnelv1alpha1.Endpoint) error {
	log := log.FromContext(ctx)
	targetEndpointName := endpoint.Name
	loadpatternConfig := FindLoadConfig(targetEndpointName, exp)
	if loadpatternConfig == nil {
		exp.Status.ExperimentState = "Error: No LoadPatternConfig found"
		log.Info("Cannot find Load Configuration for endpoint: " + targetEndpointName)
		return r.Status().Update(ctx, exp)
	}
	load := &windtunnelv1alpha1.LoadPattern{}

	loadName := types.NamespacedName{
		Namespace: GetLoadPatternNamespace(exp, &loadpatternConfig.LoadPatternRef),
		Name:      loadpatternConfig.LoadPatternRef.Name,
	}
	if err := r.Get(ctx, loadName, load); err != nil {
		exp.Status.ExperimentState = "Error: No LoadPattern found"
		log.Error(err, "Cannot get LoadPattern: "+loadName.String())
		return r.Status().Update(ctx, exp)
	}

	if err := r.CreateConfigMap(ctx, exp, endpoint, load); err != nil {
		return err
	}
	if err := r.CreateK6(ctx, exp, targetEndpointName); err != nil {
		return err
	}

	if exp.Status.ExperimentState == ExperimentReady {
		exp.Status.ExperimentState = ExperimentRunning
		currTime := &metav1.Time{Time: time.Now()}
		exp.Status.StartTime = currTime
	}
	return r.Status().Update(ctx, exp)
}

func (r *ExperimentReconciler) CreateConfigMap(ctx context.Context, exp *windtunnelv1alpha1.Experiment, endpoint *windtunnelv1alpha1.Endpoint, load *windtunnelv1alpha1.LoadPattern) error {
	log := log.FromContext(ctx)
	testRunName := utils.GetTestRunName(exp.Name, endpoint.Name)
	configName := types.NamespacedName{Namespace: exp.Namespace, Name: testRunName}
	config := &corev1.ConfigMap{}
	if err := r.Get(ctx, configName, config); err == nil {
		return nil
	}
	configMap, err := loadgen.CreateConfigMap(testRunName, exp.Status.Protocols[endpoint.Name], exp, endpoint, load)
	if err != nil {
		exp.Status.ExperimentState = "Error: Cannot create ConfigMap Manifest"
		log.Error(err, "Cannot create ConfigMap Manifest")
		return nil
	}
	if err := ctrl.SetControllerReference(exp, configMap, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, configMap); err != nil {
		exp.Status.ExperimentState = "Error: Cannot create ConfigMap"
		log.Error(err, "Cannot create ConfigMap")
		return nil
	}
	return nil
}

func (r *ExperimentReconciler) CreateConfigMapWithDataSet(ctx context.Context, exp *windtunnelv1alpha1.Experiment, endpoint *windtunnelv1alpha1.Endpoint, load *windtunnelv1alpha1.LoadPattern, dataset *windtunnelv1alpha1.DataSet) error {
	log := log.FromContext(ctx)
	testRunName := utils.GetTestRunName(exp.Name, endpoint.Name)
	configName := types.NamespacedName{Namespace: exp.Namespace, Name: testRunName}
	config := &corev1.ConfigMap{}
	if err := r.Get(ctx, configName, config); err == nil {
		return nil
	}
	configMap, err := loadgen.CreateConfigMapWithDataSet(testRunName, exp.Status.Protocols[endpoint.Name], exp, endpoint, load, dataset)
	if err != nil {
		exp.Status.ExperimentState = "Error: Cannot create ConfigMap Manifest"
		log.Error(err, "Cannot create ConfigMap Manifest")
		return nil
	}
	if err := ctrl.SetControllerReference(exp, configMap, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, configMap); err != nil {
		exp.Status.ExperimentState = "Error: Cannot create ConfigMap"
		log.Error(err, "Cannot create ConfigMap")
		return nil
	}
	return nil
}

func (r *ExperimentReconciler) CopyConfigMap(ctx context.Context, dataset *windtunnelv1alpha1.DataSet, configMap *corev1.ConfigMap, exp *windtunnelv1alpha1.Experiment) error {
	log := log.FromContext(ctx)
	pod := loadgen.CreateCopyPod(dataset, configMap)
	podName := types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}
	if err := r.Get(ctx, podName, pod); err == nil {
		if pod.Status.Phase == corev1.PodSucceeded {
			return nil
		}
		return fmt.Errorf("copying pod is running")
	}
	if err := ctrl.SetControllerReference(exp, pod, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, pod); err != nil {
		exp.Status.ExperimentState = "Error: Cannot copy files from ConfigMap to PVC"
		log.Error(err, "Cannot copy files from ConfigMap to PVC")
	}
	// TODO: Handle waiting for copying to be successed
	return fmt.Errorf("copying pod is running")
}

func (r *ExperimentReconciler) CreateK6(ctx context.Context, exp *windtunnelv1alpha1.Experiment, endpointName string) error {
	log := log.FromContext(ctx)
	testRunName := utils.GetTestRunName(exp.Name, endpointName)
	k6Name := types.NamespacedName{Namespace: exp.Namespace, Name: testRunName}
	k6 := &k6v1alpha1.K6{}
	if err := r.Get(ctx, k6Name, k6); err == nil {
		return nil
	}
	testRun := loadgen.CreateTestRunManifest(exp.Namespace, exp.Name, endpointName)

	testRun.Spec.Script = k6v1alpha1.K6Script{
		ConfigMap: k6v1alpha1.K6Configmap{
			Name: testRunName,
			File: config.GetString("k6.config.script"),
		},
	}
	if err := ctrl.SetControllerReference(exp, testRun, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, testRun); err != nil {
		exp.Status.ExperimentState = "Error: Cannot create K6 TestRun"
		log.Error(err, "Cannot create K6 TestRun: "+testRunName)
		return err
	}
	return nil
}

func (r *ExperimentReconciler) CreateK6WithDataSet(ctx context.Context, exp *windtunnelv1alpha1.Experiment, dataset *windtunnelv1alpha1.DataSet, endpointName string) error {
	log := log.FromContext(ctx)
	k6Name := types.NamespacedName{Namespace: exp.Namespace, Name: exp.Name}
	k6 := &k6v1alpha1.K6{}
	if err := r.Get(ctx, k6Name, k6); err == nil {
		return nil
	}
	testRun := loadgen.CreateTestRunManifest(exp.Namespace, exp.Name, endpointName)

	testRun.Spec.Script = k6v1alpha1.K6Script{
		VolumeClaim: k6v1alpha1.K6VolumeClaim{
			Name: utils.GetPVCName(dataset.Name, dataset.Generation),
			File: config.GetString("k6.config.script"), // TODO: Streamline the reading from configuration file.
		},
	}
	if err := ctrl.SetControllerReference(exp, testRun, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, testRun); err != nil {
		exp.Status.ExperimentState = "Error: Cannot create K6 TestRun"
		log.Error(err, "Cannot create K6 TestRun: "+exp.Name)
		return nil
	}
	return nil
}

// Check if the test run is finished
func (r *ExperimentReconciler) CheckFinished(ctx context.Context, exp *windtunnelv1alpha1.Experiment) error {
	if exp.Status.ExperimentState != ExperimentRunning {
		return nil
	}
	log := log.FromContext(ctx)
	successCounter := 0
	for _, config := range exp.Spec.LoadPatterns {
		testRun := &k6v1alpha1.K6{}
		testRunName := types.NamespacedName{Namespace: exp.Namespace, Name: utils.GetTestRunName(exp.Name, config.EndpointName)}
		if err := r.Get(ctx, testRunName, testRun); err != nil {
			exp.Status.ExperimentState = "Error: Cannot get K6 TestRun"
			log.Error(err, "Cannot get K6 TestRun: "+testRunName.String())
		}
		if testRun.Status.Stage == "finished" {
			successCounter += 1
		}
	}

	if successCounter == len(exp.Spec.LoadPatterns) {
		exp.Status.ExperimentState = ExperimentFinished
	}

	if err := r.Status().Update(ctx, exp); err != nil {
		log.Error(err, "Cannot update the status of Experiment")
		return err
	}

	if exp.Status.ExperimentState == ExperimentFinished {
		pipeline := &windtunnelv1alpha1.Pipeline{}
		pipelineName := types.NamespacedName{Namespace: exp.Spec.PipelineRef.Namespace, Name: exp.Spec.PipelineRef.Name}
		if err := r.Get(ctx, pipelineName, pipeline); err == nil {
			return r.UpdatePipeline(ctx, pipeline, corev1.ObjectReference{})
		}
	}
	return nil
}

func (r *ExperimentReconciler) UpdatePipeline(ctx context.Context, pipeline *windtunnelv1alpha1.Pipeline, expRef corev1.ObjectReference) error {
	log := log.FromContext(ctx)
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		newPipeline := &windtunnelv1alpha1.Pipeline{}
		pipelineName := types.NamespacedName{Namespace: pipeline.Namespace, Name: pipeline.Name}
		if err := r.Get(ctx, pipelineName, newPipeline); err != nil {
			log.Error(err, "Cannot get Pipeline while updating the Pipeline status: "+pipelineName.String())
			return err
		}
		newPipeline.Spec.ExperimentRef = expRef
		return r.Update(ctx, newPipeline)
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExperimentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.Experiment{}).
		Owns(&k6v1alpha1.K6{}).
		Complete(r)
}
