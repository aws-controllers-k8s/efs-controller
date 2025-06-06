// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MountTargetSpec defines the desired state of MountTarget.
type MountTargetSpec struct {

	// The ID of the file system for which to create the mount target.
	//
	// Regex Pattern: `^(arn:aws[-a-z]*:elasticfilesystem:[0-9a-z-:]+:file-system/fs-[0-9a-f]{8,40}|fs-[0-9a-f]{8,40})$`
	FileSystemID  *string                                  `json:"fileSystemID,omitempty"`
	FileSystemRef *ackv1alpha1.AWSResourceReferenceWrapper `json:"fileSystemRef,omitempty"`
	// Valid IPv4 address within the address range of the specified subnet.
	//
	// Regex Pattern: `^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`
	IPAddress         *string                                    `json:"ipAddress,omitempty"`
	SecurityGroupRefs []*ackv1alpha1.AWSResourceReferenceWrapper `json:"securityGroupRefs,omitempty"`
	// Up to five VPC security group IDs, of the form sg-xxxxxxxx. These must be
	// for the same VPC as subnet specified.
	SecurityGroups []*string `json:"securityGroups,omitempty"`
	// The ID of the subnet to add the mount target in. For One Zone file systems,
	// use the subnet that is associated with the file system's Availability Zone.
	//
	// Regex Pattern: `^subnet-[0-9a-f]{8,40}$`
	SubnetID  *string                                  `json:"subnetID,omitempty"`
	SubnetRef *ackv1alpha1.AWSResourceReferenceWrapper `json:"subnetRef,omitempty"`
}

// MountTargetStatus defines the observed state of MountTarget
type MountTargetStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRs managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// The unique and consistent identifier of the Availability Zone that the mount
	// target resides in. For example, use1-az1 is an AZ ID for the us-east-1 Region
	// and it has the same location in every Amazon Web Services account.
	// +kubebuilder:validation:Optional
	AvailabilityZoneID *string `json:"availabilityZoneID,omitempty"`
	// The name of the Availability Zone in which the mount target is located. Availability
	// Zones are independently mapped to names for each Amazon Web Services account.
	// For example, the Availability Zone us-east-1a for your Amazon Web Services
	// account might not be the same location as us-east-1a for another Amazon Web
	// Services account.
	//
	// Regex Pattern: `^.+$`
	// +kubebuilder:validation:Optional
	AvailabilityZoneName *string `json:"availabilityZoneName,omitempty"`
	// Lifecycle state of the mount target.
	// +kubebuilder:validation:Optional
	LifeCycleState *string `json:"lifeCycleState,omitempty"`
	// System-assigned mount target ID.
	//
	// Regex Pattern: `^fsmt-[0-9a-f]{8,40}$`
	// +kubebuilder:validation:Optional
	MountTargetID *string `json:"mountTargetID,omitempty"`
	// The ID of the network interface that Amazon EFS created when it created the
	// mount target.
	// +kubebuilder:validation:Optional
	NetworkInterfaceID *string `json:"networkInterfaceID,omitempty"`
	// Amazon Web Services account ID that owns the resource.
	//
	// Regex Pattern: `^(\d{12})|(\d{4}-\d{4}-\d{4})$`
	// +kubebuilder:validation:Optional
	OwnerID *string `json:"ownerID,omitempty"`
	// The virtual private cloud (VPC) ID that the mount target is configured in.
	// +kubebuilder:validation:Optional
	VPCID *string `json:"vpcID,omitempty"`
}

// MountTarget is the Schema for the MountTargets API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="FILESYSTEMID",type=string,priority=0,JSONPath=`.spec.fileSystemID`
// +kubebuilder:printcolumn:name="IPADDRESS",type=string,priority=0,JSONPath=`.spec.ipAddress`
// +kubebuilder:printcolumn:name="MOUNTTARGETID",type=string,priority=0,JSONPath=`.status.mountTargetID`
// +kubebuilder:printcolumn:name="SUBNETID",type=string,priority=0,JSONPath=`.spec.subnetID`
// +kubebuilder:printcolumn:name="VPCID",type=string,priority=1,JSONPath=`.status.vpcID`
// +kubebuilder:printcolumn:name="AVAILABILITYZONEID",type=string,priority=1,JSONPath=`.status.availabilityZoneID`
// +kubebuilder:printcolumn:name="AVAILABILITYZONENAME",type=string,priority=1,JSONPath=`.status.availabilityZoneName`
// +kubebuilder:printcolumn:name="STATE",type=string,priority=0,JSONPath=`.status.lifeCycleState`
// +kubebuilder:printcolumn:name="Synced",type="string",priority=0,JSONPath=".status.conditions[?(@.type==\"ACK.ResourceSynced\")].status"
// +kubebuilder:printcolumn:name="Age",type="date",priority=0,JSONPath=".metadata.creationTimestamp"
type MountTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MountTargetSpec   `json:"spec,omitempty"`
	Status            MountTargetStatus `json:"status,omitempty"`
}

// MountTargetList contains a list of MountTarget
// +kubebuilder:object:root=true
type MountTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MountTarget `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MountTarget{}, &MountTargetList{})
}
