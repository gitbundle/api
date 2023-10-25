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

func TestOrgFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/orgs/gogits").
		Reply(200).
		Type("application/json").
		File("testdata/organization.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Organizations.Find(context.Background(), "gogits")
	if err != nil {
		t.Error(err)
	}

	want := new(api.Organization)
	raw, _ := ioutil.ReadFile("testdata/organization.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestOrganizationFindMembership(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.Organizations.FindMembership(context.Background(), "gogits", "jcitizen")
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestOrgList(t *testing.T) {
	defer gock.Off()

	gock.New("https://example.gitbundle.com").
		Get("/api/v1/user/orgs").
		Reply(200).
		Type("application/json").
		File("testdata/organizations.json")

	client, _ := New("https://example.gitbundle.com")
	got, _, err := client.Organizations.List(context.Background(), api.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*api.Organization{}
	raw, _ := ioutil.ReadFile("testdata/organizations.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
