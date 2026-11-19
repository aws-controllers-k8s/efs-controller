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
	"strings"
	"testing"

	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
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

// efsAPI serves the EFS API with a function that gets each request and returns
// the status code, the error type for an error, and the body.
type efsAPI func(req *http.Request) (status int, errType string, body string)

func (f efsAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	status, errType, body := f(req)
	if errType != "" {
		w.Header().Set("X-Amzn-Errortype", errType)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

func TestRequeueWaitState_WithoutStateStillRequeues(t *testing.T) {
	requeue := requeueWaitState(newFileSystem(nil))
	if requeue == nil || requeue.Duration() <= 0 {
		t.Fatalf("requeueWaitState = %#v, want a requeue after a non-zero duration", requeue)
	}
}

// The status has no ID, for example because the controller lost the response
// of a Create that went through. The read finds the file system by the CR UID
// as its creation token, so a delete of the CR deletes it too.
func TestSdkFind_WithoutIDFindsTheFileSystemByCreationToken(t *testing.T) {
	const id = "fs-0123456789abcdef0"
	fileSystem := `{"FileSystems":[{"FileSystemId":"` + id + `","CreationToken":"` + testUID + `",` +
		`"LifeCycleState":"creating","PerformanceMode":"generalPurpose"}]}`
	rm := newTestManager(t, efsAPI(func(req *http.Request) (int, string, string) {
		q := req.URL.Query()
		switch {
		case req.URL.Path == "/2015-02-01/file-systems" && q.Get("CreationToken") == testUID:
			return http.StatusOK, "", fileSystem
		case req.URL.Path == "/2015-02-01/file-systems" && q.Get("FileSystemId") == id:
			return http.StatusOK, "", fileSystem
		// The policies and the replication configuration of the file system.
		case strings.HasSuffix(req.URL.Path, "/policy"):
			return http.StatusNotFound, "PolicyNotFound", `{"ErrorCode":"PolicyNotFound","Message":"none"}`
		case strings.HasSuffix(req.URL.Path, "/backup-policy"):
			return http.StatusOK, "", `{"BackupPolicy":{"Status":"DISABLED"}}`
		case strings.HasSuffix(req.URL.Path, "/lifecycle-configuration"):
			return http.StatusOK, "", `{"LifecyclePolicies":[]}`
		case strings.HasSuffix(req.URL.Path, "/replication-configurations"):
			return http.StatusOK, "", `{"Replications":[]}`
		}
		t.Errorf("unexpected request %s %s", req.Method, req.URL)
		return http.StatusBadRequest, "BadRequest", `{"ErrorCode":"BadRequest"}`
	}))

	found, err := rm.sdkFind(context.Background(), newFileSystem(nil))
	if err != nil {
		t.Fatalf("sdkFind: %v", err)
	}
	if got := found.ko.Status.FileSystemID; got == nil || *got != id {
		t.Errorf("FileSystemID = %v, want %s", got, id)
	}
}

func TestSdkFind_WithoutIDAndNoFileSystemWithTheToken_NotFound(t *testing.T) {
	var tokens []string
	rm := newTestManager(t, efsAPI(func(req *http.Request) (int, string, string) {
		tokens = append(tokens, req.URL.Query().Get("CreationToken"))
		return http.StatusOK, "", `{"FileSystems":[]}`
	}))

	if _, err := rm.sdkFind(context.Background(), newFileSystem(nil)); err != ackerr.NotFound {
		t.Errorf("sdkFind error = %v, want NotFound", err)
	}
	if len(tokens) != 1 || tokens[0] != testUID {
		t.Errorf("creation tokens looked up = %q, want the CR UID %q once", tokens, testUID)
	}
}

func TestSdkFind_WithoutIDAndUIDSendsNoRequest(t *testing.T) {
	rm := newTestManager(t, efsAPI(func(req *http.Request) (int, string, string) {
		t.Errorf("unexpected request %s %s", req.Method, req.URL)
		return http.StatusBadRequest, "BadRequest", `{"ErrorCode":"BadRequest"}`
	}))
	r := newFileSystem(nil)
	r.ko.UID = ""
	if _, err := rm.sdkFind(context.Background(), r); err != ackerr.NotFound {
		t.Errorf("sdkFind error = %v, want NotFound", err)
	}
}
