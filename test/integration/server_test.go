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

package integration

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	crclientset "k8s.io/cluster-registry/pkg/client/clientset_generated/clientset"
)

func TestClusterCRUD(t *testing.T) {
	config, tearDown := StartTestServerOrDie(t)
	defer tearDown()

	clientset, err := crclientset.NewForConfig(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	clusterName := "cluster"

	t.Run("Create", func(t *testing.T) {
		testClusterCreate(t, clientset, clusterName)
	})

	t.Run("Get", func(t *testing.T) {
		testClusterGet(t, clientset, clusterName)
	})

	t.Run("Update", func(t *testing.T) {
		testClusterUpdate(t, clientset, clusterName)
	})

	t.Run("Delete", func(t *testing.T) {
		testClusterDelete(t, clientset, clusterName)
	})
}

func testClusterCreate(t *testing.T, clientset *crclientset.Clientset, clusterName string) {
	cluster, err := clientset.ClusterregistryV1alpha1().Clusters("default").Create(&v1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName,
		},
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	} else if cluster == nil {
		t.Fatalf("Expected a cluster, got nil")
	} else if cluster.Name != clusterName {
		t.Fatalf("Expected a cluster named 'cluster', got a cluster named '%v'.", cluster.Name)
	}
}

func testClusterGet(t *testing.T, clientset *crclientset.Clientset, clusterName string) {
	cluster, err := clientset.ClusterregistryV1alpha1().Clusters("default").Get(clusterName,
		metav1.GetOptions{})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	} else if cluster == nil {
		t.Fatalf("Expected a cluster, got nil")
	} else if cluster.Name != clusterName {
		t.Fatalf("Expected a cluster named 'cluster', got a cluster named '%v'.", cluster.Name)
	}
}

func testClusterUpdate(t *testing.T, clientset *crclientset.Clientset, clusterName string) {
	authProviderName := "authProviderName"

	cluster := &v1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName,
		},
		Spec: v1alpha1.ClusterSpec{
			AuthInfo: v1alpha1.AuthInfo{
				Providers: []v1alpha1.AuthProviderConfig{
					{Name: authProviderName},
				},
			},
		},
	}

	cluster, err := clientset.ClusterregistryV1alpha1().Clusters("default").Update(cluster)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	} else if cluster == nil {
		t.Fatalf("Expected a cluster, got nil")
	} else if cluster.Name != clusterName {
		t.Fatalf("Expected a cluster named 'cluster', got a cluster named '%v'.", cluster.Name)
	} else if len(cluster.Spec.AuthInfo.Providers) != 1 || cluster.Spec.AuthInfo.Providers[0].Name != authProviderName {
		t.Fatalf("Expected a cluster auth provider named '%v', got cluster auth provider '%v'",
			authProviderName, cluster.Spec.AuthInfo.Providers[0].Name)
	}
}

func testClusterDelete(t *testing.T, clientset *crclientset.Clientset, clusterName string) {
	err := clientset.ClusterregistryV1alpha1().Clusters("default").Delete(clusterName,
		&metav1.DeleteOptions{})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// We do not expect to find the cluster we just deleted
	_, err = clientset.ClusterregistryV1alpha1().Clusters("default").Get(clusterName, metav1.GetOptions{})

	if err == nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
