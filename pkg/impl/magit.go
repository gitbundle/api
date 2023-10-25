// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package impl implements a Magit client.
package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"

	api "github.com/gitbundle/api"
)

// New returns a new Magit API client.
func New(uri string) (*api.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	client := &wrapper{new(api.Client)}
	client.BaseURL = base
	// initialize services
	client.Driver = api.DriverMagit
	client.Linker = &linker{base.String()}
	client.Contents = &contentService{client}
	client.Git = &gitService{client}
	client.Issues = &issueService{client}
	client.Milestones = &milestoneService{client}
	client.Organizations = &organizationService{client}
	client.PullRequests = &pullService{client}
	client.Repositories = &repositoryService{client}
	client.Releases = &releaseService{client}
	client.Reviews = &reviewService{client}
	client.Users = &userService{client}
	client.Webhooks = &webhookService{client}
	return client.Client, nil
}

func NewDefault() *api.Client {
	client, _ := New("http://localhost:4000")
	return client
}

// wrapper wraps the Client to provide high level helper functions
// for making http requests and unmarshaling the response.
type wrapper struct {
	*api.Client
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*api.Response, error) {
	req := &api.Request{
		Method: method,
		Path:   path,
		Header: map[string][]string{},
	}
	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(in)
		req.Header.Add("Content-Type", "application/json")
		req.Body = buf
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return res, res.Decode(ctx, out)
}
