// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	api "github.com/gitbundle/api"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

//
// issue sub-tests
//

func TestIssueFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/go-magit/magit/issues/1").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Issues.Find(context.Background(), "go-magit/magit", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(api.Issue)
	raw, _ := ioutil.ReadFile("testdata/issue.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestIssueList(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/go-magit/magit/issues").
		Reply(200).
		Type("application/json").
		File("testdata/issues.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Issues.List(context.Background(), "go-magit/magit", api.IssueListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*api.Issue{}
	raw, _ := ioutil.ReadFile("testdata/issues.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestIssueCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Post("/api/v1/repos/go-magit/magit/issues").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	input := api.IssueInput{
		Title: "Bug found",
		Body:  "I'm having a problem with this.",
	}

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Issues.Create(context.Background(), "go-magit/magit", &input)
	if err != nil {
		t.Error(err)
	}

	want := new(api.Issue)
	raw, _ := ioutil.ReadFile("testdata/issue.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestIssueClose(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.Issues.Close(context.Background(), "gogits/go-gogs-client", 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestIssueLock(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.Issues.Lock(context.Background(), "gogits/go-gogs-client", 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestIssueUnlock(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.Issues.Unlock(context.Background(), "gogits/go-gogs-client", 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

//
// issue comment sub-tests
//

func TestIssueCommentFind(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.Issues.FindComment(context.Background(), "gogits/go-gogs-client", 1, 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestIssueCommentList(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/go-magit/magit/issues/1/comments").
		Reply(200).
		Type("application/json").
		File("testdata/comments.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Issues.ListComments(context.Background(), "go-magit/magit", 1, api.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*api.Comment{}
	raw, _ := ioutil.ReadFile("testdata/comments.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestIssueCommentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Post("/api/v1/repos/go-magit/magit/issues/1/comments").
		Reply(201).
		Type("application/json").
		File("testdata/comment.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Issues.CreateComment(context.Background(), "go-magit/magit", 1, &api.CommentInput{Body: "what?"})
	if err != nil {
		t.Error(err)
	}

	want := new(api.Comment)
	raw, _ := ioutil.ReadFile("testdata/comment.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	if gock.IsPending() {
		t.Errorf("Pending API calls")
	}
}

func TestIssueCommentDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Delete("/api/v1/repos/go-magit/magit/issues/1/comments/1").
		Reply(204).
		Type("application/json")

	client, _ := New("https://example.gitbundle.com")
	_, err := client.Issues.DeleteComment(context.Background(), "go-magit/magit", 1, 1)
	if err != nil {
		t.Error(err)
	}
}
