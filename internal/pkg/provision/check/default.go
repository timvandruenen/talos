// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package check

import (
	"context"
	"time"

	"github.com/talos-systems/talos/internal/pkg/conditions"
	"github.com/talos-systems/talos/internal/pkg/provision"
	"github.com/talos-systems/talos/pkg/config/types/v1alpha1/generate"
)

// DefaultClusterChecks returns a set of default Talos cluster readiness checks.
func DefaultClusterChecks() []ClusterCheck {
	return []ClusterCheck{
		// wait for etcd to be healthy on all control plane nodes
		func(cluster provision.ClusterAccess) conditions.Condition {
			return conditions.PollingCondition("etcd to be healthy", func(ctx context.Context) error {
				return ServiceHealthAssertion(ctx, cluster, "etcd", WithNodeTypes(generate.TypeInit, generate.TypeControlPlane))
			}, 5*time.Minute, 5*time.Second)
		},
		// wait for bootkube to finish on init node
		func(cluster provision.ClusterAccess) conditions.Condition {
			return conditions.PollingCondition("bootkube to finish", func(ctx context.Context) error {
				return ServiceStateAssertion(ctx, cluster, "bootkube", "Finished")
			}, 5*time.Minute, 5*time.Second)
		},
		// wait for apid to be ready on all the nodes
		func(cluster provision.ClusterAccess) conditions.Condition {
			return conditions.PollingCondition("apid to be ready", func(ctx context.Context) error {
				return ApidReadyAssertion(ctx, cluster)
			}, 2*time.Minute, 5*time.Second)
		},
		// wait for all the nodes to report in at k8s level
		func(cluster provision.ClusterAccess) conditions.Condition {
			return conditions.PollingCondition("all k8s nodes to report", func(ctx context.Context) error {
				return K8sAllNodesReportedAssertion(ctx, cluster)
			}, 5*time.Minute, 5*time.Second)
		},
		// wait for all the nodes to report ready at k8s level
		func(cluster provision.ClusterAccess) conditions.Condition {
			return conditions.PollingCondition("all k8s nodes to report ready", func(ctx context.Context) error {
				return K8sAllNodesReadyAssertion(ctx, cluster)
			}, 10*time.Minute, 5*time.Second)
		},
		// wait for HA k8s control plane
		func(cluster provision.ClusterAccess) conditions.Condition {
			return conditions.PollingCondition("all master nodes to be part of k8s control plane", func(ctx context.Context) error {
				return K8sFullControlPlaneAssertion(ctx, cluster)
			}, 2*time.Minute, 5*time.Second)
		},
		// wait for kube-proxy to report ready
		func(cluster provision.ClusterAccess) conditions.Condition {
			return conditions.PollingCondition("kube-proxy to report ready", func(ctx context.Context) error {
				return K8sPodReadyAssertion(ctx, cluster, "kube-system", "k8s-app=kube-proxy")
			}, 3*time.Minute, 5*time.Second)
		},
		// wait for kube-dns to report ready
		func(cluster provision.ClusterAccess) conditions.Condition {
			return conditions.PollingCondition("kube-dns to report ready", func(ctx context.Context) error {
				return K8sPodReadyAssertion(ctx, cluster, "kube-system", "k8s-app=kube-dns")
			}, 3*time.Minute, 5*time.Second)
		},
	}
}