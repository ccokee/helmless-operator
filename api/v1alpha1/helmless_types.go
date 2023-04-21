package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HelmLessSpec defines the desired state of HelmLess
type HelmLessSpec struct {
	ChartRepo    string `json:"chartRepo,omitempty"`
	ChartName    string `json:"chartName,omitempty"`
	ChartVersion string `json:"chartVersion,omitempty"`
	Namespace    string `json:"namespace,omitempty"`
	PublicGist   string `json:"publicGist,omitempty"`
}

// HelmLessStatus defines the observed state of HelmLess
type HelmLessStatus struct {
	Deployed       bool                 `json:"deployed,omitempty"`
	Message        string               `json:"message,omitempty"`
	ReleaseName    string               `json:"releaseName,omitempty"`
	Namespace      string               `json:"namespace,omitempty"`
	ChartInfo      HelmChartInfo        `json:"chartInfo,omitempty"`
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
