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

package file_system

import (
	"context"
	"fmt"
	"strings"
	"time"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/efs"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/efs-controller/pkg/tags"
	"github.com/aws/aws-sdk-go/aws"
)

// getIdempotencyToken returns a unique string to be used in certain API calls
// to ensure no replay of the call.
func getIdempotencyToken() string {
	return fmt.Sprintf("%d", time.Now().UTC().Nanosecond())
}

// requeueWaitState returns a `ackrequeue.RequeueNeededAfter` struct
// explaining the filesystem cannot be modified until it reaches an active status.
func requeueWaitState(r *resource) *ackrequeue.RequeueNeededAfter {
	if r.ko.Status.LifeCycleState == nil {
		return nil
	}
	status := *r.ko.Status.LifeCycleState
	return ackrequeue.NeededAfter(
		fmt.Errorf("filesystem in '%s' state, requeuing until filesystem is '%s'",
			status, svcapitypes.LifeCycleState_available),
		time.Second*3,
	)
}

// filesystemActive returns true if the supplied filesystem is in an active status
func filesystemActive(r *resource) bool {
	if r.ko.Status.LifeCycleState == nil {
		return false
	}
	cs := *r.ko.Status.LifeCycleState
	lifeCycleState := string(svcsdktypes.LifeCycleStateAvailable)
	return cs == lifeCycleState
}

var (
	// Requeue variables for different states
	requeueWaitReplicationConfiguration = ackrequeue.NeededAfter(
		fmt.Errorf("replication configuration is inactive, waiting for active state"),
		15*time.Second,
	)
)

// replicationBusy returns true if any replication destination is in a transitional state
func replicationConfigurationActive(r *resource) bool {
	if r.ko.Status.ReplicationConfigurationStatus == nil {
		return true
	}
	for _, dest := range r.ko.Status.ReplicationConfigurationStatus {
		if dest.Status != nil {
			status := *dest.Status
			if status == string(svcapitypes.ReplicationStatus_ENABLING) ||
				status == string(svcapitypes.ReplicationStatus_DELETING) {
				return false
			}
		}
	}
	return true
}

// replicationConfigurationExists returns true if replication configuration exists
func replicationConfigurationExists(r *resource) bool {
	return len(r.ko.Status.ReplicationConfigurationStatus) > 0
}

// Ideally, a part of this code needs to be generated.. However since the
// tags packge is not imported, we can't call it directly from sdk.go. We
// have to do this Go-fu to make it work.
var syncTags = tags.SyncTags

// setResourceDefaults queries the EFS API for the current state of the
// fields that are not returned by the ReadOne or List APIs.
func (rm *resourceManager) setResourceAdditionalFields(ctx context.Context, r *svcapitypes.FileSystem) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setResourceAdditionalFields")
	defer func() { exit(err) }()

	// Set the FS policy
	r.Spec.Policy, err = rm.getPolicy(ctx, r)
	if err != nil {
		return err
	}

	// Set backup policy
	r.Spec.BackupPolicy, err = rm.getBackupPolicy(ctx, r)
	if err != nil {
		return err
	}

	// Set lifecycle policies
	r.Spec.LifecyclePolicies, err = rm.getLifecyclePolicies(ctx, r)
	if err != nil {
		return err
	}

	// Set replication configuration
	err = rm.setReplicationConfiguration(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

// getBackupPolicy returns the backup policy for the file system
func (rm *resourceManager) getBackupPolicy(ctx context.Context, r *svcapitypes.FileSystem) (_ *svcapitypes.BackupPolicy, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getBackupPolicy")
	defer func() { exit(err) }()

	var output *svcsdk.DescribeBackupPolicyOutput
	output, err = rm.sdkapi.DescribeBackupPolicy(
		ctx,
		&svcsdk.DescribeBackupPolicyInput{
			FileSystemId: r.Status.FileSystemID,
		},
	)
	rm.metrics.RecordAPICall("GET", "DescribeBackupPolicy", err)
	if err != nil {
		if strings.Contains(err.Error(), "PolicyNotFound") {
			return nil, nil
		}
		return nil, err
	}
	backupPolicy := &svcapitypes.BackupPolicy{
		Status: aws.String(string(output.BackupPolicy.Status)),
	}

	return backupPolicy, nil
}

// getPolicy returns the policy for the file system
func (rm *resourceManager) getPolicy(ctx context.Context, r *svcapitypes.FileSystem) (_ *string, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getPolicy")
	defer func() { exit(err) }()

	var output *svcsdk.DescribeFileSystemPolicyOutput
	output, err = rm.sdkapi.DescribeFileSystemPolicy(
		ctx,
		&svcsdk.DescribeFileSystemPolicyInput{
			FileSystemId: r.Status.FileSystemID,
		},
	)
	rm.metrics.RecordAPICall("GET", "DescribeFileSystemPolicy", err)
	if err != nil {
		if strings.Contains(err.Error(), "PolicyNotFound") {
			return nil, nil
		}
		return nil, err
	}
	return output.Policy, nil
}

// getLifecyclePolicies returns the lifecycle policies for the file system
func (rm *resourceManager) getLifecyclePolicies(ctx context.Context, r *svcapitypes.FileSystem) (lps []*svcapitypes.LifecyclePolicy, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getLifecyclePolicies")
	defer func() { exit(err) }()

	var output *svcsdk.DescribeLifecycleConfigurationOutput
	output, err = rm.sdkapi.DescribeLifecycleConfiguration(
		ctx,
		&svcsdk.DescribeLifecycleConfigurationInput{
			FileSystemId: r.Status.FileSystemID,
		},
	)
	rm.metrics.RecordAPICall("GET", "DescribeLifecycleConfiguration", err)
	if err != nil {
		return nil, err
	}

	for _, lp := range output.LifecyclePolicies {
		lps = append(lps, &svcapitypes.LifecyclePolicy{
			TransitionToArchive:             aws.String(string(lp.TransitionToArchive)),
			TransitionToIA:                  aws.String(string(lp.TransitionToIA)),
			TransitionToPrimaryStorageClass: aws.String(string(lp.TransitionToPrimaryStorageClass)),
		})
	}

	return lps, nil
}

func (rm *resourceManager) syncPolicy(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncPolicy")
	defer func() { exit(err) }()

	if r.ko.Spec.Policy == nil {
		return rm.deletePolicy(ctx, r)
	}

	_, err = rm.sdkapi.PutFileSystemPolicy(
		ctx,
		&svcsdk.PutFileSystemPolicyInput{
			FileSystemId: r.ko.Status.FileSystemID,
			Policy:       r.ko.Spec.Policy,
			// BypassPolicyLockoutSafetyCheck: r.ko.Spec.BypassPolicyLockoutSafetyCheck,
		},
	)
	rm.metrics.RecordAPICall("UPDATE", "PutFileSystemPolicy", err)
	return err
}

func (rm *resourceManager) deletePolicy(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.deletePolicy")
	defer func() { exit(err) }()

	_, err = rm.sdkapi.DeleteFileSystemPolicy(
		ctx,
		&svcsdk.DeleteFileSystemPolicyInput{
			FileSystemId: r.ko.Status.FileSystemID,
		},
	)
	rm.metrics.RecordAPICall("UPDATE", "DeleteFileSystemPolicy", err)
	return err
}

func (rm *resourceManager) syncBackupPolicy(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncBackupPolicy")
	defer func() { exit(err) }()

	status := svcsdktypes.StatusDisabled
	if r.ko.Spec.BackupPolicy != nil && r.ko.Spec.BackupPolicy.Status != nil {
		status = svcsdktypes.Status(*r.ko.Spec.BackupPolicy.Status)
	}

	_, err = rm.sdkapi.PutBackupPolicy(
		ctx,
		&svcsdk.PutBackupPolicyInput{
			FileSystemId: r.ko.Status.FileSystemID,
			BackupPolicy: &svcsdktypes.BackupPolicy{
				Status: status,
			},
		},
	)
	rm.metrics.RecordAPICall("UPDATE", "PutBackupPolicy", err)
	return err
}

func (rm *resourceManager) syncLifecyclePolicies(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncLifecyclePolicies")
	defer func() { exit(err) }()
	// Check if LifecyclePolicies is nil or empty
	if len(r.ko.Spec.LifecyclePolicies) == 0 {
		rlog.Debug("No LifecyclePolicies specified, skipping sync")
		return nil
	}
	// Check if FileSystemID is nil
	if r.ko.Status.FileSystemID == nil {
		err = fmt.Errorf("FileSystemID is nil, cannot sync lifecycle policies")
		rlog.Info("FileSystemID is missing", "error", err)
		return err
	}
	// Build LifecyclePolicy list
	lps := []svcsdktypes.LifecyclePolicy{}
	for _, lp := range r.ko.Spec.LifecyclePolicies {
		l := svcsdktypes.LifecyclePolicy{}
		if lp.TransitionToArchive != nil {
			l.TransitionToArchive = svcsdktypes.TransitionToArchiveRules(*lp.TransitionToArchive)
		}
		if lp.TransitionToIA != nil {
			l.TransitionToIA = svcsdktypes.TransitionToIARules(*lp.TransitionToIA)
		}
		if lp.TransitionToPrimaryStorageClass != nil {
			l.TransitionToPrimaryStorageClass = svcsdktypes.TransitionToPrimaryStorageClassRules(*lp.TransitionToPrimaryStorageClass)
		}
		lps = append(lps, l)
	}
	// Log the API call payload for traceability
	rlog.Debug("Calling PutLifecycleConfiguration", "FileSystemID", *r.ko.Status.FileSystemID, "LifecyclePolicies", lps)
	// Make API call to update lifecycle policies
	_, err = rm.sdkapi.PutLifecycleConfiguration(
		ctx,
		&svcsdk.PutLifecycleConfigurationInput{
			FileSystemId:      r.ko.Status.FileSystemID,
			LifecyclePolicies: lps,
		},
	)
	if err != nil {
		err = fmt.Errorf("failed to sync lifecycle policies for FileSystemID '%s': %w", *r.ko.Status.FileSystemID, err)
		rlog.Info("PutLifecycleConfiguration API call failed", "error", err)
	}
	rm.metrics.RecordAPICall("UPDATE", "PutLifecycleConfiguration", err)
	return err
}

func (rm *resourceManager) syncFilesystemProtection(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncFilesystemProtection")
	defer func() { exit(err) }()

	_, err = rm.sdkapi.UpdateFileSystemProtection(
		ctx,
		&svcsdk.UpdateFileSystemProtectionInput{
			FileSystemId:                   r.ko.Status.FileSystemID,
			ReplicationOverwriteProtection: svcsdktypes.ReplicationOverwriteProtection(*r.ko.Spec.FileSystemProtection.ReplicationOverwriteProtection),
		},
	)
	rm.metrics.RecordAPICall("UPDATE", "UpdateFileSystemProtection", err)
	return err
}

// Comparison defaults
var (
	defaultEncrypted                      = false
	defaultReplicationOverwriteProtection = "ENABLED"
	defaultPerformanceMode                = "generalPurpose"
	defaultThroughputMode                 = "bursting"
)

// customPreCompare ensures that default values of nil-able types are
// appropriately replaced with empty maps or structs depending on the default
// output of the SDK.
func customPreCompare(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) {
	if a.ko.Spec.Encrypted == nil {
		a.ko.Spec.Encrypted = &defaultEncrypted
	}
	if a.ko.Spec.PerformanceMode == nil {
		a.ko.Spec.PerformanceMode = &defaultPerformanceMode
	}
	if a.ko.Spec.ThroughputMode == nil {
		a.ko.Spec.ThroughputMode = &defaultThroughputMode
	}
	if a.ko.Spec.FileSystemProtection == nil {
		a.ko.Spec.FileSystemProtection = &svcapitypes.UpdateFileSystemProtectionInput{
			ReplicationOverwriteProtection: &defaultReplicationOverwriteProtection,
		}
	}
	if a.ko.Spec.KMSKeyID == nil {
		a.ko.Spec.KMSKeyID = b.ko.Spec.KMSKeyID
	}

	// EFS replication supports two modes for destination file systems. First, implicit creation where
	// users specify only region/availabilityZone in the spec and AWS creates a new, empty destination
	// file system with the generated fileSystemID appearing only in status (not patched back to spec).
	// Second, explicit reference where users specify fileSystemID of an existing file system in the spec,
	// with the existing file system requiring replication overwrite protection to be disabled so AWS can
	// sync/overwrite the existing destination to match source data.
	//
	// The customPreCompare function handles both cases by adding AWS-generated fields (fileSystemID,
	// roleARN) to prevent unnecessary update cycles while preserving user intent for explicit versus
	// implicit destination creation.
	if len(b.ko.Spec.ReplicationConfiguration) > 0 {
		for i, aDest := range a.ko.Spec.ReplicationConfiguration {
			bDest := b.ko.Spec.ReplicationConfiguration[i]
			if aDest.FileSystemID == nil && bDest.FileSystemID != nil {
				aDest.FileSystemID = bDest.FileSystemID
			}
			if aDest.RoleARN == nil && bDest.RoleARN != nil {
				aDest.RoleARN = bDest.RoleARN
			}
		}
	}
}

// setReplicationConfiguration populates both Spec and Status replication fields for the file system
func (rm *resourceManager) setReplicationConfiguration(ctx context.Context, r *svcapitypes.FileSystem) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setReplicationConfiguration")
	defer func() { exit(err) }()

	var output *svcsdk.DescribeReplicationConfigurationsOutput
	output, err = rm.sdkapi.DescribeReplicationConfigurations(
		ctx,
		&svcsdk.DescribeReplicationConfigurationsInput{
			FileSystemId: r.Status.FileSystemID,
		},
	)
	rm.metrics.RecordAPICall("GET", "DescribeReplicationConfigurations", err)
	if err != nil {
		if strings.Contains(err.Error(), "ReplicationNotFound") {
			r.Spec.ReplicationConfiguration = nil
			r.Status.ReplicationConfigurationStatus = nil
			return nil
		}
		return err
	}

	if len(output.Replications) == 0 {
		r.Spec.ReplicationConfiguration = nil
		r.Status.ReplicationConfigurationStatus = nil
		return nil
	}

	// Convert the first replication's destinations
	replication := output.Replications[0]
	specDestinations := make([]*svcapitypes.DestinationToCreate, 0, len(replication.Destinations))
	statusDestinations := make([]*svcapitypes.Destination, 0, len(replication.Destinations))

	for _, dest := range replication.Destinations {
		specDest := &svcapitypes.DestinationToCreate{}
		if dest.Region != nil {
			specDest.Region = dest.Region
		}
		if dest.FileSystemId != nil {
			specDest.FileSystemID = dest.FileSystemId
		}
		if dest.RoleArn != nil {
			specDest.RoleARN = dest.RoleArn
		}

		specDestinations = append(specDestinations, specDest)

		// Status
		statusDest := &svcapitypes.Destination{}
		if dest.Region != nil {
			statusDest.Region = dest.Region
		}
		if dest.FileSystemId != nil {
			statusDest.FileSystemID = dest.FileSystemId
		}
		if dest.OwnerId != nil {
			statusDest.OwnerID = dest.OwnerId
		}
		if dest.LastReplicatedTimestamp != nil {
			statusDest.LastReplicatedTimestamp = &metav1.Time{Time: *dest.LastReplicatedTimestamp}
		}
		if dest.Status != "" {
			statusDest.Status = aws.String(string(dest.Status))
		}
		if dest.StatusMessage != nil {
			statusDest.StatusMessage = dest.StatusMessage
		}
		if dest.RoleArn != nil {
			statusDest.RoleARN = dest.RoleArn
		}
		statusDestinations = append(statusDestinations, statusDest)
	}

	r.Spec.ReplicationConfiguration = specDestinations
	r.Status.ReplicationConfigurationStatus = statusDestinations
	return nil
}

// syncReplicationConfiguration manages replication configuration state for the file system
func (rm *resourceManager) syncReplicationConfiguration(ctx context.Context, desired *resource, latest *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncReplicationConfiguration")
	defer func() { exit(err) }()

	desiredHasReplication := len(desired.ko.Spec.ReplicationConfiguration) == 1
	latestHasReplication := len(latest.ko.Spec.ReplicationConfiguration) == 1

	// EFS replication configuration does not currently suppport more than one destination
	if len(desired.ko.Spec.ReplicationConfiguration) > 1 {
		rlog.Info("Invalid replication configuration",
			"replicationCount", len(desired.ko.Spec.ReplicationConfiguration),
			"replicationDestinations", desired.ko.Spec.ReplicationConfiguration)
		msg := fmt.Errorf("EFS replication configuration only allows one replica, got %d replicas",
			len(desired.ko.Spec.ReplicationConfiguration))
		return ackerr.NewTerminalError(msg)
	}

	if desiredHasReplication && latestHasReplication {
		// Delete existing then create new (AWS doesn't support direct update)
		err = rm.deleteReplicationConfiguration(ctx, latest)
		if err != nil {
			return err
		}
		return requeueWaitReplicationConfiguration
	}

	// If desired has replication but latest doesn't, create it
	if desiredHasReplication && !latestHasReplication {
		err = rm.createReplicationConfiguration(ctx, desired)
		if err != nil {
			return err
		}
	}

	// If desired doesn't have replication but latest does, delete it
	if !desiredHasReplication && latestHasReplication {
		err = rm.deleteReplicationConfiguration(ctx, desired)
		if err != nil {
			return err
		}
	}

	return requeueWaitReplicationConfiguration
}

// createReplicationConfiguration creates replication configuration for the file system
func (rm *resourceManager) createReplicationConfiguration(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.createReplicationConfiguration")
	defer func() { exit(err) }()

	// Convert DestinationToCreate to Smithy SDK format
	destinations := make([]svcsdktypes.DestinationToCreate, 0, len(r.ko.Spec.ReplicationConfiguration))
	for _, dest := range r.ko.Spec.ReplicationConfiguration {
		sdkDest := svcsdktypes.DestinationToCreate{}
		if dest.Region != nil {
			sdkDest.Region = dest.Region
		}
		if dest.AvailabilityZoneName != nil {
			sdkDest.AvailabilityZoneName = dest.AvailabilityZoneName
		}
		if dest.KMSKeyID != nil {
			sdkDest.KmsKeyId = dest.KMSKeyID
		}
		if dest.FileSystemID != nil {
			sdkDest.FileSystemId = dest.FileSystemID
		}
		if dest.RoleARN != nil {
			sdkDest.RoleArn = dest.RoleARN
		}
		destinations = append(destinations, sdkDest)
	}

	_, err = rm.sdkapi.CreateReplicationConfiguration(
		ctx,
		&svcsdk.CreateReplicationConfigurationInput{
			SourceFileSystemId: r.ko.Status.FileSystemID,
			Destinations:       destinations,
		},
	)
	rm.metrics.RecordAPICall("UPDATE", "CreateReplicationConfiguration", err)
	return err
}

// deleteReplicationConfiguration deletes replication configuration for the file system
func (rm *resourceManager) deleteReplicationConfiguration(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.deleteReplicationConfiguration")
	defer func() { exit(err) }()

	_, err = rm.sdkapi.DeleteReplicationConfiguration(
		ctx,
		&svcsdk.DeleteReplicationConfigurationInput{
			SourceFileSystemId: r.ko.Status.FileSystemID,
		},
	)
	rm.metrics.RecordAPICall("UPDATE", "DeleteReplicationConfiguration", err)
	return err
}
