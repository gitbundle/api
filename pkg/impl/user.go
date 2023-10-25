// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"
	"fmt"
	"strconv"

	api "github.com/gitbundle/api"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*api.User, *api.Response, error) {
	out := new(user)
	res, err := s.client.do(ctx, "GET", "api/v1/user", nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*api.User, *api.Response, error) {
	path := fmt.Sprintf("api/v1/users/%s", login)
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *api.Response, error) {
	user, res, err := s.Find(ctx)
	return user.Email, res, err
}

func (s *userService) ListEmail(context.Context, api.ListOptions) ([]*api.Email, *api.Response, error) {
	return nil, nil, api.ErrNotSupported
}

//
// native data structures
//

type user struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
	IsAdmin  bool   `json:"is_admin"`
}

//
// native data structure conversion
//

func convertUser(src *user) *api.User {
	return &api.User{
		ID:      strconv.Itoa(src.ID),
		Login:   userLogin(src),
		Avatar:  src.Avatar,
		Email:   src.Email,
		Name:    src.Fullname,
		IsAdmin: src.IsAdmin,
	}
}

func userLogin(src *user) string {
	if src.Username != "" {
		return src.Username
	}
	return src.Login
}
