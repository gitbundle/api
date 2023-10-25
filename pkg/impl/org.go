// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/gitbundle/api"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) Find(ctx context.Context, name string) (*api.Organization, *api.Response, error) {
	path := fmt.Sprintf("api/v1/orgs/%s", name)
	out := new(org)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertOrg(out), res, err
}

func (s *organizationService) FindMembership(ctx context.Context, name, username string) (*api.Membership, *api.Response, error) {
	return nil, nil, api.ErrNotSupported
}

func (s *organizationService) List(ctx context.Context, opts api.ListOptions) ([]*api.Organization, *api.Response, error) {
	path := fmt.Sprintf("api/v1/user/orgs?%s", encodeListOptions(opts))
	out := []*org{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertOrgList(out), res, err
}

func (s *organizationService) CheckMember(ctx context.Context, org, username string) (bool, *api.Response, error) {
	path := fmt.Sprintf("api/v1/orgs/%s/members/%s", org, username)
	res, err := s.client.do(ctx, "GET", path, nil, nil)
	return res.Status == http.StatusNoContent, res, err
}

func (s *organizationService) FindTeamMember(ctx context.Context, teamId int64, username string) (*api.Member, *api.Response, error) {
	path := fmt.Sprintf("api/v1/teams/%d/members/%s", teamId, username)
	out := member{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertMember(out), res, err
}

//
// native data structures
//

type (
	org struct {
		Name   string `json:"username"`
		Avatar string `json:"avatar_url"`
	}

	member struct {
		ID         int64  `json:"id"`
		Active     bool   `json:"active"`
		Email      string `json:"email"`
		FullName   string `json:"full_name"`
		IsAdmin    bool   `json:"is_admin"`
		Login      string `json:"login"`
		Restricted bool   `json:"restricted"`
		Username   string `json:"username"`
		Visibility string `json:"visibility"`
	}
)

//
// native data structure conversion
//

func convertMember(from member) *api.Member {
	return &api.Member{
		ID:         from.ID,
		Active:     from.Active,
		Email:      from.Email,
		FullName:   from.FullName,
		IsAdmin:    from.IsAdmin,
		Login:      from.Login,
		Restricted: from.Restricted,
		Username:   from.Username,
		Visibility: from.Visibility,
	}
}

func convertOrgList(from []*org) []*api.Organization {
	to := []*api.Organization{}
	for _, v := range from {
		to = append(to, convertOrg(v))
	}
	return to
}

func convertOrg(from *org) *api.Organization {
	return &api.Organization{
		Name:   from.Name,
		Avatar: from.Avatar,
	}
}
