// Copyright 2023 GitBundle Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

type LanguageStat struct {
	Repo       string  `json:"repo"`
	IsPrimary  bool    `json:"is_primary"`
	Language   string  `json:"language"`
	Percentage float32 `json:"percentage"`
	Size       int64   `json:"size"`
	Color      string  `json:"color"`
}
