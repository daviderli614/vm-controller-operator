/*
Copyright 2022.

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

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VMClaimSpec defines the desired state of VMClaim
type VMClaimSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Image     string `json:"image"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	Envs      []corev1.EnvVar `json:"envs,omitempty"`
	ScaleUp   *ScaleUpConfig `json:"scaleUp,omitempty"`
	ScaleDown *ScaleDownConfig `json:"scaleDown,omitempty"`

	// BalanceSimilarNodeGroups enables/disables the
	// `--balance-similar-node-groups` vm-controller feature.
	// This feature will automatically identify node groups with
	// the same instance type and the same set of labels and try
	// to keep the respective sizes of those node groups balanced.
	BalanceSimilarNodeGroups *bool `json:"balanceSimilarNodeGroups,omitempty"`

	// Enables/Disables `--skip-nodes-with-local-storage` feature flag. If true vm controller will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath. true by default at controller
	SkipNodesWithLocalStorage *bool `json:"skipNodesWithLocalStorage,omitempty"`

	// Refresh cycle
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	ScanInterval *string `json:"scanInterval,omitempty"`

	Expander *string `json:"expander,omitempty"`
}

// VMClaimStatus defines the observed state of VMClaim
type VMClaimStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	appsv1.DeploymentStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VMClaim is the Schema for the vmclaims API
type VMClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMClaimSpec   `json:"spec,omitempty"`
	Status VMClaimStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VMClaimList contains a list of VMClaim
type VMClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMClaim `json:"items"`
}

type ScaleUpConfig struct {
	// Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale up
	// +kubebuilder:validation:Pattern=(0.[0-9]+)
	ScaleUpUtilizationThreshold *string `json:"scaleUpUtilizationThreshold,omitempty"`
}

type ScaleDownConfig struct {
	// Should CA scale down the cluster
	Enabled bool `json:"enabled"`

	// How long after scale up that scale down evaluation resumes
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	DelayAfterAdd *string `json:"delayAfterAdd,omitempty"`

	// How long a node should be unneeded before it is eligible for scale down
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	UnneededTime *string `json:"unneededTime,omitempty"`

	// Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down
	// +kubebuilder:validation:Pattern=(0.[0-9]+)
	ScaleDownUtilizationThreshold *string `json:"scaleDownUtilizationThreshold,omitempty"`
}

func init() {
	SchemeBuilder.Register(&VMClaim{}, &VMClaimList{})
}
