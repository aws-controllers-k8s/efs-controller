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

package file_system

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/efs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.EFS{}
	_ = &svcapitypes.FileSystem{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadManyInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newListRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DescribeFileSystemsOutput
	resp, err = rm.sdkapi.DescribeFileSystemsWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_MANY", "DescribeFileSystems", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "FileSystemNotFound" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	found := false
	for _, elem := range resp.FileSystems {
		if elem.AvailabilityZoneId != nil {
			ko.Status.AvailabilityZoneID = elem.AvailabilityZoneId
		} else {
			ko.Status.AvailabilityZoneID = nil
		}
		if elem.AvailabilityZoneName != nil {
			ko.Spec.AvailabilityZoneName = elem.AvailabilityZoneName
		} else {
			ko.Spec.AvailabilityZoneName = nil
		}
		if elem.CreationTime != nil {
			ko.Status.CreationTime = &metav1.Time{*elem.CreationTime}
		} else {
			ko.Status.CreationTime = nil
		}
		if elem.Encrypted != nil {
			ko.Spec.Encrypted = elem.Encrypted
		} else {
			ko.Spec.Encrypted = nil
		}
		if elem.FileSystemArn != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.FileSystemArn)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.FileSystemId != nil {
			ko.Status.FileSystemID = elem.FileSystemId
		} else {
			ko.Status.FileSystemID = nil
		}
		if elem.FileSystemProtection != nil {
			f6 := &svcapitypes.UpdateFileSystemProtectionInput{}
			if elem.FileSystemProtection.ReplicationOverwriteProtection != nil {
				f6.ReplicationOverwriteProtection = elem.FileSystemProtection.ReplicationOverwriteProtection
			}
			ko.Spec.FileSystemProtection = f6
		} else {
			ko.Spec.FileSystemProtection = nil
		}
		if elem.KmsKeyId != nil {
			ko.Spec.KMSKeyID = elem.KmsKeyId
		} else {
			ko.Spec.KMSKeyID = nil
		}
		if elem.LifeCycleState != nil {
			ko.Status.LifeCycleState = elem.LifeCycleState
		} else {
			ko.Status.LifeCycleState = nil
		}
		if elem.Name != nil {
			ko.Status.Name = elem.Name
		} else {
			ko.Status.Name = nil
		}
		if elem.NumberOfMountTargets != nil {
			ko.Status.NumberOfMountTargets = elem.NumberOfMountTargets
		} else {
			ko.Status.NumberOfMountTargets = nil
		}
		if elem.OwnerId != nil {
			ko.Status.OwnerID = elem.OwnerId
		} else {
			ko.Status.OwnerID = nil
		}
		if elem.PerformanceMode != nil {
			ko.Spec.PerformanceMode = elem.PerformanceMode
		} else {
			ko.Spec.PerformanceMode = nil
		}
		if elem.ProvisionedThroughputInMibps != nil {
			ko.Spec.ProvisionedThroughputInMiBps = elem.ProvisionedThroughputInMibps
		} else {
			ko.Spec.ProvisionedThroughputInMiBps = nil
		}
		if elem.SizeInBytes != nil {
			f14 := &svcapitypes.FileSystemSize{}
			if elem.SizeInBytes.Timestamp != nil {
				f14.Timestamp = &metav1.Time{*elem.SizeInBytes.Timestamp}
			}
			if elem.SizeInBytes.Value != nil {
				f14.Value = elem.SizeInBytes.Value
			}
			if elem.SizeInBytes.ValueInArchive != nil {
				f14.ValueInArchive = elem.SizeInBytes.ValueInArchive
			}
			if elem.SizeInBytes.ValueInIA != nil {
				f14.ValueInIA = elem.SizeInBytes.ValueInIA
			}
			if elem.SizeInBytes.ValueInStandard != nil {
				f14.ValueInStandard = elem.SizeInBytes.ValueInStandard
			}
			ko.Status.SizeInBytes = f14
		} else {
			ko.Status.SizeInBytes = nil
		}
		if elem.Tags != nil {
			f15 := []*svcapitypes.Tag{}
			for _, f15iter := range elem.Tags {
				f15elem := &svcapitypes.Tag{}
				if f15iter.Key != nil {
					f15elem.Key = f15iter.Key
				}
				if f15iter.Value != nil {
					f15elem.Value = f15iter.Value
				}
				f15 = append(f15, f15elem)
			}
			ko.Spec.Tags = f15
		} else {
			ko.Spec.Tags = nil
		}
		if elem.ThroughputMode != nil {
			ko.Spec.ThroughputMode = elem.ThroughputMode
		} else {
			ko.Spec.ThroughputMode = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}

	rm.setStatusDefaults(ko)
	if err := rm.setResourceAdditionalFields(ctx, ko); err != nil {
		return nil, err
	}
	if !filesystemActive(&resource{ko}) {
		return &resource{ko}, requeueWaitState(r)
	}

	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadManyInput returns true if there are any fields
// for the ReadMany Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadManyInput(
	r *resource,
) bool {
	return r.ko.Status.FileSystemID == nil

}

// newListRequestPayload returns SDK-specific struct for the HTTP request
// payload of the List API call for the resource
func (rm *resourceManager) newListRequestPayload(
	r *resource,
) (*svcsdk.DescribeFileSystemsInput, error) {
	res := &svcsdk.DescribeFileSystemsInput{}

	if r.ko.Status.FileSystemID != nil {
		res.SetFileSystemId(*r.ko.Status.FileSystemID)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.FileSystemDescription
	_ = resp
	resp, err = rm.sdkapi.CreateFileSystemWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateFileSystem", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.AvailabilityZoneId != nil {
		ko.Status.AvailabilityZoneID = resp.AvailabilityZoneId
	} else {
		ko.Status.AvailabilityZoneID = nil
	}
	if resp.AvailabilityZoneName != nil {
		ko.Spec.AvailabilityZoneName = resp.AvailabilityZoneName
	} else {
		ko.Spec.AvailabilityZoneName = nil
	}
	if resp.CreationTime != nil {
		ko.Status.CreationTime = &metav1.Time{*resp.CreationTime}
	} else {
		ko.Status.CreationTime = nil
	}
	if resp.Encrypted != nil {
		ko.Spec.Encrypted = resp.Encrypted
	} else {
		ko.Spec.Encrypted = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.FileSystemArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.FileSystemArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.FileSystemId != nil {
		ko.Status.FileSystemID = resp.FileSystemId
	} else {
		ko.Status.FileSystemID = nil
	}
	if resp.FileSystemProtection != nil {
		f6 := &svcapitypes.UpdateFileSystemProtectionInput{}
		if resp.FileSystemProtection.ReplicationOverwriteProtection != nil {
			f6.ReplicationOverwriteProtection = resp.FileSystemProtection.ReplicationOverwriteProtection
		}
		ko.Spec.FileSystemProtection = f6
	} else {
		ko.Spec.FileSystemProtection = nil
	}
	if resp.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.LifeCycleState != nil {
		ko.Status.LifeCycleState = resp.LifeCycleState
	} else {
		ko.Status.LifeCycleState = nil
	}
	if resp.Name != nil {
		ko.Status.Name = resp.Name
	} else {
		ko.Status.Name = nil
	}
	if resp.NumberOfMountTargets != nil {
		ko.Status.NumberOfMountTargets = resp.NumberOfMountTargets
	} else {
		ko.Status.NumberOfMountTargets = nil
	}
	if resp.OwnerId != nil {
		ko.Status.OwnerID = resp.OwnerId
	} else {
		ko.Status.OwnerID = nil
	}
	if resp.PerformanceMode != nil {
		ko.Spec.PerformanceMode = resp.PerformanceMode
	} else {
		ko.Spec.PerformanceMode = nil
	}
	if resp.ProvisionedThroughputInMibps != nil {
		ko.Spec.ProvisionedThroughputInMiBps = resp.ProvisionedThroughputInMibps
	} else {
		ko.Spec.ProvisionedThroughputInMiBps = nil
	}
	if resp.SizeInBytes != nil {
		f14 := &svcapitypes.FileSystemSize{}
		if resp.SizeInBytes.Timestamp != nil {
			f14.Timestamp = &metav1.Time{*resp.SizeInBytes.Timestamp}
		}
		if resp.SizeInBytes.Value != nil {
			f14.Value = resp.SizeInBytes.Value
		}
		if resp.SizeInBytes.ValueInArchive != nil {
			f14.ValueInArchive = resp.SizeInBytes.ValueInArchive
		}
		if resp.SizeInBytes.ValueInIA != nil {
			f14.ValueInIA = resp.SizeInBytes.ValueInIA
		}
		if resp.SizeInBytes.ValueInStandard != nil {
			f14.ValueInStandard = resp.SizeInBytes.ValueInStandard
		}
		ko.Status.SizeInBytes = f14
	} else {
		ko.Status.SizeInBytes = nil
	}
	if resp.Tags != nil {
		f15 := []*svcapitypes.Tag{}
		for _, f15iter := range resp.Tags {
			f15elem := &svcapitypes.Tag{}
			if f15iter.Key != nil {
				f15elem.Key = f15iter.Key
			}
			if f15iter.Value != nil {
				f15elem.Value = f15iter.Value
			}
			f15 = append(f15, f15elem)
		}
		ko.Spec.Tags = f15
	} else {
		ko.Spec.Tags = nil
	}
	if resp.ThroughputMode != nil {
		ko.Spec.ThroughputMode = resp.ThroughputMode
	} else {
		ko.Spec.ThroughputMode = nil
	}

	rm.setStatusDefaults(ko)
	// We expect the fs to be in 'creating' status since we just issued
	// the call to create it, but I suppose it doesn't hurt to check here.
	if filesystemCreating(&resource{ko}) {
		// Setting resource synced condition to false will trigger a requeue of
		// the resource. No need to return a requeue error here.
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)
		return &resource{ko}, nil
	}

	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateFileSystemInput, error) {
	res := &svcsdk.CreateFileSystemInput{}

	if r.ko.Spec.AvailabilityZoneName != nil {
		res.SetAvailabilityZoneName(*r.ko.Spec.AvailabilityZoneName)
	}
	if r.ko.Spec.Backup != nil {
		res.SetBackup(*r.ko.Spec.Backup)
	}
	if r.ko.Spec.Encrypted != nil {
		res.SetEncrypted(*r.ko.Spec.Encrypted)
	}
	if r.ko.Spec.KMSKeyID != nil {
		res.SetKmsKeyId(*r.ko.Spec.KMSKeyID)
	}
	if r.ko.Spec.PerformanceMode != nil {
		res.SetPerformanceMode(*r.ko.Spec.PerformanceMode)
	}
	if r.ko.Spec.ProvisionedThroughputInMiBps != nil {
		res.SetProvisionedThroughputInMibps(*r.ko.Spec.ProvisionedThroughputInMiBps)
	}
	if r.ko.Spec.Tags != nil {
		f6 := []*svcsdk.Tag{}
		for _, f6iter := range r.ko.Spec.Tags {
			f6elem := &svcsdk.Tag{}
			if f6iter.Key != nil {
				f6elem.SetKey(*f6iter.Key)
			}
			if f6iter.Value != nil {
				f6elem.SetValue(*f6iter.Value)
			}
			f6 = append(f6, f6elem)
		}
		res.SetTags(f6)
	}
	if r.ko.Spec.ThroughputMode != nil {
		res.SetThroughputMode(*r.ko.Spec.ThroughputMode)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()

	if delta.DifferentAt("Spec.Tags") {
		err := syncTags(
			ctx, rm.sdkapi, rm.metrics,
			string(*desired.ko.Status.ACKResourceMetadata.ARN),
			desired.ko.Spec.Tags, latest.ko.Spec.Tags,
		)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.Policy") {
		err := rm.syncPolicy(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.BackupPolicy") {
		err := rm.syncBackupPolicy(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.LifecyclePolicies") {
		err := rm.syncLifecyclePolicies(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.FileSystemProtection") {
		err := rm.syncFilesystemProtection(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	// To trigger to normal update we need to make sure that at least
	// one of the following fields are different.
	if !delta.DifferentAt("Spec.ProvisionedThroughputInMiBps") && !delta.DifferentAt("Spec.ThroughputMode") {
		return desired, nil
	}
	input, err := rm.newUpdateRequestPayload(ctx, desired, delta)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.UpdateFileSystemOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateFileSystemWithContext(ctx, input)
	return desired, nil
	rm.metrics.RecordAPICall("UPDATE", "UpdateFileSystem", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.AvailabilityZoneId != nil {
		ko.Status.AvailabilityZoneID = resp.AvailabilityZoneId
	} else {
		ko.Status.AvailabilityZoneID = nil
	}
	if resp.AvailabilityZoneName != nil {
		ko.Spec.AvailabilityZoneName = resp.AvailabilityZoneName
	} else {
		ko.Spec.AvailabilityZoneName = nil
	}
	if resp.CreationTime != nil {
		ko.Status.CreationTime = &metav1.Time{*resp.CreationTime}
	} else {
		ko.Status.CreationTime = nil
	}
	if resp.Encrypted != nil {
		ko.Spec.Encrypted = resp.Encrypted
	} else {
		ko.Spec.Encrypted = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.FileSystemArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.FileSystemArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.FileSystemId != nil {
		ko.Status.FileSystemID = resp.FileSystemId
	} else {
		ko.Status.FileSystemID = nil
	}
	if resp.FileSystemProtection != nil {
		f7 := &svcapitypes.UpdateFileSystemProtectionInput{}
		if resp.FileSystemProtection.ReplicationOverwriteProtection != nil {
			f7.ReplicationOverwriteProtection = resp.FileSystemProtection.ReplicationOverwriteProtection
		}
		ko.Spec.FileSystemProtection = f7
	} else {
		ko.Spec.FileSystemProtection = nil
	}
	if resp.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.LifeCycleState != nil {
		ko.Status.LifeCycleState = resp.LifeCycleState
	} else {
		ko.Status.LifeCycleState = nil
	}
	if resp.Name != nil {
		ko.Status.Name = resp.Name
	} else {
		ko.Status.Name = nil
	}
	if resp.NumberOfMountTargets != nil {
		ko.Status.NumberOfMountTargets = resp.NumberOfMountTargets
	} else {
		ko.Status.NumberOfMountTargets = nil
	}
	if resp.OwnerId != nil {
		ko.Status.OwnerID = resp.OwnerId
	} else {
		ko.Status.OwnerID = nil
	}
	if resp.PerformanceMode != nil {
		ko.Spec.PerformanceMode = resp.PerformanceMode
	} else {
		ko.Spec.PerformanceMode = nil
	}
	if resp.ProvisionedThroughputInMibps != nil {
		ko.Spec.ProvisionedThroughputInMiBps = resp.ProvisionedThroughputInMibps
	} else {
		ko.Spec.ProvisionedThroughputInMiBps = nil
	}
	if resp.SizeInBytes != nil {
		f15 := &svcapitypes.FileSystemSize{}
		if resp.SizeInBytes.Timestamp != nil {
			f15.Timestamp = &metav1.Time{*resp.SizeInBytes.Timestamp}
		}
		if resp.SizeInBytes.Value != nil {
			f15.Value = resp.SizeInBytes.Value
		}
		if resp.SizeInBytes.ValueInArchive != nil {
			f15.ValueInArchive = resp.SizeInBytes.ValueInArchive
		}
		if resp.SizeInBytes.ValueInIA != nil {
			f15.ValueInIA = resp.SizeInBytes.ValueInIA
		}
		if resp.SizeInBytes.ValueInStandard != nil {
			f15.ValueInStandard = resp.SizeInBytes.ValueInStandard
		}
		ko.Status.SizeInBytes = f15
	} else {
		ko.Status.SizeInBytes = nil
	}
	if resp.Tags != nil {
		f16 := []*svcapitypes.Tag{}
		for _, f16iter := range resp.Tags {
			f16elem := &svcapitypes.Tag{}
			if f16iter.Key != nil {
				f16elem.Key = f16iter.Key
			}
			if f16iter.Value != nil {
				f16elem.Value = f16iter.Value
			}
			f16 = append(f16, f16elem)
		}
		ko.Spec.Tags = f16
	} else {
		ko.Spec.Tags = nil
	}
	if resp.ThroughputMode != nil {
		ko.Spec.ThroughputMode = resp.ThroughputMode
	} else {
		ko.Spec.ThroughputMode = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
	delta *ackcompare.Delta,
) (*svcsdk.UpdateFileSystemInput, error) {
	res := &svcsdk.UpdateFileSystemInput{}

	if r.ko.Status.FileSystemID != nil {
		res.SetFileSystemId(*r.ko.Status.FileSystemID)
	}
	if r.ko.Spec.ProvisionedThroughputInMiBps != nil {
		res.SetProvisionedThroughputInMibps(*r.ko.Spec.ProvisionedThroughputInMiBps)
	}
	if r.ko.Spec.ThroughputMode != nil {
		res.SetThroughputMode(*r.ko.Spec.ThroughputMode)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteFileSystemOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteFileSystemWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteFileSystem", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteFileSystemInput, error) {
	res := &svcsdk.DeleteFileSystemInput{}

	if r.ko.Status.FileSystemID != nil {
		res.SetFileSystemId(*r.ko.Status.FileSystemID)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.FileSystem,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	// No terminal_errors specified for this resource in generator config
	return false
}
