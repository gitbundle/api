// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth1

import (
	"context"
	"testing"

	api "github.com/gitbundle/api"
)

func TestContextTokenSource(t *testing.T) {
	source := ContextTokenSource()
	want := new(api.Token)

	ctx := context.Background()
	ctx = context.WithValue(ctx, api.TokenKey{}, want)
	got, err := source.Token(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if got != want {
		t.Errorf("Expect token retrieved from Context")
	}
}

func TestContextTokenSource_Nil(t *testing.T) {
	source := ContextTokenSource()

	ctx := context.Background()
	token, err := source.Token(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if token != nil {
		t.Errorf("Expect nil token from Context")
	}
}
