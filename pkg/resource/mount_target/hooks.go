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

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/efs"
	corev1 "k8s.io/api/core/v1"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
)

// errMountTargetInErrorState is terminal: the mounttarget has to be recreated.
var errMountTargetInErrorState = fmt.Errorf(
	"mounttarget is in '%s' state and cannot be modified; delete and recreate it",
	svcapitypes.LifeCycleState_error,
)

// lifeCycleState returns the mounttarget's observed lifecycle state, or
// "unknown" when the state has not been observed yet.
func lifeCycleState(r *resource) string {
	if r == nil || r.ko.Status.LifeCycleState == nil {
		return "unknown"
	}
	return *r.ko.Status.LifeCycleState
}

// requeueWaitState returns a `ackrequeue.RequeueNeededAfter` struct
// explaining the mounttarget cannot be modified until it reaches an active status.
func requeueWaitState(r *resource) *ackrequeue.RequeueNeededAfter {
	return ackrequeue.NeededAfter(
		fmt.Errorf("mounttarget in '%s' state, requeuing until mounttarget is '%s'",
			lifeCycleState(r), svcapitypes.LifeCycleState_available),
		ackrequeue.DefaultRequeueAfterDuration,
	)
}

// mounttargetActive returns true if the supplied mounttarget is in an active status
func mountTargetActive(r *resource) bool {
	return lifeCycleState(r) == string(svcapitypes.LifeCycleState_available)
}

// mountTargetInErrorState returns true if the supplied mounttarget is in a state
// it cannot recover from in place.
func mountTargetInErrorState(r *resource) bool {
	return lifeCycleState(r) == string(svcapitypes.LifeCycleState_error)
}

// setResourceDefaults queries the EFS API for the current state of the
// fields that are not returned by the ReadOne or List APIs.
func (rm *resourceManager) setResourceAdditionalFields(ctx context.Context, r *svcapitypes.MountTarget) error {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setResourceAdditionalFields")
	defer exit(nil)

	securityGroups, err := rm.getSecurityGroups(ctx, r)
	if err != nil {
		exit(err)
		return err
	}

	r.Spec.SecurityGroups = make([]*string, len(securityGroups))
	for i := range securityGroups {
		securityGroup := securityGroups[i]
		r.Spec.SecurityGroups[i] = &securityGroup
	}

	return nil
}

// getSecurityGroups returns the security groups for the mount target
func (rm *resourceManager) getSecurityGroups(ctx context.Context, r *svcapitypes.MountTarget) (_ []string, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getSecurityGroups")
	defer func() { exit(err) }()

	var output *svcsdk.DescribeMountTargetSecurityGroupsOutput
	output, err = rm.sdkapi.DescribeMountTargetSecurityGroups(
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

	securityGroups := make([]string, 0, len(r.ko.Spec.SecurityGroups))
	for _, sg := range r.ko.Spec.SecurityGroups {
		securityGroups = append(securityGroups, *sg)
	}

	_, err = rm.sdkapi.ModifyMountTargetSecurityGroups(
		ctx,
		&svcsdk.ModifyMountTargetSecurityGroupsInput{
			MountTargetId:  r.ko.Status.MountTargetID,
			SecurityGroups: securityGroups,
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
	// An 'error' mounttarget never becomes modifiable again, so fail fast
	// instead of requeuing forever.
	if mountTargetInErrorState(latest) {
		return nil, ackerr.NewTerminalError(errMountTargetInErrorState)
	}

	updated = rm.concreteResource(desired.DeepCopy())
	updated.SetStatus(latest)
	if !mountTargetActive(updated) {
		msg := fmt.Sprintf("mounttarget cannot be modified until it is '%s'",
			svcapitypes.LifeCycleState_available)
		reason := lifeCycleState(updated)
		ackcondition.SetSynced(updated, corev1.ConditionFalse, &msg, &reason)
		return updated, requeueWaitState(updated)
	}

	if delta.DifferentAt("Spec.SecurityGroups") {
		err := rm.putSecurityGroups(ctx, updated)
		if err != nil {
			return nil, err
		}
	}

	return updated, nil
}
