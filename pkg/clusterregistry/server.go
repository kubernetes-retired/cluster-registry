/*
Copyright 2014 The Kubernetes Authors.

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

package clusterregistry

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/install"
	clusterregistryv1alpha1 "k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	clientset "k8s.io/cluster-registry/pkg/client/clientset_generated/internalclientset"
	informers "k8s.io/cluster-registry/pkg/client/informers_generated/internalversion"
	"k8s.io/cluster-registry/pkg/clusterregistry/options"
)

// Run runs the cluster registry API server. It only returns if stopCh is closed
// or one of the ports cannot be listened on initially.
func Run(s *options.ServerRunOptions, stopCh <-chan struct{}) error {
	err := NonBlockingRun(s, stopCh)
	if err != nil {
		return err
	}
	<-stopCh
	return nil
}

// NonBlockingRun runs the cluster registry API server and configures it to
// stop with the given channel.
func NonBlockingRun(s *options.ServerRunOptions, stopCh <-chan struct{}) error {
	server, err := CreateServer(s)
	if err != nil {
		return err
	}

	return server.PrepareRun().NonBlockingRun(stopCh)
}

// CreateServer creates a cluster registry API server.
func CreateServer(s *options.ServerRunOptions) (*genericapiserver.GenericAPIServer, error) {
	// set defaults
	if err := s.GenericServerRunOptions.DefaultAdvertiseAddress(s.SecureServing); err != nil {
		return nil, err
	}

	if err := s.SecureServing.MaybeDefaultWithSelfSignedCerts(s.GenericServerRunOptions.AdvertiseAddress.String(), nil, nil); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	if errs := s.Validate(); len(errs) != 0 {
		return nil, utilerrors.NewAggregate(errs)
	}

	genericConfig := genericapiserver.NewConfig(install.Codecs)
	if err := s.GenericServerRunOptions.ApplyTo(genericConfig); err != nil {
		return nil, err
	}
	if err := s.SecureServing.ApplyTo(genericConfig); err != nil {
		return nil, err
	}
	if err := s.Authentication.ApplyTo(genericConfig); err != nil {
		return nil, err
	}
	if err := s.Authorization.ApplyTo(genericConfig); err != nil {
		return nil, err
	}
	if err := s.Audit.ApplyTo(genericConfig); err != nil {
		return nil, err
	}
	if err := s.Features.ApplyTo(genericConfig); err != nil {
		return nil, err
	}

	resourceConfig := defaultResourceConfig()

	if s.Etcd.StorageConfig.DeserializationCacheSize == 0 {
		// When size of cache is not explicitly set, set it to 50000
		s.Etcd.StorageConfig.DeserializationCacheSize = 50000
	}

	storageFactory := serverstorage.NewDefaultStorageFactory(
		s.Etcd.StorageConfig, s.Etcd.DefaultStorageMediaType, install.Codecs,
		serverstorage.NewDefaultResourceEncodingConfig(install.Registry),
		resourceConfig, nil,
	)

	for _, override := range s.Etcd.EtcdServersOverrides {
		tokens := strings.Split(override, "#")
		if len(tokens) != 2 {
			glog.Errorf("invalid value of etcd server overrides: %s", override)
			continue
		}

		apiresource := strings.Split(tokens[0], "/")
		if len(apiresource) != 2 {
			glog.Errorf("invalid resource definition: %s", tokens[0])
			continue
		}
		group := apiresource[0]
		resource := apiresource[1]
		groupResource := schema.GroupResource{Group: group, Resource: resource}

		servers := strings.Split(tokens[1], ";")
		storageFactory.SetEtcdLocation(groupResource, servers)
	}
	if err := s.Etcd.ApplyWithStorageFactoryTo(storageFactory, genericConfig); err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(genericConfig.LoopbackClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	genericConfig.Version = &version.Info{
		Major: "0",
		Minor: "1",
	}

	var securityDefinitions *spec.SecurityDefinitions
	if s.UseDelegatedAuth {
		config, err := s.Authentication.ToDelegatedAuthenticationConfig()
		if err != nil {
			return nil, err
		}
		authenticator, innerSecurityDefinitions, err := config.New()
		if err != nil {
			return nil, err
		}
		genericConfig.Authenticator = authenticator
		securityDefinitions = innerSecurityDefinitions

		authorizationConfig, err := s.Authorization.ToDelegatedAuthorizationConfig()
		if err != nil {
			return nil, err
		}
		authorizer, err := authorizationConfig.New()
		if err != nil {
			return nil, err
		}
		genericConfig.Authorizer = authorizer
	} else {
		authenticator, innerSecurityDefinitions, err := s.Authentication.ToStandaloneAuthenticationConfig().New()
		if err != nil {
			return nil, err
		}
		genericConfig.Authenticator = authenticator
		securityDefinitions = innerSecurityDefinitions

		authorizer, err := s.Authorization.ToStandaloneAuthorizationConfig().New()
		if err != nil {
			return nil, err
		}
		genericConfig.Authorizer = authorizer
	}

	genericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(clusterregistryv1alpha1.GetOpenAPIDefinitions, install.Scheme)
	genericConfig.OpenAPIConfig.Info.Title = "Cluster Registry"
	genericConfig.OpenAPIConfig.Info.Version = fmt.Sprintf("v%s.%s", genericConfig.Version.Major, genericConfig.Version.Minor)
	genericConfig.OpenAPIConfig.SecurityDefinitions = securityDefinitions

	m, err := genericConfig.Complete(nil).New("clusterregistry", genericapiserver.EmptyDelegate)
	if err != nil {
		return nil, err
	}

	apiResourceConfigSource := storageFactory.APIResourceConfigSource
	installClusterAPIs(m, genericConfig.RESTOptionsGetter, apiResourceConfigSource)

	sharedInformers := informers.NewSharedInformerFactory(client, genericConfig.LoopbackClientConfig.Timeout)
	m.AddPostStartHook("start-informers", func(context genericapiserver.PostStartHookContext) error {
		sharedInformers.Start(context.StopCh)
		return nil
	})
	return m, nil
}

func defaultResourceConfig() *serverstorage.ResourceConfig {
	rc := serverstorage.NewResourceConfig()
	rc.EnableVersions(clusterregistryv1alpha1.SchemeGroupVersion)
	return rc
}
