/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcilePersistentVolume reconciles PersistentVolumes
type reconcilePersistentVolume struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &reconcilePersistentVolume{}

func (r *reconcilePersistentVolume) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := log.FromContext(ctx)

	labels := prometheus.Labels{
		"name": request.NamespacedName.Name,
	}

	// Fetch the PersistentVolume from the cache
	pv := &corev1.PersistentVolume{}
	if err := r.client.Get(ctx, request.NamespacedName, pv); err != nil {
		if errors.IsNotFound(err) {
			log.Info("removing metric")

			log.Info(
				fmt.Sprintf("%d metrics was deleted",
					metric.DeletePartialMatch(labels)))

			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, fmt.Errorf("could not fetch PersistentVolume: %+v", err)
	}

	// Print the PersistentVolume
	log.Info("adding metric")

	if csi := pv.Spec.CSI; csi != nil {
		labels["csi_driver"] = csi.Driver
		labels["csi_volume_handle"] = csi.VolumeHandle
		labels["csi_fs_type"] = csi.FSType
	}
	metric.With(labels).Set(1)

	return reconcile.Result{}, nil
}
