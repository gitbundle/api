// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	api "github.com/gitbundle/api"
	"github.com/gitbundle/api/pkg/structs"
)

type repositoryService struct {
	client *wrapper
}

func (s *repositoryService) Find(ctx context.Context, repo string) (*api.Repository, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out), res, err
}

func (s *repositoryService) FindHook(ctx context.Context, repo string, id string) (*api.Hook, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/hooks/%s", repo, id)
	out := new(hook)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertHook(out), res, err
}

func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*api.Perm, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out).Perm, res, err
}

func (s *repositoryService) List(ctx context.Context, opts api.ListOptions) ([]*api.Repository, *api.Response, error) {
	path := fmt.Sprintf("api/v1/user/repos?%s", encodeListOptions(opts))
	out := []*repository{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertRepositoryList(out), res, err
}

func (s *repositoryService) ListHooks(ctx context.Context, repo string, opts api.ListOptions) ([]*api.Hook, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/hooks?%s", repo, encodeListOptions(opts))
	out := []*hook{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertHookList(out), res, err
}

func (s *repositoryService) ListStatus(ctx context.Context, repo string, ref string, opts api.ListOptions) ([]*api.Status, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/statuses/%s?%s", repo, ref, encodeListOptions(opts))
	out := []*status{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertStatusList(out), res, err
}

func (s *repositoryService) CreateHook(ctx context.Context, repo string, input *api.HookInput) (*api.Hook, *api.Response, error) {
	target, err := url.Parse(input.Target)
	if err != nil {
		return nil, nil, err
	}
	params := target.Query()
	params.Set("secret", input.Secret)
	target.RawQuery = params.Encode()

	path := fmt.Sprintf("api/v1/repos/%s/hooks", repo)
	in := new(hook)
	in.Type = "gitbundle"
	in.Active = true
	in.Config.Secret = input.Secret
	in.Config.ContentType = "json"
	in.Config.URL = target.String()
	in.Events = append(
		input.NativeEvents,
		convertHookEvent(input.Events)...,
	)
	out := new(hook)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertHook(out), res, err
}

func (s *repositoryService) CreateStatus(ctx context.Context, repo string, ref string, input *api.StatusInput) (*api.Status, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/statuses/%s", repo, ref)
	in := &statusInput{
		State:       convertFromState(input.State),
		Context:     input.Label,
		Description: input.Desc,
		TargetURL:   input.Target,
	}
	out := new(status)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertStatus(out), res, err
}

func (s *repositoryService) UpdateHook(ctx context.Context, repo, id string, input *api.HookInput) (*api.Hook, *api.Response, error) {
	target, err := url.Parse(input.Target)
	if err != nil {
		return nil, nil, err
	}
	params := target.Query()
	params.Set("secret", input.Secret)
	target.RawQuery = params.Encode()

	path := fmt.Sprintf("api/v1/repos/%s/hooks/%s", repo, id)
	in := new(hook)
	in.Type = "gitea"
	in.Active = true
	in.Config.Secret = input.Secret
	in.Config.ContentType = "json"
	in.Config.URL = target.String()
	in.Events = append(
		input.NativeEvents,
		convertHookEvent(input.Events)...,
	)
	out := new(hook)
	res, err := s.client.do(ctx, "PATCH", path, in, out)
	return convertHook(out), res, err
}

func (s *repositoryService) DeleteHook(ctx context.Context, repo string, id string) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/hooks/%s", repo, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *repositoryService) CheckCollaborator(ctx context.Context, repo string, collaborator string) (bool, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/collaborators/%s", repo, collaborator)
	res, err := s.client.do(ctx, "GET", path, nil, nil)
	return res.Status == http.StatusNoContent, res, err
}

func (s *repositoryService) CheckCollaboratorPermission(ctx context.Context, repo string, collaborator string) (*structs.RepoCollaboratorPermission, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/collaborators/%s/permission", repo, collaborator)
	out := new(structs.RepoCollaboratorPermission)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return out, res, err
}

func (s *repositoryService) ListTeams(ctx context.Context, repo string) ([]*api.Team, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/teams", repo)
	out := []*team{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertTeamList(out), res, err
}

func (s *repositoryService) FindRequires(ctx context.Context, repo string) (*structs.Requirement, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/requires", repo)
	out := new(structs.Requirement)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return out, res, err
}

func (s *repositoryService) ListClusters(ctx context.Context, repo string, opt api.QueryOption) ([]string, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/cluster-list?%s", repo, encodeQueryOption(opt))
	out := make([]string, 0, 8)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return out, res, err
}

func (s *repositoryService) ListSimplePullRequests(ctx context.Context, repo string, opt api.ListOptions) ([]*structs.SimplePullRequest, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/simple-pull-requests?%s", repo, encodeListOptions(opt))
	out := make([]*structs.SimplePullRequest, 0, 8)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return out, res, err
}

//func (s *repositoryService) FindDeployMetainfo(ctx context.Context, repo string) ([]*structs.DeployMetainfo, *api.Response, error) {
//	path := fmt.Sprintf("api/v1/repos/%s/deploy/metainfo", repo)
//	out := make([]*structs.DeployMetainfo, 0, 8)
//	res, err := s.client.do(ctx, "GET", path, nil, &out)
//	return out, res, err
//}

func (s *repositoryService) CreateDeploy(ctx context.Context, repo string, in *structs.CreateDeployOption) error {
	path := fmt.Sprintf("api/v1/repos/%s/deploy", repo)
	_, err := s.client.do(ctx, "POST", path, in, nil)
	return err
}

func (s *repositoryService) CancelDebug(ctx context.Context, repo string, in *structs.TempDeployPayload) error {
	path := fmt.Sprintf("api/v1/repos/%s/deploy/cancel_debug", repo)
	_, err := s.client.do(ctx, "POST", path, in, nil)
	return err
}

func (s *repositoryService) CancelVerify(ctx context.Context, repo string, in *structs.TempDeployPayload) error {
	path := fmt.Sprintf("api/v1/repos/%s/deploy/cancel_verify", repo)
	_, err := s.client.do(ctx, "POST", path, in, nil)
	return err
}

func (s *repositoryService) BindCI(ctx context.Context, repo string, in *structs.CIBindDeployOption) error {
	path := fmt.Sprintf("api/v1/repos/%s/deploy/bind_ci", repo)
	_, err := s.client.do(ctx, "POST", path, in, nil)
	return err
}

func (s *repositoryService) UpdateDeployment(ctx context.Context, repo, clusterName string, replicas int32) error {
	path := fmt.Sprintf("api/v1/repos/%s/deploy/%s?replicas=%d", repo, clusterName, replicas)
	_, err := s.client.do(ctx, "PUT", path, nil, nil)
	return err
}

func (s *repositoryService) ListPods(ctx context.Context, repo, cluster string, opt api.QueryOption) (*structs.PodList, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/pods/%s?%s", repo, cluster, encodeQueryOption(opt))
	out := new(structs.PodList)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return out, res, err
}

func (s *repositoryService) ListSimplePods(ctx context.Context, repo, cluster string, opt api.QueryOption) ([]string, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/pods/simple/%s?%s", repo, cluster, encodeQueryOption(opt))
	out := make([]string, 0, 30)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return out, res, err
}

func (s *repositoryService) ListPodEvents(ctx context.Context, repo, cluster, pod string, opt api.QueryOption) (*structs.EventList, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/pods/events/%s/%s?%s", repo, cluster, pod, encodeQueryOption(opt))
	out := new(structs.EventList)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return out, res, err
}

func (s *repositoryService) PodTerminalWs(ctx context.Context, repo string, opt api.QueryOption) error {
	path := fmt.Sprintf("api/v1/repos/%s/pods/terminal/ws?%s", repo, encodeQueryOption(opt))
	_, err := s.client.do(ctx, "GET", path, nil, nil)
	return err
}

func (s *repositoryService) PodLogs(ctx context.Context, repo string, opt api.QueryOption) error {
	path := fmt.Sprintf("api/v1/repos/%s/pods/logs?%s", repo, encodeQueryOption(opt))
	_, err := s.client.do(ctx, "GET", path, nil, nil)
	return err
}

func (s *repositoryService) ListMetricsPromData(ctx context.Context, repo, cluster string, in *structs.PromQueryOption) ([]byte, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/metrics/prom-data/%s", repo, cluster)
	out := make([]byte, 0, 1024)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return out, res, err
}

func (s *repositoryService) ListRepoQuotas(ctx context.Context, repo string) (*structs.ResourceQuotas, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/quotas", repo)
	out := &structs.ResourceQuotas{
		RepoQuotas: make([]*structs.RepoQuota, 0, 10),
		OrgQuotas:  make(structs.OrgQuotas, 10),
	}
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return out, res, err
}

func (s *repositoryService) CreateRepoQuota(ctx context.Context, repo string, in *structs.RepoQuota) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/quotas", repo)
	return s.client.do(ctx, "POST", path, in, nil)
}

func (s *repositoryService) UpdateRepoQuota(ctx context.Context, repo string, in *structs.RepoQuota) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/quotas", repo)
	return s.client.do(ctx, "PATCH", path, in, nil)
}

func (s *repositoryService) DeleteRepoQuota(ctx context.Context, repo string) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/quotas", repo)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *repositoryService) ListKubeConfigs(ctx context.Context, repo string) ([]*structs.KubeConfig, *api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/kubeconfigs", repo)
	out := make([]*structs.KubeConfig, 0, 10)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return out, res, err
}

func (s *repositoryService) PutKubeHPA(ctx context.Context, repo string, in *structs.DeployMetainfo) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/hpa", repo)
	return s.client.do(ctx, "PUT", path, in, nil)
}

func (s *repositoryService) DeleteKubeHPA(ctx context.Context, repo string, in *structs.DeployMetainfo) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/hpa", repo)
	return s.client.do(ctx, "DELETE", path, in, nil)
}

func (s *repositoryService) PutKubeVPA(ctx context.Context, repo string, in *structs.DeployMetainfo) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/vpa", repo)
	return s.client.do(ctx, "PUT", path, in, nil)
}

func (s *repositoryService) DeleteKubeVPA(ctx context.Context, repo string, in *structs.DeployMetainfo) (*api.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/vpa", repo)
	return s.client.do(ctx, "DELETE", path, in, nil)
}

func encodeQueryOption(o api.QueryOption) string {
	query := url.Values{}
	if o.Pod != "" {
		query.Set("pod", o.Pod)
	}
	if o.Container != "" {
		query.Set("container", o.Container)
	}
	if o.Debugging != "" {
		query.Set("debugging", o.Debugging)
	}
	if o.Cluster != "" {
		query.Set("cluster", o.Cluster)
	}
	if o.Extend != "" {
		query.Set("continue", o.Extend)
	}
	return query.Encode()
}

//
// native data structures
//

type (
	// gitea repository resource.
	repository struct {
		ID            int       `json:"id"`
		Owner         user      `json:"owner"`
		Name          string    `json:"name"`
		FullName      string    `json:"full_name"`
		Private       bool      `json:"private"`
		Fork          bool      `json:"fork"`
		HTMLURL       string    `json:"html_url"`
		SSHURL        string    `json:"ssh_url"`
		CloneURL      string    `json:"clone_url"`
		DefaultBranch string    `json:"default_branch"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Permissions   perm      `json:"permissions"`
		Archived      bool      `json:"archived"`
	}

	// gitea permissions details.
	perm struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	}

	// gitea hook resource.
	hook struct {
		ID     int        `json:"id"`
		Type   string     `json:"type"`
		Events []string   `json:"events"`
		Active bool       `json:"active"`
		Config hookConfig `json:"config"`
	}

	// gitea hook configuration details.
	hookConfig struct {
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
		Secret      string `json:"secret"`
	}

	// gitea status resource.
	status struct {
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		State       string    `json:"status"`
		TargetURL   string    `json:"target_url"`
		Description string    `json:"description"`
		Context     string    `json:"context"`
	}

	// gitea status creation request.
	statusInput struct {
		State       string `json:"state"`
		TargetURL   string `json:"target_url"`
		Description string `json:"description"`
		Context     string `json:"context"`
	}

	team struct {
		ID                      int64    `json:"id"`
		CanCreateOrgRepo        bool     `json:"can_create_org_repo"`
		Description             string   `json:"description"`
		IncludesAllRepositories bool     `json:"includes_all_repositories"`
		Name                    string   `json:"name"`
		Permission              string   `json:"permission"`
		Units                   []string `json:"units"`
		UnitsMap                unitsMap `json:"units_map"`
	}
	unitsMap struct {
		Code      string `json:"repo.code"`
		ExtIssues string `json:"repo.ext_issues"`
		ExtWiki   string `json:"repo.ext_wiki"`
		Issues    string `json:"repo.issues"`
		Packages  string `json:"repo.packages"`
		Projects  string `json:"repo.projects"`
		Pulls     string `json:"repo.pulls"`
		Releases  string `json:"repo.releases"`
		Wiki      string `json:"repo.wiki"`
	}
)

//
// native data structure conversion
//

func convertUnitsMap(from unitsMap) api.UnitsMap {
	return api.UnitsMap{
		Code:      from.Code,
		ExtIssues: from.ExtIssues,
		ExtWiki:   from.ExtWiki,
		Issues:    from.Issues,
		Packages:  from.Packages,
		Projects:  from.Projects,
		Pulls:     from.Pulls,
		Releases:  from.Releases,
		Wiki:      from.Wiki,
	}
}

func convertTeam(from *team) *api.Team {
	return &api.Team{
		ID:                      from.ID,
		CanCreateOrgRepo:        from.CanCreateOrgRepo,
		Description:             from.Description,
		IncludesAllRepositories: from.IncludesAllRepositories,
		Name:                    from.Name,
		Permission:              from.Permission,
		Units:                   from.Units,
		UnitsMap:                convertUnitsMap(from.UnitsMap),
	}
}

func convertTeamList(src []*team) []*api.Team {
	var dst []*api.Team
	for _, v := range src {
		dst = append(dst, convertTeam(v))
	}
	return dst
}

func convertRepositoryList(src []*repository) []*api.Repository {
	var dst []*api.Repository
	for _, v := range src {
		dst = append(dst, convertRepository(v))
	}
	return dst
}

func convertRepository(src *repository) *api.Repository {
	return &api.Repository{
		ID:        strconv.Itoa(src.ID),
		Namespace: userLogin(&src.Owner),
		Name:      src.Name,
		Perm:      convertPerm(src.Permissions),
		Branch:    src.DefaultBranch,
		Private:   src.Private,
		Clone:     src.CloneURL,
		CloneSSH:  src.SSHURL,
		Link:      src.HTMLURL,
		Archived:  src.Archived,
	}
}

func convertPerm(src perm) *api.Perm {
	return &api.Perm{
		Push:  src.Push,
		Pull:  src.Pull,
		Admin: src.Admin,
	}
}

func convertHookList(src []*hook) []*api.Hook {
	var dst []*api.Hook
	for _, v := range src {
		dst = append(dst, convertHook(v))
	}
	return dst
}

func convertHook(from *hook) *api.Hook {
	return &api.Hook{
		ID:     strconv.Itoa(from.ID),
		Active: from.Active,
		Target: from.Config.URL,
		Events: from.Events,
	}
}

func convertHookEvent(from api.HookEvents) []string {
	var events []string
	if from.PullRequest {
		events = append(events, "pull_request")
	}
	if from.Issue {
		events = append(events, "issues")
	}
	if from.IssueComment || from.PullRequestComment {
		events = append(events, "issue_comment")
	}
	if from.Branch || from.Tag {
		events = append(events, "create")
		events = append(events, "delete")
	}
	if from.Push {
		events = append(events, "push")
	}
	return events
}

func convertStatusList(src []*status) []*api.Status {
	var dst []*api.Status
	for _, v := range src {
		dst = append(dst, convertStatus(v))
	}
	return dst
}

func convertStatus(from *status) *api.Status {
	return &api.Status{
		State:  convertState(from.State),
		Label:  from.Context,
		Desc:   from.Description,
		Target: from.TargetURL,
	}
}

func convertState(from string) api.State {
	switch from {
	case "error":
		return api.StateError
	case "failure":
		return api.StateFailure
	case "pending":
		return api.StatePending
	case "success":
		return api.StateSuccess
	default:
		return api.StateUnknown
	}
}

func convertFromState(from api.State) string {
	switch from {
	case api.StatePending, api.StateRunning:
		return "pending"
	case api.StateSuccess:
		return "success"
	case api.StateFailure:
		return "failure"
	default:
		return "error"
	}
}
