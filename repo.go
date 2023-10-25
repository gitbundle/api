// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"context"
	"time"

	"github.com/gitbundle/api/pkg/structs"
)

type (
	// Repository represents a git repository.
	Repository struct {
		ID         string
		Namespace  string
		Name       string
		Perm       *Perm
		Branch     string
		Archived   bool
		Private    bool
		Visibility Visibility
		Clone      string
		CloneSSH   string
		Link       string
		Created    time.Time
		Updated    time.Time
	}

	// Perm represents a user's repository permissions.
	Perm struct {
		Pull  bool
		Push  bool
		Admin bool
	}

	// Hook represents a repository hook.
	Hook struct {
		ID         string
		Name       string
		Target     string
		Events     []string
		Active     bool
		SkipVerify bool
	}

	// HookInput provides the input fields required for
	// creating or updating repository webhooks.
	HookInput struct {
		Name       string
		Target     string
		Secret     string
		Events     HookEvents
		SkipVerify bool

		// NativeEvents are used to create hooks with
		// provider-specific event types that cannot be
		// abstracted or represented in HookEvents.
		NativeEvents []string
	}

	// HookEvents represents supported hook events.
	HookEvents struct {
		Branch             bool
		Deployment         bool
		Issue              bool
		IssueComment       bool
		PullRequest        bool
		PullRequestComment bool
		Push               bool
		ReviewComment      bool
		Tag                bool
	}

	// Status represents a commit status.
	Status struct {
		State  State
		Label  string
		Desc   string
		Target string

		// TODO(bradrydzewski) this field is only used
		// by Bitbucket which requires a user-defined
		// key (label), title and description. We need
		// to cleanup this abstraction.
		Title string
	}

	// StatusInput provides the input fields required for
	// creating or updating commit statuses.
	StatusInput struct {
		State  State
		Label  string
		Title  string
		Desc   string
		Target string
	}

	// DeployStatus represents a deployment status.
	DeployStatus struct {
		Number         int64
		State          State
		Desc           string
		Target         string
		Environment    string
		EnvironmentURL string
	}

	Team struct {
		ID                      int64
		CanCreateOrgRepo        bool
		Description             string
		IncludesAllRepositories bool
		Name                    string
		Permission              string
		Units                   []string
		UnitsMap                UnitsMap
	}
	UnitsMap struct {
		Code      string
		ExtIssues string
		ExtWiki   string
		Issues    string
		Packages  string
		Projects  string
		Pulls     string
		Releases  string
		Wiki      string
	}

	QueryOption struct {
		Pod       string
		Container string
		Debugging string
		Cluster   string
		Extend    string
	}

	// RepositoryService provides access to repository resources.
	RepositoryService interface {
		// Find returns a repository by name.
		Find(context.Context, string) (*Repository, *Response, error)

		// FindHook returns a repository hook.
		FindHook(context.Context, string, string) (*Hook, *Response, error)

		// FindPerms returns repository permissions.
		FindPerms(context.Context, string) (*Perm, *Response, error)

		// List returns a list of repositories.
		List(context.Context, ListOptions) ([]*Repository, *Response, error)

		// ListHooks returns a list or repository hooks.
		ListHooks(context.Context, string, ListOptions) ([]*Hook, *Response, error)

		// ListStatus returns a list of commit statuses.
		ListStatus(context.Context, string, string, ListOptions) ([]*Status, *Response, error)

		// CreateHook creates a new repository hook.
		CreateHook(context.Context, string, *HookInput) (*Hook, *Response, error)

		// CreateStatus creates a new commit status.
		CreateStatus(context.Context, string, string, *StatusInput) (*Status, *Response, error)

		// UpdateHook updates an existing repository hook.
		UpdateHook(context.Context, string, string, *HookInput) (*Hook, *Response, error)

		// DeleteHook deletes a repository hook.
		DeleteHook(context.Context, string, string) (*Response, error)

		CheckCollaborator(context.Context, string, string) (bool, *Response, error)
		CheckCollaboratorPermission(ctx context.Context, repo string, collaborator string) (*structs.RepoCollaboratorPermission, *Response, error)

		ListTeams(context.Context, string) ([]*Team, *Response, error)

		// deploy + build
		FindRequires(ctx context.Context, repo string) (*structs.Requirement, *Response, error)
		ListClusters(ctx context.Context, repo string, opt QueryOption) ([]string, *Response, error)
		ListSimplePullRequests(ctx context.Context, repo string, opt ListOptions) ([]*structs.SimplePullRequest, *Response, error)
		//FindDeployMetainfo(ctx context.Context, repo string) ([]*structs.DeployMetainfo, *Response, error)
		CreateDeploy(ctx context.Context, repo string, in *structs.CreateDeployOption) error
		CancelDebug(ctx context.Context, repo string, in *structs.TempDeployPayload) error
		CancelVerify(ctx context.Context, repo string, in *structs.TempDeployPayload) error
		// BindCI(ctx context.Context, repo string, in *structs.CIBindDeployOption) error
		UpdateDeployment(ctx context.Context, repo, ClusterName string, replicas int32) error

		// pod
		ListPods(ctx context.Context, repo, cluster string, opt QueryOption) (*structs.PodList, *Response, error)
		ListSimplePods(ctx context.Context, repo, cluster string, opt QueryOption) ([]string, *Response, error)
		ListPodEvents(ctx context.Context, repo, cluster, pod string, opt QueryOption) (*structs.EventList, *Response, error)
		PodTerminalWs(ctx context.Context, repo string, opt QueryOption) error
		PodLogs(ctx context.Context, repo string, opt QueryOption) error

		// metrics
		ListMetricsPromData(ctx context.Context, repo, cluster string, in *structs.PromQueryOption) ([]byte, *Response, error)

		// user quotas
		ListRepoQuotas(ctx context.Context, repo string) (*structs.ResourceQuotas, *Response, error)
		CreateRepoQuota(ctx context.Context, repo string, in *structs.RepoQuota) (*Response, error)
		UpdateRepoQuota(ctx context.Context, repo string, in *structs.RepoQuota) (*Response, error)
		DeleteRepoQuota(ctx context.Context, repo string) (*Response, error)

		// kube config
		ListKubeConfigs(ctx context.Context, repo string) ([]*structs.KubeConfig, *Response, error)

		// kubernetes hpa
		PutKubeHPA(ctx context.Context, repo string, in *structs.DeployMetainfo) (*Response, error)
		DeleteKubeHPA(ctx context.Context, repo string, in *structs.DeployMetainfo) (*Response, error)

		// kubernetes vpa
		PutKubeVPA(ctx context.Context, repo string, in *structs.DeployMetainfo) (*Response, error)
		DeleteKubeVPA(ctx context.Context, repo string, in *structs.DeployMetainfo) (*Response, error)
	}
)
