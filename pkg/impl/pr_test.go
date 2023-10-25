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
// pull request sub-tests
//

func TestPullRequestFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/jcitizen/my-repo/pulls/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.PullRequests.Find(context.Background(), "jcitizen/my-repo", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(api.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestList(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/jcitizen/my-repo/pulls").
		Reply(200).
		Type("application/json").
		File("testdata/prs.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.PullRequests.List(context.Background(), "jcitizen/my-repo", api.PullRequestListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*api.PullRequest{}
	raw, _ := ioutil.ReadFile("testdata/prs.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Post("/api/v1/repos/jcitizen/my-repo/pulls").
		Reply(201).
		Type("application/json").
		File("testdata/pr.json")

	input := api.PullRequestInput{
		Title:  "Add License File",
		Body:   "Using a BSD License",
		Source: "feature",
		Target: "master",
	}

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.PullRequests.Create(context.Background(), "jcitizen/my-repo", &input)
	if err != nil {
		t.Error(err)
	}

	want := new(api.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestClose(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.PullRequests.Close(context.Background(), "go-magit/magit", 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestPullRequestMerge(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Post("/api/v1/repos/go-magit/magit/pulls/1").
		Reply(204).
		Type("application/json")

	client, _ := New("https://example.gitbundle.com")
	_, err := client.PullRequests.Merge(context.Background(), "go-magit/magit", 1)
	if err != nil {
		t.Error(err)
	}
}

//
// pull request change sub-tests
//

func TestPullRequestChanges(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.PullRequests.ListChanges(context.Background(), "go-magit/magit", 1, api.ListOptions{})
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

//
// pull request comment sub-tests
//

func TestPullRequestCommentFind(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.PullRequests.FindComment(context.Background(), "go-magit/magit", 1, 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestPullRequestCommentList(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.PullRequests.ListComments(context.Background(), "go-magit/magit", 1, api.ListOptions{})
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestPullRequestCommentCreate(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.PullRequests.CreateComment(context.Background(), "go-magit/magit", 1, &api.CommentInput{})
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestPullRequestCommentDelete(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.PullRequests.DeleteComment(context.Background(), "go-magit/magit", 1, 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestPullListCommits(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.PullRequests.ListCommits(context.Background(), "go-magit/magit", 1, api.ListOptions{})
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
