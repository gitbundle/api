// Copyright 2023 GitBundle Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

type Requirement struct {
	HasOwnerPerm     bool `json:"has_owner_perm"`
	HasAdminPerm     bool `json:"has_admin_perm"`
	HasCodeWritePerm bool `json:"has_code_write_perm"`
	HasCodeReadPerm  bool `json:"has_code_read_perm"`
}
