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

package cluster

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	clusterregistry "k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
)

func validNewCluster() *clusterregistry.Cluster {
	return &clusterregistry.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "foo",
			ResourceVersion: "4",
			Labels: map[string]string{
				"name": "foo",
			},
		},
	}
}

func invalidNewCluster() *clusterregistry.Cluster {
	// Create a cluster with empty ServerAddressByClientCIDRs (which is a required field).
	return &clusterregistry.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "foo2",
			ResourceVersion: "5",
		},
	}
}

func TestClusterStrategy(t *testing.T) {
	ctx := genericapirequest.NewDefaultContext()
	if Strategy.NamespaceScoped() {
		t.Errorf("Cluster should not be namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("Cluster should not allow create on update")
	}

	cluster := validNewCluster()
	Strategy.PrepareForCreate(ctx, cluster)
	errs := Strategy.Validate(ctx, cluster)
	if len(errs) != 0 {
		t.Errorf("Unexpected error validating %v", errs)
	}

	invalidCluster := invalidNewCluster()
	Strategy.PrepareForUpdate(ctx, invalidCluster, cluster)
	errs = Strategy.ValidateUpdate(ctx, invalidCluster, cluster)
	if len(errs) == 0 {
		t.Errorf("Expected a validation error")
	}
	if cluster.ResourceVersion != "4" {
		t.Errorf("Incoming resource version on update should not be mutated")
	}
}

func TestMatchCluster(t *testing.T) {
	testFieldMap := map[bool][]fields.Set{
		true: {
			{"metadata.name": "foo"},
		},
		false: {
			{"foo": "bar"},
		},
	}

	for expectedResult, fieldSet := range testFieldMap {
		for _, field := range fieldSet {
			m := MatchCluster(labels.Everything(), field.AsSelector())
			_, matchesSingle := m.MatchesSingle()
			if e, a := expectedResult, matchesSingle; e != a {
				t.Errorf("%+v: expected %v, got %v", fieldSet, e, a)
			}
		}
	}
}
