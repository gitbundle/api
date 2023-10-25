// Copyright 2023 GitBundle Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

type PromQueryOption struct {
	Category     string `json:"category"`
	Nodes        string `json:"nodes"`
	Pods         string `json:"pods"`
	Pvc          string `json:"pvc"`
	Selector     string `json:"selector"`
	Ingress      string `json:"ingress"`
	QueryName    string `json:"query_name"`
	RateAccuracy string `json:"rate_accuracy"`
	Start        int64  `json:"start"`
	End          int64  `json:"end"`
}
