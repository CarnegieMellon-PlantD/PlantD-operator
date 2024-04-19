package v1alpha1

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ComponentStatusText defines the status of a component.
type ComponentStatusText string

const (
	ComponentReady    ComponentStatusText = "Ready"
	ComponentNotReady ComponentStatusText = "Not Ready"
	ComponentSkipped  ComponentStatusText = "Skipped"
)

// DeploymentConfig defines the desired state of a component deployed as Deployment.
type DeploymentConfig struct {
	// Number of replicas.
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`
	// Container image to use.
	Image string `json:"image,omitempty"`
	// Resources requirements.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// StatefulSetConfig defines the desired state of a component deployed as StatefulSet.
type StatefulSetConfig struct {
	// Number of replicas.
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`
	// Container image to use.
	Image string `json:"image,omitempty"`
	// Resources requirements.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Storage size.
	// Should not be modified once set.
	StorageSize *resource.Quantity `json:"storageSize,omitempty"`
}

// PrometheusConfig defines the desired state of a Prometheus component.
type PrometheusConfig struct {
	// Number of replicas.
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`
	// Interval at which Prometheus scrapes metrics.
	ScrapeInterval monitoringv1.Duration `json:"scrapeInterval,omitempty"`
	// Resources requirements.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// ThanosConfig defines the desired state of a Thanos component.
type ThanosConfig struct {
	// Thanos image to use. Must be synced with the `version` field.
	Image string `json:"image,omitempty"`
	// Thanos version to use. Must be synced with the `image` field.
	Version string `json:"version,omitempty"`
	// Object store configuration for Thanos.
	// Set this field will enable upload in Thanos-Sidecar, and deploy Thanos-Store and Thanos-Compactor.
	ObjectStoreConfig *corev1.SecretKeySelector `json:"objectStoreConfig,omitempty"`
	// Thanos-Sidecar configuration.
	// The `sidecar.replicas`, `sidecar.image` and `sidecar.storageSize` fields are always ignored.
	SidecarConfig StatefulSetConfig `json:"sidecar,omitempty"`
	// Thanos-Store configuration.
	// The `store.image` field is always ignored.
	// This field is ignored if `objectStoreConfig` is not set.
	StoreConfig StatefulSetConfig `json:"store,omitempty"`
	// Thanos-Compactor configuration.
	// The `compactor.image` field is always ignored.
	// This field is ignored if `objectStoreConfig` is not set.
	CompactorConfig StatefulSetConfig `json:"compactor,omitempty"`
	// Thanos-Querier configuration.
	// The `querier.image` and `querier.storageSize` fields are always ignored.
	QuerierConfig StatefulSetConfig `json:"querier,omitempty"`
}

// OpenCostConfig defines the desired state of an OpenCost component.
type OpenCostConfig struct {
	// Number of replicas.
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`
	// Container image to use for OpenCost.
	Image string `json:"image,omitempty"`
	// Resources requirements for OpenCost.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Container image to use for OpenCost-UI.
	UIImage string `json:"uiImage,omitempty"`
	// Resources requirements for OpenCost-UI.
	UIResources corev1.ResourceRequirements `json:"uiResources,omitempty"`
}

// ComponentStatus defines the status of a component.
type ComponentStatus struct {
	// Component status string.
	Text ComponentStatusText `json:"text,omitempty"`
	// Number of ready replicas.
	NumReady int32 `json:"numReady,omitempty"`
	// Number of desired replicas.
	NumDesired int32 `json:"numDesired,omitempty"`
}

// PlantDCoreSpec defines the desired state of PlantDCore.
type PlantDCoreSpec struct {
	// PlantD-Proxy configuration.
	ProxyConfig DeploymentConfig `json:"proxy,omitempty"`
	// PlantD-Studio configuration.
	StudioConfig DeploymentConfig `json:"studio,omitempty"`
	// Prometheus configuration.
	PrometheusConfig PrometheusConfig `json:"prometheus,omitempty"`
	// Thanos configuration.
	ThanosConfig ThanosConfig `json:"thanos,omitempty"`
	// Redis configuration.
	RedisConfig StatefulSetConfig `json:"redis,omitempty"`
	// OpenCost configuration.
	OpenCostConfig OpenCostConfig `json:"opencost,omitempty"`
}

// PlantDCoreStatus defines the observed state of PlantDCore.
type PlantDCoreStatus struct {
	// PlantD-Proxy status.
	ProxyStatus ComponentStatus `json:"proxyStatus,omitempty"`
	// PlantD-Studio status.
	StudioStatus ComponentStatus `json:"studioStatus,omitempty"`
	// Prometheus status.
	PrometheusStatus ComponentStatus `json:"prometheusStatus,omitempty"`
	// Thanos-Store status.
	ThanosStoreStatus ComponentStatus `json:"thanosStoreStatus,omitempty"`
	// Thanos-Compactor status.
	ThanosCompactorStatus ComponentStatus `json:"thanosCompactorStatus,omitempty"`
	// Thanos-Querier status.
	ThanosQuerierStatus ComponentStatus `json:"thanosQuerierStatus,omitempty"`
	// Redis status.
	RedisStatus ComponentStatus `json:"redisStatus,omitempty"`
	// OpenCost status.
	OpenCostStatus ComponentStatus `json:"opencostStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Proxy",type="string",JSONPath=".status.proxyStatus.text"
//+kubebuilder:printcolumn:name="Studio",type="string",JSONPath=".status.studioStatus.text"
//+kubebuilder:printcolumn:name="Prometheus",type="string",JSONPath=".status.prometheusStatus.text"
//+kubebuilder:printcolumn:name="Thanos-Store",type="string",JSONPath=".status.thanosStoreStatus.text"
//+kubebuilder:printcolumn:name="Thanos-Compactor",type="string",JSONPath=".status.thanosCompactorStatus.text"
//+kubebuilder:printcolumn:name="Thanos-Querier",type="string",JSONPath=".status.thanosQuerierStatus.text"
//+kubebuilder:printcolumn:name="Redis",type="string",JSONPath=".status.redisStatus.text"
//+kubebuilder:printcolumn:name="OpenCost",type="string",JSONPath=".status.opencostStatus.text"

// PlantDCore is the Schema for the plantdcores API
type PlantDCore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlantDCoreSpec   `json:"spec,omitempty"`
	Status PlantDCoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PlantDCoreList contains a list of PlantDCore
type PlantDCoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlantDCore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PlantDCore{}, &PlantDCoreList{})
}
