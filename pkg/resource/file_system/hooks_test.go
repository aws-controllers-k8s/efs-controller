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

// createFileSystemAPI serves CreateFileSystem with the given status code and
// body, and records the creation token of each request.
type createFileSystemAPI struct {
	status int
	header http.Header
	body   string
	tokens []string
}

func (a *createFileSystemAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var in struct{ CreationToken string }
	b, _ := io.ReadAll(req.Body)
	_ = json.Unmarshal(b, &in)
	a.tokens = append(a.tokens, in.CreationToken)
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

func newFileSystem(id *string) *resource {
	ko := &svcapitypes.FileSystem{}
	ko.UID = testUID
	ko.Status.FileSystemID = id
	return &resource{ko}
}

func TestSdkCreate_SendsTheCRUIDAsCreationToken(t *testing.T) {
	api := &createFileSystemAPI{
		status: http.StatusCreated,
		body:   `{"FileSystemId":"fs-0123456789abcdef0","LifeCycleState":"creating"}`,
	}
	rm := newTestManager(t, api)

	for range 2 {
		if _, err := rm.sdkCreate(context.Background(), newFileSystem(nil)); err != nil {
			t.Fatalf("sdkCreate: %v", err)
		}
	}
	if len(api.tokens) != 2 || api.tokens[0] != testUID || api.tokens[1] != testUID {
		t.Errorf("creation tokens = %q, want the CR UID %q twice", api.tokens, testUID)
	}
}

func TestSdkCreate_WithoutUIDLeavesTheTokenToTheSDK(t *testing.T) {
	api := &createFileSystemAPI{
		status: http.StatusCreated,
		body:   `{"FileSystemId":"fs-0123456789abcdef0","LifeCycleState":"creating"}`,
	}
	rm := newTestManager(t, api)

	r := newFileSystem(nil)
	r.ko.UID = ""
	if _, err := rm.sdkCreate(context.Background(), r); err != nil {
		t.Fatalf("sdkCreate: %v", err)
	}
	if len(api.tokens) != 1 || api.tokens[0] == "" {
		t.Errorf("creation tokens = %q, want one token filled in by the SDK", api.tokens)
	}
}

func TestSdkCreate_FileSystemAlreadyCreated(t *testing.T) {
	api := &createFileSystemAPI{
		status: http.StatusConflict,
		header: http.Header{"X-Amzn-Errortype": {"FileSystemAlreadyExists"}},
		body: `{"ErrorCode":"FileSystemAlreadyExists",` +
			`"Message":"File system 'fs-0123456789abcdef0' already exists with creation token '` + testUID + `'",` +
			`"FileSystemId":"fs-0123456789abcdef0"}`,
	}
	rm := newTestManager(t, api)

	created, err := rm.sdkCreate(context.Background(), newFileSystem(nil))
	if err != nil {
		t.Fatalf("sdkCreate: %v", err)
	}
	if id := created.ko.Status.FileSystemID; id == nil || *id != "fs-0123456789abcdef0" {
		t.Errorf("FileSystemID = %v, want fs-0123456789abcdef0", id)
	}
}

func TestSdkCreate_OtherErrorsAreReturned(t *testing.T) {
	api := &createFileSystemAPI{
		status: http.StatusBadRequest,
		header: http.Header{"X-Amzn-Errortype": {"BadRequest"}},
		body:   `{"ErrorCode":"BadRequest","Message":"invalid"}`,
	}
	rm := newTestManager(t, api)

	created, err := rm.sdkCreate(context.Background(), newFileSystem(nil))
	if err == nil {
		t.Fatalf("sdkCreate returned %v and no error, want the BadRequest error", created)
	}
}
