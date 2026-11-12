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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/efs"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
)

const testUID = "8b5e3e0a-6f0c-4a7e-9a55-0f7d5c1d2e3f"

// createAccessPointAPI serves CreateAccessPoint with the given status code and
// body, and records the client token of each request.
type createAccessPointAPI struct {
	status int
	header http.Header
	body   string
	tokens []string
}

func (a *createAccessPointAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var in struct{ ClientToken string }
	b, _ := io.ReadAll(req.Body)
	_ = json.Unmarshal(b, &in)
	a.tokens = append(a.tokens, in.ClientToken)
	for k, v := range a.header {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(a.status)
	_, _ = io.WriteString(w, a.body)
}

func newTestManager(t *testing.T, api http.Handler) *resourceManager {
	srv := httptest.NewServer(api)
	t.Cleanup(srv.Close)
	return &resourceManager{
		metrics: ackmetrics.NewMetrics("efs"),
		sdkapi: svcsdk.New(svcsdk.Options{
			Region:       "us-west-2",
			BaseEndpoint: aws.String(srv.URL),
			Credentials:  aws.AnonymousCredentials{},
			HTTPClient:   srv.Client(),
		}),
	}
}

func newAccessPoint(id *string) *resource {
	ko := &svcapitypes.AccessPoint{}
	ko.UID = testUID
	ko.Spec.FileSystemID = aws.String("fs-0123456789abcdef0")
	ko.Status.AccessPointID = id
	return &resource{ko}
}

func TestSdkCreate_SendsTheCRUIDAsClientToken(t *testing.T) {
	api := &createAccessPointAPI{
		status: http.StatusCreated,
		body:   `{"AccessPointId":"fsap-0123456789abcdef0","FileSystemId":"fs-0123456789abcdef0","LifeCycleState":"creating"}`,
	}
	rm := newTestManager(t, api)

	for range 2 {
		if _, err := rm.sdkCreate(context.Background(), newAccessPoint(nil)); err != nil {
			t.Fatalf("sdkCreate: %v", err)
		}
	}
	if len(api.tokens) != 2 || api.tokens[0] != testUID || api.tokens[1] != testUID {
		t.Errorf("client tokens = %q, want the CR UID %q twice", api.tokens, testUID)
	}
}

func TestSdkCreate_AccessPointAlreadyCreated(t *testing.T) {
	api := &createAccessPointAPI{
		status: http.StatusConflict,
		header: http.Header{"X-Amzn-Errortype": {"AccessPointAlreadyExists"}},
		body: `{"ErrorCode":"AccessPointAlreadyExists",` +
			`"Message":"Access point 'fsap-0123456789abcdef0' already exists with client token '` + testUID + `'",` +
			`"AccessPointId":"fsap-0123456789abcdef0"}`,
	}
	rm := newTestManager(t, api)

	created, err := rm.sdkCreate(context.Background(), newAccessPoint(nil))
	if err != nil {
		t.Fatalf("sdkCreate: %v", err)
	}
	if id := created.ko.Status.AccessPointID; id == nil || *id != "fsap-0123456789abcdef0" {
		t.Errorf("AccessPointID = %v, want fsap-0123456789abcdef0", id)
	}
}

func TestSdkCreate_OtherErrorsAreReturned(t *testing.T) {
	api := &createAccessPointAPI{
		status: http.StatusBadRequest,
		header: http.Header{"X-Amzn-Errortype": {"BadRequest"}},
		body:   `{"ErrorCode":"BadRequest","Message":"invalid"}`,
	}
	rm := newTestManager(t, api)

	created, err := rm.sdkCreate(context.Background(), newAccessPoint(nil))
	if err == nil {
		t.Fatalf("sdkCreate returned %v and no error, want the BadRequest error", created)
	}
}
