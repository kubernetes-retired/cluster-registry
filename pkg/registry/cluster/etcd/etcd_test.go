/*
Copyright 2016 The Kubernetes Authors.

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

package etcd

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	etcdtesting "k8s.io/apiserver/pkg/storage/etcd/testing"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/install"
	clusterregistry "k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	registrytest "k8s.io/cluster-registry/pkg/registry/cluster/etcd/registrytestutil"
)

func newStorage(t *testing.T) (*REST, *etcdtesting.EtcdTestServer) {
	storageConfig, server := registrytest.NewEtcdStorage(t, clusterregistry.GroupName)
	restOptions := generic.RESTOptions{
		StorageConfig:           storageConfig,
		Decorator:               generic.UndecoratedStorage,
		DeleteCollectionWorkers: 1,
		ResourcePrefix:          "clusters",
	}
	storage, _ := NewREST(restOptions, install.Scheme)
	return storage, server
}

func validNewCluster() *clusterregistry.Cluster {
	return &clusterregistry.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
			Labels: map[string]string{
				"name": "foo",
			},
		},
	}
}

func TestCreate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Store).ClusterScope()
	cluster := validNewCluster()
	cluster.ObjectMeta = metav1.ObjectMeta{GenerateName: "foo"}
	test.TestCreate(
		cluster,
		&clusterregistry.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: "-a123-a_"},
		},
	)
}

func TestUpdate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Store).ClusterScope()
	test.TestUpdate(
		// valid
		validNewCluster(),
		// updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*clusterregistry.Cluster)
			object.ObjectMeta.Annotations = map[string]string{"foo": "bar"}
			return object
		},
	)
}

func TestDelete(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Store).ClusterScope().ReturnDeletedObject()
	test.TestDelete(validNewCluster())
}

func TestGet(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Store).ClusterScope()
	test.TestGet(validNewCluster())
}

func TestList(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Store).ClusterScope()
	test.TestList(validNewCluster())
}

//func TestWatch(t *testing.T) {
//	storage, server := newStorage(t)
//	defer server.Terminate(t)
//	test := registrytest.New(t, storage.Store).ClusterScope()
//	test.TestWatch(
//		validNewCluster(),
//		// matching labels
//		[]labels.Set{
//			{"name": "foo"},
//		},
//		// not matching labels
//		[]labels.Set{
//			{"name": "bar"},
//			{"foo": "bar"},
//		},
//		// matching fields
//		[]fields.Set{
//			{"metadata.name": "foo"},
//		},
//		// not matching fields
//		[]fields.Set{
//			{"metadata.name": "bar"},
//		},
//	)
//}
