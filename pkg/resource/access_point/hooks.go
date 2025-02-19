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
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"

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

// requeueWaitState returns a `ackrequeue.RequeueNeededAfter` struct
// explaining the accesspoint cannot be modified until it reaches an active status.
func requeueWaitState(r *resource) *ackrequeue.RequeueNeededAfter {
	if r.ko.Status.LifeCycleState == nil {
		return nil
	}
	status := *r.ko.Status.LifeCycleState
	return ackrequeue.NeededAfter(
		fmt.Errorf("accesspoint in '%s' state, requeuing until accesspoint is '%s'",
			status, svcapitypes.LifeCycleState_available),
		time.Second*10,
	)
}

// accessPointActive returns true if the supplied accessPoint is in an active status
func accessPointActive(r *resource) bool {
	if r.ko.Status.LifeCycleState == nil {
		return false
	}
	cs := *r.ko.Status.LifeCycleState
	lifeCycleState := string(svcapitypes.LifeCycleState_available)
	return cs == lifeCycleState
}

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

	return desired, nil
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
