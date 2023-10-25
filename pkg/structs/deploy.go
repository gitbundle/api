// Copyright 2023 GitBundle Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

type DeployMetainfo struct {
	KubeConfig

	//ClusterMode          string                      `json:"cluster_mode"`
	ClusterReplicas      int64                       `json:"cluster_replicas"`
	ClusterAutoDeploy    bool                        `json:"cluster_auto_deploy"`
	ClusterDebug         bool                        `json:"cluster_debug"`
	ClusterDebugReplicas int64                       `json:"cluster_debug_replicas"`
	ClusterReplicasTemp  []DeployClusterReplicasTemp `json:"cluster_replicas_temp,omitempty"`

	ClusterMaxCpu                  int32 `json:"cluster_max_cpu,omitempty"`
	ClusterContainerRequestsCpu    int32 `json:"cluster_container_requestsCpu,omitempty"`
	ClusterContainerLimitsCpu      int32 `json:"cluster_container_limitsCpu,omitempty"`
	ClusterMaxMemory               int32 `json:"cluster_max_memory,omitempty"`
	ClusterContainerRequestsMemory int32 `json:"cluster_container_requestsMemory,omitempty"`
	ClusterContainerLimitsMemory   int32 `json:"cluster_container_limitsMemory,omitempty"`

	ClusterHPA *string `form:"cluster_hpa" json:"cluster_hpa,omitempty"`
	ClusterVPA *string `form:"cluster_vpa" json:"cluster_vpa,omitempty"`
	// ClusterHPASpec                 *hpav2.HorizontalPodAutoscalerSpec `form:"cluster_hpa_spec" json:"cluster_hpa_spec,omitempty"`
	// ClusterVPASpec                 *vpav1.VerticalPodAutoscalerSpec   `form:"cluster_vpa_spec" json:"cluster_vpa_spec,omitempty"`

	// OrgRemainingCpu    int32 `form:"org_remaining_cpu" json:"org_remaining_cpu"`
	// OrgRemainingMemory int32 `form:"org_remaining_memory" json:"org_remaining_memory"`
}

type DeployClusterReplicasTemp struct {
	Value int64  `json:"value"`
	Key   string `json:"key"`
}

type KubeClusterName string

func (s KubeClusterName) String() string {
	return string(s)
}

type CreateDeployOption struct {
	ClusterName     KubeClusterName `json:"cluster_name"`
	ClusterReplicas int64           `json:"cluster_replicas"`
	Sha             string          `json:"sha"`
	IsDebug         bool            `json:"is_debug"`
	IsVerify        bool            `json:"is_verify"`
	DeployID        int64           `json:"deploy_id"`
	ClusterConfig   DeployMetainfo  `json:"cluster_config"`
	Port            int64           `json:"port"`
	Protocol        Protocol        `json:"protocol"`
}

type TempDeployPayload struct {
	ClusterName KubeClusterName `json:"cluster_name"`
	Repo        string          `json:"repo"`
	DeploySha   string          `json:"deploy_sha"`
	DeployID    int64           `json:"deploy_id"`
	DoerName    string          `json:"doer_name"`
}

type CIBindDeployOption struct {
	BuildID  int64 `json:"build_id"`
	DeployID int64 `json:"deploy_id"`
}

type Protocol string

const DefaultProtocol = Tcp

const (
	Tcp   Protocol = "TCP" // default
	Udp   Protocol = "UDP"
	Sctp  Protocol = "SCTP"
	Http  Protocol = "HTTP"
	Proxy Protocol = "PROXY"
)

func (p Protocol) String() string {
	return string(p)
}

func Protocols() []Protocol {
	return []Protocol{
		Tcp, Udp, Sctp, Http, Proxy,
	}
}
