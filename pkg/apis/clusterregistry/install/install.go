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

// Package install registers the clusterregistry API group and adds its types to
// a scheme.
package install

import (
	"k8s.io/cluster-registry/pkg/apis/clusterregistry"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"

	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	// Scheme is the scheme used for the clusterregistry API.
	Scheme = runtime.NewScheme()
	// Registry is the registry used for the clusterregistry API.
	Registry = registered.NewOrDie("")
	// Codecs is the set of codecs for the clusterregistry API.
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	Install(make(announced.APIGroupFactoryRegistry), Registry, Scheme)
}

// Install registers the clusterregistry API group and adds its types to a
// scheme.
func Install(groupFactoryRegistry announced.APIGroupFactoryRegistry, registry *registered.APIRegistrationManager, scheme *runtime.Scheme) {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  clusterregistry.GroupName,
			VersionPreferenceOrder:     []string{v1alpha1.SchemeGroupVersion.Version},
			AddInternalObjectsToScheme: clusterregistry.AddToScheme,
		},
		announced.VersionToSchemeFunc{
			v1alpha1.SchemeGroupVersion.Version: v1alpha1.AddToScheme,
		},
	).Announce(groupFactoryRegistry).RegisterAndEnable(registry, scheme); err != nil {
		panic(err)
	}
}
