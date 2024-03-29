// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package structs

import (
	"time"
)

// EditReactionOption contain the reaction type
type EditReactionOption struct {
	Reaction string `json:"content"`
}

// Reaction contain one reaction
type Reaction struct {
	User     *User  `json:"user"`
	Reaction string `json:"content"`
	// swagger:strfmt date-time
	Created time.Time `json:"created_at"`
}
