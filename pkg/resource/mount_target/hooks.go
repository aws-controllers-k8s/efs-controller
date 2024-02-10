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

package mount_target

import (
	"context"
	"fmt"
	"time"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/efs"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
)

var (
	// TerminalStatuses are the status strings that are terminal states for a
	// mounttarget.
	TerminalStatuses = []string{
		string(svcapitypes.LifeCycleState_error),
		string(svcapitypes.LifeCycleState_deleted),
		string(svcapitypes.LifeCycleState_deleting),
	}
)

// requeueWaitState returns a `ackrequeue.RequeueNeededAfter` struct
// explaining the mounttarget cannot be modified until it reaches an active status.
func requeueWaitState(r *resource) *ackrequeue.RequeueNeededAfter {
	if r.ko.Status.LifeCycleState == nil {
		return nil
	}
	status := *r.ko.Status.LifeCycleState
	return ackrequeue.NeededAfter(
		fmt.Errorf("mounttarget in '%s' state, requeuing until mounttarget is '%s'",
			status, svcapitypes.LifeCycleState_available),
		time.Second*10,
	)
}

// mounttargetActive returns true if the supplied mounttarget is in an active status
func mountTargetActive(r *resource) bool {
	if r.ko.Status.LifeCycleState == nil {
		return false
	}
	cs := *r.ko.Status.LifeCycleState
	return cs == string(svcapitypes.LifeCycleState_available)
}

// mounttargetCreating returns true if the supplied mounttarget is in the process of
// being created
func mountTargetCreating(r *resource) bool {
	if r.ko.Status.LifeCycleState == nil {
		return false
	}
	cs := *r.ko.Status.LifeCycleState
	return cs == string(svcapitypes.LifeCycleState_creating)
}

// setResourceDefaults queries the EFS API for the current state of the
// fields that are not returned by the ReadOne or List APIs.
func (rm *resourceManager) setResourceAdditionalFields(ctx context.Context, r *svcapitypes.MountTarget) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setResourceAdditionalFields")
	defer func() { exit(err) }()

	r.Spec.SecurityGroups, err = rm.getSecurityGroups(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

// getSecurityGroups returns the security groups for the mount target
func (rm *resourceManager) getSecurityGroups(ctx context.Context, r *svcapitypes.MountTarget) (_ []*string, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getSecurityGroups")
	defer func() { exit(err) }()

	var output *svcsdk.DescribeMountTargetSecurityGroupsOutput
	output, err = rm.sdkapi.DescribeMountTargetSecurityGroupsWithContext(
		ctx,
		&svcsdk.DescribeMountTargetSecurityGroupsInput{
			MountTargetId: r.Status.MountTargetID,
		},
	)
	rm.metrics.RecordAPICall("GET", "DescribeMountTargetSecurityGroups", err)
	if err != nil {
		return nil, err
	}

	return output.SecurityGroups, nil
}

// putSecurityGroups updates the security groups for the mount target
func (rm *resourceManager) putSecurityGroups(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncPolicy")
	defer func() { exit(err) }()

	_, err = rm.sdkapi.ModifyMountTargetSecurityGroupsWithContext(
		ctx,
		&svcsdk.ModifyMountTargetSecurityGroupsInput{
			MountTargetId:  r.ko.Status.MountTargetID,
			SecurityGroups: r.ko.Spec.SecurityGroups,
		},
	)
	rm.metrics.RecordAPICall("UPDATE", "ModifyMountTargetSecurityGroups", err)
	return err
}

// customUpdateMountTarget updates the mount target security groups
func (rm *resourceManager) customUpdateMountTarget(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() { exit(err) }()

	if delta.DifferentAt("Spec.SecurityGroups") {
		err := rm.putSecurityGroups(ctx, desired)
		if err != nil {
			return nil, err
		}
	}

	return desired, nil
}
