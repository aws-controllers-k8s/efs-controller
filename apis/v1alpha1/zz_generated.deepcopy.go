//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessPointDescription) DeepCopyInto(out *AccessPointDescription) {
	*out = *in
	if in.AccessPointID != nil {
		in, out := &in.AccessPointID, &out.AccessPointID
		*out = new(string)
		**out = **in
	}
	if in.FileSystemID != nil {
		in, out := &in.FileSystemID, &out.FileSystemID
		*out = new(string)
		**out = **in
	}
	if in.LifeCycleState != nil {
		in, out := &in.LifeCycleState, &out.LifeCycleState
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]*Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessPointDescription.
func (in *AccessPointDescription) DeepCopy() *AccessPointDescription {
	if in == nil {
		return nil
	}
	out := new(AccessPointDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupPolicy) DeepCopyInto(out *BackupPolicy) {
	*out = *in
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupPolicy.
func (in *BackupPolicy) DeepCopy() *BackupPolicy {
	if in == nil {
		return nil
	}
	out := new(BackupPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Destination) DeepCopyInto(out *Destination) {
	*out = *in
	if in.FileSystemID != nil {
		in, out := &in.FileSystemID, &out.FileSystemID
		*out = new(string)
		**out = **in
	}
	if in.LastReplicatedTimestamp != nil {
		in, out := &in.LastReplicatedTimestamp, &out.LastReplicatedTimestamp
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Destination.
func (in *Destination) DeepCopy() *Destination {
	if in == nil {
		return nil
	}
	out := new(Destination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DestinationToCreate) DeepCopyInto(out *DestinationToCreate) {
	*out = *in
	if in.AvailabilityZoneName != nil {
		in, out := &in.AvailabilityZoneName, &out.AvailabilityZoneName
		*out = new(string)
		**out = **in
	}
	if in.FileSystemID != nil {
		in, out := &in.FileSystemID, &out.FileSystemID
		*out = new(string)
		**out = **in
	}
	if in.KMSKeyID != nil {
		in, out := &in.KMSKeyID, &out.KMSKeyID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DestinationToCreate.
func (in *DestinationToCreate) DeepCopy() *DestinationToCreate {
	if in == nil {
		return nil
	}
	out := new(DestinationToCreate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSystem) DeepCopyInto(out *FileSystem) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSystem.
func (in *FileSystem) DeepCopy() *FileSystem {
	if in == nil {
		return nil
	}
	out := new(FileSystem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FileSystem) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSystemDescription) DeepCopyInto(out *FileSystemDescription) {
	*out = *in
	if in.AvailabilityZoneID != nil {
		in, out := &in.AvailabilityZoneID, &out.AvailabilityZoneID
		*out = new(string)
		**out = **in
	}
	if in.AvailabilityZoneName != nil {
		in, out := &in.AvailabilityZoneName, &out.AvailabilityZoneName
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = (*in).DeepCopy()
	}
	if in.Encrypted != nil {
		in, out := &in.Encrypted, &out.Encrypted
		*out = new(bool)
		**out = **in
	}
	if in.FileSystemARN != nil {
		in, out := &in.FileSystemARN, &out.FileSystemARN
		*out = new(string)
		**out = **in
	}
	if in.FileSystemID != nil {
		in, out := &in.FileSystemID, &out.FileSystemID
		*out = new(string)
		**out = **in
	}
	if in.FileSystemProtection != nil {
		in, out := &in.FileSystemProtection, &out.FileSystemProtection
		*out = new(FileSystemProtectionDescription)
		(*in).DeepCopyInto(*out)
	}
	if in.KMSKeyID != nil {
		in, out := &in.KMSKeyID, &out.KMSKeyID
		*out = new(string)
		**out = **in
	}
	if in.LifeCycleState != nil {
		in, out := &in.LifeCycleState, &out.LifeCycleState
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.NumberOfMountTargets != nil {
		in, out := &in.NumberOfMountTargets, &out.NumberOfMountTargets
		*out = new(int64)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.PerformanceMode != nil {
		in, out := &in.PerformanceMode, &out.PerformanceMode
		*out = new(string)
		**out = **in
	}
	if in.ProvisionedThroughputInMiBps != nil {
		in, out := &in.ProvisionedThroughputInMiBps, &out.ProvisionedThroughputInMiBps
		*out = new(float64)
		**out = **in
	}
	if in.SizeInBytes != nil {
		in, out := &in.SizeInBytes, &out.SizeInBytes
		*out = new(FileSystemSize)
		(*in).DeepCopyInto(*out)
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]*Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.ThroughputMode != nil {
		in, out := &in.ThroughputMode, &out.ThroughputMode
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSystemDescription.
func (in *FileSystemDescription) DeepCopy() *FileSystemDescription {
	if in == nil {
		return nil
	}
	out := new(FileSystemDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSystemList) DeepCopyInto(out *FileSystemList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FileSystem, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSystemList.
func (in *FileSystemList) DeepCopy() *FileSystemList {
	if in == nil {
		return nil
	}
	out := new(FileSystemList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FileSystemList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSystemProtectionDescription) DeepCopyInto(out *FileSystemProtectionDescription) {
	*out = *in
	if in.ReplicationOverwriteProtection != nil {
		in, out := &in.ReplicationOverwriteProtection, &out.ReplicationOverwriteProtection
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSystemProtectionDescription.
func (in *FileSystemProtectionDescription) DeepCopy() *FileSystemProtectionDescription {
	if in == nil {
		return nil
	}
	out := new(FileSystemProtectionDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSystemSize) DeepCopyInto(out *FileSystemSize) {
	*out = *in
	if in.Timestamp != nil {
		in, out := &in.Timestamp, &out.Timestamp
		*out = (*in).DeepCopy()
	}
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(int64)
		**out = **in
	}
	if in.ValueInArchive != nil {
		in, out := &in.ValueInArchive, &out.ValueInArchive
		*out = new(int64)
		**out = **in
	}
	if in.ValueInIA != nil {
		in, out := &in.ValueInIA, &out.ValueInIA
		*out = new(int64)
		**out = **in
	}
	if in.ValueInStandard != nil {
		in, out := &in.ValueInStandard, &out.ValueInStandard
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSystemSize.
func (in *FileSystemSize) DeepCopy() *FileSystemSize {
	if in == nil {
		return nil
	}
	out := new(FileSystemSize)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSystemSpec) DeepCopyInto(out *FileSystemSpec) {
	*out = *in
	if in.AvailabilityZoneName != nil {
		in, out := &in.AvailabilityZoneName, &out.AvailabilityZoneName
		*out = new(string)
		**out = **in
	}
	if in.Backup != nil {
		in, out := &in.Backup, &out.Backup
		*out = new(bool)
		**out = **in
	}
	if in.BackupPolicy != nil {
		in, out := &in.BackupPolicy, &out.BackupPolicy
		*out = new(BackupPolicy)
		(*in).DeepCopyInto(*out)
	}
	if in.Encrypted != nil {
		in, out := &in.Encrypted, &out.Encrypted
		*out = new(bool)
		**out = **in
	}
	if in.FileSystemProtection != nil {
		in, out := &in.FileSystemProtection, &out.FileSystemProtection
		*out = new(UpdateFileSystemProtectionInput)
		(*in).DeepCopyInto(*out)
	}
	if in.KMSKeyID != nil {
		in, out := &in.KMSKeyID, &out.KMSKeyID
		*out = new(string)
		**out = **in
	}
	if in.KMSKeyRef != nil {
		in, out := &in.KMSKeyRef, &out.KMSKeyRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.LifecyclePolicies != nil {
		in, out := &in.LifecyclePolicies, &out.LifecyclePolicies
		*out = make([]*LifecyclePolicy, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(LifecyclePolicy)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.PerformanceMode != nil {
		in, out := &in.PerformanceMode, &out.PerformanceMode
		*out = new(string)
		**out = **in
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = new(string)
		**out = **in
	}
	if in.ProvisionedThroughputInMiBps != nil {
		in, out := &in.ProvisionedThroughputInMiBps, &out.ProvisionedThroughputInMiBps
		*out = new(float64)
		**out = **in
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]*Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.ThroughputMode != nil {
		in, out := &in.ThroughputMode, &out.ThroughputMode
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSystemSpec.
func (in *FileSystemSpec) DeepCopy() *FileSystemSpec {
	if in == nil {
		return nil
	}
	out := new(FileSystemSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSystemStatus) DeepCopyInto(out *FileSystemStatus) {
	*out = *in
	if in.ACKResourceMetadata != nil {
		in, out := &in.ACKResourceMetadata, &out.ACKResourceMetadata
		*out = new(corev1alpha1.ResourceMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*corev1alpha1.Condition, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.Condition)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.AvailabilityZoneID != nil {
		in, out := &in.AvailabilityZoneID, &out.AvailabilityZoneID
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = (*in).DeepCopy()
	}
	if in.FileSystemID != nil {
		in, out := &in.FileSystemID, &out.FileSystemID
		*out = new(string)
		**out = **in
	}
	if in.LifeCycleState != nil {
		in, out := &in.LifeCycleState, &out.LifeCycleState
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.NumberOfMountTargets != nil {
		in, out := &in.NumberOfMountTargets, &out.NumberOfMountTargets
		*out = new(int64)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.SizeInBytes != nil {
		in, out := &in.SizeInBytes, &out.SizeInBytes
		*out = new(FileSystemSize)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSystemStatus.
func (in *FileSystemStatus) DeepCopy() *FileSystemStatus {
	if in == nil {
		return nil
	}
	out := new(FileSystemStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LifecyclePolicy) DeepCopyInto(out *LifecyclePolicy) {
	*out = *in
	if in.TransitionToArchive != nil {
		in, out := &in.TransitionToArchive, &out.TransitionToArchive
		*out = new(string)
		**out = **in
	}
	if in.TransitionToIA != nil {
		in, out := &in.TransitionToIA, &out.TransitionToIA
		*out = new(string)
		**out = **in
	}
	if in.TransitionToPrimaryStorageClass != nil {
		in, out := &in.TransitionToPrimaryStorageClass, &out.TransitionToPrimaryStorageClass
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LifecyclePolicy.
func (in *LifecyclePolicy) DeepCopy() *LifecyclePolicy {
	if in == nil {
		return nil
	}
	out := new(LifecyclePolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MountTarget) DeepCopyInto(out *MountTarget) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MountTarget.
func (in *MountTarget) DeepCopy() *MountTarget {
	if in == nil {
		return nil
	}
	out := new(MountTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MountTarget) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MountTargetDescription) DeepCopyInto(out *MountTargetDescription) {
	*out = *in
	if in.AvailabilityZoneID != nil {
		in, out := &in.AvailabilityZoneID, &out.AvailabilityZoneID
		*out = new(string)
		**out = **in
	}
	if in.AvailabilityZoneName != nil {
		in, out := &in.AvailabilityZoneName, &out.AvailabilityZoneName
		*out = new(string)
		**out = **in
	}
	if in.FileSystemID != nil {
		in, out := &in.FileSystemID, &out.FileSystemID
		*out = new(string)
		**out = **in
	}
	if in.IPAddress != nil {
		in, out := &in.IPAddress, &out.IPAddress
		*out = new(string)
		**out = **in
	}
	if in.LifeCycleState != nil {
		in, out := &in.LifeCycleState, &out.LifeCycleState
		*out = new(string)
		**out = **in
	}
	if in.MountTargetID != nil {
		in, out := &in.MountTargetID, &out.MountTargetID
		*out = new(string)
		**out = **in
	}
	if in.NetworkInterfaceID != nil {
		in, out := &in.NetworkInterfaceID, &out.NetworkInterfaceID
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.SubnetID != nil {
		in, out := &in.SubnetID, &out.SubnetID
		*out = new(string)
		**out = **in
	}
	if in.VPCID != nil {
		in, out := &in.VPCID, &out.VPCID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MountTargetDescription.
func (in *MountTargetDescription) DeepCopy() *MountTargetDescription {
	if in == nil {
		return nil
	}
	out := new(MountTargetDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MountTargetList) DeepCopyInto(out *MountTargetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MountTarget, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MountTargetList.
func (in *MountTargetList) DeepCopy() *MountTargetList {
	if in == nil {
		return nil
	}
	out := new(MountTargetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MountTargetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MountTargetSpec) DeepCopyInto(out *MountTargetSpec) {
	*out = *in
	if in.FileSystemID != nil {
		in, out := &in.FileSystemID, &out.FileSystemID
		*out = new(string)
		**out = **in
	}
	if in.FileSystemRef != nil {
		in, out := &in.FileSystemRef, &out.FileSystemRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.IPAddress != nil {
		in, out := &in.IPAddress, &out.IPAddress
		*out = new(string)
		**out = **in
	}
	if in.SecurityGroupRefs != nil {
		in, out := &in.SecurityGroupRefs, &out.SecurityGroupRefs
		*out = make([]*corev1alpha1.AWSResourceReferenceWrapper, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.AWSResourceReferenceWrapper)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.SecurityGroups != nil {
		in, out := &in.SecurityGroups, &out.SecurityGroups
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.SubnetID != nil {
		in, out := &in.SubnetID, &out.SubnetID
		*out = new(string)
		**out = **in
	}
	if in.SubnetRef != nil {
		in, out := &in.SubnetRef, &out.SubnetRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MountTargetSpec.
func (in *MountTargetSpec) DeepCopy() *MountTargetSpec {
	if in == nil {
		return nil
	}
	out := new(MountTargetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MountTargetStatus) DeepCopyInto(out *MountTargetStatus) {
	*out = *in
	if in.ACKResourceMetadata != nil {
		in, out := &in.ACKResourceMetadata, &out.ACKResourceMetadata
		*out = new(corev1alpha1.ResourceMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*corev1alpha1.Condition, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.Condition)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.AvailabilityZoneID != nil {
		in, out := &in.AvailabilityZoneID, &out.AvailabilityZoneID
		*out = new(string)
		**out = **in
	}
	if in.AvailabilityZoneName != nil {
		in, out := &in.AvailabilityZoneName, &out.AvailabilityZoneName
		*out = new(string)
		**out = **in
	}
	if in.LifeCycleState != nil {
		in, out := &in.LifeCycleState, &out.LifeCycleState
		*out = new(string)
		**out = **in
	}
	if in.MountTargetID != nil {
		in, out := &in.MountTargetID, &out.MountTargetID
		*out = new(string)
		**out = **in
	}
	if in.NetworkInterfaceID != nil {
		in, out := &in.NetworkInterfaceID, &out.NetworkInterfaceID
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.VPCID != nil {
		in, out := &in.VPCID, &out.VPCID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MountTargetStatus.
func (in *MountTargetStatus) DeepCopy() *MountTargetStatus {
	if in == nil {
		return nil
	}
	out := new(MountTargetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicationConfigurationDescription) DeepCopyInto(out *ReplicationConfigurationDescription) {
	*out = *in
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = (*in).DeepCopy()
	}
	if in.OriginalSourceFileSystemARN != nil {
		in, out := &in.OriginalSourceFileSystemARN, &out.OriginalSourceFileSystemARN
		*out = new(string)
		**out = **in
	}
	if in.SourceFileSystemARN != nil {
		in, out := &in.SourceFileSystemARN, &out.SourceFileSystemARN
		*out = new(string)
		**out = **in
	}
	if in.SourceFileSystemID != nil {
		in, out := &in.SourceFileSystemID, &out.SourceFileSystemID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicationConfigurationDescription.
func (in *ReplicationConfigurationDescription) DeepCopy() *ReplicationConfigurationDescription {
	if in == nil {
		return nil
	}
	out := new(ReplicationConfigurationDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tag) DeepCopyInto(out *Tag) {
	*out = *in
	if in.Key != nil {
		in, out := &in.Key, &out.Key
		*out = new(string)
		**out = **in
	}
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tag.
func (in *Tag) DeepCopy() *Tag {
	if in == nil {
		return nil
	}
	out := new(Tag)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UpdateFileSystemProtectionInput) DeepCopyInto(out *UpdateFileSystemProtectionInput) {
	*out = *in
	if in.ReplicationOverwriteProtection != nil {
		in, out := &in.ReplicationOverwriteProtection, &out.ReplicationOverwriteProtection
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UpdateFileSystemProtectionInput.
func (in *UpdateFileSystemProtectionInput) DeepCopy() *UpdateFileSystemProtectionInput {
	if in == nil {
		return nil
	}
	out := new(UpdateFileSystemProtectionInput)
	in.DeepCopyInto(out)
	return out
}
