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

package access_point

import (
	"context"
	"fmt"
	"reflect"
	"time"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	corev1 "k8s.io/api/core/v1"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/efs-controller/pkg/tags"
)

// getIdempotencyToken returns a unique string to be used in certain API calls
// to ensure no replay of the call.
func getIdempotencyToken() string {
	t := time.Now().UTC()
	return t.Format("20060102150405000000")
}

// Ideally, a part of this code needs to be generated.. However since the
// tags packge is not imported, we can't call it directly from sdk.go. We
// have to do this Go-fu to make it work.
var syncTags = tags.SyncTags

// lifeCycleState returns the accesspoint's observed lifecycle state, or
// "unknown" when the state has not been observed yet.
func lifeCycleState(r *resource) string {
	if r == nil || r.ko.Status.LifeCycleState == nil {
		return "unknown"
	}
	return *r.ko.Status.LifeCycleState
}

// requeueWaitState returns a `ackrequeue.RequeueNeededAfter` struct
// explaining the accesspoint cannot be modified until it reaches an active status.
func requeueWaitState(r *resource) *ackrequeue.RequeueNeededAfter {
	return ackrequeue.NeededAfter(
		fmt.Errorf("accesspoint in '%s' state, requeuing until accesspoint is '%s'",
			lifeCycleState(r), svcapitypes.LifeCycleState_available),
		ackrequeue.DefaultRequeueAfterDuration,
	)
}

// accessPointActive returns true if the supplied accessPoint is in an active status
func accessPointActive(r *resource) bool {
	return lifeCycleState(r) == string(svcapitypes.LifeCycleState_available)
}

// accessPointInErrorState returns true if the supplied accessPoint is in a state
// it cannot recover from in place.
func accessPointInErrorState(r *resource) bool {
	return lifeCycleState(r) == string(svcapitypes.LifeCycleState_error)
}

// errAccessPointInErrorState is terminal: the accesspoint has to be recreated.
var errAccessPointInErrorState = fmt.Errorf(
	"accesspoint is in '%s' state and cannot be modified; delete and recreate it",
	svcapitypes.LifeCycleState_error,
)

// customUpdateAccessPoint updates the access point
func (rm *resourceManager) customUpdateAccessPoint(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() { exit(err) }()

	// An 'error' accesspoint never becomes modifiable again, so fail fast
	// instead of requeuing forever.
	if accessPointInErrorState(latest) {
		return nil, ackerr.NewTerminalError(errAccessPointInErrorState)
	}

	updated = rm.concreteResource(desired.DeepCopy())
	updated.SetStatus(latest)
	if !accessPointActive(updated) {
		msg := fmt.Sprintf("accesspoint cannot be modified until it is '%s'",
			svcapitypes.LifeCycleState_available)
		reason := lifeCycleState(updated)
		ackcondition.SetSynced(updated, corev1.ConditionFalse, &msg, &reason)
		return updated, requeueWaitState(updated)
	}

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

	return updated, nil
}

var (
	defaultRootDirectory = svcapitypes.RootDirectory{
		Path: aws.String("/"),
	}
)

func customPreCompare(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) {
	if a.ko.Spec.RootDirectory == nil {
		a.ko.Spec.RootDirectory = &defaultRootDirectory
	}
	// PosixUser.SecondaryGIDs is not generated, we need to compare it manually
	if ackcompare.HasNilDifference(a.ko.Spec.PosixUser, b.ko.Spec.PosixUser) {
		delta.Add("Spec.PosixUser", a.ko.Spec.PosixUser, b.ko.Spec.PosixUser)
	} else if a.ko.Spec.PosixUser != nil && b.ko.Spec.PosixUser != nil {
		if len(a.ko.Spec.PosixUser.SecondaryGIDs) != len(b.ko.Spec.PosixUser.SecondaryGIDs) {
			delta.Add("Spec.PosixUser.SecondaryGIDs", a.ko.Spec.PosixUser.GID, b.ko.Spec.PosixUser.GID)
		} else if len(a.ko.Spec.PosixUser.SecondaryGIDs) > 0 && !reflect.DeepEqual(a.ko.Spec.PosixUser.SecondaryGIDs, b.ko.Spec.PosixUser.SecondaryGIDs) {
			delta.Add("Spec.PosixUser.GID", a.ko.Spec.PosixUser.GID, b.ko.Spec.PosixUser.GID)
		}
	}
}
