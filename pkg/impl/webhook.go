// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	api "github.com/gitbundle/api"
	"github.com/gitbundle/api/internal/hmac"
)

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn api.SecretFunc) (api.Webhook, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}

	var hook api.Webhook
	switch req.Header.Get("X-Gitea-Event") {
	case "push":
		hook, err = s.parsePushHook(data)
	case "create":
		hook, err = s.parseCreateHook(data)
	case "delete":
		hook, err = s.parseDeleteHook(data)
	case "issues":
		hook, err = s.parseIssueHook(data)
	case "issue_comment":
		hook, err = s.parseIssueCommentHook(data)
	case "pull_request":
		hook, err = s.parsePullRequestHook(data)
	default:
		return nil, api.ErrUnknownEvent
	}
	if err != nil {
		return nil, err
	}

	// get the magit signature key to verify the payload
	// signature. If no key is provided, no validation
	// is performed.
	key, err := fn(hook)
	if err != nil {
		return hook, err
	} else if key == "" {
		return hook, nil
	}

	secret := req.FormValue("secret")
	signature := req.Header.Get("X-Gitea-Signature")

	// fail if no signature passed
	if signature == "" && secret == "" {
		return hook, api.ErrSignatureInvalid
	}

	// test signature if header not set and secret is in payload
	if signature == "" && secret != "" && secret != key {
		return hook, api.ErrSignatureInvalid
	}

	// test signature using header
	if signature != "" && !hmac.Validate(sha256.New, data, []byte(key), signature) {
		return hook, api.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePushHook(data []byte) (api.Webhook, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	return convertPushHook(dst), err
}

func (s *webhookService) parseCreateHook(data []byte) (api.Webhook, error) {
	dst := new(createHook)
	err := json.Unmarshal(data, dst)
	switch dst.RefType {
	case "tag":
		return convertTagHook(dst, api.ActionCreate), err
	case "branch":
		return convertBranchHook(dst, api.ActionCreate), err
	default:
		return nil, api.ErrUnknownEvent
	}
}

func (s *webhookService) parseDeleteHook(data []byte) (api.Webhook, error) {
	dst := new(createHook)
	err := json.Unmarshal(data, dst)
	switch dst.RefType {
	case "tag":
		return convertTagHook(dst, api.ActionDelete), err
	case "branch":
		return convertBranchHook(dst, api.ActionDelete), err
	default:
		return nil, api.ErrUnknownEvent
	}
}

func (s *webhookService) parseIssueHook(data []byte) (api.Webhook, error) {
	dst := new(issueHook)
	err := json.Unmarshal(data, dst)
	return convertIssueHook(dst), err
}

func (s *webhookService) parseIssueCommentHook(data []byte) (api.Webhook, error) {
	dst := new(issueHook)
	err := json.Unmarshal(data, dst)
	if dst.Issue.PullRequest != nil {
		return convertPullRequestCommentHook(dst), err
	}
	return convertIssueCommentHook(dst), err
}

func (s *webhookService) parsePullRequestHook(data []byte) (api.Webhook, error) {
	dst := new(pullRequestHook)
	err := json.Unmarshal(data, dst)
	return convertPullRequestHook(dst), err
}

//
// native data structures
//

type (
	// magit push webhook payload
	pushHook struct {
		Ref        string     `json:"ref"`
		Before     string     `json:"before"`
		After      string     `json:"after"`
		Compare    string     `json:"compare_url"`
		Commits    []commit   `json:"commits"`
		Repository repository `json:"repository"`
		Pusher     user       `json:"pusher"`
		Sender     user       `json:"sender"`
	}

	// magit create webhook payload
	createHook struct {
		Ref           string     `json:"ref"`
		RefType       string     `json:"ref_type"`
		Sha           string     `json:"sha"`
		DefaultBranch string     `json:"default_branch"`
		Repository    repository `json:"repository"`
		Sender        user       `json:"sender"`
	}

	// magit issue webhook payload
	issueHook struct {
		Action     string       `json:"action"`
		Issue      issue        `json:"issue"`
		Comment    issueComment `json:"comment"`
		Repository repository   `json:"repository"`
		Sender     user         `json:"sender"`
	}

	// magit pull request webhook payload
	pullRequestHook struct {
		Action      string     `json:"action"`
		Number      int        `json:"number"`
		PullRequest pr         `json:"pull_request"`
		Repository  repository `json:"repository"`
		Sender      user       `json:"sender"`
	}
)

//
// native data structure conversion
//

func convertTagHook(dst *createHook, action api.Action) *api.TagHook {
	return &api.TagHook{
		Action: action,
		Ref: api.Reference{
			Name: dst.Ref,
			Sha:  dst.Sha,
		},
		Repo:   *convertRepository(&dst.Repository),
		Sender: *convertUser(&dst.Sender),
	}
}

func convertBranchHook(dst *createHook, action api.Action) *api.BranchHook {
	return &api.BranchHook{
		Action: action,
		Ref: api.Reference{
			Name: dst.Ref,
		},
		Repo:   *convertRepository(&dst.Repository),
		Sender: *convertUser(&dst.Sender),
	}
}

func convertPushHook(dst *pushHook) *api.PushHook {
	if len(dst.Commits) > 0 {
		var commits []api.Commit
		for _, c := range dst.Commits {
			commits = append(commits,
				api.Commit{
					Sha:     c.ID,
					Message: c.Message,
					Link:    c.URL,
					Author: api.Signature{
						Login: c.Author.Username,
						Email: c.Author.Email,
						Name:  c.Author.Name,
						Date:  c.Timestamp,
					},
					Committer: api.Signature{
						Login: c.Committer.Username,
						Email: c.Committer.Email,
						Name:  c.Committer.Name,
						Date:  c.Timestamp,
					},
				})
		}

		return &api.PushHook{
			Ref:    dst.Ref,
			Before: dst.Before,
			Commit: api.Commit{
				Sha:     dst.After,
				Message: dst.Commits[0].Message,
				Link:    dst.Compare,
				Author: api.Signature{
					Login: dst.Commits[0].Author.Username,
					Email: dst.Commits[0].Author.Email,
					Name:  dst.Commits[0].Author.Name,
					Date:  dst.Commits[0].Timestamp,
				},
				Committer: api.Signature{
					Login: dst.Commits[0].Committer.Username,
					Email: dst.Commits[0].Committer.Email,
					Name:  dst.Commits[0].Committer.Name,
					Date:  dst.Commits[0].Timestamp,
				},
			},
			Commits: commits,
			Repo:    *convertRepository(&dst.Repository),
			Sender:  *convertUser(&dst.Sender),
		}
	} else {
		return &api.PushHook{
			Ref: dst.Ref,
			Commit: api.Commit{
				Sha:  dst.After,
				Link: dst.Compare,
				Author: api.Signature{
					Login: dst.Pusher.Login,
					Email: dst.Pusher.Email,
					Name:  dst.Pusher.Fullname,
				},
				Committer: api.Signature{
					Login: dst.Pusher.Login,
					Email: dst.Pusher.Email,
					Name:  dst.Pusher.Fullname,
				},
			},
			Repo:   *convertRepository(&dst.Repository),
			Sender: *convertUser(&dst.Sender),
		}
	}
}

func convertPullRequestHook(dst *pullRequestHook) *api.PullRequestHook {
	return &api.PullRequestHook{
		Action: convertAction(dst.Action),
		PullRequest: api.PullRequest{
			Number: dst.PullRequest.Number,
			Title:  dst.PullRequest.Title,
			Body:   dst.PullRequest.Body,
			Closed: dst.PullRequest.State == "closed",
			Author: api.User{
				Login:  dst.PullRequest.User.Login,
				Email:  dst.PullRequest.User.Email,
				Avatar: dst.PullRequest.User.Avatar,
			},
			Merged: dst.PullRequest.Merged,
			// Created: nil,
			// Updated: nil,
			Source: dst.PullRequest.Head.Name,
			Target: dst.PullRequest.Base.Name,
			Fork:   dst.PullRequest.Head.Repo.FullName,
			Link:   dst.PullRequest.HTMLURL,
			Ref:    fmt.Sprintf("refs/pull/%d/head", dst.PullRequest.Number),
			Sha:    dst.PullRequest.Head.Sha,
		},
		Repo:   *convertRepository(&dst.Repository),
		Sender: *convertUser(&dst.Sender),
	}
}

func convertPullRequestCommentHook(dst *issueHook) *api.PullRequestCommentHook {
	return &api.PullRequestCommentHook{
		Action:      convertAction(dst.Action),
		PullRequest: *convertPullRequestFromIssue(&dst.Issue),
		Comment:     *convertIssueComment(&dst.Comment),
		Repo:        *convertRepository(&dst.Repository),
		Sender:      *convertUser(&dst.Sender),
	}
}

func convertIssueHook(dst *issueHook) *api.IssueHook {
	return &api.IssueHook{
		Action: convertAction(dst.Action),
		Issue:  *convertIssue(&dst.Issue),
		Repo:   *convertRepository(&dst.Repository),
		Sender: *convertUser(&dst.Sender),
	}
}

func convertIssueCommentHook(dst *issueHook) *api.IssueCommentHook {
	return &api.IssueCommentHook{
		Action:  convertAction(dst.Action),
		Issue:   *convertIssue(&dst.Issue),
		Comment: *convertIssueComment(&dst.Comment),
		Repo:    *convertRepository(&dst.Repository),
		Sender:  *convertUser(&dst.Sender),
	}
}

func convertAction(src string) (action api.Action) {
	switch src {
	case "create", "created":
		return api.ActionCreate
	case "delete", "deleted":
		return api.ActionDelete
	case "update", "updated", "edit", "edited":
		return api.ActionUpdate
	case "open", "opened":
		return api.ActionOpen
	case "reopen", "reopened":
		return api.ActionReopen
	case "close", "closed":
		return api.ActionClose
	case "label", "labeled":
		return api.ActionLabel
	case "unlabel", "unlabeled":
		return api.ActionUnlabel
	case "merge", "merged":
		return api.ActionMerge
	case "synchronize", "synchronized":
		return api.ActionSync
	default:
		return
	}
}
