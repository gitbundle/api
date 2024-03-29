// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package structs

import (
	"time"
)

// Package represents a package
type Package struct {
	ID         int64       `json:"id"`
	Owner      *User       `json:"owner"`
	Repository *Repository `json:"repository"`
	Creator    *User       `json:"creator"`
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	Version    string      `json:"version"`
	// swagger:strfmt date-time
	CreatedAt time.Time `json:"created_at"`
}

// PackageFile represents a package file
type PackageFile struct {
	ID         int64 `json:"id"`
	Size       int64
	Name       string `json:"name"`
	HashMD5    string `json:"md5"`
	HashSHA1   string `json:"sha1"`
	HashSHA256 string `json:"sha256"`
	HashSHA512 string `json:"sha512"`
}
