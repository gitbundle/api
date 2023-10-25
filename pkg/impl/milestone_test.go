// Copyright 2023 GitBundle Inc. All rights reserved.
package impl

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	api "github.com/gitbundle/api"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestMilestoneFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/jcitizen/my-repo/milestones/1").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Milestones.Find(context.Background(), "jcitizen/my-repo", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(api.Milestone)
	raw, _ := ioutil.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestMilestoneList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/jcitizen/my-repo/milestones").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/milestones.json")

	client, _ := New("https://example.gitbundle.com")
	got, res, err := client.Milestones.List(context.Background(), "jcitizen/my-repo", api.MilestoneListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*api.Milestone{}
	raw, _ := ioutil.ReadFile("testdata/milestones.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestMilestoneCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Post("/api/v1/repos/jcitizen/my-repo/milestones").
		File("testdata/milestone_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://example.gitbundle.com")
	dueDate, _ := time.Parse(api.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &api.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "open",
		DueDate:     dueDate,
	}
	got, _, err := client.Milestones.Create(context.Background(), "jcitizen/my-repo", input)
	if err != nil {
		t.Error(err)
	}

	want := new(api.Milestone)
	raw, _ := ioutil.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestMilestoneUpdate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Patch("/api/v1/repos/jcitizen/my-repo/milestones").
		File("testdata/milestone_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://example.gitbundle.com")
	dueDate, _ := time.Parse(api.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &api.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "open",
		DueDate:     dueDate,
	}
	got, _, err := client.Milestones.Update(context.Background(), "jcitizen/my-repo", 1, input)
	if err != nil {
		t.Error(err)
	}

	want := new(api.Milestone)
	raw, _ := ioutil.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestMilestoneDelete(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Delete("/api/v1/repos/jcitizen/my-repo/milestones/1").
		Reply(200).
		Type("application/json")

	client, _ := New("https://example.gitbundle.com")
	_, err := client.Milestones.Delete(context.Background(), "jcitizen/my-repo", 1)
	if err != nil {
		t.Error(err)
	}
}

var mockPageHeaders = map[string]string{
	"Link": `<https://example.gitbundle.com/v1/resource?page=2>; rel="next",` +
		`<https://example.gitbundle.com/v1/resource?page=1>; rel="prev",` +
		`<https://example.gitbundle.com/v1/resource?page=1>; rel="first",` +
		`<https://example.gitbundle.com/v1/resource?page=5>; rel="last"`,
}

func mockServerVersion() {
	gock.New("https://example.gitbundle.com").
		Get("/api/v1/version").
		Reply(200).
		Type("application/json").
		File("testdata/version.json")
}
