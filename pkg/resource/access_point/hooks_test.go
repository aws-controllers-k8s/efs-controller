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
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
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

func TestSdkCreate_WithoutUIDLeavesTheTokenToTheSDK(t *testing.T) {
	api := &createAccessPointAPI{
		status: http.StatusCreated,
		body:   `{"AccessPointId":"fsap-0123456789abcdef0","FileSystemId":"fs-0123456789abcdef0","LifeCycleState":"creating"}`,
	}
	rm := newTestManager(t, api)

	r := newAccessPoint(nil)
	r.ko.UID = ""
	if _, err := rm.sdkCreate(context.Background(), r); err != nil {
		t.Fatalf("sdkCreate: %v", err)
	}
	if len(api.tokens) != 1 || api.tokens[0] == "" {
		t.Errorf("client tokens = %q, want one token filled in by the SDK", api.tokens)
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

func accessPointJSON(id, token, state string) string {
	return `{"AccessPointId":"` + id + `","FileSystemId":"fs-0123456789abcdef0",` +
		`"ClientToken":"` + token + `","LifeCycleState":"` + state + `"}`
}

// The response of a Create that went through was lost, the retry gets
// AccessPointAlreadyExists, and the access point is still being created: the
// read right after Create must ask for a requeue that waits, not for none.
func TestReadOne_AfterAccessPointAlreadyExistsWhileCreating_Requeues(t *testing.T) {
	const id = "fsap-0123456789abcdef0"
	rm := newTestManager(t, efsAPI(func(req *http.Request) (int, string, string) {
		switch {
		case req.Method == http.MethodPost:
			return http.StatusConflict, "AccessPointAlreadyExists",
				`{"ErrorCode":"AccessPointAlreadyExists","Message":"exists","AccessPointId":"` + id + `"}`
		case req.URL.Query().Get("AccessPointId") == id:
			return http.StatusOK, "", `{"AccessPoints":[` + accessPointJSON(id, testUID, "creating") + `]}`
		}
		t.Errorf("unexpected request %s %s", req.Method, req.URL)
		return http.StatusBadRequest, "BadRequest", `{"ErrorCode":"BadRequest"}`
	}))
	ctx := context.Background()

	created, err := rm.Create(ctx, newAccessPoint(nil))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	_, err = rm.ReadOne(ctx, created)

	var requeue *ackrequeue.RequeueNeededAfter
	if !errors.As(err, &requeue) || requeue == nil || requeue.Duration() <= 0 {
		t.Fatalf("ReadOne error = %#v, want a requeue after a non-zero duration", err)
	}
	if !strings.Contains(requeue.Error(), "'creating'") {
		t.Errorf("requeue message = %q, want the observed state 'creating'", requeue.Error())
	}
}

func TestRequeueWaitState_WithoutStateStillRequeues(t *testing.T) {
	requeue := requeueWaitState(newAccessPoint(nil))
	if requeue == nil || requeue.Duration() <= 0 {
		t.Fatalf("requeueWaitState = %#v, want a requeue after a non-zero duration", requeue)
	}
}

// The status has no ID, for example because the controller lost the response
// of a Create that went through. The read finds the access point by the CR UID
// as its client token, so a delete of the CR deletes it too.
func TestSdkFind_WithoutIDFindsTheAccessPointByClientToken(t *testing.T) {
	const id = "fsap-0123456789abcdef0"
	rm := newTestManager(t, efsAPI(func(req *http.Request) (int, string, string) {
		q := req.URL.Query()
		switch {
		case q.Get("FileSystemId") == "fs-0123456789abcdef0":
			return http.StatusOK, "", `{"AccessPoints":[` +
				accessPointJSON("fsap-0aaaaaaaaaaaaaaaa", "another-token", "available") + `,` +
				accessPointJSON(id, testUID, "available") + `]}`
		case q.Get("AccessPointId") == id:
			return http.StatusOK, "", `{"AccessPoints":[` + accessPointJSON(id, testUID, "available") + `]}`
		}
		t.Errorf("unexpected request %s %s", req.Method, req.URL)
		return http.StatusBadRequest, "BadRequest", `{"ErrorCode":"BadRequest"}`
	}))

	found, err := rm.sdkFind(context.Background(), newAccessPoint(nil))
	if err != nil {
		t.Fatalf("sdkFind: %v", err)
	}
	if got := found.ko.Status.AccessPointID; got == nil || *got != id {
		t.Errorf("AccessPointID = %v, want %s", got, id)
	}
}

func TestSdkFind_WithoutIDAndNoAccessPointWithTheToken_NotFound(t *testing.T) {
	for name, answer := range map[string]struct {
		status        int
		errType, body string
	}{
		"other tokens only": {http.StatusOK, "", `{"AccessPoints":[` +
			accessPointJSON("fsap-0aaaaaaaaaaaaaaaa", "another-token", "available") + `]}`},
		"file system gone": {http.StatusNotFound, "FileSystemNotFound",
			`{"ErrorCode":"FileSystemNotFound","Message":"gone"}`},
	} {
		t.Run(name, func(t *testing.T) {
			rm := newTestManager(t, efsAPI(func(req *http.Request) (int, string, string) {
				return answer.status, answer.errType, answer.body
			}))
			_, err := rm.sdkFind(context.Background(), newAccessPoint(nil))
			if err != ackerr.NotFound {
				t.Errorf("sdkFind error = %v, want NotFound", err)
			}
		})
	}
}

func TestSdkFind_WithoutIDAndUIDSendsNoRequest(t *testing.T) {
	rm := newTestManager(t, efsAPI(func(req *http.Request) (int, string, string) {
		t.Errorf("unexpected request %s %s", req.Method, req.URL)
		return http.StatusBadRequest, "BadRequest", `{"ErrorCode":"BadRequest"}`
	}))
	r := newAccessPoint(nil)
	r.ko.UID = ""
	if _, err := rm.sdkFind(context.Background(), r); err != ackerr.NotFound {
		t.Errorf("sdkFind error = %v, want NotFound", err)
	}
}
