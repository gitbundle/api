// Copyright 2023 GitBundle Inc. All rights reserved.
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

func TestConvertAPIURLToHTMLURL(t *testing.T) {

	got := convertAPIURLToHTMLURL("https://try.magit.com/api/v1/repos/octocat/Hello-World/123", "v1.0.0")
	want := "https://try.magit.com/octocat/Hello-World/releases/tag/v1.0.0"

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		t.Log("got:")
		t.Log(string(got))
	}

}

func TestConvertAPIURLToHTMLURLEmptyLinkWhenURLParseFails(t *testing.T) {

	broken := []string{"http s://try.magit.com/api/v1/repos/octocat/Hello-World/123", "https://try.magit.com/api/v1/repos/octocat/Hello-World"}
	for _, url := range broken {

		got := convertAPIURLToHTMLURL(url, "v1.0.0")
		want := ""

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Unexpected Results")
			t.Log(diff)

			t.Log("got:")
			t.Log(string(got))
		}
	}

}

func TestReleaseFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Get("/repos/octocat/hello-world/releases/1").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	client, err := New("https://example.gitbundle.com")
	if err != nil {
		t.Error(err)
		return
	}
	got, _, err := client.Releases.Find(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(api.Release)
	raw, _ := ioutil.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}
}

func TestReleaseFindByTag(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Get("/repos/octocat/hello-world/releases/tags/v1.0.0").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	client, err := New("https://example.gitbundle.com")
	if err != nil {
		t.Error(err)
		return
	}
	got, _, err := client.Releases.FindByTag(context.Background(), "octocat/hello-world", "v1.0.0")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(api.Release)
	raw, _ := ioutil.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}
}

func TestReleaseList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Get("/repos/octocat/hello-world/releases").
		MatchParam("page", "1").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/releases.json")

	client, err := New("https://example.gitbundle.com")
	if err != nil {
		t.Error(err)
		return
	}

	got, _, err := client.Releases.List(context.Background(), "octocat/hello-world", api.ReleaseListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*api.Release{}
	raw, _ := ioutil.ReadFile("testdata/releases.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}

}

func TestReleaseCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Post("/repos/octocat/hello-world/releases").
		File("testdata/release_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	client, err := New("https://example.gitbundle.com")
	if err != nil {
		t.Error(err)
		return
	}
	input := &api.ReleaseInput{
		Title:       "v1.0.0",
		Description: "Description of the release",
		Tag:         "v1.0.0",
		Commitish:   "master",
		Draft:       false,
		Prerelease:  false,
	}

	got, _, err := client.Releases.Create(context.Background(), "octocat/hello-world", input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(api.Release)
	raw, _ := ioutil.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}

}

func TestReleaseUpdate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Patch("/repos/octocat/hello-world/releases/1").
		File("testdata/release_update.json").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	client, err := New("https://example.gitbundle.com")
	if err != nil {
		t.Error(err)
		return
	}
	input := &api.ReleaseInput{
		Title:       "v1.0.0",
		Description: "Description of the release",
		Tag:         "v1.0.0",
		Commitish:   "master",
		Draft:       false,
		Prerelease:  false,
	}
	got, _, err := client.Releases.Update(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(api.Release)
	raw, _ := ioutil.ReadFile("testdata/release.json.golden")
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

func TestReleaseUpdateByTag(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Get("/repos/octocat/hello-world/releases/tags/v1.0.0").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	gock.New("https://example.gitbundle.com").
		Patch("/repos/octocat/hello-world/releases/1").
		File("testdata/release_update.json").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	client, err := New("https://example.gitbundle.com")
	if err != nil {
		t.Error(err)
		return
	}
	input := &api.ReleaseInput{
		Title:       "v1.0.0",
		Description: "Description of the release",
		Tag:         "v1.0.0",
		Commitish:   "master",
		Draft:       false,
		Prerelease:  false,
	}
	got, _, err := client.Releases.UpdateByTag(context.Background(), "octocat/hello-world", "v1.0.0", input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(api.Release)
	raw, _ := ioutil.ReadFile("testdata/release.json.golden")
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

func TestReleaseDelete(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Delete("/repos/octocat/hello-world/releases/1").
		Reply(200).
		Type("application/json")

	client, err := New("https://example.gitbundle.com")
	_, err = client.Releases.Delete(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestReleaseDeleteByTag(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://example.gitbundle.com").
		Delete("/repos/octocat/hello-world/releases/tags/v1.0.0").
		Reply(200).
		Type("application/json")

	client, err := New("https://example.gitbundle.com")
	_, err = client.Releases.DeleteByTag(context.Background(), "octocat/hello-world", "v1.0.0")
	if err != nil {
		t.Error(err)
		return
	}

}
