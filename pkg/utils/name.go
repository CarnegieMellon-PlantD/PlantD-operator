package utils

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetPVCName(ownerName string, Generation int64) string {
	return fmt.Sprintf("%s-%d-pvc", ownerName, Generation)
}

func GetJobName(ownerName string, funcName string, Generation int64) string {
	return fmt.Sprintf("%s-%s-%d-job", ownerName, funcName, Generation)
}

func GetVolumeName(ownerName string) string {
	return fmt.Sprintf("%s-volume", ownerName)
}

func GetNamespacedName(obj client.Object) string {
	return fmt.Sprintf("%s-%s", obj.GetNamespace(), obj.GetName())
}

func GetTestRunName(expName string, endpointName string) string {
	return fmt.Sprintf("%s-%s", expName, endpointName)
}
