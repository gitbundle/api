// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package structs

import (
	"time"
)

// CommitStatus holds a single status of a single Commit
type CommitStatus struct {
	ID          int64             `json:"id"`
	State       CommitStatusState `json:"status"`
	TargetURL   string            `json:"target_url"`
	Description string            `json:"description"`
	URL         string            `json:"url"`
	Context     string            `json:"context"`
	Creator     *User             `json:"creator"`
	// swagger:strfmt date-time
	Created time.Time `json:"created_at"`
	// swagger:strfmt date-time
	Updated time.Time `json:"updated_at"`
}

// CombinedStatus holds the combined state of several statuses for a single commit
type CombinedStatus struct {
	State      CommitStatusState `json:"state"`
	SHA        string            `json:"sha"`
	TotalCount int               `json:"total_count"`
	Statuses   []*CommitStatus   `json:"statuses"`
	Repository *Repository       `json:"repository"`
	CommitURL  string            `json:"commit_url"`
	URL        string            `json:"url"`
}

// CreateStatusOption holds the information needed to create a new CommitStatus for a Commit
type CreateStatusOption struct {
	State       CommitStatusState `json:"state"`
	TargetURL   string            `json:"target_url"`
	Description string            `json:"description"`
	Context     string            `json:"context"`
}
