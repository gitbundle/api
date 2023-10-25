// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"
	"testing"

	api "github.com/gitbundle/api"
)

func TestReviewFind(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.Reviews.Find(context.Background(), "go-magit/magit", 1, 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewList(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.Reviews.List(context.Background(), "go-magit/magit", 1, api.ListOptions{})
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewCreate(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, _, err := client.Reviews.Create(context.Background(), "go-magit/magit", 1, nil)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewDelete(t *testing.T) {
	client, _ := New("https://example.gitbundle.com")
	_, err := client.Reviews.Delete(context.Background(), "go-magit/magit", 1, 1)
	if err != api.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
