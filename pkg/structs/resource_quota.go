// Copyright 2023 GitBundle Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

type RepoQuota struct {
	ClusterName KubeClusterName `json:"cluster_name,omitempty"`
	MaxCpu      int32           `json:"max_cpu"`
	MaxMemory   int32           `json:"max_memory"`
}

// For kubernetes, cpu_max_quota means allow organization to use the max cpu cores (Kubernetes Cpus) in a kubernetes cluster
// For kubernetes, memory_max_quota means allow organization to use the max memory in a kubernetes cluster
// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#meaning-of-cpu
type OrgQuota struct {
	MaxCpu          int32 `json:"max_cpu,omitempty"` // m units: millicpu, 0 means no cpus allocated, that is no deployments allowed in this organization
	RemainingCpu    int32 `json:"remaining_cpu,omitempty"`
	MaxMemory       int32 `json:"max_memory,omitempty"` // M, Mi units: megabytes, 0 means no memories allocated, that is no deployments allowed in this organization
	RemainingMemory int32 `json:"remaining_memory,omitempty"`
}

type OrgQuotas map[KubeClusterName]*OrgQuota

type ResourceQuotas struct {
	RepoQuotas []*RepoQuota `json:"repo_quotas"`
	OrgQuotas  OrgQuotas    `json:"org_quotas"`
}
