// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"bytes"
	"context"
	"fmt"

	api "github.com/gitbundle/api"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*api.Content, *api.Response, error) {
	endpoint := fmt.Sprintf("api/v1/repos/%s/raw/%s/%s", repo, api.TrimRef(ref), path)
	out := new(bytes.Buffer)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return &api.Content{
		Path: path,
		Data: out.Bytes(),
	}, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *api.ContentParams) (*api.Response, error) {
	return nil, api.ErrNotSupported
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *api.ContentParams) (*api.Response, error) {
	return nil, api.ErrNotSupported
}

func (s *contentService) Delete(ctx context.Context, repo, path string, params *api.ContentParams) (*api.Response, error) {
	return nil, api.ErrNotSupported
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, _ api.ListOptions) ([]*api.ContentInfo, *api.Response, error) {
	endpoint := fmt.Sprintf("api/v1/repos/%s/contents/%s?ref=%s", repo, path, ref)
	out := []*content{}
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertContentInfoList(out), res, err
}

type content struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
}

func convertContentInfoList(from []*content) []*api.ContentInfo {
	to := []*api.ContentInfo{}
	for _, v := range from {
		to = append(to, convertContentInfo(v))
	}
	return to
}

func convertContentInfo(from *content) *api.ContentInfo {
	to := &api.ContentInfo{Path: from.Path, BlobID: from.Sha}
	switch from.Type {
	case "file":
		to.Kind = api.ContentKindFile
	case "dir":
		to.Kind = api.ContentKindDirectory
	case "symlink":
		to.Kind = api.ContentKindSymlink
	case "submodule":
		to.Kind = api.ContentKindGitlink
	default:
		to.Kind = api.ContentKindUnsupported
	}
	return to
}
