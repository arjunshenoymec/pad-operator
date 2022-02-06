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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PadSpec defines the desired state of Pad
type PadSpec struct {
	// +kubebuilder:default:="quay.io/aicoe/prometheus-anomaly-detector:latest"
	Image string `json:"image,omitempty"`
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`
	// +kubebuilder:default:="http://demo.robustperception.io:9090/"
	Source string `json:"source,omitempty"`
	// +kubebuilder:default:="up"
	Metrics string `json:"metrics,omitempty"`
	// +kubebuilder:default:="15"
	Retraining_interval string `json:"retraining_interval,omitempty"`
	// +kubebuilder:default:="24h"
	Training_window_size string `json:"training_window_size,omitempty"`
}

// PadStatus defines the observed state of Pad
type PadStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Empty
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// Pad is the Schema for the pads API
type Pad struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PadSpec   `json:"spec,omitempty"`
	Status PadStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PadList contains a list of Pad
type PadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pad `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pad{}, &PadList{})
}
