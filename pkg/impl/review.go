// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"

	api "github.com/gitbundle/api"
)

type reviewService struct {
	client *wrapper
}

func (s *reviewService) Find(ctx context.Context, repo string, number, id int) (*api.Review, *api.Response, error) {
	return nil, nil, api.ErrNotSupported
}

func (s *reviewService) List(ctx context.Context, repo string, number int, opts api.ListOptions) ([]*api.Review, *api.Response, error) {
	return nil, nil, api.ErrNotSupported
}

func (s *reviewService) Create(ctx context.Context, repo string, number int, input *api.ReviewInput) (*api.Review, *api.Response, error) {
	return nil, nil, api.ErrNotSupported
}

func (s *reviewService) Delete(ctx context.Context, repo string, number, id int) (*api.Response, error) {
	return nil, api.ErrNotSupported
}
