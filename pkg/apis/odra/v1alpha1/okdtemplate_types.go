package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type OKDTemplateConditionType string

const (
	OKDTemplateNone      OKDTemplateConditionType = ""
	OKDTemplateNew       OKDTemplateConditionType = "New"
	OKDTemplateReconcile OKDTemplateConditionType = "Reconcile"
	OKDTemplateReady     OKDTemplateConditionType = "Ready"
	OKDTemplateError     OKDTemplateConditionType = "Error"
	OKDTemplateDelete    OKDTemplateConditionType = "Delete"
)

type TemplateSource struct {
	Url        string `json:"url"`
	Local      string `json:"local"`
}

// OKDTemplateSpec defines the desired state of OKDTemplate
type OKDTemplateSpec struct {
	Source *TemplateSource `json:"source"`
	Parameters map[string]string `json:"parameters"`
}

// OKDTemplateStatus defines the observed state of OKDTemplate
type OKDTemplateStatus struct {
	Type    OKDTemplateConditionType `json:"type"`
	Reason  *string `json:"reason"`
	Message *string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OKDTemplate is the Schema for the okdtemplates API
// +k8s:openapi-gen=true
type OKDTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OKDTemplateSpec   `json:"spec,omitempty"`
	Status OKDTemplateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OKDTemplateList contains a list of OKDTemplate
type OKDTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OKDTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OKDTemplate{}, &OKDTemplateList{})
}
