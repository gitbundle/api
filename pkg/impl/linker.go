// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"
	"fmt"

	api "github.com/gitbundle/api"
)

type linker struct {
	base string
}

// Resource returns a link to the resource.
func (l *linker) Resource(ctx context.Context, repo string, ref api.Reference) (string, error) {
	switch {
	case api.IsTag(ref.Path):
		t := api.TrimRef(ref.Path)
		return fmt.Sprintf("%s%s/src/tag/%s", l.base, repo, t), nil
	case api.IsPullRequest(ref.Path):
		d := api.ExtractPullRequest(ref.Path)
		return fmt.Sprintf("%s%s/pulls/%d", l.base, repo, d), nil
	case ref.Sha == "":
		t := api.TrimRef(ref.Path)
		return fmt.Sprintf("%s%s/src/branch/%s", l.base, repo, t), nil
	default:
		return fmt.Sprintf("%s%s/commit/%s", l.base, repo, ref.Sha), nil
	}
}

// Diff returns a link to the diff.
func (l *linker) Diff(ctx context.Context, repo string, source, target api.Reference) (string, error) {
	if api.IsPullRequest(target.Path) {
		d := api.ExtractPullRequest(target.Path)
		return fmt.Sprintf("%s%s/pulls/%d/files", l.base, repo, d), nil
	}

	s := source.Sha
	t := target.Sha
	if s == "" {
		s = api.TrimRef(source.Path)
	}
	if t == "" {
		t = api.TrimRef(target.Path)
	}

	return fmt.Sprintf("%s%s/compare/%s...%s", l.base, repo, s, t), nil
}
