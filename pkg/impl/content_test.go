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

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/go-magit/magit/raw/f05f642b892d59a0a9ef6a31f6c905a24b5db13a/README.md").
		Reply(200).
		Type("plain/text").
		BodyString("Hello World\n")

	client, _ := New("https://example.gitbundle.com")
	result, _, err := client.Contents.Find(
		context.Background(),
		"go-magit/magit",
		"README.md",
		"f05f642b892d59a0a9ef6a31f6c905a24b5db13a",
	)
	if err != nil {
		t.Error(err)
	}

	if got, want := result.Path, "README.md"; got != want {
		t.Errorf("Want file Path %q, got %q", want, got)
	}
	if got, want := string(result.Data), "Hello World\n"; got != want {
		t.Errorf("Want file Data %q, got %q", want, got)
	}
}

func TestContentCreate(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.Contents.Create(context.Background(), "go-magit/magit", "README.md", nil)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentUpdate(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.Contents.Update(context.Background(), "go-magit/magit", "README.md", nil)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentDelete(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.Contents.Delete(context.Background(), "go-magit/magit", "README.md", &api.ContentParams{})
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/repos/go-magit/magit/contents/docs/content/doc").
		MatchParam("ref", "master").
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Contents.List(
		context.Background(),
		"go-magit/magit",
		"docs/content/doc",
		"master",
		api.ListOptions{},
	)
	if err != nil {
		t.Error(err)
	}

	want := []*api.ContentInfo{}
	raw, _ := ioutil.ReadFile("testdata/content_list.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
