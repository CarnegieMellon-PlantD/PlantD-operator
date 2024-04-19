package core

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
)

// ServiceAccount
var (
	serviceAccountControllerManager = config.GetString("rbac.serviceAccount.controllerManager")
	serviceAccountPrometheus        = config.GetString("rbac.serviceAccount.prometheus")
	serviceAccountOpenCost          = config.GetString("rbac.serviceAccount.openCost")
)

// getResources accepts a resource requirements and returns a new resource requirements with the default values.
func getResources(in *corev1.ResourceRequirements, defaultCPURequest, defaultMemoryRequest, defaultCPULimit, defaultMemoryLimit string) corev1.ResourceRequirements {
	out := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(defaultCPURequest),
			corev1.ResourceMemory: resource.MustParse(defaultMemoryRequest),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(defaultCPULimit),
			corev1.ResourceMemory: resource.MustParse(defaultMemoryLimit),
		},
	}

	if in == nil {
		return out
	}

	if in.Requests != nil {
		if !in.Requests.Cpu().IsZero() {
			out.Requests[corev1.ResourceCPU] = in.Requests[corev1.ResourceCPU]
		}
		if !in.Requests.Memory().IsZero() {
			out.Requests[corev1.ResourceMemory] = in.Requests[corev1.ResourceMemory]
		}
	}
	if in.Limits != nil {
		if !in.Limits.Cpu().IsZero() {
			out.Limits[corev1.ResourceCPU] = in.Limits[corev1.ResourceCPU]
		}
		if !in.Limits.Memory().IsZero() {
			out.Limits[corev1.ResourceMemory] = in.Limits[corev1.ResourceMemory]
		}
	}

	return out
}
