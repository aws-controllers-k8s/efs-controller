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

package tags

import (
	"context"

	"github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"

	svcsdk "github.com/aws/aws-sdk-go-v2/service/efs"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
)

// Ideally, a part of this code needs to be generated, the other part
// needs to be implemented in a different repository (runtime or pkg).
//
// Few things to node:
// - Some AWS APIs support map[string]string for tags, while others support []*Tag.
// - Some AWS APIs have a limit on the number of tags that can be associated with a resource.
// - We can call a few different names and ways to tag and untag resources.
//   - Some allow to add/remove one tag at a time, while others allow to add/remove multiple tags at once.
//   - Some have a seperate API to list tags, while others return tags as part of the describe response.
//   - Even when the API model states that a Describe response will contain tags, the actual response may
//     not contain any tags. And users are expected to call a seperate ListTags API to get the tags.
//
// - Noting a few diffrent API names:
//   - TagResource, UntagResource, ListTagsForResource
//   - CreateTags, DeleteTags, ListTags
//   - AddTags, RemoveTags, ListTags

// Below are some abstractions that can be used to abstract the implementation details
// of tagging and untagging resources.

type metricsRecorder interface {
	RecordAPICall(opType string, opID string, err error)
}

type tagsClient interface {
	TagResource(context.Context, *svcsdk.TagResourceInput, ...func(*svcsdk.Options)) (*svcsdk.TagResourceOutput, error)
	ListTagsForResource(context.Context, *svcsdk.ListTagsForResourceInput, ...func(*svcsdk.Options)) (*svcsdk.ListTagsForResourceOutput, error)
	UntagResource(context.Context, *svcsdk.UntagResourceInput, ...func(*svcsdk.Options)) (*svcsdk.UntagResourceOutput, error)
}

// syncTags examines the Tags in the supplied Resource and calls the
// TagResource and UntagResource APIs to ensure that the set of
// associated Tags stays in sync with the Resource.Spec.Tags
func SyncTags(
	ctx context.Context,
	client tagsClient,
	mr metricsRecorder,
	resourceID string,
	aTags []*v1alpha1.Tag,
	bTags []*v1alpha1.Tag,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTags")
	defer func() { exit(err) }()

	desiredTags := map[string]*string{}
	for _, t := range aTags {
		desiredTags[*t.Key] = t.Value
	}
	existingTags := map[string]*string{}
	for _, t := range bTags {
		existingTags[*t.Key] = t.Value
	}

	toAdd := map[string]*string{}
	toDelete := []string{}

	for k, v := range desiredTags {
		if ev, found := existingTags[k]; !found || *ev != *v {
			toAdd[k] = v
		}
	}

	for k := range existingTags {
		if _, found := desiredTags[k]; !found {
			deleteKey := k
			toDelete = append(toDelete, deleteKey)
		}
	}

	if len(toAdd) > 0 {
		for k, v := range toAdd {
			rlog.Debug("adding tag to resource", "key", k, "value", *v)
		}
		if err = addTags(
			ctx,
			client,
			mr,
			resourceID,
			toAdd,
		); err != nil {
			return err
		}
	}
	if len(toDelete) > 0 {
		for _, k := range toDelete {
			rlog.Debug("removing tag from resource", "key", k)
		}
		if err = removeTags(
			ctx,
			client,
			mr,
			resourceID,
			toDelete,
		); err != nil {
			return err
		}
	}

	return nil
}

func addTags(
	ctx context.Context,
	client tagsClient,
	mr metricsRecorder,
	resourceID string,
	tags map[string]*string,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.addTag")
	defer func() { exit(err) }()

	sdkTags := toSdkTags(tags)
	input := &svcsdk.TagResourceInput{
		ResourceId: &resourceID,
		Tags:       sdkTags,
	}

	_, err = client.TagResource(ctx, input)
	mr.RecordAPICall("UPDATE", "TagResource", err)
	return err
}

func removeTags(
	ctx context.Context,
	client tagsClient,
	mr metricsRecorder,
	resourceID string,
	tagKeys []string,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.removeTag")
	defer func() { exit(err) }()

	input := &svcsdk.UntagResourceInput{
		ResourceId: &resourceID,
		TagKeys:    tagKeys,
	}
	_, err = client.UntagResource(ctx, input)
	mr.RecordAPICall("UPDATE", "UntagResource", err)
	return err
}

func toSdkTags(tags map[string]*string) []svcsdktypes.Tag {
	sdkTags := []svcsdktypes.Tag{}
	for key, val := range tags {
		sdkTags = append(sdkTags, svcsdktypes.Tag{Key: &key, Value: val})
	}
	return sdkTags
}
