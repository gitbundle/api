// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2022 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traverse

import (
	"context"

	api "github.com/gitbundle/api"
)

// Repos returns the full repository list, traversing and
// combining paginated responses if necessary.
func Repos(ctx context.Context, client *api.Client) ([]*api.Repository, error) {
	list := []*api.Repository{}
	opts := api.ListOptions{Size: 100}
	for {
		result, meta, err := client.Repositories.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, src := range result {
			if src != nil {
				list = append(list, src)
			}
		}
		opts.Page = meta.Page.Next
		opts.URL = meta.Page.NextURL

		if opts.Page == 0 && opts.URL == "" {
			break
		}
	}
	return list, nil
}
