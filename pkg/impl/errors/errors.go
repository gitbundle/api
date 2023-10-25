// Copyright 2023 GitBundle Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"net/http"
)

type APIStatus interface {
	Status() StatusError
}

var knownReasons = map[StatusReason]struct{}{
	// StatusReasonUnknown : {}
	StatusReasonUnauthorized:          {},
	StatusReasonForbidden:             {},
	StatusReasonNotFound:              {},
	StatusReasonAlreadyExists:         {},
	StatusReasonConflict:              {},
	StatusReasonGone:                  {},
	StatusReasonInvalid:               {},
	StatusReasonServerTimeout:         {},
	StatusReasonTimeout:               {},
	StatusReasonTooManyRequests:       {},
	StatusReasonBadRequest:            {},
	StatusReasonMethodNotAllowed:      {},
	StatusReasonNotAcceptable:         {},
	StatusReasonRequestEntityTooLarge: {},
	StatusReasonUnsupportedMediaType:  {},
	StatusReasonInternalError:         {},
	StatusReasonExpired:               {},
	StatusReasonServiceUnavailable:    {},
}

func IsInvalid(err error) bool {
	reason, code := reasonAndCodeForError(err)
	if reason == StatusReasonInvalid {
		return true
	}
	if _, ok := knownReasons[reason]; !ok && code == http.StatusUnprocessableEntity {
		return true
	}
	return false
}

func reasonAndCodeForError(err error) (StatusReason, int32) {
	if status := APIStatus(nil); errors.As(err, &status) {
		return status.Status().Reason, status.Status().Code
	}
	return StatusReasonUnknown, 0
}
