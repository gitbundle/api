// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package structs

// AddCollaboratorOption options when adding a user as a collaborator of a repository
type AddCollaboratorOption struct {
	Permission *string `json:"permission"`
}

// RepoCollaboratorPermission to get repository permission for a collaborator
type RepoCollaboratorPermission struct {
	Permission  string `json:"permission"`
	RoleName    string `json:"role_name"`
	IsRepoAdmin bool   `json:"is_repo_admin"`
	User        *User  `json:"user"`
}
