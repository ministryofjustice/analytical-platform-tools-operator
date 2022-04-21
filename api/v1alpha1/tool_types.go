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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ToolSpec defines the desired state of Tool
type ToolSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the name of the tool to be deployed
	Name string `json:"Name"`
	// ImageVersion refers to the version of the pushed image
	ImageVersion string `json:"Version,omitempty"`
	// Username of the initialed user
	Username string `json:"User"`
	// IamRole assigned at initialisation
	IamRole string `json:"Iamrole"`
	// ProxyClientID is the client id of the auth proxy
	ProxyClientID string `json:"ProxyClientID,omitempty"`
	// ProxyClientSecret is the client secret of the auth proxy
	ProxyClientSecret string `json:"ProxyClientSecret,omitempty"`
	// ProxyDomain is the domain of the auth ProxyClientID
	ProxyDomain string `json:"ProxyDomain,omitempty"`
	// ProxyCookies is the cookies of the auth proxy
	ProxyCookies string `json:"ProxyCookies,omitempty"`
	// ProxyImageVersion is the version of the auth proxy
	ProxyImageVersion string `json:"ProxyVersion,omitempty"`
}

// ToolStatus defines the observed state of Tool
type ToolStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// DateLastUsed is the date the tool was last used by the user
	DateLastUsed string `json:"DateLastUsed,omitempty"`
	// PreviousVersions is a list of previous versions of the tool used by the user
	PreviousVersions []string `json:"PreviousVersions,omitempty"`
	// Url is the url of the tool
	Url string `json:"Url,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Tool is the Schema for the tools API
type Tool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ToolSpec   `json:"spec,omitempty"`
	Status ToolStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ToolList contains a list of Tool
type ToolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tool{}, &ToolList{})
}
