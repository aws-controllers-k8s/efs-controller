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
	"testing"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
)

const testFileSystemARN = "arn:aws:elasticfilesystem:us-west-2:111122223333:file-system/fs-0b07fc1c44310f532"

func policyResource(policy *string, arn *string) *resource {
	ko := &svcapitypes.FileSystem{}
	ko.Spec.Policy = policy
	if arn != nil {
		a := ackv1alpha1.AWSResourceName(*arn)
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{ARN: &a}
	}
	return &resource{ko: ko}
}

func TestComparePolicy(t *testing.T) {
	desired := `{"Version":"2012-10-17","Statement":[{"Sid":"lt2","Effect":"Allow","Principal":{"AWS":"*"},"Action":["elasticfilesystem:ClientMount"],"Condition":{"Bool":{"elasticfilesystem:AccessedViaMountTarget":"true"}}}]}`
	observed := `{
  "Version" : "2012-10-17",
  "Statement" : [ {
    "Sid" : "lt2",
    "Effect" : "Allow",
    "Principal" : {
      "AWS" : "*"
    },
    "Action" : "elasticfilesystem:ClientMount",
    "Resource" : "` + testFileSystemARN + `",
    "Condition" : {
      "Bool" : {
        "elasticfilesystem:AccessedViaMountTarget" : "true"
      }
    }
  } ]
}`
	otherResource := `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Principal":{"AWS":"*"},"Action":"elasticfilesystem:ClientMount","Resource":"arn:aws:elasticfilesystem:us-west-2:111122223333:file-system/fs-other"}}`
	singleStatement := `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Principal":{"AWS":"*"},"Action":"elasticfilesystem:ClientMount"}}`
	singleObserved := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":"elasticfilesystem:ClientMount","Resource":"` + testFileSystemARN + `"}]}`

	tests := []struct {
		name      string
		desired   *string
		latest    *string
		latestARN *string
		wantDiff  bool
	}{
		{"both nil", nil, nil, aws.String(testFileSystemARN), false},
		{"desired nil", nil, aws.String(observed), aws.String(testFileSystemARN), true},
		{"latest nil", aws.String(desired), nil, aws.String(testFileSystemARN), true},
		{"omitted resource defaults to file system ARN", aws.String(desired), aws.String(observed), aws.String(testFileSystemARN), false},
		{"single statement object", aws.String(singleStatement), aws.String(singleObserved), aws.String(testFileSystemARN), false},
		{"explicit resource differs", aws.String(otherResource), aws.String(singleObserved), aws.String(testFileSystemARN), true},
		{"no ARN in status", aws.String(desired), aws.String(observed), nil, true},
		{"invalid desired policy", aws.String("not-json"), aws.String(observed), aws.String(testFileSystemARN), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := policyResource(tt.desired, nil)
			b := policyResource(tt.latest, tt.latestARN)
			delta := ackcompare.NewDelta()
			comparePolicy(delta, a, b)
			if got := delta.DifferentAt("Spec.Policy"); got != tt.wantDiff {
				t.Errorf("DifferentAt(Spec.Policy) = %v, want %v", got, tt.wantDiff)
			}
			if tt.desired != nil && a.ko.Spec.Policy != tt.desired {
				t.Errorf("desired policy was mutated")
			}
		})
	}
}
