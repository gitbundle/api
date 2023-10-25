// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth1

import (
	"context"

	api "github.com/gitbundle/api"
)

// StaticTokenSource returns a TokenSource that always
// returns the same token. Because the provided token t
// is never refreshed, StaticTokenSource is only useful
// for tokens that never expire.
func StaticTokenSource(t *api.Token) api.TokenSource {
	return staticTokenSource{t}
}

type staticTokenSource struct {
	token *api.Token
}

func (s staticTokenSource) Token(context.Context) (*api.Token, error) {
	return s.token, nil
}

// ContextTokenSource returns a TokenSource that returns
// a token from the http.Request context.
func ContextTokenSource() api.TokenSource {
	return contextTokenSource{}
}

type contextTokenSource struct {
}

func (s contextTokenSource) Token(ctx context.Context) (*api.Token, error) {
	token, _ := ctx.Value(api.TokenKey{}).(*api.Token)
	return token, nil
}
