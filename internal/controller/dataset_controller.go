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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/datagen"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"
)

// DataSetReconciler reconciles a DataSet object
type DataSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	errorTypeLogContainer = "container"
	errorTypeLogPod       = "pod"
)
const (
	CREATING   = "Creating"
	GENERATING = "Generating"
	RUNNING    = "Running"
	FAILED     = "Failed"
	SUCCESS    = "Success"
	UNKNOWN    = "Unknown"
)

//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=datasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=datasets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=datasets/finalizers,verbs=update
//+kubebuilder:rbac:groups=windtunnel.plantd.org,resources=schemas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods/log,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *DataSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	/*
		### 1: Get the DataSet by name
		We'll fetch the CronJob using our client.  All client methods take a
		context (to allow for cancellation) as their first argument, and the object
		in question as their last.  Get is a bit special, in that it takes a
		[`NamespacedName`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client?tab=doc#ObjectKey)
		as the middle argument.
		Many client methods also take variadic options at the end.
	*/
	dataSet := &windtunnelv1alpha1.DataSet{}
	if err := r.Get(ctx, req.NamespacedName, dataSet); err != nil {
		if apierrors.IsNotFound(err) {
			// Custom resource not found, perform cleanup tasks here.
			// Created objects are automatically garbage collected.
			log.Info("DataSet is deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to fetch DataSet.")
		return ctrl.Result{}, err
	}

	schemaMap := make(map[string]*windtunnelv1alpha1.Schema, len(dataSet.Spec.Schemas))
	for _, schema := range dataSet.Spec.Schemas {
		s := &windtunnelv1alpha1.Schema{}
		schemaName := types.NamespacedName{Namespace: req.Namespace, Name: schema.Name}
		if err := r.Get(ctx, schemaName, s); err != nil {
			log.Error(err, "Cannot get Schema: "+schemaName.String())
			r.updateErrors(dataSet, errorTypeLogContainer, err.Error())
		}
		schemaMap[schema.Name] = s
	}

	if dataSet.GetDeletionTimestamp() != nil {
		log.Info("Delete DataSet")
		return ctrl.Result{}, nil
	}

	if dataSet.Status.StartTime == nil {
		dataSet.Status.StartTime = &metav1.Time{}
	}
	if dataSet.Status.JobStatus == "" {
		dataSet.Status.JobStatus = UNKNOWN
		dataSet.Status.PVCStatus = UNKNOWN
	}
	if err := r.Status().Update(ctx, dataSet); err != nil {
		log.Error(err, "Cannot update the status of DataSet.")
	}

	/*
		### 2: Create or update pvc and indexedJob
		dataSet.Generation is used to track the change in dataSet.Spec
		once the spec is created/updated, we create new pvc & indexJob, delete the old pvs & indexJob
	*/
	if dataSet.Status.LastGeneration != dataSet.Generation {
		log.Info("Create/Update DataSet")

		// Create a new pvc
		pvcName := utils.GetPVCName(dataSet.Name, dataSet.Generation)
		log.Info("Create a new pvc")
		newPVC := datagen.CreatePVC(types.NamespacedName{Name: pvcName, Namespace: dataSet.Namespace})
		if err := ctrl.SetControllerReference(dataSet, newPVC, r.Scheme); err != nil {
			log.Error(err, "Cannot set owner controller.")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, newPVC); err != nil {
			log.Error(err, "Cannot create persistent volume claim.")
			r.updateErrors(dataSet, errorTypeLogContainer, err.Error())
		}
		dataSet.Status.PVCStatus = string(newPVC.Status.Phase)

		// Create a new indexedJob
		jobName := utils.GetJobName(dataSet.Name, "parallel-gen-data", dataSet.Generation)
		newIndexedJob, err := datagen.CreateJobByDataSet(jobName, pvcName, dataSet, schemaMap)
		if err != nil {
			log.Error(err, "Cannot create job.")
			r.updateErrors(dataSet, errorTypeLogContainer, err.Error())

		}
		if err := ctrl.SetControllerReference(dataSet, newIndexedJob, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, newIndexedJob); err != nil {
			log.Error(err, "Cannot create indexed job.")
			r.updateErrors(dataSet, errorTypeLogContainer, err.Error())
		}
		log.Info("Create an indexed parallel job.")
		dataSet.Status.StartTime = &newIndexedJob.ObjectMeta.CreationTimestamp
		dataSet.Status.JobStatus = GENERATING

		// Delete the job of last generation if exists one
		lastIndexedJob := &kbatch.Job{}
		lastJobName := utils.GetJobName(dataSet.Name, "parallel-gen-data", dataSet.Status.LastGeneration)
		if err := r.Get(ctx, client.ObjectKey{Name: lastJobName, Namespace: dataSet.Namespace}, lastIndexedJob); err == nil {
			if lastIndexedJob.GetDeletionTimestamp() == nil {
				log.Info("Job of last generation still alive, proceed to kill.")
				propagationPolicy := metav1.DeletePropagationBackground
				if err := r.Delete(ctx, lastIndexedJob, &client.DeleteOptions{PropagationPolicy: &propagationPolicy}); err != nil {
					log.Error(err, "Cannot kill the current running jobs.")
					return ctrl.Result{}, err
				}
				log.Info("Job of last generation deleted.")
			}
		}

		// Delete the pvc of last generation if exists one
		lastPvc := &corev1.PersistentVolumeClaim{}
		lastPvcName := utils.GetPVCName(dataSet.Name, dataSet.Status.LastGeneration)
		if err := r.Get(ctx, client.ObjectKey{Namespace: dataSet.Namespace, Name: lastPvcName}, lastPvc); err == nil {
			var refPods int64
			for refPods > 0 {
				refPods = 0
				podList := &corev1.PodList{}
				if err := r.List(ctx, podList, client.InNamespace(dataSet.Namespace)); err != nil {
					log.Error(err, "Cannot list pods using pvc.")
					r.updateErrors(dataSet, errorTypeLogContainer, err.Error())
				}
				for _, pod := range podList.Items {
					for _, volume := range pod.Spec.Volumes {
						if volume.VolumeSource.PersistentVolumeClaim != nil && volume.VolumeSource.PersistentVolumeClaim.ClaimName == lastPvc.Name {
							log.Info("There's still pod referencing the pvc")
							refPods = refPods + 1
							break
						}
					}
					if refPods > 0 {
						break
					}
				}
			}
			if err := r.Delete(ctx, lastPvc); err != nil {
				log.Error(err, "Cannot delete pvc.")
				r.updateErrors(dataSet, errorTypeLogContainer, err.Error())
			}
			log.Info("PVC of last generation deleted.")
		}

		dataSet.Status.LastGeneration = dataSet.Generation

		if err := r.Status().Update(ctx, dataSet); err != nil {
			log.Error(err, "Cannot update the status of DataSet.")
			r.updateErrors(dataSet, errorTypeLogContainer, err.Error())
		}
	}

	/*
		### 3: Fetch the current pvc & indexJob and update status
	*/

	jobName := utils.GetJobName(dataSet.Name, "parallel-gen-data", dataSet.Generation)
	indexedJob := &kbatch.Job{}
	if err := r.Get(ctx, client.ObjectKey{Name: jobName, Namespace: dataSet.Namespace}, indexedJob); err != nil {
		log.Error(err, "Cannot get indexed job")
		return ctrl.Result{}, err
	}

	pvcName := utils.GetPVCName(dataSet.Name, dataSet.Generation)
	pvc := &corev1.PersistentVolumeClaim{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: dataSet.Namespace, Name: pvcName}, pvc); err != nil {
		log.Error(err, "Cannot get pvc")
		return ctrl.Result{}, err
	}
	dataSet.Status.PVCStatus = string(pvc.Status.Phase)
	if err := r.Status().Update(ctx, dataSet); err != nil {
		log.Error(err, "Cannot update the status of DataSet 3.")
		return ctrl.Result{}, err
	}

	/*
		### 4: Check if the job has completed
	*/
	ok, conditionType := IsJobFinished(indexedJob)
	if ok {
		switch conditionType {
		case kbatch.JobComplete:
			dataSet.Status.JobStatus = SUCCESS
			log.Info("Job Completed.")

		case kbatch.JobFailed:
			log.Info("Job failed.")
			dataSet.Status.JobStatus = FAILED
			// Get pod logs
			podLogs, err := r.getPodLogs(ctx, indexedJob)
			if err != nil {
				log.Error(err, "Cannot get pod logs")
				r.updateErrors(dataSet, errorTypeLogContainer, err.Error())
			} else {
				// Set errorString to pod logs
				r.updateErrors(dataSet, errorTypeLogPod, podLogs)
			}
		}
		dataSet.Status.CompletionTime = indexedJob.Status.CompletionTime
		if err := r.Status().Update(ctx, dataSet); err != nil {
			log.Error(err, "Cannot update the status of DataSet.")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	log.Info("Job is running.")
	return ctrl.Result{Requeue: true, RequeueAfter: time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DataSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windtunnelv1alpha1.DataSet{}).
		Complete(r)
}

func (r *DataSetReconciler) updateErrors(dataset *windtunnelv1alpha1.DataSet, errorType string, err string) {
	if dataset.Status.Errors == nil {
		dataset.Status.Errors = make(map[string][]string)
	}

	// Check if the error already exists in the list
	if !containsError(dataset.Status.Errors[errorType], err) {
		dataset.Status.ErrorCount++
		dataset.Status.Errors[errorType] = append(dataset.Status.Errors[errorType], err)
	}
}

// Helper function to check if an error already exists in the list
func containsError(errors []string, err string) bool {
	for _, e := range errors {
		if e == err {
			return true
		}
	}
	return false
}

func IsJobFinished(job *kbatch.Job) (bool, kbatch.JobConditionType) {
	for _, c := range job.Status.Conditions {
		if (c.Type == kbatch.JobComplete || c.Type == kbatch.JobFailed) && c.Status == corev1.ConditionTrue {
			return true, c.Type
		}
	}
	return false, ""
}

func (r *DataSetReconciler) getPodLogs(ctx context.Context, job *kbatch.Job) (string, error) {
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(job.Namespace)); err != nil {
		return "", fmt.Errorf("failed to list pods: %w", err)
	}

	for _, pod := range podList.Items {
		// Check if the pod belongs to the job
		if metav1.IsControlledBy(&pod, job) {
			// Create a new context to ensure it is not affected by the parent context cancellation
			// this is to avoid any rate limiting issues from kubernetes API
			podLogCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			// Get the pod logs using the new context
			logs, err := r.getPodLogsByPodName(podLogCtx, pod.Namespace, pod.Name)
			if err != nil {
				return "", fmt.Errorf("failed to get pod logs: %w", err)
			}
			return logs, nil
		}
	}

	return "", fmt.Errorf("no pod found for job: %s/%s", job.Namespace, job.Name)
}

func (r *DataSetReconciler) getPodLogsByPodName(ctx context.Context, namespace, podName string) (string, error) {
	// pod := &corev1.Pod{}
	// err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: podName}, pod)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to get pod: %w", err)
	// }

	// if len(pod.Spec.Containers) == 0 {
	// 	return "", fmt.Errorf("no containers found in pod")
	// }

	// containerName := pod.Spec.Containers[0].Name

	// podLogOpts := &corev1.PodLogOptions{
	// 	Container: containerName,
	// }

	// podLogRequest := r.K8sClient.CoreV1().Pods(namespace).GetLogs(podName, podLogOpts)
	// podLogs, err := podLogRequest.Stream(ctx)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to open log stream for pod: %w", err)
	// }
	// defer podLogs.Close()

	// buf := new(bytes.Buffer)
	// _, err = io.Copy(buf, podLogs)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to read pod logs: %w", err)
	// }

	// return buf.String(), nil
	return "", nil
}
