/*
Copyright 2017 The Kubernetes Authors.

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

package testing

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	crclientset "k8s.io/cluster-registry/pkg/client/clientset_generated/clientset"
)

func TestCreateAndGet(t *testing.T) {
	config, tearDown := StartTestServerOrDie(t)
	defer tearDown()

	clientset, err := crclientset.NewForConfig(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	_, err = clientset.ClusterregistryV1alpha1().Clusters().Create(&v1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
	})

	cluster, err := clientset.ClusterregistryV1alpha1().Clusters().Get("cluster", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cluster == nil {
		t.Fatalf("Expected a cluster, got nil")
	}
	if cluster.Name != "cluster" {
		t.Fatalf("Expected a cluster named 'cluster', got a cluster named '%v.", cluster.Name)
	}
}
