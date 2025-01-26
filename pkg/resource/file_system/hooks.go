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
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/efs"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/efs/types"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/efs-controller/pkg/tags"
	"github.com/aws/aws-sdk-go/aws"
)

// getIdempotencyToken returns a unique string to be used in certain API calls
// to ensure no replay of the call.
func getIdempotencyToken() string {
	t := time.Now().UTC()
	return t.Format("20060102150405000000")
}

var (
	// TerminalStatuses are the status strings that are terminal states for a
	// filesystem.
	TerminalStatuses = []string{
		string(svcapitypes.LifeCycleState_error),
		string(svcapitypes.LifeCycleState_deleted),
		string(svcapitypes.LifeCycleState_deleting),
	}
)

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
		time.Second*10,
	)
}

// filesystemHasTerminalStatus returns whether the supplied filesystem is in a
// terminal state
func filesystemHasTerminalStatus(r *resource) bool {
	if r.ko.Status.LifeCycleState == nil {
		return false
	}
	cs := *r.ko.Status.LifeCycleState
	for _, s := range TerminalStatuses {
		if cs == s {
			return true
		}
	}
	return false
}

// filesystemActive returns true if the supplied filesystem is in an active status
func filesystemActive(r *resource) bool {
	if r.ko.Status.LifeCycleState == nil {
		return false
	}
	cs := *r.ko.Status.LifeCycleState
	lifeCycleState := string(svcapitypes.LifeCycleState_available)
	return cs == lifeCycleState
}

// filesystemCreating returns true if the supplied filesystem is in the process of
// being created
func filesystemCreating(r *resource) bool {
	if r.ko.Status.LifeCycleState == nil {
		return false
	}
	cs := *r.ko.Status.LifeCycleState
	lifeCycleState := string(svcapitypes.LifeCycleState_creating)
	return cs == lifeCycleState
}

// filesystemDeleting returns true if the supplied filesystem is in the process of
// being deleted
func filesystemDeleting(r *resource) bool {
	if r.ko.Status.LifeCycleState == nil {
		return false
	}
	cs := *r.ko.Status.LifeCycleState
	lifeCycleState := string(svcapitypes.LifeCycleState_deleting)
	return cs == lifeCycleState
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

func (rm *resourceManager) syncBackupPolicy(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncBackupPolicy")
	defer func() { exit(err) }()

	_, err = rm.sdkapi.PutBackupPolicy(
		ctx,
		&svcsdk.PutBackupPolicyInput{
			FileSystemId: r.ko.Status.FileSystemID,
			BackupPolicy: &svcsdktypes.BackupPolicy{
				Status: svcsdktypes.Status(*r.ko.Spec.BackupPolicy.Status),
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
}
