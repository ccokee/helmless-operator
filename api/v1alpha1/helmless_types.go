package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HelmLessSpec defines the desired state of HelmLess
type HelmLessSpec struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ChartRepo string `json:"chartRepo,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ChartName string `json:"chartName,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ChartReleaseName string `json:"chartReleaseName,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ChartVersion string `json:"chartVersion,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Namespace string `json:"namespace,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ValuesUrl string `json:"valuesUrl,omitempty"`
}

// HelmLessStatus defines the observed state of HelmLess
type HelmLessStatus struct {
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Deployed bool `json:"deployed,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Message string `json:"message,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ReleaseName string `json:"releaseName,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Namespace string `json:"namespace,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ChartInfo HelmChartInfo `json:"chartInfo,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=status
	DeploymentInfo []HelmDeploymentInfo `json:"deploymentInfo,omitempty"`
}

type HelmChartInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type HelmDeploymentInfo struct {
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HelmLess is the Schema for the helmless API
type HelmLess struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmLessSpec   `json:"spec,omitempty"`
	Status HelmLessStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HelmLessList contains a list of HelmLess
type HelmLessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmLess `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelmLess{}, &HelmLessList{})
}
