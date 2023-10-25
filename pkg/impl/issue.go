// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"
	"fmt"
	"time"

	api "github.com/gitbundle/api"
)

type issueService struct {
	client *wrapper
}

func (s *issueService) Find(ctx context.Context, repo string, number int) (*api.Issue, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues/%d", repo, number)
	out := new(issue)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertIssue(out), res, err
}

func (s *issueService) FindComment(ctx context.Context, repo string, index, id int) (*api.Comment, *api.Response, error) {
	return nil, nil, api.ErrNotSupported
}

func (s *issueService) List(ctx context.Context, repo string, opts api.IssueListOptions) ([]*api.Issue, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues?%s", repo, encodeIssueListOptions(opts))
	out := []*issue{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertIssueList(out), res, err
}

func (s *issueService) ListComments(ctx context.Context, repo string, index int, opts api.ListOptions) ([]*api.Comment, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues/%d/comments?%s", repo, index, encodeListOptions(opts))
	out := []*issueComment{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertIssueCommentList(out), res, err
}

func (s *issueService) Create(ctx context.Context, repo string, input *api.IssueInput) (*api.Issue, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues", repo)
	in := &issueInput{
		Title: input.Title,
		Body:  input.Body,
	}
	out := new(issue)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertIssue(out), res, err
}

func (s *issueService) CreateComment(ctx context.Context, repo string, index int, input *api.CommentInput) (*api.Comment, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues/%d/comments", repo, index)
	in := &issueCommentInput{
		Body: input.Body,
	}
	out := new(issueComment)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertIssueComment(out), res, err
}

func (s *issueService) DeleteComment(ctx context.Context, repo string, index, id int) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues/%d/comments/%d", repo, index, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *issueService) Close(ctx context.Context, repo string, number int) (*api.Response, error) {
	return nil, api.ErrNotSupported
}

func (s *issueService) Lock(ctx context.Context, repo string, number int) (*api.Response, error) {
	return nil, api.ErrNotSupported
}

func (s *issueService) Unlock(ctx context.Context, repo string, number int) (*api.Response, error) {
	return nil, api.ErrNotSupported
}

//
// native data structures
//

type (
	// magit issue response object.
	issue struct {
		ID          int       `json:"id"`
		Number      int       `json:"number"`
		User        user      `json:"user"`
		Title       string    `json:"title"`
		Body        string    `json:"body"`
		State       string    `json:"state"`
		Labels      []string  `json:"labels"`
		Comments    int       `json:"comments"`
		Created     time.Time `json:"created_at"`
		Updated     time.Time `json:"updated_at"`
		PullRequest *struct {
			Merged   bool        `json:"merged"`
			MergedAt interface{} `json:"merged_at"`
		} `json:"pull_request"`
	}

	// magit issue request object.
	issueInput struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	// magit issue comment response object.
	issueComment struct {
		ID        int       `json:"id"`
		HTMLURL   string    `json:"html_url"`
		User      user      `json:"user"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// magit issue comment request object.
	issueCommentInput struct {
		Body string `json:"body"`
	}
)

//
// native data structure conversion
//

func convertIssueList(from []*issue) []*api.Issue {
	to := []*api.Issue{}
	for _, v := range from {
		to = append(to, convertIssue(v))
	}
	return to
}

func convertIssue(from *issue) *api.Issue {
	return &api.Issue{
		Number:  from.Number,
		Title:   from.Title,
		Body:    from.Body,
		Link:    "", // TODO construct the link to the issue.
		Closed:  from.State == "closed",
		Author:  *convertUser(&from.User),
		Created: from.Created,
		Updated: from.Updated,
	}
}

func convertIssueCommentList(from []*issueComment) []*api.Comment {
	to := []*api.Comment{}
	for _, v := range from {
		to = append(to, convertIssueComment(v))
	}
	return to
}

func convertIssueComment(from *issueComment) *api.Comment {
	return &api.Comment{
		ID:      from.ID,
		Body:    from.Body,
		Author:  *convertUser(&from.User),
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
	}
}
