// Copyright 2023 GitBundle Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

import (
	hpav2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
)

type KubeConfig struct {
	ClusterName                KubeClusterName          `binding:"Required" form:"cluster_name" json:"cluster_name"`
	ClusterPriority            int64                    `binding:"-" form:"-" json:"cluster_priority"`
	ClusterAllowDebug          bool                     `binding:"-" form:"-" json:"cluster_allow_debug"`
	ClusterMetricSourceTypes   []hpav2.MetricSourceType `form:"cluster_metric_source_types" json:"cluster_metric_source_types"`
	ClusterControlledResources []corev1.ResourceName    `form:"cluster_controlled_resources" json:"cluster_controlled_resources"`
	ClusterLabels              []DeployKV               `form:"cluster_labels" json:"cluster_labels,omitempty"`
	ClusterAnnotations         []DeployKV               `form:"cluster_annotations" json:"cluster_annotations,omitempty"`
}

type DeployKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
