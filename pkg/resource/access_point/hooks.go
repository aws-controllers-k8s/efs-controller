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
	"errors"
	"fmt"
	"reflect"
	"time"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/efs/types"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/efs-controller/pkg/tags"
)

// getIdempotencyToken returns the ClientToken for the access point of r.
//
// The token is the UID of the CR, so that every Create for the CR sends the
// same one. When the controller loses the response of a Create that went
// through, for example because it was restarted before the access point ID
// reached the status, the next Create gets AccessPointAlreadyExists with that
// ID instead of creating a second access point (see accessPointAlreadyCreated).
//
// EFS frees the token once the access point is deleted, so the Create that
// replaces an access point deleted outside of the controller makes a new one.
func getIdempotencyToken(r *resource) string {
	return string(r.ko.UID)
}

// accessPointAlreadyCreated returns desired with the ID of its access point in
// the status when err is the AccessPointAlreadyExists error, and nil for any
// other error.
//
// CreateAccessPoint returns that error for a client token it has already
// seen. The token comes from getIdempotencyToken, so the access point was
// created for this CR, by an earlier Create whose response was lost. The
// returned resource carries only the ID; the reconciler reads the rest of the
// access point right after Create.
func accessPointAlreadyCreated(desired *resource, err error) *resource {
	var exists *svcsdktypes.AccessPointAlreadyExists
	if !errors.As(err, &exists) || exists.AccessPointId == nil {
		return nil
	}
	ko := desired.ko.DeepCopy()
	ko.Status.AccessPointID = exists.AccessPointId
	return &resource{ko}
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
