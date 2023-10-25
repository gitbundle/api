// Copyright 2023 GitBundle Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PodList struct {
	Data     []*Pod `json:"data"`
	Continue string `json:"continue"`
}

type Pod struct {
	UID             types.UID `json:"uid"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	Ready           string    `json:"ready"`
	ReadyCount      int       `json:"ready_count"`
	ReadyTotal      int       `json:"ready_total"`
	ReadyPercentage int       `json:"ready_percentage"`
	Restarts        string    `json:"restarts"`
	Status          string    `json:"status"`
	Age             string    `json:"age"`

	Cpu           string `json:"cpu"`
	CpuRequests   int64  `json:"cpu_requests"`
	CpuLimits     int64  `json:"cpu_limits"`
	CpuPercentage int    `json:"cpu_percentage"`
	Mem           string `json:"mem"`
	MemRequests   int64  `json:"mem_requests"`
	MemLimits     int64  `json:"mem_limits"`
	MemPercentage int    `json:"mem_percentage"`

	Raw corev1.Pod `json:"raw"`
}

type EventList struct {
	Data     Events `json:"data"`
	Continue string `json:"continue"`
}

type EventMeta struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type EventObject struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type Event struct {
	Metadata      EventMeta   `json:"metadata"`
	Type          string      `json:"type"`
	Reason        string      `json:"reason"`
	Object        EventObject `json:"object"`
	Message       string      `json:"message"`
	LastSeen      string      `json:"last_seen"`
	LastTimestamp time.Time   `json:"last_timestamp"`
}

type Events []*Event

func (s Events) Len() int           { return len(s) }
func (s Events) Less(i, j int) bool { return s[i].LastTimestamp.After(s[j].LastTimestamp) }
func (s Events) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
